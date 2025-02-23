package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	MailHost     string
	MailPort     int
	MailUser     string
	MailPassword string
)

func init() {
	configPath := Initenv()

	// Load the .env file
	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatalf("Problem loading .env file: %v", err)
		os.Exit(-1)
	}

	mailport := os.Getenv("MAIL_PORT")
	MailHost = os.Getenv("MAIL_HOST")
	MailPort, err = strconv.Atoi(mailport)
	if err != nil {
		log.Fatalf("Problem parsing MAIL_PORT: %v", err)
	}
	MailUser = os.Getenv("MAIL_USER")
	MailPassword = os.Getenv("MAIL_PASSWORD")
}
