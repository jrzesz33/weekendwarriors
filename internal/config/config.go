package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	Port        int
	Environment string
	DatabaseURL string
	CORSOrigins []string
	LogLevel    string
}

// Load loads configuration from environment variables with sensible defaults
func Load() *Config {
	cfg := &Config{
		Port:        getEnvAsInt("PORT", 8080),
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseURL: getEnv("DATABASE_URL", "data/golf_gamez.db"),
		CORSOrigins: getEnvAsSlice("CORS_ORIGINS", []string{
			"http://localhost:3000",
			"http://localhost:8000",
			"http://localhost:8080",
			"https://golfgamez.com",
		}),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvAsSlice gets an environment variable as slice or returns a default value
func getEnvAsSlice(key string, defaultVal []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultVal
}
