package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"uplytics/backend/auth"
	"uplytics/backend/handlers"
	"uplytics/backend/worker"
	"uplytics/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Custom file server that sets correct MIME types
func staticFileServer(dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(dir, r.URL.Path)

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		// Set correct MIME type based on file extension
		ext := filepath.Ext(filePath)
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
		case ".ico":
			w.Header().Set("Content-Type", "image/x-icon")
		case ".woff":
			w.Header().Set("Content-Type", "font/woff")
		case ".woff2":
			w.Header().Set("Content-Type", "font/woff2")
		}

		// Serve the file
		http.ServeFile(w, r, filePath)
	})
}

func main() {
	// Ensure correct MIME types are registered
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".woff", "font/woff")
	mime.AddExtensionType(".woff2", "font/woff2")

	conn := db.OpenDB()
	defer conn.Close()

	// Ping database to ensure connection
	if err := db.PingDB(conn); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	appHandlers := handlers.NewHandler(conn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Initialize custom authentication
	auth.NewAuth()

	// Start health checker worker (check every 30 seconds)
	healthChecker := worker.NewHealthChecker(conn, 30*time.Second)
	go healthChecker.Start()
	log.Println("✅ Health checker worker started (checking every 30 seconds)")

	// --- API routes (must come first) ---
	r.Route("/api", func(r chi.Router) {
		r.With(auth.AuthMiddleware).Get("/start-onboarding", handlers.StartOnboardingHandler)
		r.With(auth.AuthMiddleware).Post("/go-to-dashboard", appHandlers.GoToDashboardHandler)
		r.With(auth.AuthMiddleware).Get("/user-status", appHandlers.GetUserStatusHandler)
		r.With(auth.AuthMiddleware).Get("/latest-status", appHandlers.GetLatestStatusHandler)
		r.With(auth.AuthMiddleware).Post("/update-theme", appHandlers.UpdateThemeHandler)

		// Public API - no authentication required
		r.Get("/public/status/{slug}", appHandlers.GetPublicStatusHandler)
		r.Get("/public/ping/{slug}", appHandlers.GetCurrentResponseTimeHandler)
	})

	//--- OAuth Auth Routes (must come before catch-all) ---
	r.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}", handlers.BeginAuthHandler)
		r.Get("/{provider}/callback", appHandlers.GetAuthHandler)
		r.Get("/logout", handlers.LogoutHandler)
		r.Get("/check-session", handlers.CheckSessionHandler)
	})

	// --- Serve React build files ---
	workDir, _ := os.Getwd()
	reactBuildDir := filepath.Join(workDir, "frontend", "dist")

	log.Printf("Serving static files from: %s", reactBuildDir)

	// Serve static assets
	r.Handle("/assets/*", http.StripPrefix("/assets/", staticFileServer(filepath.Join(reactBuildDir, "assets"))))

	// Serve favicon
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		faviconPath := filepath.Join(reactBuildDir, "favicon.ico")
		if _, err := os.Stat(faviconPath); err == nil {
			w.Header().Set("Content-Type", "image/x-icon")
			http.ServeFile(w, r, faviconPath)
		} else {
			http.NotFound(w, r)
		}
	})

	// --- Serve React app for frontend routes ---
	r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join(reactBuildDir, "index.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, indexPath)
	})

	// Catch-all route to serve React app for client-side routing
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join(reactBuildDir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			log.Printf("React build not found at %s", indexPath)
			http.Error(w, "Frontend not built", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, indexPath)
	})

	log.Println("Server starting on http://localhost:3333")
	if err := http.ListenAndServe(":3333", r); err != nil {
		log.Fatal(err)
	}
}
