package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	FullName      string         `gorm:"size:100;not null" json:"full_name"`
	Email         string         `gorm:"unique;not null" json:"email"`
	PasswordHash  string         `gorm:"not null" json:"-"`
	Role          string         `gorm:"default:user" json:"role"`          // user, admin
	Plan          string         `gorm:"default:free" json:"plan"`          // free, trial, premium
	PlanStatus    string         `gorm:"default:active" json:"plan_status"` // active, expired
	FreeTrialUsed bool           `gorm:"default:false" json:"free_trial_used"`
	Templates     []Template     `gorm:"foreignKey:UserID"`
	Calculations  []Calculation  `gorm:"foreignKey:UserID"`
	PlanExpiry    time.Time      `json:"plan_expiry"`
	PlanStart     time.Time      `json:"plan_start"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
