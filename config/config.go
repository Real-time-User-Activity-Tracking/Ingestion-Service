package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Server ServerConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	port := getEnv("PORT", "9094")
	host := getEnv("HOST", "0.0.0.0")

	return &Config{
		Server: ServerConfig{
			Port: port,
			Host: host,
		},
	}
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
