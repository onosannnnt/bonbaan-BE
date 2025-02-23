package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Port string
var BucketName string
var BucketKey string

func init() {
	configPath := Initenv()

	// Load the .env file
	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatalf("Problem loading .env file: %v", err)
		os.Exit(-1)
	}

	Port = os.Getenv("PORT")
	BucketName = os.Getenv("BUCKET_NAME")
	BucketKey = os.Getenv("BUCKET_KEY")
}
