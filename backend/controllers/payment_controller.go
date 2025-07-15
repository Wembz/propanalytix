package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

func CreateCheckoutSession(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	plan := c.Query("plan")

	if plan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing plan type"})
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	successURL := os.Getenv("FRONTEND_URL") + "/payment_success"
	cancelURL := os.Getenv("FRONTEND_URL") + "/payment_cancel"

	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(getPriceIDForPlan(plan)),
				Quantity: stripe.Int64(1),
			},
		},

		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String(cancelURL),
		ClientReferenceID: stripe.String(string(rune(userID))),
	}

	s, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"checkout_url": s.URL})
}

func getPriceIDForPlan(plan string) string {
	switch plan {
	case "pro":
		return os.Getenv("STRIPE_PRICE_ID_PRO")
	case "premium":
		return os.Getenv("STRIPE_PRICE_ID_PREMIUM")
	default:
		return ""
	}
}
