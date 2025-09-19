package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"uplytics/backend/utils"
)

type OnboardingRequest struct {
	Name     string `json:"name"`
	Homepage string `json:"homepage"`
	Alerts   string `json:"alerts"`
}

func StartOnboardingHandler(w http.ResponseWriter, r *http.Request) {
	//TODO ADD LOGIN VERIFICATION BEFORE REDIRECTING USER
	log.Println("Reaching the endpoint")
	http.Redirect(w, r, "/onboarding", http.StatusFound)
}

func GoToDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req OnboardingRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !utils.CheckUsername(req.Name) {
		http.Error(w, "Invalid username format", http.StatusBadRequest)
		return
	}

	if !utils.CheckURLFormat(req.Homepage) {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	if !utils.CheckAlerts(req.Alerts) {
		http.Error(w, "Invalid alerts format", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}
