package controllers

import (
	"net/http"
	"time"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
    "github.com/Wembz/propanalytix/backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input *models.RegisterInput 

	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error securing password"})
		return
	} 

	user := models.User {
		FullName: input.Name,
		Email: input.Email,
		PasswordHash: string(hashPassword),
		Plan: "trial",
		PlanExpiry: time.Now().AddDate(0, 0, 7),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	token , err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Login(c *gin.Context) {
	var input  *models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
		if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	if user.PlanExpiry.Before(time.Now()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Subscription expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
