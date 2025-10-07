package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"statusframe/backend/auth"
	"statusframe/backend/utils"
	"statusframe/db"

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

	// Validate app name and slug are not empty
	if req.AppName == "" {
		http.Error(w, "App name is required", http.StatusBadRequest)
		return
	}

	if req.Slug == "" {
		http.Error(w, "Slug is required", http.StatusBadRequest)
		return
	}

	// Set default theme if not provided
	if req.Theme == "" {
		req.Theme = "cyberpunk"
	}

	conn := h.conn

	user, err := db.GetUserFromContext(conn, r.Context())
	if err != nil {
		log.Println("Error getting user from context", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check plan limit before creating app
	plan, _ := db.GetUserPlan(conn, user.Id)
	planLimit := db.GetPlanLimit(plan)
	appCount, err := db.GetAppCount(conn, user.Id)
	if err != nil {
		log.Println("Error getting app count", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if appCount >= planLimit {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "plan_limit_reached",
			"message": fmt.Sprintf("You've reached your %s plan limit (%d apps)", plan, planLimit),
			"plan":    plan,
			"limit":   planLimit,
		})
		return
	}

	// Create new app
	log.Printf("Creating app: user_id=%d, app_name=%s, slug=%s, health_url=%s, theme=%s, alerts=%s", 
		user.Id, req.AppName, req.Slug, req.Homepage, req.Theme, req.Alerts)
	
	appId, err := db.CreateApp(conn, user.Id, req.AppName, req.Slug, req.Homepage, req.Theme, req.Alerts)
	if err != nil {
		log.Println("Error creating app in GoToDashboardHandler", err)
		if strings.Contains(err.Error(), "duplicate") {
			http.Error(w, "Slug already exists", http.StatusConflict)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully created app with ID: %d", appId)

	// Return success JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "App created successfully",
		"app_id":  appId,
		"slug":    req.Slug,
	})
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
		// Check if existing user needs onboarding (no apps created yet)
		appCount, err := db.GetAppCount(conn, id)
		if err != nil {
			appCount = 0
		}
		needsOnboarding = (appCount == 0)
		log.Printf("Existing user found: ID=%d, AppCount=%d, NeedsOnboarding=%t", id, appCount, needsOnboarding)
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

	// Get user from database to ensure they exist
	conn := h.conn
	_, err := db.GetUserById(conn, userID.(int))
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Fetch user apps to determine default slug/app info in multi-app world
	apps, err := db.GetUserApps(conn, userID.(int))
	if err != nil {
		log.Printf("Error getting user apps: %v", err)
		http.Error(w, "Failed to load apps", http.StatusInternalServerError)
		return
	}

	var defaultSlug string
	var defaultHomepage string
	var defaultAppName string
	var defaultTheme string

	if len(apps) > 0 {
		defaultSlug = apps[0].Slug
		defaultHomepage = apps[0].HealthUrl
		defaultAppName = apps[0].AppName
		defaultTheme = apps[0].Theme
	}

	// Default theme if no apps or app has no theme
	if defaultTheme == "" {
		defaultTheme = "cyberpunk"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"userId":        userID,
		"userName":      userName,
		"theme":         defaultTheme,
		"homepage":      defaultHomepage,
		"slug":          defaultSlug,
		"appName":       defaultAppName,
		"apps":          apps,
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

// GetCurrentResponseTimeHandler pings an endpoint and returns real-time response time
// Does NOT store the result - just for live display on status page
func (h *Handler) GetCurrentResponseTimeHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Slug required", http.StatusBadRequest)
		return
	}

	// Get app by slug to find the health URL
	app, err := db.GetAppBySlug(h.conn, slug)
	if err != nil {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	if app.HealthUrl == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":         "Health URL not configured",
			"response_time": 0,
		})
		return
	}

	// Ping the endpoint and measure response time
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	startTime := time.Now()
	resp, err := client.Get(app.HealthUrl)
	responseTime := time.Since(startTime).Milliseconds()

	statusCode := 0
	if err == nil {
		defer resp.Body.Close()
		statusCode = resp.StatusCode
	}

	// Return real-time response data (not stored in database)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response_time": responseTime,
		"status_code":   statusCode,
		"timestamp":     time.Now().UTC(),
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

	// Get user info to verify they exist
	_, err := db.GetUserById(conn, userId.(int))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get latest status check from database (this is now per-app, this handler is deprecated)
	// TODO: Remove this handler and use app-specific status endpoints
	latestStatus, err := db.GetLatestStatusByUser(conn, userId.(int))
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
	uptime, err := db.GetUptimePercentage(conn, userId.(int), 24)
	if err != nil {
		log.Printf("Error calculating uptime: %v", err)
		uptime = 0
	}

	response := map[string]interface{}{
		"status_code": latestStatus.StatusCode,
		"status":      latestStatus.Status,
		"checked_at":  latestStatus.CheckedAt,
		"uptime_24h":  uptime,
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

	// Get app by slug (now using apps table)
	app, err := db.GetAppBySlug(conn, slug)
	if err != nil {
		log.Printf("App not found for slug %s: %v", slug, err)
		http.Error(w, "Status page not found", http.StatusNotFound)
		return
	}

	// Get latest status check from database using app_id
	query := `
		SELECT status_code, checked_at 
		FROM user_status 
		WHERE app_id = $1 
		ORDER BY checked_at DESC 
		LIMIT 1
	`
	var statusCode int
	var checkedAt string
	err = conn.QueryRow(query, app.Id).Scan(&statusCode, &checkedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// No status checks yet - return pending state
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"app_name":    app.AppName,
				"theme":       app.Theme,
				"endpoint":    app.HealthUrl,
				"status":      "pending",
				"status_code": 0,
				"message":     "Waiting for first health check",
				"user_id":     app.UserId,
			})
			return
		}
		log.Printf("Error getting status for app %s: %v", slug, err)
		http.Error(w, "Error fetching status", http.StatusInternalServerError)
		return
	}

	// Derive status from status code
	status := db.GetStatusFromCode(statusCode)

	// Get uptime percentage for this app
	uptimeQuery := `
		SELECT 
			ROUND(
				CAST(COUNT(*) FILTER (WHERE status_code >= 200 AND status_code < 300) AS NUMERIC) / 
				NULLIF(COUNT(*), 0) * 100, 
				2
			) as uptime_24h
		FROM user_status
		WHERE app_id = $1 AND checked_at > NOW() - INTERVAL '24 hours'
	`
	var uptime float64
	err = conn.QueryRow(uptimeQuery, app.Id).Scan(&uptime)
	if err != nil {
		log.Printf("Error calculating uptime: %v", err)
		uptime = 0
	}

	// Get 30-day uptime history for bar graph
	historyQuery := `
		SELECT 
			DATE(checked_at) as date,
			COUNT(*) as total_checks,
			COUNT(*) FILTER (WHERE status_code >= 200 AND status_code < 300) as successful_checks
		FROM user_status
		WHERE app_id = $1
		AND checked_at > NOW() - INTERVAL '1 day' * $2
		GROUP BY DATE(checked_at)
		ORDER BY date DESC
	`
	rows, err := conn.Query(historyQuery, app.Id, 30)
	if err != nil {
		log.Printf("Error getting uptime history: %v", err)
	}
	defer rows.Close()

	var uptimeHistory []db.DailyUptime
	if rows != nil {
		for rows.Next() {
			var daily db.DailyUptime
			err := rows.Scan(&daily.Date, &daily.TotalChecks, &daily.SuccessfulChecks)
			if err != nil {
				log.Printf("Error scanning uptime history: %v", err)
				continue
			}
			if daily.TotalChecks > 0 {
				daily.UptimePercentage = float64(daily.SuccessfulChecks) / float64(daily.TotalChecks) * 100
			}
			uptimeHistory = append(uptimeHistory, daily)
		}
	}

	// Return public status data (no sensitive info)
	response := map[string]interface{}{
		"app_name":       app.AppName,
		"theme":          app.Theme,
		"status_code":    statusCode,
		"status":         status,
		"checked_at":     checkedAt,
		"uptime_24h":     uptime,
		"uptime_history": uptimeHistory,
		"user_id":        app.UserId, // Include user ID for owner detection
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateThemeHandler allows authenticated users to update their app theme
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
		Slug  string `json:"slug"`  // Add slug to identify which app to update
		AppId int    `json:"app_id"` // Alternative: use app_id
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
	
	// Get the app to update - either by slug or app_id
	var app *db.App
	if req.Slug != "" {
		app, err = db.GetAppBySlug(conn, req.Slug)
		if err != nil {
			http.Error(w, "App not found", http.StatusNotFound)
			return
		}
	} else if req.AppId != 0 {
		app, err = db.GetAppById(conn, req.AppId)
		if err != nil {
			http.Error(w, "App not found", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "Either slug or app_id must be provided", http.StatusBadRequest)
		return
	}

	// Verify the user owns this app
	if app.UserId != userId.(int) {
		http.Error(w, "Unauthorized - you don't own this app", http.StatusForbidden)
		return
	}

	// Update the app's theme
	err = db.UpdateAppTheme(conn, app.Id, req.Theme)
	if err != nil {
		log.Printf("Error updating app theme: %v", err)
		http.Error(w, "Failed to update theme", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"theme":   req.Theme,
	})
}

// ========== MULTI-APP DASHBOARD HANDLERS ==========

// GetUserAppsHandler returns all apps for the authenticated user with their status
func (h *Handler) GetUserAppsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	apps, err := db.GetUserAppsWithStatus(h.conn, user.Id)
	if err != nil {
		log.Printf("Error fetching user apps: %v", err)
		http.Error(w, "Failed to fetch apps", http.StatusInternalServerError)
		return
	}

	// Get user's plan info
	plan, _ := db.GetUserPlan(h.conn, user.Id)
	planLimit := db.GetPlanLimit(plan)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"apps":       apps,
		"plan":       plan,
		"plan_limit": planLimit,
		"app_count":  len(apps),
	})
}

// DeleteAppHandler deletes an app
func (h *Handler) DeleteAppHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appId := chi.URLParam(r, "appId")
	if appId == "" {
		http.Error(w, "App ID required", http.StatusBadRequest)
		return
	}

	var id int
	_, err = fmt.Sscanf(appId, "%d", &id)
	if err != nil {
		http.Error(w, "Invalid app ID", http.StatusBadRequest)
		return
	}

	err = db.DeleteApp(h.conn, id, user.Id)
	if err != nil {
		log.Printf("Error deleting app: %v", err)
		http.Error(w, "Failed to delete app", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "App deleted successfully",
	})
}

// CheckPlanLimitHandler checks if user can add more apps
func (h *Handler) CheckPlanLimitHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plan, _ := db.GetUserPlan(h.conn, user.Id)
	planLimit := db.GetPlanLimit(plan)
	appCount, err := db.GetAppCount(h.conn, user.Id)
	if err != nil {
		log.Printf("Error getting app count: %v", err)
		http.Error(w, "Failed to check limit", http.StatusInternalServerError)
		return
	}

	canAdd := appCount < planLimit

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"can_add":    canAdd,
		"plan":       plan,
		"plan_limit": planLimit,
		"app_count":  appCount,
		"remaining":  planLimit - appCount,
	})
}

// GetPlanFeaturesHandler returns all features for user's current plan
func (h *Handler) GetPlanFeaturesHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plan, _ := db.GetUserPlan(h.conn, user.Id)
	features := db.GetPlanFeatures(plan)
	appCount, _ := db.GetAppCount(h.conn, user.Id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"plan":                plan,
		"max_monitors":        features.MaxMonitors,
		"min_check_interval":  features.MinCheckInterval,
		"webhooks":            features.Webhooks,
		"custom_domain":       features.CustomDomain,
		"ssl_monitoring":      features.SSLMonitoring,
		"api_access":          features.APIAccess,
		"email_alerts":        features.EmailAlerts,
		"max_alerts_per_day":  features.MaxAlertsPerDay,
		"current_app_count":   appCount,
		"remaining_monitors":  features.MaxMonitors - appCount,
	})
}
