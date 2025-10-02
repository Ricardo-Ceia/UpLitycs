package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	Name      string
	AvatarUrl string
	Email     string
	HealthUrl string
	Theme     string
	Alerts    string
	Slug      string
	AppName   string
	Id        int
}

type LatestStatus struct {
	Status      string
	Status_code int
	CheckedAt   string
}

func OpenDB() *sql.DB {
	connStr := "user=postgres dbname=uplytics password=example host=localhost port=5432 sslmode=disable"

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
		"INSERT INTO users (username, avatar_url, email, health_url, theme, alerts, slug, app_name) VALUES ($1, $2, $3, '', 'cyberpunk', '', '', '') ON CONFLICT (email) DO UPDATE SET username = EXCLUDED.username, avatar_url = EXCLUDED.avatar_url RETURNING id",
		name,
		avatarUrl,
		email,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateUser(conn *sql.DB, id int, health_url, theme, alerts, slug, appName string) error {
	_, err := conn.Exec(
		"UPDATE users SET health_url = $2, theme = $3, alerts = $4, slug = $5, app_name = $6 WHERE id = $1",
		id,
		health_url,
		theme,
		alerts,
		slug,
		appName,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetUserFromContext(conn *sql.DB, ctx context.Context) (User, error) {
	user, ok := ctx.Value("user").(string)
	if !ok {
		return User{}, sql.ErrNoRows
	}

	var u User
	var homepage, theme, alerts, slug, appName sql.NullString

	err := conn.QueryRow("SELECT id, username, health_url, theme, alerts, slug, app_name FROM users WHERE username=$1", user).Scan(&u.Id, &u.Name, &homepage, &theme, &alerts, &slug, &appName)
	if err != nil {
		return User{}, err
	}

	// Handle NULL values
	u.HealthUrl = homepage.String
	u.Theme = theme.String
	u.Alerts = alerts.String
	u.Slug = slug.String
	u.AppName = appName.String

	return u, nil
}

func GetUserByEmail(conn *sql.DB, email string) (User, error) {
	var u User
	var homepage, theme, alerts, slug, appName sql.NullString

	err := conn.QueryRow("SELECT id, username, health_url, theme, alerts, email, avatar_url, slug, app_name FROM users WHERE email=$1", email).Scan(&u.Id, &u.Name, &homepage, &theme, &alerts, &u.Email, &u.AvatarUrl, &slug, &appName)
	if err != nil {
		return User{}, err
	}

	// Handle NULL values
	u.HealthUrl = homepage.String
	u.Theme = theme.String
	u.Alerts = alerts.String
	u.Slug = slug.String
	u.AppName = appName.String

	return u, nil
}

func GetUserById(conn *sql.DB, id int) (User, error) {
	var u User
	var homepage, theme, alerts, slug, appName sql.NullString

	err := conn.QueryRow("SELECT id, username, health_url, theme, alerts, email, avatar_url, slug, app_name FROM users WHERE id=$1", id).Scan(&u.Id, &u.Name, &homepage, &theme, &alerts, &u.Email, &u.AvatarUrl, &slug, &appName)
	if err != nil {
		return User{}, err
	}

	// Handle NULL values
	u.HealthUrl = homepage.String
	u.Theme = theme.String
	u.Alerts = alerts.String
	u.Slug = slug.String
	u.AppName = appName.String

	return u, nil
}

func GetUserBySlug(conn *sql.DB, slug string) (User, error) {
	var u User
	var homepage, theme, alerts, slugVal, appName sql.NullString

	err := conn.QueryRow("SELECT id, username, health_url, theme, alerts, email, avatar_url, slug, app_name FROM users WHERE slug=$1", slug).Scan(&u.Id, &u.Name, &homepage, &theme, &alerts, &u.Email, &u.AvatarUrl, &slugVal, &appName)
	if err != nil {
		return User{}, err
	}

	// Handle NULL values
	u.HealthUrl = homepage.String
	u.Theme = theme.String
	u.Alerts = alerts.String
	u.Slug = slugVal.String
	u.AppName = appName.String

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

func GetAllUsers(conn *sql.DB) ([]struct {
	Username  string
	HealthUrl string
	Alerts    string
	Id        int
}, error) {
	// Define a slice of the anonymous struct
	users := []struct {
		Username  string
		HealthUrl string
		Alerts    string
		Id        int
	}{}

	rows, err := conn.Query(`SELECT id, username, health_url, alerts FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u struct {
			Username  string
			HealthUrl string
			Alerts    string
			Id        int
		}
		var healthUrl, alerts sql.NullString

		err := rows.Scan(&u.Id, &u.Username, &healthUrl, &alerts)
		if err != nil {
			return nil, err
		}

		// Handle NULL values
		u.HealthUrl = healthUrl.String
		u.Alerts = alerts.String

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserPage(conn *sql.DB, id int) (string, error) {
	var homepage sql.NullString

	err := conn.QueryRow("SELECT health_url FROM users where id=$1", id).Scan(&homepage)

	if err != nil {
		return "", err
	}

	return homepage.String, nil
}

func InsertStatus(conn *sql.DB, userID int, page string, status string, status_code int) error {
	_, err := conn.Exec("INSERT INTO user_status (user_id, page, status, status_code) VALUES ($1, $2, $3, $4)", userID, page, status, status_code)
	if err != nil {
		return err
	}
	return nil
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
	Id             int    `json:"id"`
	UserId         int    `json:"user_id"`
	Endpoint       string `json:"endpoint"`
	StatusCode     int    `json:"status_code"`
	Status         string `json:"status"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	CheckedAt      string `json:"checked_at"`
}

// InsertStatusCheck records a health check result
func InsertStatusCheck(conn *sql.DB, userId int, endpoint string, statusCode int, status string, responseTime int64) error {
	query := `
		INSERT INTO user_status (user_id, page, status_code, status, response_time_ms, checked_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err := conn.Exec(query, userId, endpoint, statusCode, status, responseTime)
	return err
}

// GetLatestStatusByUser gets the most recent status check for a user
func GetLatestStatusByUser(conn *sql.DB, userId int) (*StatusCheck, error) {
	var check StatusCheck
	query := `
		SELECT id, user_id, page, status_code, status, COALESCE(response_time_ms, 0), checked_at
		FROM user_status
		WHERE user_id = $1
		ORDER BY checked_at DESC
		LIMIT 1
	`
	err := conn.QueryRow(query, userId).Scan(
		&check.Id,
		&check.UserId,
		&check.Endpoint,
		&check.StatusCode,
		&check.Status,
		&check.ResponseTimeMs,
		&check.CheckedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &check, err
}

// GetLatestStatusBySlug gets the most recent status check for a user by slug
func GetLatestStatusBySlug(conn *sql.DB, slug string) (*StatusCheck, error) {
	var check StatusCheck
	query := `
		SELECT sc.id, sc.user_id, sc.page, sc.status_code, sc.status, 
		       COALESCE(sc.response_time_ms, 0), sc.checked_at
		FROM user_status sc
		JOIN users u ON sc.user_id = u.id
		WHERE u.slug = $1
		ORDER BY sc.checked_at DESC
		LIMIT 1
	`
	err := conn.QueryRow(query, slug).Scan(
		&check.Id,
		&check.UserId,
		&check.Endpoint,
		&check.StatusCode,
		&check.Status,
		&check.ResponseTimeMs,
		&check.CheckedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &check, err
}

// GetStatusHistory gets recent status checks for a user
func GetStatusHistory(conn *sql.DB, userId int, limit int) ([]StatusCheck, error) {
	query := `
		SELECT id, user_id, page, status_code, status, COALESCE(response_time_ms, 0), checked_at
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
			&check.Endpoint,
			&check.StatusCode,
			&check.Status,
			&check.ResponseTimeMs,
			&check.CheckedAt,
		)
		if err != nil {
			return nil, err
		}
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

// GetAllUsersForHealthCheck gets all users with their health URLs for the health checker
func GetAllUsersForHealthCheck(conn *sql.DB) ([]User, error) {
	query := `SELECT id, username, email, avatar_url, health_url, theme, slug, app_name, alerts FROM users WHERE health_url IS NOT NULL AND health_url != ''`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var healthUrl, theme, slug, appName, alerts sql.NullString
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.AvatarUrl,
			&healthUrl,
			&theme,
			&slug,
			&appName,
			&alerts,
		)
		if err != nil {
			return nil, err
		}
		user.HealthUrl = healthUrl.String
		user.Theme = theme.String
		if !theme.Valid || theme.String == "" {
			user.Theme = "cyberpunk"
		}
		user.Slug = slug.String
		user.AppName = appName.String
		user.Alerts = alerts.String
		users = append(users, user)
	}
	return users, nil
}

// InsertAlert records an alert
func InsertAlert(conn *sql.DB, userId int, status string, statusCode int) error {
	query := `
		INSERT INTO alerts (user_id, status, status_code, sent_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := conn.Exec(query, userId, status, statusCode)
	return err
}

// GetLastAlert gets the time of the last alert for a user
func GetLastAlert(conn *sql.DB, userId int) (string, error) {
	var sentAt string
	query := `SELECT sent_at FROM alerts WHERE user_id = $1 ORDER BY sent_at DESC LIMIT 1`
	err := conn.QueryRow(query, userId).Scan(&sentAt)
	return sentAt, err
}

// DailyUptime represents uptime data for a single day
type DailyUptime struct {
	Date            string  `json:"date"`
	UptimePercentage float64 `json:"uptime_percentage"`
	TotalChecks     int     `json:"total_checks"`
	SuccessfulChecks int    `json:"successful_checks"`
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
