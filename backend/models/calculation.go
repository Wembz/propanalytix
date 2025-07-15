package models

import (
	"time"

	"gorm.io/gorm"
)

type Calculation struct {
	ID             uint		`gorm:"primaryKey" json:"id"`
	UserID         uint		`gorm:"not null" json:"user_id"`
	TemplateID     *uint	`json:"template_id"`
	InputData      string	`gorm:"type:jsonb;not null" json:"input_data" `
	Results        string	`gorm:"type:jsonb;not null" json:"results"`
	Exported       bool		`gorm:"default:false" json:"exported"`
	ExportedFormat string	`gorm:"size:20" json:"exported_format"`

	CreatedAt time.Time		`json:"created_at"`
	UpdatedAt time.Time		`json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`	
}
