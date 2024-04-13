package helpers_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	println("TestMain is now running.")

	// Perform setup specific to tests, if any.
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Run all tests.
	code := m.Run()

	// Perform teardown, if necessary.

	// Exit with the code from the test run.
	os.Exit(code)
}
