package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"statusframe/backend/auth"
	"statusframe/backend/utils"
	"statusframe/db"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	ADMIN_EMAIL = "ricardoceia.sete@gmail.com"
)

type AdminUser struct {
	ID                   int       `json:"id"`
	Username             string    `json:"username"`
	Email                string    `json:"email"`
	AvatarURL            string    `json:"avatar_url"`
	Plan                 string    `json:"plan"`
	PlanStartedAt        time.Time `json:"plan_started_at"`
	StripeCustomerID     *string   `json:"stripe_customer_id"`
	StripeSubscriptionID *string   `json:"stripe_subscription_id"`
	CreatedAt            time.Time `json:"created_at"`
	AppCount             int       `json:"app_count"`
	TotalChecks          int       `json:"total_checks"`
}

type AdminStats struct {
	TotalUsers        int `json:"total_users"`
	FreeUsers         int `json:"free_users"`
	ProUsers          int `json:"pro_users"`
	BusinessUsers     int `json:"business_users"`
	TotalApps         int `json:"total_apps"`
	TotalStatusChecks int `json:"total_status_checks"`
	ActiveSubscribers int `json:"active_subscribers"`
}

// AdminMiddleware checks if the logged-in user is the admin
func (h *Handler) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.Store.Get(r, "auth-session")
		if err != nil {
			http.Error(w, "Unauthorized - Please login", http.StatusUnauthorized)
			return
		}

		userId, ok := session.Values["userId"].(int)
		if !ok {
			http.Error(w, "Unauthorized - Please login", http.StatusUnauthorized)
			return
		}

		// Get user from database to check email
		user, err := db.GetUserById(h.conn, userId)
		if err != nil {
			http.Error(w, "Unauthorized - User not found", http.StatusUnauthorized)
			return
		}

		// Check if user email matches admin email
		if user.Email != ADMIN_EMAIL {
			log.Printf("⚠️  Non-admin user attempted to access admin panel: %s", user.Email)
			http.Error(w, "Forbidden - Admin access required", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "admin", true)
		ctx = context.WithValue(ctx, "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminCheckSessionHandler checks if current user is admin
func (h *Handler) AdminCheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, "auth-session")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	userId, ok := session.Values["userId"].(int)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	// Get user from database
	user, err := db.GetUserById(h.conn, userId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	// Check if user is admin
	isAdmin := user.Email == ADMIN_EMAIL

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": isAdmin,
		"email":         user.Email,
		"username":      user.Name,
	})
}

// GetAllUsersHandler returns all users with their subscription details
func (h *Handler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			u.id,
			u.username,
			u.email,
			COALESCE(u.avatar_url, '') as avatar_url,
			u.plan,
			u.plan_started_at,
			u.stripe_customer_id,
			u.stripe_subscription_id,
			u.created_at,
			COUNT(DISTINCT a.id) as app_count,
			COUNT(us.id) as total_checks
		FROM users u
		LEFT JOIN apps a ON u.id = a.user_id
		LEFT JOIN user_status us ON a.id = us.app_id
		GROUP BY u.id, u.username, u.email, u.avatar_url, u.plan, u.plan_started_at, 
		         u.stripe_customer_id, u.stripe_subscription_id, u.created_at
		ORDER BY u.created_at DESC
	`

	rows, err := h.conn.Query(query)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []AdminUser
	for rows.Next() {
		var user AdminUser
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.AvatarURL,
			&user.Plan,
			&user.PlanStartedAt,
			&user.StripeCustomerID,
			&user.StripeSubscriptionID,
			&user.CreatedAt,
			&user.AppCount,
			&user.TotalChecks,
		)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) UploadLogoHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form with 5MB max size
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}

	// Get the file from form
	file, fileHeader, err := r.FormFile("logo")
	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type based on header
	contentType := fileHeader.Header.Get("Content-Type")
	validTypes := map[string]bool{
		"image/png":     true,
		"image/jpeg":    true,
		"image/jpg":     true,
		"image/svg+xml": true,
		"image/webp":    true,
	}

	if !validTypes[contentType] {
		http.Error(w, "Invalid file type. Only images allowed (PNG, JPG, SVG, WebP)", http.StatusBadRequest)
		return
	}

	// Validate file size
	if fileHeader.Size > 5<<20 { // 5MB
		http.Error(w, "File size exceeds 5MB limit", http.StatusBadRequest)
		return
	}

	// Create AWS session
	sess, err := utils.CreateAWSSession()
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
		return
	}

	// Get bucket name from environment
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	if bucketName == "" {
		log.Printf("AWS_S3_BUCKET_NAME not configured")
		http.Error(w, "S3 bucket name not configured", http.StatusInternalServerError)
		return
	}

	// Create S3 client
	s3Client := s3.New(sess)

	// Generate unique file name with original extension
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".png" // default extension
	}
	fileName := fmt.Sprintf("logo_%d_%d%s", time.Now().Unix(), time.Now().Nanosecond(), ext)
	s3Path := "logos/" + fileName

	// Upload directly to S3 - streaming the file without reading into memory
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(s3Path),
		Body:          file, // Stream directly from multipart.File
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fileHeader.Size),
		// ACL removed - bucket must have public read policy configured
	})
	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
		return
	}

	// Construct the public file URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		bucketName,
		os.Getenv("AWS_REGION"),
		s3Path,
	)

	// Return the URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": fileURL,
	})
}

// GetAdminStatsHandler returns overall platform statistics
func (h *Handler) GetAdminStatsHandler(w http.ResponseWriter, r *http.Request) {
	var stats AdminStats

	// Get user counts by plan
	err := h.conn.QueryRow(`
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE plan = 'free') as free,
			COUNT(*) FILTER (WHERE plan = 'pro') as pro,
			COUNT(*) FILTER (WHERE plan = 'business') as business,
			COUNT(*) FILTER (WHERE stripe_subscription_id IS NOT NULL) as active_subscribers
		FROM users
	`).Scan(&stats.TotalUsers, &stats.FreeUsers, &stats.ProUsers, &stats.BusinessUsers, &stats.ActiveSubscribers)

	if err != nil {
		log.Printf("Error fetching user stats: %v", err)
		http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	// Get total apps
	err = h.conn.QueryRow(`SELECT COUNT(*) FROM apps`).Scan(&stats.TotalApps)
	if err != nil {
		log.Printf("Error fetching app count: %v", err)
	}

	// Get total status checks
	err = h.conn.QueryRow(`SELECT COUNT(*) FROM user_status`).Scan(&stats.TotalStatusChecks)
	if err != nil {
		log.Printf("Error fetching status check count: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
