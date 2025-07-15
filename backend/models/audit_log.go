package models

import (
	"time"

	"gorm.io/gorm"
)

type Auditlog struct {
	ID 	uint	`gorm:"primaryKey" json:"id"`
	UserID	uint	`gorm:"index;not null" json:"user_id"`
	Action string	`gorm:"size:100;not null" json:"action"`
	Description string	`gorm:"size:255" json:"description"`
	IPAddress string	`gorm:"size:45" json:"ip_address"`
	CreatedAt time.Time	`json:"created_at"`
	DeletedAt gorm.DeletedAt	`gorm:"index" json:"-"`
	
}