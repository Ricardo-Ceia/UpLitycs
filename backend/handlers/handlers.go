package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"uplytics/backend/auth"
	"uplytics/backend/utils"
	"uplytics/db"

	"github.com/markbates/goth/gothic"
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

	conn := h.conn

	user, err := db.GetUserFromContext(conn, r.Context())
	err = db.UpdateUser(conn, user.Id, req.Homepage, req.Alerts)

	if err != nil {
		log.Println("Error updating user on the GoToDashboardHandler", err)
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

func BeginAuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (h *Handler) GetAuthHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		log.Printf("Gothic authentication failed:%v", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	conn := h.conn
	mode := r.URL.Query().Get("mode")

	existingUser, err := db.GetUserByEmail(conn, user.Email)

	if mode == "login" {
		if err != nil {
			log.Printf("User not found during login: %v", err)
			http.Redirect(w, r, "/auth?error=user_not_found", http.StatusFound)
			return
		}
		id := existingUser.Id

		session, _ := auth.Store.Get(r, "auth-session")
		session.Values["user"] = user.Name
		session.Values["userId"] = id
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", http.StatusFound)
	} else {
		var id int
		if err != nil {
			id, err = db.InsertUser(conn, user.Name, user.AvatarURL, user.Email)
			if err != nil {
				log.Printf("Error creating user:%v", err)
				http.Redirect(w, r, "/auth?error=signup_failed", http.StatusTemporaryRedirect)
			}
		} else {
			id = existingUser.Id
		}
		session, _ := auth.Store.Get(r, "auth-session")
		session.Values["user"] = user.Name
		session.Values["userId"] = id
		session.Save(r, w)

		//check if user needs onboarding
		dbUser, _ := db.GetUserById(conn, id)
		if dbUser.Homepage == "" {
			http.Redirect(w, r, "/onboarding", http.StatusFound)
		} else {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		}
	}
}

func (h *Handler) GetUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	// This is called by ProtectedRoute to check authentication
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"authenticated": true}`))
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, "auth-session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
