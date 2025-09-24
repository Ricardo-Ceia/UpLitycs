package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"uplytics/backend/utils"
	"uplytics/db"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
)

type OnboardingRequest struct {
	Name     string `json:"name"`
	Homepage string `json:"homepage"`
	Alerts   string `json:"alerts"`
}

type Handler struct {
	conn *sql.DB
}

func NewHandler(conn *sql.DB) *Handler {
	return &Handler{conn: conn}
}

func StartOnboardingHandler(w http.ResponseWriter, r *http.Request) {
	//TODO ADD LOGIN VERIFICATION BEFORE REDIRECTING USER
	http.Redirect(w, r, "/onboarding", http.StatusFound)
}

func (h *Handler) GoToDashboardHandler(w http.ResponseWriter, r *http.Request) {
	//TODO ADD LOGIN VERIFICATION BEFORE REDIRECTING USER
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

	err = db.InsertUser(h.conn, req.Name, req.Homepage, req.Alerts)

	if err != nil {
		log.Println("Error inserting user on the GoToDashboardHandler", err)
		if err == sql.ErrNoRows {
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}
	// Return success JSON instead of redirect
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Onboarding completed successfully"}`))
}

func (h *Handler) LatestDataStatusHandler(w http.ResponseWriter, r *http.Request) {
	//TODO:check if user is logged in
	conn := h.conn
	user, err := db.GetUserFromContext(conn, r.Context())

	if err != nil {
		log.Println("Error getting user from context in LatestDataStatusHandler:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := db.GetUserIdFromUser(conn, user)

	if err != nil {
		log.Println("Error getting user ID from user in LatestDataStatusHandler:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	LatestStatus, err := db.GetLatestStatus(conn, id, user.Homepage)

	if err != nil {
		log.Println("Error getting latest status in LatestDataStatusHandler:", err)
		http.Error(w, "Error fetching latest status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LatestStatus)
}


func GetAuthHandler (w http.ResponseWriter,r *http.Request){
	user,	err := gothic.CompleteUserAuth(w,r)

	if err != nil{
		log.Println("Gothic authentication failed:%v",err)
		return
	}

	err = InsertUser(db,user.Name,user.AvatarURL,user.Email)



	http.Redirect(w,r,"http://localhost:3333",http.StatusFound)
}


