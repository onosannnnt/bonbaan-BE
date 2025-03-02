package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func Initenv() string {
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
	return configPath
}
