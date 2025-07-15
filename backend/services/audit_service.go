package services

import (
	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
)

func LogAction(userID uint, action, description, ip string) {
	audit := models.Auditlog{
		UserID: userID,
		Action: action,
		Description: description,
		IPAddress: ip,
	}

	config.DB.Create(&audit)
}

