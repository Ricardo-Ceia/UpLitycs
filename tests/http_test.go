package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"uplytics/backend/status_checker"
)

func TestGetPageStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	status, err := status_checker.GetPageStatus(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != "up" {
		t.Errorf("expected status up, got %s", status)
	}
}
