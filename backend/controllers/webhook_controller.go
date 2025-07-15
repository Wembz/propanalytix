package controllers

import (
	"encoding/json"

	"io"
	"log"
	"net/http"
	"os"

	"github.com/Wembz/propanalytix/backend/services"
	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/webhook"
)

type WebhookController struct {
	SubscriptionService *services.SubscriptionService
}

func NewWebhookController(subService *services.SubscriptionService) *WebhookController {
	return &WebhookController{SubscriptionService: subService}
}

func (wc *WebhookController) HandleStripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Request too large"})
		return
	}

	sig := c.GetHeader("Stripe-Signature")
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")


	event, err := webhook.ConstructEvent(payload, sig, endpointSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("Failed to parse session: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session data"})
			return
		}	

		userID := session.Metadata["user_id"]
		plan := session.Metadata["plan"]
		log.Printf("Checkout completed: User %s bought plan %s", userID, plan)
		
	case "invoice.payment_succeeded": 
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice data"})
			return
		}	
		log.Printf("Payment succeeded for customer %v", invoice.Customer.ID)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			log.Printf("Failed to parse subscription %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription data"})
			return
		}	
		log.Printf("Subscription cancelled for customer %v", subscription.Customer.ID)

	default: 
		log.Printf("Unhandled event type %s", event.Type)	
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}