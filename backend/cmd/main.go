package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Wembz/propanalytix/backend/clients"
	"github.com/mehanizm/airtable"
	"github.com/Wembz/propanalytix/backend/clients/stripe"
	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
	"github.com/Wembz/propanalytix/backend/routes"
	"github.com/Wembz/propanalytix/backend/utils"
    "gopkg.in/gomail.v2"
	"github.com/joho/godotenv"
)

func main () {

	var (
		emailClient *gomail.Dialer
		airtableClient *airtable.Client 
	)

	//Utils logger
	utils.InitLogger()
	utils.Log.Println("Server is starting...")


	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with environment variables...")
	}
	
	// Connect to PostgreSQL
	config.InitDB()

	//Utils Validator
	utils.InitValidator()

	 // Auto migrate models
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Calculation{},
		&models.Template{},
		&models.Subscription{},
		&models.Auditlog{},
	)
	if err != nil {
		log.Fatalf("Migation failed: %v", err)
	}

	// 🔐 Initialize external services
	stripe.InitStripe()
	emailClient = clients.NewEmailClient()
	airtableClient = clients.NewAirtableClient()

	// Set up router
	router := routes.SetupRouter(emailClient, airtableClient)
	router.Run("8080")



	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}