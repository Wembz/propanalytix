package controllers

import (
	"net/http"
	"time"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
	"github.com/Wembz/propanalytix/backend/services"
	"github.com/gin-gonic/gin"
)

// UpgradeMembership handles user plan upgrade request

type SubscriptionController struct {
	SubscriptionService *services.SubscriptionService
}

func NewSubscriptionController() *SubscriptionController {
	return &SubscriptionController{
		SubscriptionService: services.NewSubscriptionService(),
	}
}

// PUT /subscription/upgrade
func (sc *SubscriptionController) UpgradeMembership(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req struct {
		Plan string `json:"plan"`
		Expiry string `json:"expiry,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Plan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plan type is required"})
		return
	}

	var expiry time.Time
	var err error
	if req.Expiry != "" {
		expiry, err = time.Parse(time.RFC3339, req.Expiry)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid  expiry date format"})
			return
		}
	}else {
		expiry = time.Now().AddDate(0, 1, 0)
	}

	err = sc.SubscriptionService.UpgradeUserPlan(userID, req.Plan, expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Membership upgraded successfully"})
}

func (sc *SubscriptionController) GetPlanStatus(c *gin.Context) {
	userID := c.MustGet("user_ID").(uint)
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plan":  user.Plan,
		"plan_start": user.PlanStart,
		"plan_expiry": user.PlanExpiry,
		"is_expired": time.Now().After(user.PlanExpiry),
	})
}