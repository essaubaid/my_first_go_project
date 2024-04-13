package testing_config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
