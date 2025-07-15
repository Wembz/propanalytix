package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	PlanName         string    `gorm:"not null" json:"plan_name"`
	Status           string    `gorm:"not null" json:"status"`
	StripeSubID      string    `gorm:"uniqueIndex" json:"stripe_sub_id"`
	CurrentPeriodEnd time.Time `json:"current_period_end"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
