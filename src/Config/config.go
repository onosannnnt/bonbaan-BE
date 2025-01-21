package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Port string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	Port = os.Getenv("PORT")
}
