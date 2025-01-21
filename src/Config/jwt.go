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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	JwtSecret = os.Getenv("JWT_SECRET")

}
