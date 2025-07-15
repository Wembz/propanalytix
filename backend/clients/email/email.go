package email

import (
	"log"
	"os"
	"gopkg.in/gomail.v2"
	
)


// NewEmailClient returns a preconfigured SMTP email client (e.g. Gmail or SendGrid)
func NewEmailClient() *gomail.Dialer {
	host := os.Getenv("EMAIL_HOST")
	port := 587
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	
	if host == "" || username == "" || password == "" {
		log.Println("Email credentials not set correctly in environment")
	}

	return gomail.NewDialer(host, port, username, password)
}
