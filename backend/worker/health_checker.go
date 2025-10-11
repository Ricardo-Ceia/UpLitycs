package worker

import (
	"database/sql"
	"log"
	"net/http"
	"statusframe/db"
	"time"
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

	// Start cleanup routine (runs every 24 hours)
	go hc.startCleanupRoutine()

	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	// Run immediately on start
	hc.checkAllUsers()

	for range ticker.C {
		hc.checkAllUsers()
	}
}

// startCleanupRoutine runs data retention cleanup every 24 hours
func (hc *HealthChecker) startCleanupRoutine() {
	// Run cleanup immediately on start
	log.Println("🧹 Starting data retention cleanup routine")
	db.CleanupOldStatusChecks(hc.conn)

	// Then run every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("🧹 Running scheduled data retention cleanup")
		db.CleanupOldStatusChecks(hc.conn)
	}
}

func (hc *HealthChecker) checkAllUsers() {
	// Get apps that are due for checking based on their plan's check interval
	query := `
		SELECT a.id, a.user_id, a.app_name, a.slug, a.health_url, u.plan
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
		var appName, slug, healthUrl, plan string

		err := rows.Scan(&appId, &userId, &appName, &slug, &healthUrl, &plan)
		if err != nil {
			log.Printf("❌ Error scanning app: %v", err)
			continue
		}

		appCount++
		go hc.checkAppHealth(appId, userId, appName, slug, healthUrl, plan)
	}

	if appCount == 0 {
		// Don't log every tick if no apps are due
		return
	}

	log.Printf("🔍 Checking health for %d app(s) due now", appCount)
}

func (hc *HealthChecker) checkAppHealth(appId, userId int, appName, slug, healthUrl, plan string) {
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
}
