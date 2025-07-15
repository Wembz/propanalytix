package services

import (
	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
)

func RecordPayment(userID uint, amount float64, eventID string, method string) error {
	payment := models.Payment{
		UserID: userID,
		Amount: amount,
		Status: "success",
		StripeEventID: eventID,
		PaymentMethod: method,
	}

	return config.DB.Create(&payment).Error
}