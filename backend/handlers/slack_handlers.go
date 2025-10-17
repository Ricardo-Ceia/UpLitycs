package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"statusframe/backend/auth"
	"statusframe/db"
	"strings"
	"time"
)

// SlackIntegration represents Slack integration for a user
type SlackIntegration struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	SlackTeamID      string `json:"slack_team_id"`
	SlackTeamName    string `json:"slack_team_name"`
	SlackChannelID   string `json:"slack_channel_id"`
	SlackChannelName string `json:"slack_channel_name"`
	IsEnabled        bool   `json:"is_enabled"`
	CreatedAt        string `json:"created_at"`
}

// SlackAuthRequest handles OAuth from Slack
type SlackAuthRequest struct {
	Code string `json:"code"`
}

// SlackTeamInfo captures team metadata returned by Slack OAuth
type SlackTeamInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SlackOAuthResponse from Slack after auth
type SlackOAuthResponse struct {
	OK          bool           `json:"ok"`
	Error       string         `json:"error"`
	AccessToken string         `json:"access_token"`
	BotUserID   string         `json:"bot_user_id"`
	AppID       string         `json:"app_id"`
	TeamID      string         `json:"team_id"`
	TeamName    string         `json:"team_name"`
	Team        *SlackTeamInfo `json:"team"`
}

// SlackChannel for channel info
type SlackChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SlackChannelsResponse from Slack API
type SlackChannelsResponse struct {
	OK       bool           `json:"ok"`
	Channels []SlackChannel `json:"channels"`
	Error    string         `json:"error"`
}

// IncidentAlert for sending to Slack
type IncidentAlert struct {
	AppID      int
	AppName    string
	Status     string
	StatusCode int
	Message    string
	Timestamp  time.Time
}

// StartSlackAuthHandler initiates Slack OAuth flow
func (h *Handler) StartSlackAuthHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user has Pro or Business plan
	plan, err := db.GetUserPlan(h.conn, user.Id)
	if err != nil {
		log.Printf("Error fetching plan for user %d during Slack auth start: %v", user.Id, err)
		http.Error(w, "Unable to verify subscription for Slack integration", http.StatusInternalServerError)
		return
	}
	if plan != "pro" && plan != "business" {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": "Slack integration requires Pro or Business plan",
		})
		return
	}

	clientID := os.Getenv("SLACK_CLIENT_ID")
	redirectURI := os.Getenv("SLACK_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/api/slack/callback"
	}

	if clientID == "" {
		http.Error(w, "Slack client ID not configured", http.StatusInternalServerError)
		return
	}

	state, err := generateSlackState()
	if err != nil {
		log.Printf("Error generating Slack OAuth state: %v", err)
		http.Error(w, "Failed to start Slack authentication", http.StatusInternalServerError)
		return
	}

	session, err := auth.Store.Get(r, "auth-session")
	if err != nil {
		log.Printf("Error retrieving session for Slack OAuth: %v", err)
		http.Error(w, "Failed to start Slack authentication", http.StatusInternalServerError)
		return
	}

	session.Values["slack_oauth_state"] = state
	session.Values["slack_oauth_user_id"] = user.Id
	if err := session.Save(r, w); err != nil {
		log.Printf("Error saving Slack OAuth state to session: %v", err)
		http.Error(w, "Failed to start Slack authentication", http.StatusInternalServerError)
		return
	}

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("scope", "chat:write,channels:read")
	params.Set("redirect_uri", redirectURI)
	params.Set("state", state)

	oauthURL := fmt.Sprintf("https://slack.com/oauth/v2/authorize?%s", params.Encode())

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"oauth_url": oauthURL,
	})
}

// SlackCallbackHandler handles Slack OAuth callback
func (h *Handler) SlackCallbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if errParam := query.Get("error"); errParam != "" {
		log.Printf("Slack OAuth returned error: %s", errParam)
		redirectToSettings(w, r, map[string]string{
			"error": "Slack authorization was cancelled",
		})
		return
	}

	code := query.Get("code")
	state := query.Get("state")
	if code == "" || state == "" {
		redirectToSettings(w, r, map[string]string{
			"error": "Missing authorization data from Slack",
		})
		return
	}

	session, err := auth.Store.Get(r, "auth-session")
	if err != nil {
		log.Printf("Error retrieving session during Slack callback: %v", err)
		redirectToSettings(w, r, map[string]string{
			"error": "Session expired. Please try connecting again",
		})
		return
	}

	storedState, _ := session.Values["slack_oauth_state"].(string)
	if storedState == "" || storedState != state {
		log.Printf("Slack OAuth state mismatch: expected %s, got %s", storedState, state)
		redirectToSettings(w, r, map[string]string{
			"error": "Invalid Slack authorization state",
		})
		return
	}

	userID, ok := session.Values["slack_oauth_user_id"].(int)
	if !ok {
		log.Printf("Slack OAuth missing user ID in session")
		redirectToSettings(w, r, map[string]string{
			"error": "Session missing user information. Please try again",
		})
		return
	}

	plan, err := db.GetUserPlan(h.conn, userID)
	if err != nil {
		log.Printf("Error fetching plan for user %d during Slack callback: %v", userID, err)
		redirectToSettings(w, r, map[string]string{
			"error": "Unable to verify subscription for Slack integration",
		})
		return
	}

	if plan != "pro" && plan != "business" {
		redirectToSettings(w, r, map[string]string{
			"error": "Slack integration requires Pro or Business plan",
		})
		return
	}

	// Exchange code for token
	token, teamID, teamName, channelID, channelName, err := h.exchangeSlackCode(code)
	if err != nil {
		log.Printf("Error exchanging Slack code for user %d: %v", userID, err)
		redirectToSettings(w, r, map[string]string{
			"error": "Unable to connect to Slack. Please try again",
		})
		return
	}

	if channelID == "" {
		log.Printf("Slack OAuth completed for user %d but no channel could be resolved", userID)
		redirectToSettings(w, r, map[string]string{
			"error": "Could not find a Slack channel to post alerts. Ensure a public channel exists and try again",
		})
		return
	}

	if _, err := db.SaveSlackIntegration(h.conn, userID, token, teamID, teamName, channelID, channelName); err != nil {
		log.Printf("Error saving Slack integration for user %d: %v", userID, err)
		redirectToSettings(w, r, map[string]string{
			"error": "Failed to save Slack integration",
		})
		return
	}

	delete(session.Values, "slack_oauth_state")
	delete(session.Values, "slack_oauth_user_id")
	if err := session.Save(r, w); err != nil {
		log.Printf("Error clearing Slack OAuth session values: %v", err)
	}

	redirectToSettings(w, r, map[string]string{
		"success": "Slack",
	})
}

// SaveSlackIntegrationHandler saves Slack integration for user
func (h *Handler) SaveSlackIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check plan
	plan, err := db.GetUserPlan(h.conn, user.Id)
	if err != nil {
		log.Printf("Error fetching plan for user %d while saving Slack integration: %v", user.Id, err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to verify subscription for Slack integration",
		})
		return
	}
	if plan != "pro" && plan != "business" {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": "Slack integration requires Pro or Business plan",
		})
		return
	}

	var req struct {
		BotToken    string `json:"bot_token"`
		TeamID      string `json:"team_id"`
		TeamName    string `json:"team_name"`
		ChannelID   string `json:"channel_id"`
		ChannelName string `json:"channel_name"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
		return
	}

	// Save to database
	integration, err := db.SaveSlackIntegration(h.conn, user.Id, req.BotToken, req.TeamID, req.TeamName, req.ChannelID, req.ChannelName)
	if err != nil {
		log.Printf("Error saving Slack integration: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to save integration",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":     true,
		"integration": integration,
		"message":     "Slack integration saved successfully",
	})
}

// GetSlackIntegrationHandler retrieves user's Slack integration
func (h *Handler) GetSlackIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plan, err := db.GetUserPlan(h.conn, user.Id)
	if err != nil {
		log.Printf("Error fetching plan for user %d while loading Slack integration: %v", user.Id, err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to verify subscription for Slack integration",
		})
		return
	}
	if plan != "pro" && plan != "business" {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": "Slack integration requires Pro or Business plan",
		})
		return
	}

	integration, err := db.GetSlackIntegration(h.conn, user.Id)
	if err != nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"integration": nil,
			"message":     "No Slack integration found",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"integration": integration,
	})
}

// DisableSlackIntegrationHandler disables Slack integration
func (h *Handler) DisableSlackIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plan, planErr := db.GetUserPlan(h.conn, user.Id)
	if planErr != nil {
		log.Printf("Error fetching plan for user %d while disabling Slack integration: %v", user.Id, planErr)
	}
	allowedPlan := plan == "pro" || plan == "business"

	err = db.DisableSlackIntegration(h.conn, user.Id)
	if err != nil {
		log.Printf("Error disabling Slack integration: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to disable integration",
		})
		return
	}

	message := "Slack integration disabled"
	if !allowedPlan {
		message = "Slack integration disabled. Upgrade to reconnect."
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": message,
	})
}

// SendSlackAlert sends an incident alert to Slack
func (h *Handler) SendSlackAlert(alert IncidentAlert) error {
	// Get user's Slack integration
	integration, err := db.GetSlackIntegrationByAppID(h.conn, alert.AppID)
	if err != nil || integration == nil || !integration.IsEnabled {
		log.Printf("No active Slack integration for app %d", alert.AppID)
		return nil // Not an error, just no integration
	}

	// Prepare Slack message
	color := "#36a64f" // Green
	switch alert.Status {
	case "down", "error":
		color = "#ff0000" // Red
	case "degraded", "client_error":
		color = "#ffaa00" // Orange
	}

	payload := map[string]interface{}{
		"channel": integration.SlackChannelID,
		"attachments": []map[string]interface{}{
			{
				"fallback": fmt.Sprintf("%s is %s", alert.AppName, alert.Status),
				"color":    color,
				"title":    alert.AppName,
				"fields": []map[string]interface{}{
					{
						"title": "Status",
						"value": alert.Status,
						"short": true,
					},
					{
						"title": "Status Code",
						"value": fmt.Sprintf("%d", alert.StatusCode),
						"short": true,
					},
					{
						"title": "Details",
						"value": alert.Message,
						"short": false,
					},
					{
						"title": "Timestamp",
						"value": alert.Timestamp.Format("2006-01-02 15:04:05 MST"),
						"short": false,
					},
				},
			},
		},
	}

	// Send to Slack
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", integration.SlackBotToken))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Slack message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Slack API error: %s", string(body))
	}

	// Log notification
	db.LogIncidentNotification(h.conn, alert.AppID, "slack", "sent")

	log.Printf("✅ Slack alert sent for app %s", alert.AppName)
	return nil
}

// exchangeSlackCode exchanges auth code for token
func (h *Handler) exchangeSlackCode(code string) (string, string, string, string, string, error) {
	clientID := os.Getenv("SLACK_CLIENT_ID")
	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	redirectURI := os.Getenv("SLACK_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/api/slack/callback"
	}

	if clientID == "" || clientSecret == "" {
		return "", "", "", "", "", fmt.Errorf("Slack credentials not configured")
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", "https://slack.com/api/oauth.v2.access", strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", "", "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", "", "", fmt.Errorf("Slack OAuth HTTP %d: %s", resp.StatusCode, string(body))
	}

	var oauthResp SlackOAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		return "", "", "", "", "", err
	}

	if !oauthResp.OK {
		return "", "", "", "", "", fmt.Errorf("Slack OAuth failed: %s", oauthResp.Error)
	}

	teamID := oauthResp.TeamID
	teamName := oauthResp.TeamName
	if oauthResp.Team != nil {
		if teamID == "" {
			teamID = oauthResp.Team.ID
		}
		if teamName == "" {
			teamName = oauthResp.Team.Name
		}
	}

	channels, err := h.getSlackChannels(oauthResp.AccessToken)
	if err != nil {
		return "", "", "", "", "", err
	}

	var channelID, channelName string
	if len(channels) > 0 {
		for _, ch := range channels {
			if ch.Name == "general" {
				channelID = ch.ID
				channelName = ch.Name
				break
			}
		}
		if channelID == "" {
			channelID = channels[0].ID
			channelName = channels[0].Name
		}
	}

	return oauthResp.AccessToken, teamID, teamName, channelID, channelName, nil
}

// getSlackChannels retrieves list of Slack channels
func (h *Handler) getSlackChannels(token string) ([]SlackChannel, error) {
	req, err := http.NewRequest("GET", "https://slack.com/api/conversations.list", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("types", "public_channel")
	q.Add("limit", "100")
	q.Add("exclude_archived", "true")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Slack channels API error: %s", string(body))
	}

	var channelResp SlackChannelsResponse
	err = json.NewDecoder(resp.Body).Decode(&channelResp)
	if err != nil {
		return nil, err
	}

	if !channelResp.OK {
		return nil, fmt.Errorf("failed to get channels: %s", channelResp.Error)
	}

	return channelResp.Channels, nil
}

// respondJSON writes a JSON response with the provided status code
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// generateSlackState creates a cryptographically secure state parameter for Slack OAuth
func generateSlackState() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer), nil
}

// redirectToSettings sends the user back to the settings page with optional query parameters
func redirectToSettings(w http.ResponseWriter, r *http.Request, params map[string]string) {
	values := url.Values{}
	values.Set("tab", "integrations")
	for key, value := range params {
		if value == "" {
			continue
		}
		values.Set(key, value)
	}

	redirectURL := fmt.Sprintf("/settings?%s", values.Encode())
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
