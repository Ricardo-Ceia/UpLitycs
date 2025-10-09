package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"statusframe/backend/worker"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHealthChecker_RunImmediateCheck_NoAlert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	hc := worker.NewHealthChecker(db, 1*time.Hour) // long ticker so only immediate run executes

	// HTTP test server that returns 200 OK
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	appID := 1
	userID := 42

	// Expect the apps selection query (apps due for check)
	mock.ExpectQuery("SELECT .*FROM apps").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "app_name", "slug", "health_url", "alerts", "plan"}).
			AddRow(appID, userID, "Test App", "test-slug", ts.URL, "n", "free"))

	// Expect insert into user_status with status 200
	mock.ExpectExec("INSERT INTO user_status").
		WithArgs(appID, 200).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect update next_check_at for the app
	mock.ExpectExec("UPDATE apps SET next_check_at").
		WithArgs(sqlmock.AnyArg(), appID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Start the checker (runs initial check immediately)
	go hc.Start()

	// wait a short time for the immediate run to complete
	time.Sleep(200 * time.Millisecond)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestHealthChecker_RunImmediateCheck_WithAlert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	hc := worker.NewHealthChecker(db, 1*time.Hour) // long ticker so only immediate run executes

	// HTTP test server that returns 500
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	appID := 2
	userID := 99

	// apps selection returns one app due
	mock.ExpectQuery("SELECT .*FROM apps").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "app_name", "slug", "health_url", "alerts", "plan"}).
			AddRow(appID, userID, "Down App", "down-slug", ts.URL, "y", "free"))

	// Expect insert into user_status with status 500
	mock.ExpectExec("INSERT INTO user_status").
		WithArgs(appID, 500).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect update next_check_at for the app
	mock.ExpectExec("UPDATE apps SET next_check_at").
		WithArgs(sqlmock.AnyArg(), appID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// checkAndSendAppAlert flow:
	// - SELECT sent_at FROM alerts -> return zero rows so alert proceeds
	mock.ExpectQuery("SELECT sent_at FROM alerts").
		WithArgs(appID).
		WillReturnRows(sqlmock.NewRows([]string{"sent_at"})) // zero rows

	// - SELECT email FROM users WHERE id = $1 -> return user email
	mock.ExpectQuery("SELECT email FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow("owner@example.test"))

	// - INSERT INTO alerts (...) -> record alert
	mock.ExpectExec("INSERT INTO alerts").
		WithArgs(appID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Start the checker (runs initial check immediately)
	go hc.Start()

	// wait a short time for the immediate run to complete
	time.Sleep(300 * time.Millisecond)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
