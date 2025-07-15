package stripe

import (
	"log"
	"os"

	"github.com/stripe/stripe-go"
)

// InitStripe sets the Stripe API key from environment
func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		log.Println("STRIPE_SECRET_KEY not set in environment")
	}
}
