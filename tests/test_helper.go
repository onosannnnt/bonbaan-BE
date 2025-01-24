package tests

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func SetupTestEnvironment() {
	err := godotenv.Load("../.env.test")
	if err != nil {
		println("Error loading .env.test file")
	}
	// Add any other setup code here
}

func TeardownTestEnvironment() {
	// Add any cleanup code here
}

func TestMain(m *testing.M) {
	SetupTestEnvironment()
	code := m.Run()
	TeardownTestEnvironment()
	os.Exit(code)
}
