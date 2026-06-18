package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	APP_ENV          string
	APP_LOCALE       string
	APP_SECRET       string
	APP_JWT_DURATION string
	DB_URL           string
	API_URL          string
	APP_WEB_URL      string
	APP_PORT         string
	APP_LOG          bool
)

func setup() {
	APP_ENV = getEnv("APP_ENV", "local")
	APP_LOCALE = getEnv("APP_LOCALE", "en")
	APP_SECRET = getEnv("APP_SECRET", "secret")
	APP_JWT_DURATION = getEnv("APP_JWT_DURATION", "1d")
	DB_URL = getEnv("DB_URL", "postgresql://postgres:password@localhost:5432/portfolio")
	APP_WEB_URL = getEnv("APP_WEB_URL", "http://localhost:3000")
	API_URL = getEnv("API_URL", "http://localhost:8000")
	APP_PORT = getEnv("APP_PORT", "8000")
	APP_LOG = getEnv("APP_LOG", "true") == "true"
}

func InitConfigApp() {
	godotenv.Load(".env")
	setup()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
