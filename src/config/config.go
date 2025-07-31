package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      AppConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path string
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	DefaultPageSize int
	MaxPageSize     int
	Environment     string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "./src/db/todo.db"),
		},
		App: AppConfig{
			DefaultPageSize: getEnvAsInt("DEFAULT_PAGE_SIZE", 10),
			MaxPageSize:     getEnvAsInt("MAX_PAGE_SIZE", 100),
			Environment:     getEnv("ENV", "development"),
		},
	}
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets environment variable as integer with default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
