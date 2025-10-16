package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"statusframe/backend/auth"
	"statusframe/backend/handlers"
	"statusframe/backend/stripe_config"
	"statusframe/backend/worker"
	"statusframe/db"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
		case ".mp4":
			w.Header().Set("Content-Type", "video/mp4")
		case ".webm":
			w.Header().Set("Content-Type", "video/webm")
		case ".ogg":
			w.Header().Set("Content-Type", "video/ogg")
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
	mime.AddExtensionType(".mp4", "video/mp4")
	mime.AddExtensionType(".webm", "video/webm")
	mime.AddExtensionType(".ogg", "video/ogg")

	conn := db.OpenDB()
	defer conn.Close()

	// Initialize Stripe configuration
	stripe_config.Initialize()

	// Ping database to ensure connection
	if err := db.PingDB(conn); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	appHandlers := handlers.NewHandler(conn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Add CORS middleware to allow credentials (cookies)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize custom authentication
	auth.NewAuth()

	// Start health checker worker (check every 30 seconds)
	healthChecker := worker.NewHealthChecker(conn, 30*time.Second)
	go healthChecker.Start()
	log.Println("✅ Health checker worker started (checking every 30 seconds)")

	// Start SSL certificate checker (check daily)
	sslChecker := worker.NewSSLChecker(conn)
	go sslChecker.Start()
	log.Println("✅ SSL certificate checker started (checking daily)")

	// --- API routes (must come first) ---
	r.Route("/api", func(r chi.Router) {
		r.With(auth.AuthMiddleware).Get("/start-onboarding", handlers.StartOnboardingHandler)
		r.With(auth.AuthMiddleware).Post("/go-to-dashboard", appHandlers.GoToDashboardHandler)
		r.With(auth.AuthMiddleware).Get("/user-status", appHandlers.GetUserStatusHandler)
		r.With(auth.AuthMiddleware).Get("/latest-status", appHandlers.GetLatestStatusHandler)
		r.With(auth.AuthMiddleware).Post("/update-theme", appHandlers.UpdateThemeHandler)

		// Multi-app dashboard routes
		r.With(auth.AuthMiddleware).Get("/user-apps", appHandlers.GetUserAppsHandler)
		r.With(auth.AuthMiddleware).Delete("/apps/{appId}", appHandlers.DeleteAppHandler)
		r.With(auth.AuthMiddleware).Get("/check-plan-limit", appHandlers.CheckPlanLimitHandler)
		r.With(auth.AuthMiddleware).Get("/plan-features", appHandlers.GetPlanFeaturesHandler)

		// Stripe payment routes
		r.With(auth.AuthMiddleware).Post("/create-checkout-session", appHandlers.CreateCheckoutSessionHandler)
		r.With(auth.AuthMiddleware).Post("/create-portal-session", appHandlers.CreateCustomerPortalSessionHandler)
		r.With(auth.AuthMiddleware).Get("/stripe-success", appHandlers.StripeSuccessHandler) // Handle successful payment
		r.Post("/stripe-webhook", appHandlers.StripeWebhookHandler)                          // No auth - Stripe signs the request

		// Public API - no authentication required
		r.Get("/public/status/{slug}", appHandlers.GetPublicStatusHandler)
		r.Get("/public/ping/{slug}", appHandlers.GetCurrentResponseTimeHandler)
		r.Get("/badge/{slug}", appHandlers.GetUptimeBadgeHandler) // Public uptime badge

		// Admin routes - check if logged-in user is admin email
		r.Route("/admin", func(r chi.Router) {
			r.Get("/check-session", appHandlers.AdminCheckSessionHandler)

			// Protected admin routes - must be logged in with admin email
			r.Group(func(r chi.Router) {
				r.Use(appHandlers.AdminMiddleware)
				r.Get("/users", appHandlers.GetAllUsersHandler)
				r.Get("/stats", appHandlers.GetAdminStatsHandler)
			})
		})
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

	// Serve video files and other root static files
	r.Get("/demo.mp4", func(w http.ResponseWriter, r *http.Request) {
		videoPath := filepath.Join(reactBuildDir, "demo.mp4")
		if _, err := os.Stat(videoPath); err == nil {
			w.Header().Set("Content-Type", "video/mp4")
			http.ServeFile(w, r, videoPath)
		} else {
			http.NotFound(w, r)
		}
	})

	r.Get("/demo-poster.jpg", func(w http.ResponseWriter, r *http.Request) {
		posterPath := filepath.Join(reactBuildDir, "demo-poster.jpg")
		if _, err := os.Stat(posterPath); err == nil {
			w.Header().Set("Content-Type", "image/jpeg")
			http.ServeFile(w, r, posterPath)
		} else {
			http.NotFound(w, r)
		}
	})

	r.Get("/vite.svg", func(w http.ResponseWriter, r *http.Request) {
		svgPath := filepath.Join(reactBuildDir, "vite.svg")
		if _, err := os.Stat(svgPath); err == nil {
			w.Header().Set("Content-Type", "image/svg+xml")
			http.ServeFile(w, r, svgPath)
		} else {
			http.NotFound(w, r)
		}
	})

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

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
