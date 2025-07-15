package repository

import (
	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
)

type SubscriptionRepository struct {}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{}
} 

// UpsertSubscription creates or updates a user's subscription
func (r *SubscriptionRepository) UpsertSubscription(sub *models.Subscription) error {
	return config.DB. 
	Where("user_id = ?", sub.UserID). 
	Assign(sub). 
	FirstOrCreate(sub).Error
}

// GetSubscriptionByUserID fetches a subscription by user ID
func (r *SubscriptionRepository)  GetSubscriptionByUserID(userID uint) (*models.Subscription, error) {
	var sub models.Subscription
	err := config.DB.Where("user_id = ?", userID).First(&sub).Error 
	if err != nil {
		return nil, err
	}

	return &sub, err
}