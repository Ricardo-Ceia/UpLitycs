package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"uplytics/backend/handlers"
	"uplytics/backend/status_checker"
	"uplytics/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func testFunc(intervalSeconds int) {
	log.Println("Starting periodic status checks every", intervalSeconds, "seconds")
	conn := db.OpenDB()
	defer conn.Close()

	users, err := db.GetAllUsers(conn)
	if err != nil {
		log.Println("GetAllUsers error:", err)
		return
	}
	var urls []string
	for _, u := range users {
		hp, err := db.GetURLOfMainPage(conn, u)
		if err != nil {
			log.Println("GetURLOfMainPage error:", err)
			continue
		}
		urls = append(urls, hp)
	}
	if len(urls) == 0 {
		log.Println("No homepages to check yet.")
	}

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		statuses, err := status_checker.GetPagesStatus(urls)
		if err != nil {
			log.Println("status check error:", err)
		} else {
			log.Println("Status codes:", statuses)
		}
		<-ticker.C
	}
}

func main() {
	// Ensure correct MIME types are registered
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".json", "application/json")

	r := chi.NewRouter()

	// Add some useful middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- API routes ---
	r.Route("/api", func(r chi.Router) {
		r.Get("/start-onboarding", handlers.StartOnboardingHandler)
		r.Post("/go-to-dashboard", handlers.GoToDashboardHandler)
	})

	// --- React build directory ---
	workDir, _ := os.Getwd()
	reactDir := filepath.Join(workDir, "frontend", "dist")

	// Create a file server for static assets
	fileServer := http.FileServer(http.Dir(reactDir))

	// Handle static assets with proper MIME types
	r.Handle("/assets/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set proper content type based on file extension
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		fileServer.ServeHTTP(w, r)
	}))

	// Handle other static files
	r.Handle("/favicon.ico", fileServer)

	// Catch-all: serve index.html for React Router
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Skip API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// For any other route, serve the React app
		indexPath := filepath.Join(reactDir, "index.html")

		// Check if index.html exists
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			http.Error(w, "React build not found. Run 'npm run build' in your frontend directory.", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, indexPath)
	})

	// Start background checks
	go testFunc(30)

	log.Println("Server starting on http://localhost:3333")
	if err := http.ListenAndServe(":3333", r); err != nil {
		log.Fatal(err)
	}
}
