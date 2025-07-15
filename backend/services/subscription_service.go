package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
)

type SubscriptionService struct{}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{}
}


func(s *SubscriptionService) UpgradeUserPlan(userID uint, plan string, expiry time.Time) error {
	if plan != "trial" && plan != "premium" {
		return errors.New("unsupported plan type")
	}

	var periodEnd time.Time
	switch plan {
	case "trial":
		periodEnd = time.Now().AddDate(0, 0, 14)	
	case "premium":
		periodEnd = time.Now().AddDate(0, 1, 0)	
	}

	sub := models.Subscription{
		UserID: userID,
		PlanName: plan,
		Status: "active",
		CurrentPeriodEnd: periodEnd,
	}

	// Save to subscriptions table
	if err := config.DB.
		Where("user_id = ?", userID).
		Assign(sub). 
		FirstOrCreate(&sub).Error; err != nil {
			return err
		}
	
	// Update user record with plan details
	if err := config.DB. 
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"plan": plan,
			"plan_start": time.Now(),
			"plan_expiry": periodEnd,
		}).Error; err != nil {
			return err
		}

		return nil
}

func (s *SubscriptionService) IsSubscriptionExpired(userID uint) bool {
	var sub models.Subscription
	err := config.DB.Where("user_id = ?", userID).First(&sub).Error
	if err != nil {
		return true
	}

	return sub.Status != "acive" || sub.CurrentPeriodEnd.Before(time.Now())
}

// CountMonthlyCalculations counts how many calculations a user has made this month
func (s *SubscriptionService) CountMonthlyCalculations(userID uint) (int64, error) {
	var count int64
	startOfMonth := time.Now().Truncate(24 *time.Hour).AddDate(0, 0, -time.Now().Day()+1)

	err := config.DB.Model(&models.Calculation{}).
		Where("user_id = ? AND created_at >= ?", userID, startOfMonth).
		Count(&count).Error

		return count, err
}

// DeleteCalculation removes a saved calculation
func (s *SubscriptionService) DeleteCalculation(userID, calcID uint) error {
	return config.DB. 
		Where("id = ? AND user_id = ?", calcID, userID). 
		Delete(&models.Calculation{}).Error
}

// GetRecentCalculations retrieves the N most recent calculations
func (s *SubscriptionService) GetRecentCalculations(userID uint, limit int) ([]models.Calculation, error) {
	var recent []models.Calculation
	err := config.DB.Where("user_id = ?", userID). 
		Order("created_at desc"). 
		Limit(limit). 
		Find(&recent).Error
	return recent, err	
}

func (s *SubscriptionService) HandleSubscriptionSuccess(userID uint, newPlan string) error {
	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	
	// Update plan and set expiry (e.g., 30 days from now)
	user.Plan = newPlan
	user.PlanExpiry = time.Now().Add(30 * 24 * time.Hour)

	if err := config.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update user plan: %w", err)
	}

	return nil
}

// FindUserIDByStripeCustomer returns the user ID linked to a Stripe customer
func (s *SubscriptionService) FindUserIDByStripeCustomer(CustomerID string) uint {
	var user models.User
	if err := config.DB.Where("stripe_customer_id = ?", CustomerID).First(&user).Error; err == nil {
		return user.ID
	}
	return 0
}

func (s *SubscriptionService) HandleSubscriptionCancel(userID uint) error {
	// Deactivate subscription
	if err := config.DB.Model(&models.Subscription{}). 
		Where("user_id = ?", userID). 
		Update("status", "cancelled").Error; err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}

	// Downgrade user
	if err := config.DB.Model(&models.User{}). 
	Where("id = ?", userID). 
	Updates(map[string]interface{}{
		"plan": 	"trial",
		"plan_expiry": time.Now().AddDate(0, 0, 7),
	}).Error; err != nil {
		return fmt.Errorf("failed to downgrade user plan: %w", err)
	}	

	return nil
}