package auth

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "uplyticskey12345"
	MaxAge = 86400 * 30
	IsProd = false
)

var Store *sessions.CookieStore

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "auth-session")

		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		userId, ok := session.Values["UserId"].(int)

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loadinng .env file")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if googleClientId == "" || googleClientSecret == "" {
		log.Fatal("Google client ID or secret not set in environment variables")
	}

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:3333/auth/google/callback", "email", "profile"),
	)
}
