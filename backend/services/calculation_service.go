package services

import (
	"errors"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
	
)
type CalculationService struct {}

func NewCalculationService()  *CalculationService {
	return &CalculationService{}
}

func(cs *CalculationService) SaveCalculations(userID uint, inputData, results string, templateID *uint, exportFormat string) error {
	calculation := models.Calculation{
		UserID: userID,
		InputData: inputData,
		Results: results,
		Exported: exportFormat != "",
		ExportedFormat: exportFormat,
		TemplateID: templateID,
	}
	return config.DB.Create(&calculation).Error

}

// GetCalculationHistory retrieves all past calculations for a user
func(cs *CalculationService) GetCalculationHistory(userID uint)([]models.Calculation, error) {
	var history []models.Calculation
	err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&history).Error
	return history, err
}

// GetCalculationByID returns a specific saved calculation
func(cs *CalculationService) GetCalculationByID(userID, calcID uint) (*models.Calculation, error) {
	var calc models.Calculation
	err := config.DB.Where("id = ? AND user_id = ?", calcID, userID).First(&calc).Error
	if err != nil {
		return nil, errors.New("calculation not found")
	}

	return &calc, nil
}

// SaveCalculation stores the result of a user's analysis
func(cs *CalculationService) SaveCalculation(userID uint, inputData, result string, templateID *uint, exportFormat string)error{
	calculation := models.Calculation{
		UserID: userID,
		InputData: inputData,
		Results: result,
		Exported: exportFormat != "",
		ExportedFormat: exportFormat,
		TemplateID: templateID,
	}
	return config.DB.Create(&calculation).Error
}

func (cs *CalculationService) DeleteCalculation(userID, calculationID uint) error {
	return config.DB.
	Where("id = ? AND user_id = ?", calculationID, userID). 
	Delete(&models.Calculation{}).
	Error
}

func (cs *CalculationService) GetRecentCalculations(userID uint, limit int) ([]models.Calculation, error) {
	var calculations []models.Calculation
	err := config.DB. 
	Where("user_id = ?", userID).
	Order("created_at DESC"). 
	Limit(5). 
	Find(&calculations).Error

	return calculations, err
}