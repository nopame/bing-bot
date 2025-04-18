package config

import (
	"log"
	"os"
	"strconv"

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
	OpenBrowser = getEnvBool("OPEN_BROWSER", false)
}

// คืนค่า string จาก ENV
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// คืนค่า bool จาก ENV เช่น "true", "1" = true
func getEnvBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return b
}