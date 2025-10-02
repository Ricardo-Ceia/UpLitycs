package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"uplytics/backend/auth"
	"uplytics/backend/utils"
	"uplytics/db"

	"github.com/go-chi/chi/v5"
)

type OnboardingRequest struct {
	Name     string `json:"name"`
	Homepage string `json:"homepage"`
	Alerts   string `json:"alerts"`
	Theme    string `json:"theme"`
	Slug     string `json:"slug"`
	AppName  string `json:"appName"`
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

	// Set default theme if not provided
	if req.Theme == "" {
		req.Theme = "cyberpunk"
	}

	conn := h.conn

	user, err := db.GetUserFromContext(conn, r.Context())
	err = db.UpdateUser(conn, user.Id, req.Homepage, req.Theme, req.Alerts, req.Slug, req.AppName)

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

func BeginAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Get provider from URL parameter
	provider := chi.URLParam(r, "provider")

	// Fallback: parse path (e.g. /auth/google)
	if provider == "" {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) >= 2 && parts[0] == "auth" {
			provider = parts[1]
		}
	}

	if provider == "" {
		http.Error(w, "you must select a provider", http.StatusBadRequest)
		return
	}

	// Currently only supporting Google
	if provider == "google" {
		auth.BeginGoogleAuth(w, r)
	} else {
		http.Error(w, "unsupported provider", http.StatusBadRequest)
	}
}

// GetAuthHandler handles the OAuth callback
func (h *Handler) GetAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Get provider from URL
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) >= 3 && parts[0] == "auth" {
			provider = parts[1]
		}
	}

	if provider != "google" {
		log.Printf("Unsupported provider: %s", provider)
		http.Redirect(w, r, "/auth?error=unsupported_provider", http.StatusTemporaryRedirect)
		return
	}

	// Handle Google OAuth callback
	userInfo, err := auth.HandleGoogleCallback(w, r)
	if err != nil {
		log.Printf("Google authentication failed: %v", err)
		http.Redirect(w, r, "/auth?error=auth_failed", http.StatusTemporaryRedirect)
		return
	}

	conn := h.conn

	// Try to get existing user by email
	existingUser, err := db.GetUserByEmail(conn, userInfo.Email)

	var id int
	var needsOnboarding bool

	if err != nil {
		// User doesn't exist, create new user
		log.Printf("User not found, creating new user: %s", userInfo.Email)
		id, err = db.InsertUser(conn, userInfo.Name, userInfo.Picture, userInfo.Email)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			http.Redirect(w, r, "/auth?error=signup_failed", http.StatusTemporaryRedirect)
			return
		}
		needsOnboarding = true // New users always need onboarding
		log.Printf("Created new user with ID: %d", id)
	} else {
		// User exists, use existing ID
		id = existingUser.Id
		// Check if existing user needs onboarding (homepage is empty)
		needsOnboarding = (existingUser.HealthUrl == "")
		log.Printf("Existing user found: ID=%d, Homepage='%s', NeedsOnboarding=%t", id, existingUser.HealthUrl, needsOnboarding)
	}

	// Set session for both new and existing users
	session, _ := auth.Store.Get(r, "auth-session")
	session.Values["user"] = userInfo.Name
	session.Values["userId"] = id
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
	}

	// Redirect based on onboarding status
	if needsOnboarding {
		log.Printf("Redirecting to onboarding for user ID: %d", id)
		http.Redirect(w, r, "/onboarding", http.StatusFound)
	} else {
		log.Printf("Redirecting to dashboard for user ID: %d", id)
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}
}

// GetUserStatusHandler checks if user is authenticated
func (h *Handler) GetUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId")
	userName := r.Context().Value("user")

	log.Printf("GetUserStatusHandler: userId=%v, userName=%v", userID, userName)

	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database to fetch theme
	conn := h.conn
	user, err := db.GetUserById(conn, userID.(int))
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"userId":        userID,
		"userName":      userName,
		"theme":         user.Theme,
		"homepage":      user.HealthUrl,
		"slug":          user.Slug,
		"appName":       user.AppName,
	})
}

func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId, ok := session.Values["userId"].(int)

	if !ok || userId == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"userId":        userId,
	})
}

// LogoutHandler clears the user session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, "auth-session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) GetLatestStatusHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")

	if userId == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn := h.conn

	// Get user info
	user, err := db.GetUserById(conn, userId.(int))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.HealthUrl == "" {
		http.Error(w, "Health URL not configured", http.StatusBadRequest)
		return
	}

	// Get latest status check from database
	latestStatus, err := db.GetLatestStatusByUser(conn, user.Id)
	if err != nil {
		log.Printf("Error getting latest status: %v", err)
		http.Error(w, "Error fetching status", http.StatusInternalServerError)
		return
	}

	if latestStatus == nil {
		// No status checks yet - return pending state
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "pending",
			"status_code": 0,
			"checked_at":  "",
			"message":     "Waiting for first health check (runs every 30 seconds)",
		})
		return
	}

	// Get uptime percentage
	uptime, err := db.GetUptimePercentage(conn, user.Id, 24)
	if err != nil {
		log.Printf("Error calculating uptime: %v", err)
		uptime = 0
	}

	response := map[string]interface{}{
		"status_code":      latestStatus.StatusCode,
		"status":           latestStatus.Status,
		"checked_at":       latestStatus.CheckedAt,
		"response_time_ms": latestStatus.ResponseTimeMs,
		"endpoint":         latestStatus.Endpoint,
		"uptime_24h":       uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPublicStatusHandler returns public status page data by slug (NO AUTH REQUIRED)
func (h *Handler) GetPublicStatusHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if slug == "" {
		http.Error(w, "Slug required", http.StatusBadRequest)
		return
	}

	conn := h.conn

	// Get user by slug
	user, err := db.GetUserBySlug(conn, slug)
	if err != nil {
		log.Printf("User not found for slug %s: %v", slug, err)
		http.Error(w, "Status page not found", http.StatusNotFound)
		return
	}

	// Get latest status check from database
	latestStatus, err := db.GetLatestStatusBySlug(conn, slug)
	if err != nil {
		log.Printf("Error getting status for slug %s: %v", slug, err)
		http.Error(w, "Error fetching status", http.StatusInternalServerError)
		return
	}

	if latestStatus == nil {
		// No status checks yet - return pending state
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"app_name":    user.AppName,
			"theme":       user.Theme,
			"endpoint":    user.HealthUrl,
			"status":      "pending",
			"status_code": 0,
			"message":     "Waiting for first health check",
		})
		return
	}

	// Get uptime percentage
	uptime, err := db.GetUptimePercentage(conn, user.Id, 24)
	if err != nil {
		log.Printf("Error calculating uptime: %v", err)
		uptime = 0
	}

	// Get 30-day uptime history for bar graph
	uptimeHistory, err := db.GetDailyUptimeHistoryBySlug(conn, slug, 30)
	if err != nil {
		log.Printf("Error getting uptime history: %v", err)
		uptimeHistory = []db.DailyUptime{}
	}

	// Return public status data (no sensitive info)
	response := map[string]interface{}{
		"app_name":         user.AppName,
		"theme":            user.Theme,
		"endpoint":         latestStatus.Endpoint,
		"status_code":      latestStatus.StatusCode,
		"status":           latestStatus.Status,
		"checked_at":       latestStatus.CheckedAt,
		"response_time_ms": latestStatus.ResponseTimeMs,
		"uptime_24h":       uptime,
		"uptime_history":   uptimeHistory,
		"user_id":          user.Id, // Include user ID for owner detection
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateThemeHandler allows authenticated users to update their theme
func (h *Handler) UpdateThemeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" && r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userId := r.Context().Value("userId")
	if userId == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Theme string `json:"theme"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate theme
	validThemes := []string{"cyberpunk", "matrix", "retro", "minimal"}
	isValid := false
	for _, t := range validThemes {
		if req.Theme == t {
			isValid = true
			break
		}
	}

	if !isValid {
		http.Error(w, "Invalid theme. Valid themes: cyberpunk, matrix, retro, minimal", http.StatusBadRequest)
		return
	}

	conn := h.conn
	user, err := db.GetUserById(conn, userId.(int))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update only the theme, keep other fields the same
	err = db.UpdateUser(conn, user.Id, user.HealthUrl, req.Theme, user.Alerts, user.Slug, user.AppName)
	if err != nil {
		log.Printf("Error updating theme: %v", err)
		http.Error(w, "Failed to update theme", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"theme":   req.Theme,
	})
}
