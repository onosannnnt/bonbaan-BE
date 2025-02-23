package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	OmisePublicKey string
	OmiseSecretKey string
)

func init() {
	configPath := Initenv()

	// Load the .env file
	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatalf("Problem loading .env file: %v", err)
		os.Exit(-1)
	}
	OmisePublicKey = os.Getenv("OMISE_PUBLIC_KEY")
	OmiseSecretKey = os.Getenv("OMISE_SECRET_KEY")
}
