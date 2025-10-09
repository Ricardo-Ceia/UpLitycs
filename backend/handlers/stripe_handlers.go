package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"statusframe/backend/stripe_config"
	"statusframe/db"

	"github.com/stripe/stripe-go/v81"
	billingportalsession "github.com/stripe/stripe-go/v81/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/webhook"
)

// CreateCheckoutSessionHandler creates a Stripe checkout session for subscription
func (h *Handler) CreateCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user
	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req struct {
		Plan          string `json:"plan"`           // "pro" or "business"
		BillingPeriod string `json:"billing_period"` // "monthly" or "yearly"
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate plan
	if req.Plan != "pro" && req.Plan != "business" {
		http.Error(w, "Invalid plan. Must be 'pro' or 'business'", http.StatusBadRequest)
		return
	}

	// Validate billing period
	if req.BillingPeriod != "monthly" && req.BillingPeriod != "yearly" {
		http.Error(w, "Invalid billing period. Must be 'monthly' or 'yearly'", http.StatusBadRequest)
		return
	}

	// Get the price ID for the selected plan and billing period
	priceID := stripe_config.GetPriceID(req.Plan, req.BillingPeriod)
	if priceID == "" {
		log.Printf("No price ID found for plan=%s, billing=%s", req.Plan, req.BillingPeriod)
		http.Error(w, "Invalid plan configuration", http.StatusInternalServerError)
		return
	}

	// Check if user already has a Stripe customer ID
	stripeCustomerID, err := db.GetStripeCustomerId(h.conn, user.Id)
	if err != nil {
		log.Printf("Error getting Stripe customer ID: %v", err)
	}

	// Create Stripe customer if they don't have one
	if stripeCustomerID == "" {
		customerParams := &stripe.CustomerParams{
			Email: stripe.String(user.Email),
			Name:  stripe.String(user.Name),
			Metadata: map[string]string{
				"user_id": fmt.Sprintf("%d", user.Id),
			},
		}

		cust, err := customer.New(customerParams)
		if err != nil {
			log.Printf("Error creating Stripe customer: %v", err)
			http.Error(w, "Failed to create customer", http.StatusInternalServerError)
			return
		}
		stripeCustomerID = cust.ID

		// Update user with Stripe customer ID
		_, err = h.conn.Exec(
			"UPDATE users SET stripe_customer_id = $1 WHERE id = $2",
			stripeCustomerID,
			user.Id,
		)
		if err != nil {
			log.Printf("Error saving Stripe customer ID: %v", err)
		}
	}

	// Create Checkout Session
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeCustomerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(stripe_config.StripeConfig.AppURL + "/dashboard?upgrade=success"),
		CancelURL:  stripe.String(stripe_config.StripeConfig.AppURL + "/pricing?upgrade=cancelled"),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", user.Id),
			"plan":    req.Plan,
		},
	}

	sess, err := checkoutsession.New(params)
	if err != nil {
		log.Printf("Error creating checkout session: %v", err)
		http.Error(w, "Failed to create checkout session", http.StatusInternalServerError)
		return
	}

	// Return the checkout session URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"url": sess.URL,
	})
}

// StripeWebhookHandler handles Stripe webhook events
func (h *Handler) StripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Verify webhook signature
	event, err := webhook.ConstructEvent(
		payload,
		r.Header.Get("Stripe-Signature"),
		stripe_config.StripeConfig.WebhookSecret,
	)

	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Handle the event
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.handleCheckoutSessionCompleted(session)

	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.handleSubscriptionCreated(subscription)

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.handleSubscriptionUpdated(subscription)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.handleSubscriptionDeleted(subscription)

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

// handleCheckoutSessionCompleted processes successful checkout
func (h *Handler) handleCheckoutSessionCompleted(session stripe.CheckoutSession) {
	log.Printf("✅ Checkout session completed: %s", session.ID)

	// The subscription will be handled by customer.subscription.created event
	// This is just for logging
}

// handleSubscriptionCreated activates the user's subscription
func (h *Handler) handleSubscriptionCreated(subscription stripe.Subscription) {
	log.Printf("🎉 Subscription created: %s for customer %s", subscription.ID, subscription.Customer.ID)

	// Get user by Stripe customer ID
	user, err := db.GetUserByStripeCustomerId(h.conn, subscription.Customer.ID)
	if err != nil {
		log.Printf("❌ Error finding user for customer %s: %v", subscription.Customer.ID, err)
		return
	}

	// Determine plan from subscription
	plan := h.getPlanFromSubscription(&subscription)

	// Update user's subscription in database
	err = db.UpdateUserSubscription(
		h.conn,
		user.Id,
		plan,
		subscription.Customer.ID,
		subscription.ID,
	)
	if err != nil {
		log.Printf("❌ Error updating subscription: %v", err)
		return
	}

	log.Printf("✅ User %d upgraded to %s plan", user.Id, plan)
}

// handleSubscriptionUpdated processes subscription changes
func (h *Handler) handleSubscriptionUpdated(subscription stripe.Subscription) {
	log.Printf("🔄 Subscription updated: %s", subscription.ID)

	user, err := db.GetUserByStripeCustomerId(h.conn, subscription.Customer.ID)
	if err != nil {
		log.Printf("❌ Error finding user: %v", err)
		return
	}

	// Check subscription status
	if subscription.Status == "active" {
		plan := h.getPlanFromSubscription(&subscription)
		err = db.UpdateUserSubscription(
			h.conn,
			user.Id,
			plan,
			subscription.Customer.ID,
			subscription.ID,
		)
		if err != nil {
			log.Printf("❌ Error updating subscription: %v", err)
		}
	} else if subscription.Status == "canceled" || subscription.Status == "unpaid" {
		// Downgrade to free
		err = db.CancelUserSubscription(h.conn, user.Id)
		if err != nil {
			log.Printf("❌ Error canceling subscription: %v", err)
		}
		log.Printf("⬇️  User %d downgraded to free plan", user.Id)
	}
}

// handleSubscriptionDeleted downgrades user to free plan
func (h *Handler) handleSubscriptionDeleted(subscription stripe.Subscription) {
	log.Printf("❌ Subscription deleted: %s", subscription.ID)

	user, err := db.GetUserByStripeCustomerId(h.conn, subscription.Customer.ID)
	if err != nil {
		log.Printf("❌ Error finding user: %v", err)
		return
	}

	err = db.CancelUserSubscription(h.conn, user.Id)
	if err != nil {
		log.Printf("❌ Error canceling subscription: %v", err)
		return
	}

	log.Printf("⬇️  User %d downgraded to free plan", user.Id)
}

// getPlanFromSubscription determines the plan from a Stripe subscription
func (h *Handler) getPlanFromSubscription(subscription *stripe.Subscription) string {
	if len(subscription.Items.Data) == 0 {
		return "free"
	}

	priceID := subscription.Items.Data[0].Price.ID

	// Check against our configured price IDs
	if priceID == stripe_config.StripeConfig.ProMonthlyPriceID ||
		priceID == stripe_config.StripeConfig.ProYearlyPriceID {
		return "pro"
	}

	if priceID == stripe_config.StripeConfig.BusinessMonthlyPriceID ||
		priceID == stripe_config.StripeConfig.BusinessYearlyPriceID {
		return "business"
	}

	return "free"
}

// CreateCustomerPortalSessionHandler creates a Stripe customer portal session
func (h *Handler) CreateCustomerPortalSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	user, err := db.GetUserFromContext(h.conn, r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user's Stripe customer ID
	stripeCustomerID, err := db.GetStripeCustomerId(h.conn, user.Id)
	if err != nil || stripeCustomerID == "" {
		http.Error(w, "No subscription found", http.StatusNotFound)
		return
	}

	// Create customer portal session
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(stripeCustomerID),
		ReturnURL: stripe.String(stripe_config.StripeConfig.AppURL + "/dashboard"),
	}

	sess, err := billingportalsession.New(params)
	if err != nil {
		log.Printf("Error creating customer portal session: %v", err)
		http.Error(w, "Failed to create portal session", http.StatusInternalServerError)
		return
	}

	// Return the portal URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"url": sess.URL,
	})
}
