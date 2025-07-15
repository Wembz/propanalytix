package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	UserID      uint   `gorm:"not null" json:"user_id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Data        string `gorm:"type:jsonb;not null" json:"data"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
