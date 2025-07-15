package services

import (
	"errors"
	"time"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetByID(userID uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, userID).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *UserService) UpdateProfile(userID uint, email, name string) error {
	return config.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(models.User{Email: email, FullName: name}).Error
}

func (s *UserService) SetPlanStatus(userID uint, plan string, expiry time.Time) error {
	return config.DB.Model(&models.User{}).
		Where("id = ? ", userID).
		Updates(map[string]interface{}{
			"plan":        plan,
			"plan_expiry": expiry,
		}).Error
}

func (s *UserService) DeactivateAccount(userID uint) error {
	return config.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("active", false).Error
}
