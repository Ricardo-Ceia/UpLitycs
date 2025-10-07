package worker

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"uplytics/db"
)

type HealthChecker struct {
	conn     *sql.DB
	interval time.Duration
	client   *http.Client
}

func NewHealthChecker(conn *sql.DB, interval time.Duration) *HealthChecker {
	return &HealthChecker{
		conn:     conn,
		interval: interval,
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
	// Get all apps instead of users
	query := "SELECT id, user_id, app_name, slug, health_url, alerts FROM apps WHERE health_url != ''"
	rows, err := hc.conn.Query(query)
	if err != nil {
		log.Printf("❌ Error getting apps for health check: %v", err)
		return
	}
	defer rows.Close()

	appCount := 0
	for rows.Next() {
		var appId, userId int
		var appName, slug, healthUrl, alerts string

		err := rows.Scan(&appId, &userId, &appName, &slug, &healthUrl, &alerts)
		if err != nil {
			log.Printf("❌ Error scanning app: %v", err)
			continue
		}

		appCount++
		go hc.checkAppHealth(appId, userId, appName, slug, healthUrl, alerts)
	}

	if appCount == 0 {
		log.Println("⏸️  No apps with health URLs configured yet")
		return
	}

	log.Printf("🔍 Checking health for %d app(s)", appCount)
}

func (hc *HealthChecker) checkUserHealth(user db.User) {
	// This function is deprecated - apps are checked directly now
}

func (hc *HealthChecker) checkAppHealth(appId, userId int, appName, slug, healthUrl, alerts string) {
	startTime := time.Now()

	resp, err := hc.client.Get(healthUrl)

	responseTime := time.Since(startTime).Milliseconds()
	statusCode := 0

	if err != nil {
		log.Printf("❌ %s | App: %s (ID: %d) | Error: %v",
			healthUrl, appName, appId, err)
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

		log.Printf("%s %s | App: %s (ID: %d) | Status: %d (%s) | Response: %dms",
			emoji, healthUrl, appName, appId, statusCode, status, responseTime)
	}

	// Check if we need to send alert (when service goes down)
	if (status == "down" || status == "error" || statusCode >= 500) && alerts == "y" {
		hc.checkAndSendAppAlert(appId, userId, appName, healthUrl, statusCode, status)
	}
}

func (hc *HealthChecker) checkAndSendAlert(user db.User, statusCode int, status string) {
	// Deprecated - use checkAndSendAppAlert
}

func (hc *HealthChecker) checkAndSendAppAlert(appId, userId int, appName, healthUrl string, statusCode int, status string) {
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
	log.Printf("🚨 ALERT: App %s (%s) - Service %s is %s (HTTP %d)",
		appName, userEmail, healthUrl, status, statusCode)

	// Save alert record (user_id removed from schema)
	_, err = hc.conn.Exec("INSERT INTO alerts (app_id, sent_at) VALUES ($1, NOW())", appId)
	if err != nil {
		log.Printf("❌ Error saving alert for app %d: %v", appId, err)
	}

	// TODO: Implement actual email sending here
	// For now, we just log it
	log.Printf("📧 Email would be sent to: %s", userEmail)
}
