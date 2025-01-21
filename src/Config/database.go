package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbSchema   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbSchema = os.Getenv("DB_SCHEMA")

}
