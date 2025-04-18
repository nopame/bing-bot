package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	OpenBrowser       = false
	ResourceBlockLog  = false
	QueueSize         = 8
	MaxConcurrentJobs = 4
	API_URL           string
	AUTH_TOKEN        string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default values")
	}

	API_URL = getEnv("API_URL", "")
	AUTH_TOKEN = getEnv("AUTH_TOKEN", "")
}

// getEnv คืนค่าจาก os.Getenv ถ้าไม่มีให้คืนค่าดีฟอลต์
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}