package config

import (
	"os"
	"strings"
)

// Config holds the application configuration.
type Config struct {
	HTTPAddr    string
	SQLiteDSN   string
	LogLevel    string
	CORSAllowed []string
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	return &Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		SQLiteDSN:   getEnv("SQLITE_DSN", "./data/todos.db"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		CORSAllowed: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
