package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	Username string
	Homepage string
	Alerts   string
	Id       int
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

	log.Println("Connection to DATABSE ESTABLISHED✅")

	return conn
}

func PingDB(conn *sql.DB) error {
	err := conn.Ping()

	if err != nil {
		return err
	}

	return nil
}

func InsertUser(conn *sql.DB, username, homepage, alerts string) error {
	_, err := conn.Exec("INSERT INTO users (username, homepage, alerts) VALUES ($1, $2, $3)", username, homepage, alerts)
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
	err := conn.QueryRow("SELECT id, username, homepage, alerts FROM users WHERE username=$1", user).Scan(&u.Id, &u.Username, &u.Homepage, &u.Alerts)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func GetUserIdFromUser(conn *sql.DB, u User) (int, error) {
	var id int
	err := conn.QueryRow("SELECT id FROM users WHERE username=$1", u.Username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetAllUsers(conn *sql.DB) ([]struct {
	Username string
	Homepage string
	Alerts   string
	Id       int
}, error) {
	// Define a slice of the anonymous struct
	users := []struct {
		Username string
		Homepage string
		Alerts   string
		Id       int
	}{}

	rows, err := conn.Query(`SELECT id, username, homepage, alerts FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u struct {
			Username string
			Homepage string
			Alerts   string
			Id       int
		}

		err := rows.Scan(&u.Id, &u.Username, &u.Homepage, &u.Alerts)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserPage(conn *sql.DB, id int) (string, error) {
	var homepage string

	err := conn.QueryRow("SELECT homepage FROM users where id=$1", id).Scan(&homepage)

	if err != nil {
		return "", err
	}

	return homepage, nil
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
