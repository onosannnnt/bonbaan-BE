package Config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
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
	_, filename, _, _ := runtime.Caller(0)
	currentFileDir := filepath.Dir(filename)

	// Walk up the directory tree to find the .env file
	var configPath string
	for {
		configPath = filepath.Join(currentFileDir, ".env")
		if _, err := os.Stat(configPath); !os.IsNotExist(err) {
			break
		}
		parentDir := filepath.Dir(currentFileDir)
		if parentDir == currentFileDir {
			log.Fatalf(".env file not found")
			os.Exit(-1)
		}
		currentFileDir = parentDir
	}

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
