package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// KafkaService handles Kafka producer operations
type KafkaService struct {
	producer sarama.AsyncProducer
	config   KafkaConfig
	logger   *zap.Logger
	ctx      context.Context
	cancel   context.CancelFunc
}

// KafkaConfig holds Kafka configuration
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

// Message represents a Kafka message
type Message struct {
	Key   []byte
	Value []byte
	Topic string
}

// NewKafkaService creates a new Kafka service instance
func NewKafkaService(config KafkaConfig, logger *zap.Logger) (*KafkaService, error) {
	ctx, cancel := context.WithCancel(context.Background())

	service := &KafkaService{
		config: config,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}

	if err := service.initializeProducer(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize Kafka producer: %w", err)
	}

	// Start error handling goroutine
	go service.handleErrors()
	go service.handleSuccesses()

	return service, nil
}

// initializeProducer sets up the Kafka producer

func (ks *KafkaService) initializeProducer() error {
	config := sarama.NewConfig()

	// Set acknowledgment level
	switch ks.config.Acks {
	case "all":
		config.Producer.RequiredAcks = sarama.WaitForAll
	case "1":
		config.Producer.RequiredAcks = sarama.WaitForLocal
	case "0":
		config.Producer.RequiredAcks = sarama.NoResponse
	default:
		return fmt.Errorf("invalid acks value: %s", ks.config.Acks)
	}

	// Set retry configuration
	config.Producer.Retry.Max = ks.config.Retries
	config.Producer.Retry.Backoff = time.Millisecond * 100

	// Set batch configuration
	config.Producer.Flush.Bytes = ks.config.BatchSize
	config.Producer.Flush.Frequency = time.Duration(ks.config.LingerMs) * time.Millisecond

	// Set compression
	switch ks.config.Compression {
		case "none":
			config.Producer.Compression = sarama.CompressionNone
		case "gzip":
			config.Producer.Compression = sarama.CompressionGZIP
		case "snappy":
			config.Producer.Compression = sarama.CompressionSnappy
		case "lz4":
			config.Producer.Compression = sarama.CompressionLZ4
		case "zstd":
			config.Producer.Compression = sarama.CompressionZSTD
		default:
			return fmt.Errorf("invalid compression: %s", ks.config.Compression)
	}

	// Set message size limit
	config.Producer.MaxMessageBytes = ks.config.MaxMessageBytes

	// Enable idempotent producer for exactly-once semantics
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	// Create async producer
	producer, err := sarama.NewAsyncProducer(ks.config.Brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	ks.producer = producer
	ks.logger.Info("Kafka producer initialized successfully",
		zap.Strings("brokers", ks.config.Brokers),
		zap.String("topic", ks.config.Topic),
	)

	return nil
}

// PublishMessage sends a message to Kafka
func (ks *KafkaService) PublishMessage(ctx context.Context, key string, value interface{}) error {
	select {
		case <-ks.ctx.Done():
			return fmt.Errorf("kafka service is shutting down")
		default:
	}

	// Serialize value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message value: %w", err)
	}

	// Create Kafka message
	message := &sarama.ProducerMessage{
		Topic:     ks.config.Topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(jsonValue),
		Timestamp: time.Now(),
	}

	// Send message asynchronously
	select {
		case ks.producer.Input() <- message:
			ks.logger.Debug("Message sent to Kafka",
				zap.String("topic", ks.config.Topic),
				zap.String("key", key),
				zap.Int("size", len(jsonValue)),
			)
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while sending message")
		case <-ks.ctx.Done():
			return fmt.Errorf("kafka service is shutting down")
	}
}

// PublishEvent publishes an event to Kafka with proper error handling
func (ks *KafkaService) PublishEvent(ctx context.Context, event interface{}, eventID string) error {
	return ks.PublishMessage(ctx, eventID, event)
}

// handleErrors processes Kafka producer errors
func (ks *KafkaService) handleErrors() {
	for {
		select {
		case err := <-ks.producer.Errors():
			keyBytes, _ := err.Msg.Key.Encode()
			ks.logger.Error("Kafka producer error",
				zap.Error(err),
				zap.String("topic", err.Msg.Topic),
				zap.String("key", string(keyBytes)),
			)
		case <-ks.ctx.Done():
			return
		}
	}
}

// handleSuccesses processes successful Kafka messages
func (ks *KafkaService) handleSuccesses() {
	for {
		select {
		case msg := <-ks.producer.Successes():
			keyBytes, _ := msg.Key.Encode()
			ks.logger.Debug("Message successfully sent to Kafka",
				zap.String("topic", msg.Topic),
				zap.Int32("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
				zap.String("key", string(keyBytes)),
			)
		case <-ks.ctx.Done():
			return
		}
	}
}

// HealthCheck checks if the Kafka service is healthy
func (ks *KafkaService) HealthCheck() error {
	select {
	case <-ks.ctx.Done():
		return fmt.Errorf("kafka service is shutting down")
	default:
	}

	// Check if producer is still active
	if ks.producer == nil {
		return fmt.Errorf("kafka producer is not initialized")
	}

	return nil
}

// Close gracefully shuts down the Kafka service
func (ks *KafkaService) Close() error {
	ks.logger.Info("Shutting down Kafka service")

	// Cancel context to stop goroutines
	ks.cancel()

	// Close producer
	if ks.producer != nil {
		if err := ks.producer.Close(); err != nil {
			ks.logger.Error("Error closing Kafka producer", zap.Error(err))
			return fmt.Errorf("failed to close Kafka producer: %w", err)
		}
	}

	ks.logger.Info("Kafka service shut down successfully")
	return nil
}

// GetStats returns Kafka producer statistics
func (ks *KafkaService) GetStats() map[string]interface{} {
	if ks.producer == nil {
		return map[string]interface{}{
			"status": "not_initialized",
		}
	}

	return map[string]interface{}{
		"status":      "active",
		"brokers":     ks.config.Brokers,
		"topic":       ks.config.Topic,
		"compression": ks.config.Compression,
		"acks":        ks.config.Acks,
		"retries":     ks.config.Retries,
	}
}
