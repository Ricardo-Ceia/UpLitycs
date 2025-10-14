package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	Name      string
	AvatarUrl string
	Email     string
	Id        int
	Plan      string
}

type App struct {
	Id        int     `json:"id"`
	UserId    int     `json:"user_id"`
	AppName   string  `json:"app_name"`
	Slug      string  `json:"slug"`
	HealthUrl string  `json:"health_url"`
	Theme     string  `json:"theme"`
	Alerts    string  `json:"alerts"`
	LogoURL   *string `json:"logo_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type AppWithStatus struct {
	App
	Status      string  `json:"status"`
	StatusCode  int     `json:"status_code"`
	Uptime24h   float64 `json:"uptime_24h"`
	LastChecked string  `json:"last_checked"`
}

type LatestStatus struct {
	Status      string
	Status_code int
	CheckedAt   string
}

func OpenDB() *sql.DB {
	connStr := "user=postgres dbname=statusframe password=your_secure_password_here_change_this host=db port=5432 sslmode=disable"

	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to DATABASE ESTABLISHED✅")

	return conn
}

func PingDB(conn *sql.DB) error {
	err := conn.Ping()

	if err != nil {
		return err
	}

	return nil
}

func InsertUser(conn *sql.DB, name, avatarUrl, email string) (int, error) {
	var id int

	// First try to insert, if conflict occurs, get the existing user's ID
	err := conn.QueryRow(
		"INSERT INTO users (username, avatar_url, email) VALUES ($1, $2, $3) ON CONFLICT (email) DO UPDATE SET username = EXCLUDED.username, avatar_url = EXCLUDED.avatar_url RETURNING id",
		name,
		avatarUrl,
		email,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetUserFromContext(conn *sql.DB, ctx context.Context) (User, error) {
	user, ok := ctx.Value("user").(string)
	if !ok {
		return User{}, sql.ErrNoRows
	}

	var u User
	err := conn.QueryRow("SELECT id, username FROM users WHERE username=$1", user).Scan(&u.Id, &u.Name)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func GetUserByEmail(conn *sql.DB, email string) (User, error) {
	var u User
	err := conn.QueryRow("SELECT id, username, email, avatar_url FROM users WHERE email=$1", email).Scan(&u.Id, &u.Name, &u.Email, &u.AvatarUrl)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func GetUserById(conn *sql.DB, id int) (User, error) {
	var u User
	err := conn.QueryRow("SELECT id, username, email, avatar_url FROM users WHERE id=$1", id).Scan(&u.Id, &u.Name, &u.Email, &u.AvatarUrl)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func GetUserIdFromUser(conn *sql.DB, u User) (int, error) {
	var id int
	err := conn.QueryRow("SELECT id FROM users WHERE username=$1", u.Name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetAllStatuses(conn *sql.DB, userID int, page string) ([]string, error) {
	rows, err := conn.Query("SELECT status FROM user_status WHERE user_id=$1 AND page=$2 ORDER BY checked_at DESC", userID, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []string
	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return statuses, nil
}

func GetLatestStatus(conn *sql.DB, userID int, page string) (LatestStatus, error) {
	var status LatestStatus
	err := conn.QueryRow("SELECT status, status_code, checked_at FROM user_status WHERE user_id=$1 AND page=$2 ORDER BY checked_at DESC LIMIT 1", userID, page).Scan(&status.Status, &status.Status_code, &status.CheckedAt)
	if err != nil {
		return LatestStatus{}, err
	}
	return status, nil
}

// StatusCheck represents a health check result
type StatusCheck struct {
	Id         int    `json:"id"`
	UserId     int    `json:"user_id"`
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"` // Derived from status_code
	CheckedAt  string `json:"checked_at"`
}

// GetStatusFromCode derives status text from HTTP status code
func GetStatusFromCode(statusCode int) string {
	if statusCode >= 200 && statusCode < 300 {
		return "up"
	} else if statusCode >= 300 && statusCode < 400 {
		return "degraded"
	} else if statusCode >= 400 && statusCode < 500 {
		return "client_error"
	} else if statusCode >= 500 {
		return "down"
	}
	return "error" // 0 or other invalid codes
}

// InsertStatusCheck records a health check result
func InsertStatusCheck(conn *sql.DB, userId int, statusCode int) error {
	query := `
		INSERT INTO user_status (user_id, status_code, checked_at)
		VALUES ($1, $2, NOW())
	`
	_, err := conn.Exec(query, userId, statusCode)
	return err
}

// GetLatestStatusByUser gets the most recent status check for a user
func GetLatestStatusByUser(conn *sql.DB, userId int) (*StatusCheck, error) {
	var check StatusCheck
	query := `
		SELECT id, user_id, status_code, checked_at
		FROM user_status
		WHERE user_id = $1
		ORDER BY checked_at DESC
		LIMIT 1
	`
	err := conn.QueryRow(query, userId).Scan(
		&check.Id,
		&check.UserId,
		&check.StatusCode,
		&check.CheckedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	// Derive status from status code
	check.Status = GetStatusFromCode(check.StatusCode)
	return &check, err
}

// GetLatestStatusBySlug gets the most recent status check for a user by slug
func GetLatestStatusBySlug(conn *sql.DB, slug string) (*StatusCheck, error) {
	var check StatusCheck
	query := `
		SELECT sc.id, sc.user_id, sc.status_code, sc.checked_at
		FROM user_status sc
		JOIN users u ON sc.user_id = u.id
		WHERE u.slug = $1
		ORDER BY sc.checked_at DESC
		LIMIT 1
	`
	err := conn.QueryRow(query, slug).Scan(
		&check.Id,
		&check.UserId,
		&check.StatusCode,
		&check.CheckedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	// Derive status from status code
	check.Status = GetStatusFromCode(check.StatusCode)
	return &check, err
}

// GetStatusHistory gets recent status checks for a user
func GetStatusHistory(conn *sql.DB, userId int, limit int) ([]StatusCheck, error) {
	query := `
		SELECT id, user_id, status_code, checked_at
		FROM user_status
		WHERE user_id = $1
		ORDER BY checked_at DESC
		LIMIT $2
	`
	rows, err := conn.Query(query, userId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []StatusCheck
	for rows.Next() {
		var check StatusCheck
		err := rows.Scan(
			&check.Id,
			&check.UserId,
			&check.StatusCode,
			&check.CheckedAt,
		)
		if err != nil {
			return nil, err
		}
		// Derive status from status code
		check.Status = GetStatusFromCode(check.StatusCode)
		checks = append(checks, check)
	}
	return checks, nil
}

// GetUptimePercentage calculates uptime for the last N hours
func GetUptimePercentage(conn *sql.DB, userId int, hours int) (float64, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status_code >= 200 AND status_code < 300) as successful
		FROM user_status
		WHERE user_id = $1
		AND checked_at > NOW() - INTERVAL '1 hour' * $2
	`
	var total, successful int
	err := conn.QueryRow(query, userId, hours).Scan(&total, &successful)
	if err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, nil
	}
	return float64(successful) / float64(total) * 100, nil
}

// GetLastAlert gets the time of the last alert for a user
// Note: This is now deprecated, use app-specific alerts instead
func GetLastAlert(conn *sql.DB, userId int) (string, error) {
	var sentAt string
	query := `SELECT sent_at FROM alerts WHERE user_id = $1 ORDER BY sent_at DESC LIMIT 1`
	err := conn.QueryRow(query, userId).Scan(&sentAt)
	return sentAt, err
}

// DailyUptime represents uptime data for a single day
type DailyUptime struct {
	Date             string  `json:"date"`
	UptimePercentage float64 `json:"uptime_percentage"`
	TotalChecks      int     `json:"total_checks"`
	SuccessfulChecks int     `json:"successful_checks"`
}

// GetDailyUptimeHistory gets uptime percentage for each of the last N days
func GetDailyUptimeHistory(conn *sql.DB, userId int, days int) ([]DailyUptime, error) {
	query := `
		SELECT 
			DATE(checked_at) as date,
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status_code >= 200 AND status_code < 300) as successful
		FROM user_status
		WHERE user_id = $1
		AND checked_at > NOW() - INTERVAL '1 day' * $2
		GROUP BY DATE(checked_at)
		ORDER BY date DESC
	`
	rows, err := conn.Query(query, userId, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []DailyUptime
	for rows.Next() {
		var daily DailyUptime
		err := rows.Scan(
			&daily.Date,
			&daily.TotalChecks,
			&daily.SuccessfulChecks,
		)
		if err != nil {
			return nil, err
		}

		// Calculate uptime percentage
		if daily.TotalChecks > 0 {
			daily.UptimePercentage = float64(daily.SuccessfulChecks) / float64(daily.TotalChecks) * 100
		} else {
			daily.UptimePercentage = 0
		}

		history = append(history, daily)
	}
	return history, nil
}

// GetDailyUptimeHistoryBySlug gets uptime history for a user by slug
func GetDailyUptimeHistoryBySlug(conn *sql.DB, slug string, days int) ([]DailyUptime, error) {
	query := `
		SELECT 
			DATE(sc.checked_at) as date,
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE sc.status_code >= 200 AND sc.status_code < 300) as successful
		FROM user_status sc
		JOIN users u ON sc.user_id = u.id
		WHERE u.slug = $1
		AND sc.checked_at > NOW() - INTERVAL '1 day' * $2
		GROUP BY DATE(sc.checked_at)
		ORDER BY date DESC
	`
	rows, err := conn.Query(query, slug, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []DailyUptime
	for rows.Next() {
		var daily DailyUptime
		err := rows.Scan(
			&daily.Date,
			&daily.TotalChecks,
			&daily.SuccessfulChecks,
		)
		if err != nil {
			return nil, err
		}

		// Calculate uptime percentage
		if daily.TotalChecks > 0 {
			daily.UptimePercentage = float64(daily.SuccessfulChecks) / float64(daily.TotalChecks) * 100
		} else {
			daily.UptimePercentage = 0
		}

		history = append(history, daily)
	}
	return history, nil
}

// ========== APP MANAGEMENT ==========

// CreateApp creates a new app for a user
func CreateApp(conn *sql.DB, userId int, appName, slug, healthUrl, theme, alerts string) (int, error) {
	var appId int
	err := conn.QueryRow(
		"INSERT INTO apps (user_id, app_name, slug, health_url, theme, alerts) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		userId, appName, slug, healthUrl, theme, alerts,
	).Scan(&appId)

	if err != nil {
		return 0, err
	}
	return appId, nil
}

// CreateAppWithLogo creates a new app for a user with optional logo URL
func CreateAppWithLogo(conn *sql.DB, userId int, appName, slug, healthUrl, theme, alerts string, logoURL *string) (int, error) {
	var appId int
	err := conn.QueryRow(
		"INSERT INTO apps (user_id, app_name, slug, health_url, theme, alerts, logo_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		userId, appName, slug, healthUrl, theme, alerts, logoURL,
	).Scan(&appId)

	if err != nil {
		return 0, err
	}
	return appId, nil
}

// GetUserApps returns all apps for a user
func GetUserApps(conn *sql.DB, userId int) ([]App, error) {
	rows, err := conn.Query(
		"SELECT id, user_id, app_name, slug, health_url, theme, alerts, logo_url, created_at, updated_at FROM apps WHERE user_id = $1 ORDER BY created_at DESC",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []App
	for rows.Next() {
		var app App
		var updatedAt sql.NullString
		err := rows.Scan(&app.Id, &app.UserId, &app.AppName, &app.Slug, &app.HealthUrl, &app.Theme, &app.Alerts, &app.LogoURL, &app.CreatedAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			app.UpdatedAt = updatedAt.String
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// GetUserAppsWithStatus returns all apps for a user with their current status
func GetUserAppsWithStatus(conn *sql.DB, userId int) ([]AppWithStatus, error) {
	query := `
		SELECT 
			a.id, a.user_id, a.app_name, a.slug, a.health_url, a.theme, a.alerts, a.created_at, a.updated_at, a.logo_url,
			COALESCE(ls.status_code, 0) as status_code,
			ls.checked_at as last_checked,
			COALESCE(uptime.uptime_24h, 0) as uptime_24h
		FROM apps a
		LEFT JOIN LATERAL (
			SELECT status_code, checked_at 
			FROM user_status 
			WHERE app_id = a.id 
			ORDER BY checked_at DESC 
			LIMIT 1
		) ls ON true
		LEFT JOIN LATERAL (
			SELECT 
				ROUND(
					CAST(COUNT(*) FILTER (WHERE status_code >= 200 AND status_code < 300) AS NUMERIC) / 
					NULLIF(COUNT(*), 0) * 100, 
					2
				) as uptime_24h
			FROM user_status
			WHERE app_id = a.id AND checked_at > NOW() - INTERVAL '24 hours'
		) uptime ON true
		WHERE a.user_id = $1
		ORDER BY a.created_at DESC
	`

	rows, err := conn.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []AppWithStatus
	for rows.Next() {
		var app AppWithStatus
		var updatedAt, lastChecked sql.NullString
		var statusCode sql.NullInt64
		var uptime24h sql.NullFloat64

		err := rows.Scan(
			&app.Id, &app.UserId, &app.AppName, &app.Slug, &app.HealthUrl, &app.Theme, &app.Alerts,
			&app.CreatedAt, &updatedAt, &app.LogoURL, &statusCode, &lastChecked, &uptime24h,
		)
		if err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			app.UpdatedAt = updatedAt.String
		}
		if statusCode.Valid {
			app.StatusCode = int(statusCode.Int64)
			app.Status = GetStatusFromCode(app.StatusCode)
		} else {
			app.StatusCode = 0
			app.Status = "unknown"
		}
		if lastChecked.Valid {
			app.LastChecked = lastChecked.String
		}
		if uptime24h.Valid {
			app.Uptime24h = uptime24h.Float64
		}

		apps = append(apps, app)
	}
	return apps, nil
}

// GetAppBySlug returns an app by its slug
func GetAppBySlug(conn *sql.DB, slug string) (*App, error) {
	var app App
	var updatedAt sql.NullString

	err := conn.QueryRow(
		"SELECT id, user_id, app_name, slug, health_url, theme, alerts, logo_url, created_at, updated_at FROM apps WHERE slug = $1",
		slug,
	).Scan(&app.Id, &app.UserId, &app.AppName, &app.Slug, &app.HealthUrl, &app.Theme, &app.Alerts, &app.LogoURL, &app.CreatedAt, &updatedAt)

	if err != nil {
		return nil, err
	}
	if updatedAt.Valid {
		app.UpdatedAt = updatedAt.String
	}
	return &app, nil
}

// DeleteApp deletes an app and all associated data
func DeleteApp(conn *sql.DB, appId, userId int) error {
	// Verify ownership before deleting
	_, err := conn.Exec("DELETE FROM apps WHERE id = $1 AND user_id = $2", appId, userId)
	return err
}

// GetAppCount returns the number of apps a user has
func GetAppCount(conn *sql.DB, userId int) (int, error) {
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM apps WHERE user_id = $1", userId).Scan(&count)
	return count, err
}

// GetUserPlan returns the user's current plan
func GetUserPlan(conn *sql.DB, userId int) (string, error) {
	var plan string
	err := conn.QueryRow("SELECT plan FROM users WHERE id = $1", userId).Scan(&plan)
	if err != nil {
		return "free", err
	}
	return plan, nil
}

// PlanFeatures defines all features and limits for each plan
type PlanFeatures struct {
	MaxMonitors       int
	MinCheckInterval  int // in seconds
	DataRetentionDays int // how many days of historical data to keep
}

// GetPlanFeatures returns all features for a given plan
func GetPlanFeatures(plan string) PlanFeatures {
	features := map[string]PlanFeatures{
		"free": {
			MaxMonitors:       1,
			MinCheckInterval:  300, // 5 minutes
			DataRetentionDays: 7,   // 7 days
		},
		"pro": {
			MaxMonitors:       25,
			MinCheckInterval:  60, // 1 minute
			DataRetentionDays: 30, // 30 days
		},
		"business": {
			MaxMonitors:       100,
			MinCheckInterval:  30, // 30 seconds
			DataRetentionDays: 90, // 90 days
		},
	}

	feature, ok := features[plan]
	if !ok {
		return features["free"] // Default to free plan
	}
	return feature
}

// GetPlanLimit returns the app limit for a plan
func GetPlanLimit(plan string) int {
	return GetPlanFeatures(plan).MaxMonitors
}

// GetPlanCheckInterval returns minimum check interval for a plan
func GetPlanCheckInterval(plan string) int {
	return GetPlanFeatures(plan).MinCheckInterval
}

// ValidateCheckInterval ensures user isn't checking too frequently for their plan
func ValidateCheckInterval(plan string, requestedInterval int) error {
	minInterval := GetPlanCheckInterval(plan)
	if requestedInterval < minInterval {
		return fmt.Errorf("your %s plan allows minimum %d second intervals", plan, minInterval)
	}
	return nil
}

// GetAllAppsForHealthCheck returns all apps that need health checking
func GetAllAppsForHealthCheck(conn *sql.DB) ([]App, error) {
	rows, err := conn.Query(
		"SELECT id, user_id, app_name, slug, health_url, theme, alerts, created_at FROM apps WHERE health_url != '' ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []App
	for rows.Next() {
		var app App
		err := rows.Scan(&app.Id, &app.UserId, &app.AppName, &app.Slug, &app.HealthUrl, &app.Theme, &app.Alerts, &app.CreatedAt)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// GetUserEmailById returns user email by ID (for alerts)
func GetUserEmailById(conn *sql.DB, userId int) (string, error) {
	var email string
	err := conn.QueryRow("SELECT email FROM users WHERE id = $1", userId).Scan(&email)
	return email, err
}

// GetAppById returns an app by its ID
func GetAppById(conn *sql.DB, appId int) (*App, error) {
	var app App
	var updatedAt sql.NullString

	err := conn.QueryRow(
		"SELECT id, user_id, app_name, slug, health_url, theme, alerts, created_at, updated_at FROM apps WHERE id = $1",
		appId,
	).Scan(&app.Id, &app.UserId, &app.AppName, &app.Slug, &app.HealthUrl, &app.Theme, &app.Alerts, &app.CreatedAt, &updatedAt)

	if err != nil {
		return nil, err
	}
	if updatedAt.Valid {
		app.UpdatedAt = updatedAt.String
	}
	return &app, nil
}

// UpdateAppTheme updates the theme for a specific app
func UpdateAppTheme(conn *sql.DB, appId int, theme string) error {
	_, err := conn.Exec(
		"UPDATE apps SET theme = $1, updated_at = NOW() WHERE id = $2",
		theme,
		appId,
	)
	return err
}

// ========== STRIPE SUBSCRIPTION MANAGEMENT ==========

// UpdateUserSubscription updates user's Stripe subscription information
func UpdateUserSubscription(conn *sql.DB, userId int, plan, stripeCustomerId, stripeSubscriptionId string) error {
	_, err := conn.Exec(
		`UPDATE users 
		SET plan = $1, 
		    stripe_customer_id = $2, 
		    stripe_subscription_id = $3,
		    plan_started_at = NOW()
		WHERE id = $4`,
		plan,
		stripeCustomerId,
		stripeSubscriptionId,
		userId,
	)
	return err
}

// CancelUserSubscription reverts user back to free plan
func CancelUserSubscription(conn *sql.DB, userId int) error {
	_, err := conn.Exec(
		`UPDATE users 
		SET plan = 'free', 
		    stripe_subscription_id = NULL,
		    plan_started_at = NOW()
		WHERE id = $1`,
		userId,
	)
	return err
}

// GetUserByStripeCustomerId finds a user by their Stripe customer ID
func GetUserByStripeCustomerId(conn *sql.DB, stripeCustomerId string) (*User, error) {
	var user User
	err := conn.QueryRow(
		"SELECT id, username, email, avatar_url, plan FROM users WHERE stripe_customer_id = $1",
		stripeCustomerId,
	).Scan(&user.Id, &user.Name, &user.Email, &user.AvatarUrl, &user.Plan)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetStripeCustomerId gets the Stripe customer ID for a user
func GetStripeCustomerId(conn *sql.DB, userId int) (string, error) {
	var customerId sql.NullString
	err := conn.QueryRow(
		"SELECT stripe_customer_id FROM users WHERE id = $1",
		userId,
	).Scan(&customerId)

	if err != nil {
		return "", err
	}

	if customerId.Valid {
		return customerId.String, nil
	}
	return "", nil
}

// CleanupOldStatusChecks removes status checks older than the retention period for each plan
func CleanupOldStatusChecks(conn *sql.DB) error {
	// Delete old records for each plan type based on their retention policy
	queries := []struct {
		plan string
		days int
	}{
		{"free", 7},
		{"pro", 30},
		{"business", 90},
	}

	totalDeleted := 0
	for _, q := range queries {
		result, err := conn.Exec(`
			DELETE FROM user_status
			WHERE app_id IN (
				SELECT a.id FROM apps a
				JOIN users u ON a.user_id = u.id
				WHERE u.plan = $1
			)
			AND checked_at < NOW() - INTERVAL '1 day' * $2
		`, q.plan, q.days)

		if err != nil {
			log.Printf("❌ Error cleaning up %s plan data: %v", q.plan, err)
			continue
		}

		rows, _ := result.RowsAffected()
		if rows > 0 {
			log.Printf("🧹 Cleaned up %d old status checks for %s plan (>%d days)", rows, q.plan, q.days)
			totalDeleted += int(rows)
		}
	}

	if totalDeleted > 0 {
		log.Printf("✅ Total cleanup: removed %d old status checks", totalDeleted)
	}

	return nil
}
