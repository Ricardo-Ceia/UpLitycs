// backend/auth/custom_auth.go
package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	Scopes       []string
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

var googleConfig GoogleOAuthConfig

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "auth-session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		userId, ok := session.Values["userId"].(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username, ok := session.Values["user"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		ctx = context.WithValue(ctx, "user", username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewAuth() {

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	sessionSecret := os.Getenv("SESSION_SECRET")

	if googleClientId == "" || googleClientSecret == "" {
		log.Fatal("Google client ID or secret not set in environment variables")
	}

	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET environment variable not set")
	}

	Store = sessions.NewCookieStore([]byte(sessionSecret))
	Store.MaxAge(86400 * 30)
	Store.Options.Path = "/"
	Store.Options.HttpOnly = true
	Store.Options.Secure = false // Set to true in production with HTTPS
	Store.Options.SameSite = http.SameSiteLaxMode
	Store.Options.MaxAge = 86400 * 30 // 30 days - ensures cookie persists across browser sessions

	// Initialize Google OAuth config
	googleConfig = GoogleOAuthConfig{
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		RedirectURL:  googleRedirectURL,
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://oauth2.googleapis.com/token",
		UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
		Scopes:       []string{"email", "profile"},
	}
}

// Generate a random state parameter for OAuth security
func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// BeginGoogleAuth starts the Google OAuth flow
func BeginGoogleAuth(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	// Store state in session for verification
	session, _ := Store.Get(r, "auth-session")
	session.Values["oauth_state"] = state
	session.Save(r, w)

	// Build authorization URL
	params := url.Values{}
	params.Add("client_id", googleConfig.ClientID)
	params.Add("redirect_uri", googleConfig.RedirectURL)
	params.Add("scope", strings.Join(googleConfig.Scopes, " "))
	params.Add("response_type", "code")
	params.Add("state", state)
	params.Add("access_type", "offline")

	authURL := googleConfig.AuthURL + "?" + params.Encode()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleGoogleCallback handles the OAuth callback
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) (*GoogleUserInfo, error) {
	// Verify state parameter
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		return nil, fmt.Errorf("session error: %v", err)
	}

	storedState, ok := session.Values["oauth_state"].(string)
	if !ok || storedState != r.URL.Query().Get("state") {
		return nil, fmt.Errorf("invalid state parameter")
	}

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		return nil, fmt.Errorf("no authorization code received")
	}

	// Exchange code for token
	token, err := exchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	// Get user info
	userInfo, err := getUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	return userInfo, nil
}

// exchangeCodeForToken exchanges authorization code for access token
func exchangeCodeForToken(code string) (*GoogleTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", googleConfig.ClientID)
	data.Set("client_secret", googleConfig.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", googleConfig.RedirectURL)

	resp, err := http.PostForm(googleConfig.TokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var token GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

// getUserInfo fetches user information from Google
func getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequest("GET", googleConfig.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
