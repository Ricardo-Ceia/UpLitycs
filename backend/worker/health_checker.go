package worker

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"statusframe/backend/email"
	"statusframe/db"
	"time"
)

type HealthChecker struct {
	conn       *sql.DB
	interval   time.Duration
	client     *http.Client
	emailClient *email.SESClient
}

func NewHealthChecker(conn *sql.DB, interval time.Duration) *HealthChecker {
	// Initialize email client
	emailClient, err := email.NewSESClient()
	if err != nil {
		log.Printf("⚠️  Warning: Email service not available: %v", err)
		log.Println("   Email alerts will be disabled. Check your AWS SES configuration.")
		emailClient = nil
	} else {
		log.Println("✅ Email service initialized successfully")
	}

	return &HealthChecker{
		conn:       conn,
		interval:   interval,
		emailClient: emailClient,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Start begins the health checking loop
func (hc *HealthChecker) Start() {
	log.Println("🚀 Health checker started - monitoring every", hc.interval)
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	// Run immediately on start
	hc.checkAllUsers()

	for range ticker.C {
		hc.checkAllUsers()
	}
}

func (hc *HealthChecker) checkAllUsers() {
	// Get apps that are due for checking based on their plan's check interval
	query := `
		SELECT a.id, a.user_id, a.app_name, a.slug, a.health_url, a.alerts, u.plan
		FROM apps a
		JOIN users u ON a.user_id = u.id
		WHERE a.health_url != '' 
		  AND a.next_check_at <= NOW()
	`
	rows, err := hc.conn.Query(query)
	if err != nil {
		log.Printf("❌ Error getting apps for health check: %v", err)
		return
	}
	defer rows.Close()

	appCount := 0
	for rows.Next() {
		var appId, userId int
		var appName, slug, healthUrl, alerts, plan string

		err := rows.Scan(&appId, &userId, &appName, &slug, &healthUrl, &alerts, &plan)
		if err != nil {
			log.Printf("❌ Error scanning app: %v", err)
			continue
		}

		appCount++
		go hc.checkAppHealth(appId, userId, appName, slug, healthUrl, alerts, plan)
	}

	if appCount == 0 {
		// Don't log every tick if no apps are due
		return
	}

	log.Printf("🔍 Checking health for %d app(s) due now", appCount)
}

func (hc *HealthChecker) checkAppHealth(appId, userId int, appName, slug, healthUrl, alerts, plan string) {
	startTime := time.Now()

	resp, err := hc.client.Get(healthUrl)

	responseTime := time.Since(startTime).Milliseconds()
	statusCode := 0

	if err != nil {
		log.Printf("❌ %s | App: %s (ID: %d, Plan: %s) | Error: %v",
			healthUrl, appName, appId, plan, err)
		statusCode = 0
	} else {
		defer resp.Body.Close()
		statusCode = resp.StatusCode
	}

	// Derive status from status code
	status := db.GetStatusFromCode(statusCode)

	// Save to database with app_id only (user_id removed from schema)
	query := "INSERT INTO user_status (app_id, status_code, checked_at) VALUES ($1, $2, NOW())"
	_, err = hc.conn.Exec(query, appId, statusCode)
	if err != nil {
		log.Printf("❌ Error saving status check for app %s (ID: %d): %v", appName, appId, err)
	} else {
		// Log with appropriate emoji based on status
		emoji := "✅"
		if status == "down" || status == "error" {
			emoji = "🔴"
		} else if status == "degraded" || status == "client_error" {
			emoji = "🟡"
		}

		log.Printf("%s %s | App: %s (ID: %d, Plan: %s) | Status: %d (%s) | Response: %dms",
			emoji, healthUrl, appName, appId, plan, statusCode, status, responseTime)
	}

	// Update next_check_at based on plan interval
	planInterval := db.GetPlanCheckInterval(plan)
	nextCheck := time.Now().Add(time.Duration(planInterval) * time.Second)

	updateQuery := "UPDATE apps SET next_check_at = $1 WHERE id = $2"
	_, err = hc.conn.Exec(updateQuery, nextCheck, appId)
	if err != nil {
		log.Printf("❌ Error updating next_check_at for app %s: %v", appName, err)
	}

	// Check if we need to send alert (when service goes down)
	// Only send alerts if user has alerts enabled AND plan supports alerts
	planFeatures := db.GetPlanFeatures(plan)
	shouldSendAlert := (status == "down" || status == "error" || statusCode >= 500) &&
		alerts == "y" &&
		planFeatures.EmailAlerts // Check if plan has email alerts enabled

	if shouldSendAlert {
		hc.checkAndSendAppAlert(appId, userId, appName, healthUrl, statusCode, status, plan)
	}
}

func (hc *HealthChecker) checkAndSendAppAlert(appId, userId int, appName, healthUrl string, statusCode int, status string, plan string) {
	// Check if we already sent an alert recently (within last 5 minutes)
	query := "SELECT sent_at FROM alerts WHERE app_id = $1 ORDER BY sent_at DESC LIMIT 1"
	var lastAlert string
	err := hc.conn.QueryRow(query, appId).Scan(&lastAlert)

	if err == nil && lastAlert != "" {
		lastAlertTime, err := time.Parse(time.RFC3339, lastAlert)
		if err == nil && time.Since(lastAlertTime) < 5*time.Minute {
			log.Printf("⏭️  Skipping alert for app %s (ID: %d) - alert sent %s ago",
				appName, appId, time.Since(lastAlertTime).Round(time.Second))
			return // Don't spam alerts
		}
	}

	// Get user email for alert
	var userEmail string
	err = hc.conn.QueryRow("SELECT email FROM users WHERE id = $1", userId).Scan(&userEmail)
	if err != nil {
		log.Printf("❌ Error getting user email for app %s: %v", appName, err)
		return
	}

	// Log the alert (in production, send actual email)
	log.Printf("🚨 ALERT [%s plan]: App %s (%s) - Service %s is %s (HTTP %d)",
		plan, appName, userEmail, healthUrl, status, statusCode)

	// Save alert record (user_id removed from schema)
	_, err = hc.conn.Exec("INSERT INTO alerts (app_id, sent_at) VALUES ($1, NOW())", appId)
	if err != nil {
		log.Printf("❌ Error saving alert for app %d: %v", appId, err)
	}

	// Send email alert using AWS SES
	if hc.emailClient != nil {
		errorMsg := fmt.Sprintf("Service returned HTTP %d (%s)", statusCode, status)
		if statusCode == 0 {
			errorMsg = "Connection failed - service is unreachable"
		}

		alertEmail := email.AlertEmail{
			AppName:      appName,
			HealthURL:    healthUrl,
			StatusCode:   statusCode,
			Status:       status,
			ErrorMessage: errorMsg,
			Timestamp:    time.Now(),
			UserEmail:    userEmail,
			Plan:         plan,
		}

		err = hc.emailClient.SendDowntimeAlert(alertEmail)
		if err != nil {
			log.Printf("❌ Failed to send email alert to %s: %v", userEmail, err)
		} else {
			log.Printf("✅ Email alert sent successfully to %s", userEmail)
		}
	} else {
		log.Printf("⚠️  Email service not available - alert not sent to: %s", userEmail)
	}

	// Get plan features to check if webhooks are enabled
	planFeatures := db.GetPlanFeatures(plan)
	if planFeatures.Webhooks {
		log.Printf("🔗 Webhook would be triggered for business plan user: %s", userEmail)
		// TODO: Implement webhook sending here
	}
}
