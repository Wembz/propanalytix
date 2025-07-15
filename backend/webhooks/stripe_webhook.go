package webhook

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v74"
    "github.com/stripe/stripe-go/v74/webhook"

	"github.com/Wembz/propanalytix/backend/services"
	
	

	"github.com/gin-gonic/gin"
)

func HandleStripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Request too large"})
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	sighHeader := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, endpointSecret, sighHeader)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	switch event.Type {
	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err == nil {
			customerID := invoice.Customer.ID
			subService := services.NewSubscriptionService()
			// Lookup userID via Stripe Customer ID → implement logic
			userID := subService.FindUserIDByStripeCustomer(customerID)
			if userID == 0 {
				log.Println("User not found for customer ID:", customerID)
				return
			}

			err := subService.HandleSubscriptionSuccess(userID, "premium")
			if err != nil {
				log.Println("Failed to update subscriptions:", err)
			}
		}
	case "customer.subscription.delete": 
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err == nil {
			subService := services.NewSubscriptionService()

			customerID := subscription.Customer.ID
			userID := subService.FindUserIDByStripeCustomer(customerID)

			if userID != 0 {
				err := subService.HandleSubscriptionCancel(userID)
				if err != nil {
					log.Println("Failed to cancel subscription", err)
				}
			}
			
		}	
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})

}