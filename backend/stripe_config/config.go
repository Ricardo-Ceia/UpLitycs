package stripe_config

import (
	"log"
	"os"

	"github.com/stripe/stripe-go/v81"
)

type Config struct {
	SecretKey              string
	PublishableKey         string
	ProMonthlyPriceID      string
	ProYearlyPriceID       string
	BusinessMonthlyPriceID string
	BusinessYearlyPriceID  string
	WebhookSecret          string
	AppURL                 string
}

var StripeConfig *Config

// Initialize sets up Stripe configuration from environment variables
func Initialize() {
	StripeConfig = &Config{
		SecretKey:              os.Getenv("STRIPE_SECRET_KEY"),
		PublishableKey:         os.Getenv("STRIPE_PUBLISHABLE_KEY"),
		ProMonthlyPriceID:      os.Getenv("STRIPE_PRO_MONTHLY_PRICE_ID"),
		ProYearlyPriceID:       os.Getenv("STRIPE_PRO_YEARLY_PRICE_ID"),
		BusinessMonthlyPriceID: os.Getenv("STRIPE_BUSINESS_MONTHLY_PRICE_ID"),
		BusinessYearlyPriceID:  os.Getenv("STRIPE_BUSINESS_YEARLY_PRICE_ID"),
		WebhookSecret:          os.Getenv("STRIPE_WEBHOOK_SECRET"),
		AppURL:                 getAppURL(),
	}

	// Set the Stripe API key
	stripe.Key = StripeConfig.SecretKey

	// Validate required config
	if StripeConfig.SecretKey == "" {
		log.Fatal("STRIPE_SECRET_KEY environment variable is required")
	}

	log.Println("✅ Stripe configuration initialized")
	log.Printf("📦 Pro Monthly Price: %s", StripeConfig.ProMonthlyPriceID)
	log.Printf("📦 Pro Yearly Price: %s", StripeConfig.ProYearlyPriceID)
	log.Printf("📦 Business Monthly Price: %s", StripeConfig.BusinessMonthlyPriceID)
	log.Printf("📦 Business Yearly Price: %s", StripeConfig.BusinessYearlyPriceID)
}

func getAppURL() string {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3333"
	}
	return appURL
}

// GetPriceID returns the Stripe price ID for a given plan and billing period
func GetPriceID(plan string, billingPeriod string) string {
	switch plan {
	case "pro":
		if billingPeriod == "yearly" {
			return StripeConfig.ProYearlyPriceID
		}
		return StripeConfig.ProMonthlyPriceID
	case "business":
		if billingPeriod == "yearly" {
			return StripeConfig.BusinessYearlyPriceID
		}
		return StripeConfig.BusinessMonthlyPriceID
	default:
		return ""
	}
}
