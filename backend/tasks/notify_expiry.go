package tasks

import (
	"log"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/services"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.InitDB()

	notifier := services.NewNotificationService()
	if err := notifier.NotifyExpiringSubscriptions(); err != nil {
		log.Fatal("Error sending notification: %v", err)
	}

	log.Println("Notification task completed successfully")
}