package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	JwtSecret string
)

func init() {
	configPath := Initenv()

	// Load the .env file
	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatalf("Problem loading .env file: %v", err)
		os.Exit(-1)
	}

	JwtSecret = os.Getenv("JWT_SECRET")

}
