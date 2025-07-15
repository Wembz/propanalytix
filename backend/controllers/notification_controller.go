package controllers

import (
	"net/http"

	"github.com/Wembz/propanalytix/backend/models"
	"github.com/Wembz/propanalytix/backend/services"
	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationService *services.NotificationService
}

func NewNotificationController() *NotificationController {
	return &NotificationController{
		NotificationService: services.NewNotificationService(),
	}
}

// POST /notifications/summary
func(nc *NotificationController) SendSummaryEmail(c *gin.Context) {
	subject := "Your Property Summary"
	body := "Here is your latest property analysis summary"

	var user models.User

	err := nc.NotificationService.SendSummaryEmail(user.Email, subject, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send summary email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Summary email sent"})
}

// POST /notifications/expiry-alerts (admin/cron use)
func (nc *NotificationController) NotifyExpiringSubscriptions(c *gin.Context) {
	err := nc.NotificationService.NotifyExpiringSubscriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to notify expiry subscriptions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Expiry notifications sent"})
} 