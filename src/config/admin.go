package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AdminEmail string
var AdminUsername string
var AdminPassword string

func init() {
	configPath := Initenv()

	// Load the .env file
	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatalf("Problem loading .env file: %v", err)
		os.Exit(-1)
	}

	AdminEmail = os.Getenv("ADMIN_EMAIL")
	AdminUsername = os.Getenv("ADMIN_USERNAME")
	AdminPassword = os.Getenv("ADMIN_PASSWORD")
}
