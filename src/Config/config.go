package Config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

var Port string

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

	Port = os.Getenv("PORT")
}
