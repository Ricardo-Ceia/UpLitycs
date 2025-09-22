package status_checker

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"uplytics/backend/utils"
	"uplytics/db"
)

func GetPageStatus(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return utils.MapStatusCode(resp.StatusCode), nil
}

func UpdateStatuses(conn *sql.DB) {
	users, err := db.GetAllUsers(conn)
	if err != nil {
		log.Println("Error fetching users:", err)
		return
	}

	for _, user := range users {
		page, err := db.GetUserPage(conn, user.Id)
		if err != nil {
			log.Printf("Error fetching pages for user %s: %v", user, err)
			continue
		}

		status, err := GetPageStatus(page)
		if err != nil {
			log.Printf("Error checking statuses for user %s: %v", user, err)
			continue
		}

		err = db.InsertStatus(conn, user.Id, user.Homepage, status)
		if err != nil {
			log.Printf("Error inserting status for user %s, page %s: %v", user.Username, page, err)
		}
		log.Printf("User: %s, Page: %s, Status: %s", user.Username, page, status)
	}
}

func StartStatusUpdater(conn *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			UpdateStatuses(conn)
		}
	}()
}
