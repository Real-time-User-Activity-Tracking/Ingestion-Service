package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	Server ServerConfig
	Kafka  KafkaConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
}

// KafkaConfig holds Kafka-related configuration
type KafkaConfig struct {
	Brokers         []string
	Topic           string
	Acks            string
	Retries         int
	BatchSize       int
	LingerMs        int
	Compression     string
	MaxMessageBytes int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "9094"),
			Host: getEnv("HOST", "0.0.0.0"),
		},
		Kafka: KafkaConfig{
			Brokers:         parseBrokers(getEnv("KAFKA_BROKERS", "localhost:9092")),
			Topic:           getEnv("KAFKA_TOPIC", "user-activity-events"),
			Acks:            getEnv("KAFKA_ACKS", "all"),
			Retries:         getEnvAsInt("KAFKA_RETRIES", 3),
			BatchSize:       getEnvAsInt("KAFKA_BATCH_SIZE", 16384),
			LingerMs:        getEnvAsInt("KAFKA_LINGER_MS", 5),
			Compression:     getEnv("KAFKA_COMPRESSION", "snappy"),
			MaxMessageBytes: getEnvAsInt("KAFKA_MAX_MESSAGE_BYTES", 1000000),
		},
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

// validate performs configuration validation
func (c *Config) validate() error {
	if len(c.Kafka.Brokers) == 0 {
		return fmt.Errorf("at least one Kafka broker must be specified")
	}

	if c.Kafka.Topic == "" {
		return fmt.Errorf("Kafka topic must be specified")
	}

	validAcks := map[string]bool{"all": true, "1": true, "0": true}
	if !validAcks[c.Kafka.Acks] {
		return fmt.Errorf("invalid Kafka acks value: %s", c.Kafka.Acks)
	}

	if c.Kafka.Retries < 0 {
		return fmt.Errorf("Kafka retries must be non-negative")
	}

	if c.Kafka.BatchSize <= 0 {
		return fmt.Errorf("Kafka batch size must be positive")
	}

	if c.Kafka.LingerMs < 0 {
		return fmt.Errorf("Kafka linger ms must be non-negative")
	}

	validCompression := map[string]bool{"none": true, "gzip": true, "snappy": true, "lz4": true, "zstd": true}
	if !validCompression[c.Kafka.Compression] {
		return fmt.Errorf("invalid Kafka compression: %s", c.Kafka.Compression)
	}

	if c.Kafka.MaxMessageBytes <= 0 {
		return fmt.Errorf("Kafka max message bytes must be positive")
	}

	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets environment variable as integer with fallback
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// parseBrokers parses comma-separated broker list
func parseBrokers(brokers string) []string {
	if brokers == "" {
		return []string{}
	}
	return strings.Split(brokers, ",")
}
