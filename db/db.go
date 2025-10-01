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
