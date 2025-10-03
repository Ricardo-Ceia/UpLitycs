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
	users, err := db.GetAllUsersForHealthCheck(hc.conn)
	if err != nil {
		log.Printf("❌ Error getting users for health check: %v", err)
		return
	}

	if len(users) == 0 {
		log.Println("⏸️  No users with health URLs configured yet")
		return
	}

	log.Printf("🔍 Checking health for %d user(s)", len(users))

	for _, user := range users {
		if user.HealthUrl == "" {
			continue
		}
		go hc.checkUserHealth(user)
	}
}

func (hc *HealthChecker) checkUserHealth(user db.User) {
	startTime := time.Now()

	resp, err := hc.client.Get(user.HealthUrl)

	responseTime := time.Since(startTime).Milliseconds()
	statusCode := 0

	if err != nil {
		log.Printf("❌ %s | User: %s (ID: %d) | Error: %v",
			user.HealthUrl, user.Name, user.Id, err)
		statusCode = 0
	} else {
		defer resp.Body.Close()
		statusCode = resp.StatusCode
	}

	// Derive status from status code
	status := db.GetStatusFromCode(statusCode)

	// Save to database (simplified - only statusCode)
	err = db.InsertStatusCheck(hc.conn, user.Id, statusCode)
	if err != nil {
		log.Printf("❌ Error saving status check for user %s (ID: %d): %v", user.Name, user.Id, err)
	} else {
		// Log with appropriate emoji based on status
		emoji := "✅"
		if status == "down" || status == "error" {
			emoji = "🔴"
		} else if status == "degraded" || status == "client_error" {
			emoji = "🟡"
		}

		log.Printf("%s %s | User: %s (ID: %d) | Status: %d (%s) | Response: %dms",
			emoji, user.HealthUrl, user.Name, user.Id, statusCode, status, responseTime)
	}

	// Check if we need to send alert (when service goes down)
	if (status == "down" || status == "error" || statusCode >= 500) && user.Alerts == "y" {
		hc.checkAndSendAlert(user, statusCode, status)
	}
}

func (hc *HealthChecker) checkAndSendAlert(user db.User, statusCode int, status string) {
	// Check if we already sent an alert recently (within last 5 minutes)
	lastAlert, err := db.GetLastAlert(hc.conn, user.Id)
	if err == nil && lastAlert != "" {
		lastAlertTime, err := time.Parse(time.RFC3339, lastAlert)
		if err == nil && time.Since(lastAlertTime) < 5*time.Minute {
			log.Printf("⏭️  Skipping alert for user %s (ID: %d) - alert sent %s ago",
				user.Name, user.Id, time.Since(lastAlertTime).Round(time.Second))
			return // Don't spam alerts
		}
	}

	// Log the alert (in production, send actual email)
	log.Printf("🚨 ALERT: User %s (%s) - Service %s is %s (HTTP %d)",
		user.Name, user.Email, user.HealthUrl, status, statusCode)

	// Save alert record (simplified - just tracks when sent)
	err = db.InsertAlert(hc.conn, user.Id)
	if err != nil {
		log.Printf("❌ Error saving alert for user %d: %v", user.Id, err)
	}

	// TODO: Implement actual email sending here
	// For now, we just log it
	log.Printf("📧 Email would be sent to: %s", user.Email)
}
