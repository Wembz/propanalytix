package services

import (
	"errors"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
)

type TemplateService struct{}

func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

func (s *TemplateService) Save(userID uint, name string, data string) error {
	template := models.Template{
		UserID: userID,
		Name: name,
		Data: data,
	}
	return config.DB.Create(&template).Error
}

func (s *TemplateService) GetAll(userID uint) ([]models.Template, error) {
	var templates []models.Template
	err := config.DB.Where("user_id = ?", userID).Find(&templates).Error
	return templates, err
}

func (s *TemplateService) GetByID(userID, templateID uint) (*models.Template, error) {
	var tpl models.Template
	err := config.DB.Where("id = ? AND user_id = ? ", templateID, userID).First(&tpl).Error
	if err != nil {
		return nil, errors.New("template not found")
	}
	return &tpl, nil
}

func (s *TemplateService) Update(userID, templateID uint, name string, data string) error {
	return config.DB.Model(&models.Template{}). 
		Where("id = ? AND user_id = ?", userID, templateID).
		Updates(models.Template{Name: name, Data: data}).Error
}

func (s *TemplateService) Delete(UserID, templateID uint) error {
	return config.DB.Where("id = ? AND user_id = ? ", templateID, UserID).Delete(&models.Template{}).Error
}
