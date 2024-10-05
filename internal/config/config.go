package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	PaytmApiKey       string
	PaytmSecretKey    string
	PaytmApiBaseUrl   string
	PaytmLoginBaseUrl string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file; proceeding with defaults")
	}

	config := &Config{
		Port: getEnv("PORT", "8080"),
		PaytmApiKey:       getEnv("PAYTM_API_KEY", ""),
		PaytmSecretKey:    getEnv("PAYTM_CLIENT_SECRET", ""),
		PaytmApiBaseUrl:   getEnv("PAYTM_API_BASE_URL", ""),
		PaytmLoginBaseUrl: getEnv("PAYTM_LOGIN_BASE_URL", ""), // Default value
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
