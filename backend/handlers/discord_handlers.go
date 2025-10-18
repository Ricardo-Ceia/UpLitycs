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

// DiscordIntegration represents Discord integration for a user
type DiscordIntegrationHandler struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	DiscordUserID   string `json:"discord_user_id"`
	DiscordUsername string `json:"discord_username"`
	ServerID        string `json:"server_id"`
	ServerName      string `json:"server_name"`
	ChannelID       string `json:"channel_id"`
	ChannelName     string `json:"channel_name"`
	IsEnabled       bool   `json:"is_enabled"`
	CreatedAt       string `json:"created_at"`
}

// DiscordAuthRequest handles OAuth from Discord
type DiscordAuthRequest struct {
	Code string `json:"code"`
}

// DiscordOAuthResponse from Discord after auth
type DiscordOAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// DiscordUser represents the user returned by Discord API
type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// DiscordGuild represents a Discord server
type DiscordGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Permissions string `json:"permissions"`
}

// DiscordChannel represents a Discord channel
type DiscordChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

// StartDiscordAuthHandler initiates Discord OAuth flow
func (h *Handler) StartDiscordAuthHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user has Pro or Business plan
	plan, err := db.GetUserPlan(h.conn, user.Id)
	if err != nil {
		log.Printf("Error fetching plan for user %d during Discord auth start: %v", user.Id, err)
		http.Error(w, "Unable to verify subscription for Discord integration", http.StatusInternalServerError)
		return
	}
	if plan != "pro" && plan != "business" {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": "Discord integration requires Pro or Business plan",
		})
		return
	}

	clientID := os.Getenv("DISCORD_CLIENT_ID")
	redirectURI := os.Getenv("DISCORD_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/api/discord/callback"
	}

	if clientID == "" {
		log.Printf("Error: DISCORD_CLIENT_ID environment variable not set")
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Discord client ID not configured. Please set DISCORD_CLIENT_ID environment variable.",
		})
		return
	}

	log.Printf("🔍 Discord Auth Start - ClientID: %s, RedirectURI: %s", clientID, redirectURI)

	state, err := generateDiscordState()
	if err != nil {
		log.Printf("Error generating Discord OAuth state: %v", err)
		http.Error(w, "Failed to start Discord authentication", http.StatusInternalServerError)
		return
	}

	session, err := auth.Store.Get(r, "auth-session")
	if err != nil {
		log.Printf("Error retrieving session for Discord OAuth: %v", err)
		http.Error(w, "Failed to start Discord authentication", http.StatusInternalServerError)
		return
	}

	session.Values["discord_oauth_state"] = state
	session.Values["discord_oauth_user_id"] = user.Id
	if err := session.Save(r, w); err != nil {
		log.Printf("Error saving Discord OAuth state to session: %v", err)
		http.Error(w, "Failed to start Discord authentication", http.StatusInternalServerError)
		return
	}

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("scope", "identify email")
	params.Set("response_type", "code")
	params.Set("redirect_uri", redirectURI)
	params.Set("state", state)

	oauthURL := fmt.Sprintf("https://discord.com/api/oauth2/authorize?%s", params.Encode())

	log.Printf("🔗 Discord OAuth URL: %s", oauthURL)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"oauth_url": oauthURL,
	})
}

// DiscordCallbackHandler handles Discord OAuth callback
func (h *Handler) DiscordCallbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if errParam := query.Get("error"); errParam != "" {
		log.Printf("Discord OAuth returned error: %s", errParam)
		redirectToSettings(w, r, map[string]string{
			"error": "Discord authorization was cancelled",
		})
		return
	}

	code := query.Get("code")
	state := query.Get("state")
	if code == "" || state == "" {
		redirectToSettings(w, r, map[string]string{
			"error": "Missing authorization data from Discord",
		})
		return
	}

	session, err := auth.Store.Get(r, "auth-session")
	if err != nil {
		log.Printf("Error retrieving session during Discord callback: %v", err)
		redirectToSettings(w, r, map[string]string{
			"error": "Session expired. Please try connecting again",
		})
		return
	}

	storedState, _ := session.Values["discord_oauth_state"].(string)
	if storedState == "" || storedState != state {
		log.Printf("Discord OAuth state mismatch: expected %s, got %s", storedState, state)
		redirectToSettings(w, r, map[string]string{
			"error": "Invalid Discord authorization state",
		})
		return
	}

	userID, ok := session.Values["discord_oauth_user_id"].(int)
	if !ok {
		log.Printf("Discord OAuth missing user ID in session")
		redirectToSettings(w, r, map[string]string{
			"error": "Session missing user information. Please try again",
		})
		return
	}

	plan, err := db.GetUserPlan(h.conn, userID)
	if err != nil {
		log.Printf("Error fetching plan for user %d during Discord callback: %v", userID, err)
		redirectToSettings(w, r, map[string]string{
			"error": "Unable to verify subscription for Discord integration",
		})
		return
	}

	if plan != "pro" && plan != "business" {
		redirectToSettings(w, r, map[string]string{
			"error": "Discord integration requires Pro or Business plan",
		})
		return
	}

	// Exchange code for token
	log.Printf("🔄 Exchanging Discord code for user %d...", userID)
	discordUser, serverID, serverName, channelID, channelName, webhookURL, err := h.exchangeDiscordCode(code)
	if err != nil {
		log.Printf("❌ Error exchanging Discord code for user %d: %v", userID, err)
		redirectToSettings(w, r, map[string]string{
			"error": "Unable to connect to Discord. Please try again",
		})
		return
	}

	// Generate a temporary webhook URL placeholder - user will provide actual webhook URL
	if webhookURL == "" {
		webhookURL = "pending:" + discordUser.ID // Use a special prefix to indicate pending webhook setup
		log.Printf("⚠️  Webhook URL not yet provided, user will need to set it up separately")
	}

	log.Printf("💾 Saving Discord integration for user %d: user=%s, server=%s, channel=%s", userID, discordUser.Username, serverName, channelName)
	if _, err := db.SaveDiscordIntegration(h.conn, userID, discordUser.ID, discordUser.Username, webhookURL, serverID, serverName, channelID, channelName); err != nil {
		log.Printf("❌ Error saving Discord integration for user %d: %v", userID, err)
		redirectToSettings(w, r, map[string]string{
			"error": "Failed to save Discord integration",
		})
		return
	}
	log.Printf("✅ Discord integration saved successfully for user %d", userID)

	delete(session.Values, "discord_oauth_state")
	delete(session.Values, "discord_oauth_user_id")
	if err := session.Save(r, w); err != nil {
		log.Printf("Error clearing Discord OAuth session values: %v", err)
	}

	// Redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

// GetDiscordIntegrationHandler retrieves user's Discord integration
func (h *Handler) GetDiscordIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plan, err := db.GetUserPlan(h.conn, user.Id)
	if err != nil {
		log.Printf("Error fetching plan for user %d while loading Discord integration: %v", user.Id, err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to verify subscription for Discord integration",
		})
		return
	}
	if plan != "pro" && plan != "business" {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": "Discord integration requires Pro or Business plan",
		})
		return
	}

	integration, err := db.GetDiscordIntegration(h.conn, user.Id)
	if err != nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"integration": nil,
			"message":     "No Discord integration found",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"integration": integration,
	})
}

// DisableDiscordIntegrationHandler disables Discord integration
func (h *Handler) DisableDiscordIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plan, planErr := db.GetUserPlan(h.conn, user.Id)
	if planErr != nil {
		log.Printf("Error fetching plan for user %d while disabling Discord integration: %v", user.Id, planErr)
	}
	allowedPlan := plan == "pro" || plan == "business"

	err = db.DisableDiscordIntegration(h.conn, user.Id)
	if err != nil {
		log.Printf("Error disabling Discord integration: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to disable integration",
		})
		return
	}

	message := "Discord integration disabled"
	if !allowedPlan {
		message = "Discord integration disabled. Upgrade to reconnect."
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": message,
	})
}

// UpdateDiscordWebhookHandler updates the Discord webhook URL
func (h *Handler) UpdateDiscordWebhookHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plan, err := db.GetUserPlan(h.conn, user.Id)
	if err != nil {
		log.Printf("Error fetching plan for user %d while updating Discord webhook: %v", user.Id, err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to verify subscription for Discord integration",
		})
		return
	}
	if plan != "pro" && plan != "business" {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": "Discord integration requires Pro or Business plan",
		})
		return
	}

	// Parse webhook URL from request body
	var body struct {
		WebhookURL string `json:"webhook_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	// Validate webhook URL format
	if body.WebhookURL == "" || !strings.HasPrefix(body.WebhookURL, "https://discord.com/api/webhooks/") {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid Discord webhook URL format",
		})
		return
	}

	// Get existing integration
	integration, err := db.GetDiscordIntegration(h.conn, user.Id)
	if err != nil || integration == nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "No Discord integration found. Please connect Discord first.",
		})
		return
	}

	// Update webhook URL in database
	_, err = h.conn.Exec(
		"UPDATE discord_integrations SET webhook_url = $1, updated_at = NOW() WHERE user_id = $2",
		body.WebhookURL,
		user.Id,
	)
	if err != nil {
		log.Printf("Error updating Discord webhook for user %d: %v", user.Id, err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update webhook URL",
		})
		return
	}

	log.Printf("✅ Discord webhook URL updated for user %d", user.Id)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Discord webhook URL updated successfully",
	})
}

// SendDiscordAlert sends an incident alert to Discord via direct message
func (h *Handler) SendDiscordAlert(alert IncidentAlert) error {
	// Get user's Discord integration
	integration, err := db.GetDiscordIntegrationByAppID(h.conn, alert.AppID)
	if err != nil || integration == nil || !integration.IsEnabled {
		log.Printf("No active Discord integration for app %d", alert.AppID)
		return nil // Not an error, just no integration
	}

	// Skip if Discord user ID is not set (user not fully connected)
	if integration.DiscordUserID == "" {
		log.Printf("⚠️  Discord integration for app %d has no Discord user ID set", alert.AppID)
		return nil // Not an error, just not configured yet
	}

	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Printf("❌ DISCORD_BOT_TOKEN environment variable not set, cannot send DMs")
		return fmt.Errorf("Discord bot token not configured")
	}

	// Create Discord embed for DM
	statusEmoji := "🟢"
	switch alert.Status {
	case "down", "error":
		statusEmoji = "🔴"
	case "degraded", "client_error":
		statusEmoji = "🟡"
	}

	// Build message content
	messageContent := fmt.Sprintf(
		"**%s %s - %s**\n\n**App**: %s\n**Status**: %s\n**Status Code**: %d\n**Details**: %s\n**Time**: %s",
		statusEmoji,
		alert.Status,
		alert.AppName,
		alert.AppName,
		alert.Status,
		alert.StatusCode,
		alert.Message,
		alert.Timestamp.Format("2006-01-02 15:04:05 MST"),
	)

	// Create DM with user first
	dmPayload := map[string]interface{}{
		"recipient_id": integration.DiscordUserID,
	}

	dmJsonPayload, err := json.Marshal(dmPayload)
	if err != nil {
		return fmt.Errorf("error marshaling DM payload: %w", err)
	}

	// Get or create DM channel
	req, err := http.NewRequest("POST", "https://discord.com/api/v10/users/@me/channels", bytes.NewBuffer(dmJsonPayload))
	if err != nil {
		return fmt.Errorf("error creating DM request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))

	client := &http.Client{Timeout: 5 * time.Second}
	dmResp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating DM channel: %w", err)
	}
	defer dmResp.Body.Close()

	if dmResp.StatusCode != http.StatusOK && dmResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(dmResp.Body)
		return fmt.Errorf("Discord DM creation error (%d): %s", dmResp.StatusCode, string(body))
	}

	// Parse the DM channel response to get channel ID
	var dmChannel map[string]interface{}
	if err := json.NewDecoder(dmResp.Body).Decode(&dmChannel); err != nil {
		return fmt.Errorf("error parsing DM channel response: %w", err)
	}

	channelID, ok := dmChannel["id"].(string)
	if !ok {
		return fmt.Errorf("could not extract channel ID from DM response")
	}

	// Send message to the DM channel
	messagePayload := map[string]interface{}{
		"content": messageContent,
	}

	msgJsonPayload, err := json.Marshal(messagePayload)
	if err != nil {
		return fmt.Errorf("error marshaling message payload: %w", err)
	}

	msgReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID),
		bytes.NewBuffer(msgJsonPayload),
	)
	if err != nil {
		return fmt.Errorf("error creating message request: %w", err)
	}

	msgReq.Header.Set("Content-Type", "application/json")
	msgReq.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))

	msgResp, err := client.Do(msgReq)
	if err != nil {
		return fmt.Errorf("error sending Discord DM: %w", err)
	}
	defer msgResp.Body.Close()

	if msgResp.StatusCode != http.StatusOK && msgResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(msgResp.Body)
		return fmt.Errorf("Discord message send error (%d): %s", msgResp.StatusCode, string(body))
	}

	// Log notification
	db.LogIncidentNotification(h.conn, alert.AppID, "discord", "sent")

	log.Printf("✅ Discord DM alert sent to user %s for app %s", integration.DiscordUserID, alert.AppName)
	return nil
}

// exchangeDiscordCode exchanges auth code for token and gets user info
func (h *Handler) exchangeDiscordCode(code string) (*DiscordUser, string, string, string, string, string, error) {
	clientID := os.Getenv("DISCORD_CLIENT_ID")
	clientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
	redirectURI := os.Getenv("DISCORD_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/api/discord/callback"
	}

	if clientID == "" || clientSecret == "" {
		return nil, "", "", "", "", "", fmt.Errorf("Discord credentials not configured")
	}

	// Exchange code for token
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", "https://discord.com/api/v10/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, "", "", "", "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error exchanging Discord code: %v", err)
		return nil, "", "", "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ Discord token exchange failed with status %d: %s", resp.StatusCode, string(body))
		return nil, "", "", "", "", "", fmt.Errorf("Discord OAuth HTTP %d: %s", resp.StatusCode, string(body))
	}

	var oauthResp DiscordOAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		log.Printf("❌ Error decoding Discord OAuth response: %v", err)
		return nil, "", "", "", "", "", err
	}
	log.Printf("✅ Discord token exchange successful, access token: %s", oauthResp.AccessToken[:20]+"...")

	// Get user info
	discordUser, err := h.getDiscordUser(oauthResp.AccessToken)
	if err != nil {
		log.Printf("❌ Error fetching Discord user: %v", err)
		return nil, "", "", "", "", "", err
	}
	log.Printf("✅ Discord user fetched: %s (%s)", discordUser.Username, discordUser.ID)

	// For Discord integration, we store user info but don't automatically create webhooks
	// Users will need to create webhooks manually in Discord or through future flow
	// This is more secure and gives users full control
	serverID := "0"
	serverName := "Manual Setup"
	channelID := "0"
	channelName := "See instructions"
	// Use a placeholder that will be replaced by user-provided webhook URL
	webhookURL := ""

	return discordUser, serverID, serverName, channelID, channelName, webhookURL, nil
}

// getDiscordUser gets the current user's info
func (h *Handler) getDiscordUser(token string) (*DiscordUser, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Discord user API error: %s", string(body))
	}

	var user DiscordUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// getDiscordGuildAndChannel gets the user's guilds and finds the first suitable channel
func (h *Handler) getDiscordGuildAndChannel(token string) (string, string, string, string, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me/guilds", nil)
	if err != nil {
		return "", "", "", "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", "", fmt.Errorf("Discord guilds API error: %s", string(body))
	}

	var guilds []DiscordGuild
	err = json.NewDecoder(resp.Body).Decode(&guilds)
	if err != nil {
		return "", "", "", "", err
	}

	if len(guilds) == 0 {
		return "", "", "", "", fmt.Errorf("user has no Discord servers")
	}

	// Use the first guild
	serverID := guilds[0].ID
	serverName := guilds[0].Name

	// Get channels from that guild
	channelID, channelName, err := h.getDiscordChannels(token, serverID)
	if err != nil {
		return "", "", "", "", err
	}

	return serverID, serverName, channelID, channelName, nil
}

// getDiscordChannels gets text channels from a guild
func (h *Handler) getDiscordChannels(token, guildID string) (string, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/channels", guildID), nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("Discord channels API error: %s", string(body))
	}

	var channels []DiscordChannel
	err = json.NewDecoder(resp.Body).Decode(&channels)
	if err != nil {
		return "", "", err
	}

	// Find the first text channel (type 0)
	for _, ch := range channels {
		if ch.Type == 0 { // Text channel
			return ch.ID, ch.Name, nil
		}
	}

	if len(channels) > 0 {
		return channels[0].ID, channels[0].Name, nil
	}

	return "", "", fmt.Errorf("no text channels found in guild")
}

// createDiscordWebhook creates a webhook in the specified channel
func (h *Handler) createDiscordWebhook(token, guildID, channelID string) (string, error) {
	payload := map[string]interface{}{
		"name": "StatusFrame Alerts",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/webhooks", channelID), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Discord webhook creation error: %s", string(body))
	}

	var webhook map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&webhook)
	if err != nil {
		return "", err
	}

	// Construct webhook URL
	webhookID := webhook["id"].(string)
	webhookToken := webhook["token"].(string)
	webhookURL := fmt.Sprintf("https://discord.com/api/v10/webhooks/%s/%s", webhookID, webhookToken)

	return webhookURL, nil
}

// generateDiscordState creates a cryptographically secure state parameter for Discord OAuth
func generateDiscordState() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer), nil
}
