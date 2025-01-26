package Config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

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
    _, filename, _, _ := runtime.Caller(0)
    currentFileDir := filepath.Dir(filename)

    // Walk up the directory tree to find the .env file
    var configPath string
    for {
        configPath = filepath.Join(currentFileDir, ".env.test")
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

	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbSchema = os.Getenv("DB_SCHEMA")

}
