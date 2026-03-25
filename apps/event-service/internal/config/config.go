package config

import (
	"log"
	"os"
)

// Config хранит конфигурацию приложения.
type Config struct {
	Port    string
	DBHost  string
	DBPort  string
	DBName  string
	DBUser  string
	DBPass  string
	SSLMode string
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return value
}

// Load загружает конфигурацию приложения из env.
func Load() Config {
	return Config{
		Port:    getEnv("PORT", "8080"),
		DBHost:  getEnv("DB_HOST", "localhost"),
		DBPort:  getEnv("DB_PORT", "5432"),
		DBName:  getEnv("DB_NAME", "events"),
		DBUser:  mustGetEnv("DB_USER"),
		DBPass:  mustGetEnv("DB_PASSWORD"),
		SSLMode: getEnv("DB_SSLMODE", "disable"),
	}
}
