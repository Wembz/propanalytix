package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	UserID        uint    `gorm:"not null;index" json:"user_id"`
	Amount        float64 `gorm:"not null" json:"amount"`
	Currency      string  `gorm:"default:'GBP'" json:"currency"`
	PaymentMethod string  `json:"payment_method"` // card, paypal, etc.
	StripeEventID string  `gorm:"uniqueIndex" json:"stripe_event_id"`
	Status        string  `json:"status"` // success, failed, refunded

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
