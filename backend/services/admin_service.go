package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
)

func FetchAllUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.Find(&users).Error
	return users, err
} 

func GenerateRevenueReport() (map[string]interface{}, error) {
	//Placeholder for revenue logic
	return map[string]interface{}{
		"monthly_total": 2000,
		"active_subscription": 248,
	}, nil
}

func FetchAuditLogs() ([]models.Auditlog, error) {
	var logs []models.Auditlog
	err := config.DB.Order("created_at desc").Limit(100).Find(&logs).Error
	return logs, err
}

func GenerateCalculationAnalytics() (map[string]interface{}, error) {
	var count int64
	config.DB.Model(&models.Calculation{}).Count(&count)

	return map[string]interface{}{
		"total_calculations": count,
		"generated_at": time.Now(),
	},nil
}

func BanUserByID(userID string) error {
	return config.DB.Model(&models.User{}).Where("id = ?",userID).Update("status", "banned").Error
}

func UpdateUserRole(userID, role string) error {
	return config.DB.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

func DeleteUserByID(userID string) error {
	return config.DB.Delete(&models.User{}, userID).Error
}

func GeneratePlanStatistics() (map[string]int64, error) {
	var trial, premium int64
	config.DB.Model(&models.User{}).Where("plan = ?", "trial").Count(&trial)
	config.DB.Model(&models.User{}).Where("plan = ?", "premium").Count(&premium)
	return map[string]int64{
		"trial": trial,
		"premium": premium,
	}, nil
}

func ExportUserAsCSV() ([]byte, error) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	_ =writer.Write([]string{"ID","Name", "Email", "Plan", "CreatedAt"})
	for _, u := range users {
		_ = writer.Write([]string{
			fmt.Sprint(u.ID),
			u.FullName,
			u.Email,
			u.Plan,
			u.CreatedAt.Format(time.RFC3339),
		})
	}
	writer.Flush()
	return buf.Bytes(), nil
}