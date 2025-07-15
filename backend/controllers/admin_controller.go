package controllers

import (
	"net/http"
	"time"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
	"github.com/gin-gonic/gin"
)

// AdminDashboard gives a basic admin summary (e.g. usage, revenue)
func AdminDashboard(c *gin.Context) {
	var userCount int64
	var calcCount int64
	var recentUsers []models.User

	config.DB.Model(&models.User{}).Count(&userCount)
	config.DB.Model(&models.Calculation{}).Count(&calcCount)
	config.DB.Order("created_at desc").Limit(5).Find(&recentUsers)

	c.JSON(http.StatusOK, gin.H{
		"user_count": userCount,
		"calc_count": calcCount,
		"recent_users": recentUsers,
		"timestamp": time.Now(),
	})
}

// GetAllUsers lists all users for admin
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func AdminUpgradeUserPlan(c *gin.Context) {
	var input struct {
		UserID uint `json:"user_id" binding:"required"`
		Plan string `json:"plan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", input.UserID).Updates(models.User{
		Plan:     input.Plan,
		PlanStart: time.Now(),
	}).Error; err != nil { 
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User plan updated successfully"})
}
