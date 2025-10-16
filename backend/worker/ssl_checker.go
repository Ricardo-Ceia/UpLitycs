package worker

import (
	"crypto/tls"
	"database/sql"
	"log"
	"net/url"
	"statusframe/db"
	"time"
)

type SSLChecker struct {
	conn *sql.DB
}

func NewSSLChecker(conn *sql.DB) *SSLChecker {
	return &SSLChecker{
		conn: conn,
	}
}

// Start begins the daily SSL certificate check routine
func (sc *SSLChecker) Start() {
	log.Println("🔒 SSL certificate checker started - checking daily")

	// Run immediately on start
	sc.checkAllSSLCertificates()

	// Run every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		sc.checkAllSSLCertificates()
	}
}

func (sc *SSLChecker) checkAllSSLCertificates() {
	log.Println("🔍 Starting SSL certificate check for all HTTPS apps...")

	apps, err := sc.getHTTPSApps()
	if err != nil {
		log.Printf("❌ Error fetching HTTPS apps: %v", err)
		return
	}

	checkedCount := 0
	errorCount := 0

	for _, app := range apps {
		expiryDate, issuer, err := sc.checkSSLCertificate(app.HealthURL)
		if err != nil {
			log.Printf("⚠️  SSL check failed for %s (%s): %v", app.AppName, app.HealthURL, err)
			errorCount++
			// Clear SSL data on error
			db.UpdateSSLInfo(sc.conn, app.AppID, nil, nil, nil)
			continue
		}

		daysUntilExpiry := int(time.Until(expiryDate).Hours() / 24)

		// Update SSL info in database
		err = db.UpdateSSLInfo(sc.conn, app.AppID, &expiryDate, &daysUntilExpiry, &issuer)
		if err != nil {
			log.Printf("❌ Error updating SSL info for %s: %v", app.AppName, err)
			continue
		}

		checkedCount++

		// Log with appropriate emoji
		emoji := "✅"
		if daysUntilExpiry <= 7 {
			emoji = "🔴"
		} else if daysUntilExpiry <= 30 {
			emoji = "🟡"
		}

		log.Printf("%s SSL checked for %s: expires in %d days (%s) - Issuer: %s",
			emoji, app.AppName, daysUntilExpiry, expiryDate.Format("2006-01-02"), issuer)
	}

	log.Printf("🔒 SSL check complete: %d checked, %d errors", checkedCount, errorCount)
}

type httpsApp struct {
	AppID     int
	AppName   string
	HealthURL string
}

func (sc *SSLChecker) getHTTPSApps() ([]httpsApp, error) {
	query := `
		SELECT id, app_name, health_url
		FROM apps
		WHERE health_url LIKE 'https://%'
	`

	rows, err := sc.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []httpsApp
	for rows.Next() {
		var app httpsApp
		err := rows.Scan(&app.AppID, &app.AppName, &app.HealthURL)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (sc *SSLChecker) checkSSLCertificate(healthURL string) (time.Time, string, error) {
	parsedURL, err := url.Parse(healthURL)
	if err != nil {
		return time.Time{}, "", err
	}

	host := parsedURL.Host
	// Add default port if not specified
	if parsedURL.Port() == "" {
		host += ":443"
	}

	// Connect with TLS
	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: false, // Verify certificates properly
	})
	if err != nil {
		return time.Time{}, "", err
	}
	defer conn.Close()

	// Get peer certificates
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return time.Time{}, "", err
	}

	// Get the leaf certificate (first one)
	cert := certs[0]

	// Extract issuer
	issuer := "Unknown"
	if len(cert.Issuer.Organization) > 0 {
		issuer = cert.Issuer.Organization[0]
	} else if len(cert.Issuer.CommonName) > 0 {
		issuer = cert.Issuer.CommonName
	}

	return cert.NotAfter, issuer, nil
}
