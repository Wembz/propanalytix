package airtable


import (
	"log"
	"os"

	"github.com/mehanizm/airtable"
)



func NewAirtableClient() *airtable.Client {
	apiKey := os.Getenv("AIRTABLE_API_KEY not set")
	if apiKey == "" {
		log.Println("AIRTABLE_API_KEY not set")
	}
	return airtable.NewClient(apiKey)
}