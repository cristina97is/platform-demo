package config

import "os"

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

// getenvOrDefault возвращает значение переменной окружения,
// а если она не задана — значение по умолчанию.
func getenvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Load загружает конфигурацию приложения из env.
func Load() Config {
	return Config{
		Port:    getenvOrDefault("PORT", "8080"),
		DBHost:  getenvOrDefault("DB_HOST", "localhost"),
		DBPort:  getenvOrDefault("DB_PORT", "5432"),
		DBName:  getenvOrDefault("DB_NAME", "events"),
		DBUser:  getenvOrDefault("DB_USER", "events"),
		DBPass:  getenvOrDefault("DB_PASSWORD", "events"),
		SSLMode: getenvOrDefault("DB_SSLMODE", "disable"),
	}
}
