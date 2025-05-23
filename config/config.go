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
	API_PREFIX        string
	AUTH_TOKEN        string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default values")
	}

	API_URL = getEnv("API_URL", "")
	API_PREFIX = getEnv("API_PREFIX", "")
	AUTH_TOKEN = getEnv("AUTH_TOKEN", "")
	OpenBrowser = getEnvBool("OPEN_BROWSER", false)
	MaxConcurrentJobs = getEnvInt("MAX_CONCURRENT_JOBS", 4)
}

// คืนค่า string
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// คืนค่า bool เช่น "true", "1" = true
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

// คืนค่า int เช่น "4" = 4
func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}