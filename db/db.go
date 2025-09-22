package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func OpenDB() *sql.DB {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		env("DB_USER", "postgres"),
		env("DB_PASSWORD", "example"),
		env("DB_HOST", "localhost"),
		env("DB_PORT", "5432"),
		env("DB_NAME", "uplytics"),
		env("DB_SSLMODE", "disable"),
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to DATABASE ESTABLISHED ✅:", env("DB_NAME", "uplytics"))
	return conn
}

func PingDB(conn *sql.DB) error {
	return conn.Ping()
}

func GetAllUsers(conn *sql.DB) ([]string, error) {
	rows, err := conn.Query("SELECT username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []string
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func GetURLOfMainPage(conn *sql.DB, username string) (string, error) {
	var homepage string
	err := conn.QueryRow("SELECT homepage FROM users WHERE username=$1", username).Scan(&homepage)
	return homepage, err
}
