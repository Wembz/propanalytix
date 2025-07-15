package stripe

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/Wembz/propanalytix/backend/utils"
	"github.com/stripe/stripe-go"
)

var endpointSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

// HandleStripeWebhook processes incoming Stripe webhook events
func HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	 r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	 payload, err := io.ReadAll(r.Body)
	 if err != nil {
		http.Error(w, "Request too large", http.StatusServiceUnavailable)
		return
	 }


	// Verify signature
	sigHeader := r.Header.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		log.Printf("Webhook signature verification failed %v", err)
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	// Handle the event
	switch event.Type {

		
	case "checkout.session.completed": 
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Println("Failed to parse session:", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		handleCheckoutSessionCompleted(&session)
	
	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			log.Println("Failed to parse invoice:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		handleInvoicePaymentSucceeded(&invoice)
		
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

// handleCheckoutSessionCompleted upgrades user plan after checkout
func handleCheckoutSessionCompleted(session *stripe.CheckoutSession) {
	metadata := session.Metadata
	userID := metadata["user_id"]
	plan := metadata["plan"]

	if userID == "" || plan == "" {
		log.Println("Missing metadata in session")
	}

	uid, err := utils.ParseUint(userID)
	if err != nil {
		log.Printf("Invalid user ID in metatdata: %v", err)
		return
	}

	err = upgradeUserPlan(uid, plan)
	if err != nil {
		log.Printf("Failed to upgrade plan: %v", err)
	}

}


// handleInvoicePaymentSucceeded can be used to renew subscriptions
func handleInvoicePaymentSucceeded(invoice *stripe.Invoice) {
	customerID := invoice.Customer.ID

	log.Printf("Payment succeeded for customer: %s\n", customerID)
}

// upgradeUserPlan directly updates DB after payment
func upgradeUserPlan(userID uint, plan string) error {
	if plan != "premium" {
		return errors.New("only 'premium' plans are supported in Stripe webhook")
	}

	periodEnd := time.Now().AddDate(0, 1, 0)

	sub := models.Subscription{
		UserID: userID,
		PlanName: plan,
		Status: "active",
		CurrentPeriodEnd: periodEnd,
	}

	if err := config.DB. 
		Where("user_id = ?", userID).
		Assign(sub).
		FirstOrCreate(&sub).Error; err != nil {
			return err
		}
	return config.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"plan":         plan,
			"plan_start":   time.Now(),
			"plan_expiry":  periodEnd,
		}).Error	

}
