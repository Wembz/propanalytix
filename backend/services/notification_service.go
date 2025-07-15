package services

import (
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
)

type NotificationService struct {}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// NotifyExpiringSubscriptions sends reminder emails for users with expiring plans
func (s *NotificationService) NotifyExpiringSubscriptions() error {
	var users []models.User
	soon := time.Now().Add(72 * time.Hour)

	err := config.DB.Where("plan_expiry BETWEEN ? AND ? AND active = ?", time.Now(), soon, true).Find(&users).Error
	if err != nil {
		return err
	}

	for _, user := range users {
		msg := fmt.Sprintf("Subject: Your plan id expiring soon\n\nHi %s,\n\nYour current plan will expire on %s. Please upgrade to continue uninterrupted access", user.FullName, user.PlanExpiry.Format("Jan 2, 2006"))
		s.send(user.Email, msg)
	}

	return nil
}

// SendCalculationSummary sends an optional summary email after analysis
func (s *NotificationService) SendCalculationSummary(email string, summary string) error {
	msg := fmt.Sprintf("Subject: Your Calculation Summary\n\n%s", summary)
	return s.send(email, msg)
}

// Core email sender (uses basic SMTP)
func (s *NotificationService) send(to string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, []byte(body))
}

func (s *NotificationService) SendSummaryEmail(to string, subject string, body string) error {

	// Replace with your own email config
	from := "youremail@example.com"
	password := "yourpassword"

	//Mail server config
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	// Message
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject" + subject + "\n\n" +
		body

	//Auth	 
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}