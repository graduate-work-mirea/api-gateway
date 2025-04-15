package config

import (
	"os"
	"strconv"
)

// Config holds all the configuration for the application
type Config struct {
	Server     ServerConfig
	Auth       ServiceConfig
	ML         ServiceConfig
	DB         DatabaseConfig
	CacheSize  int
	JWTSecret  string
	CorsOrigin string
}

// ServerConfig holds the configuration for the API Gateway server
type ServerConfig struct {
	Port string
}

// ServiceConfig holds the configuration for external services
type ServiceConfig struct {
	Host string
	Port string
}

// DatabaseConfig holds the configuration for the database
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cacheSize, _ := strconv.Atoi(getEnv("CACHE_SIZE", "1000"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8000"),
		},
		Auth: ServiceConfig{
			Host: getEnv("AUTH_SERVICE_HOST", "localhost"),
			Port: getEnv("AUTH_SERVICE_PORT", "8080"),
		},
		ML: ServiceConfig{
			Host: getEnv("ML_SERVICE_HOST", "localhost"),
			Port: getEnv("ML_SERVICE_PORT", "6785"),
		},
		DB: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Name:     getEnv("POSTGRES_DB", "marketplace_data"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		CacheSize:  cacheSize,
		JWTSecret:  getEnv("JWT_SECRET", "your_secret_key_here"),
		CorsOrigin: getEnv("CORS_ORIGIN", "http://localhost"),
	}, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
