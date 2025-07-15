package controllers

import (
	"net/http"
	"time"

	"github.com/Wembz/propanalytix/backend/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		UserService: services.NewUserService(),
	}
}

// GET /user/profile
func (uc *UserController) GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	user, err := uc.UserService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// PUT /user/profile
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := uc.UserService.UpdateProfile(userID, req.Email, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile Updated"})
}

// DELETE /user/deactivate
func (uc *UserController) DeactivateAccount(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	err := uc.UserService.DeactivateAccount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate account"})
		return
	}
	c.JSON( http.StatusOK, gin.H{"message": "Account deactivated"})
}

// PUT /user/plan
func (uc *UserController) SetPlan(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var req struct {
		Plan string `json:"plan"`
		Expiry string `json:"expiry"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	expiry, err := time.Parse(time.RFC3339, req.Expiry)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry format"})
		return
	} 

	err = uc.UserService.SetPlanStatus(userID, req.Plan, expiry)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plan updated"})
}

