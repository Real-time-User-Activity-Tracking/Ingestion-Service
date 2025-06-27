package main

import (
	"context"
	"ingestion-service/config"
	"ingestion-service/handlers"
	"ingestion-service/router"
	"ingestion-service/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := initializeLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting ingestion service")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Kafka service
	kafkaService, err := initializeKafkaService(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Kafka service", zap.Error(err))
	}
	defer func() {
		if err := kafkaService.Close(); err != nil {
			logger.Error("Failed to close Kafka service", zap.Error(err))
		}
	}()

	// Initialize handlers
	eventHandler := handlers.NewEventHandler(kafkaService, logger)

	// Setup router with dependencies
	router := router.SetupRouter(eventHandler, logger)

	// Create HTTP server
	server := &http.Server{
		Addr:    cfg.GetServerAddress(),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", cfg.GetServerAddress()))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Display startup information
	displayStartupInfo(cfg, logger)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// initializeLogger sets up the application logger
func initializeLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	// Set log level based on environment
	if os.Getenv("LOG_LEVEL") == "debug" {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	return config.Build()
}

// initializeKafkaService creates and initializes the Kafka service
func initializeKafkaService(cfg *config.Config, logger *zap.Logger) (*services.KafkaService, error) {
	kafkaConfig := services.KafkaConfig{
		Brokers:         cfg.Kafka.Brokers,
		Topic:           cfg.Kafka.Topic,
		Acks:            cfg.Kafka.Acks,
		Retries:         cfg.Kafka.Retries,
		BatchSize:       cfg.Kafka.BatchSize,
		LingerMs:        cfg.Kafka.LingerMs,
		Compression:     cfg.Kafka.Compression,
		MaxMessageBytes: cfg.Kafka.MaxMessageBytes,
	}

	logger.Info("Initializing Kafka service",
		zap.Strings("brokers", kafkaConfig.Brokers),
		zap.String("topic", kafkaConfig.Topic),
		zap.String("compression", kafkaConfig.Compression),
	)

	return services.NewKafkaService(kafkaConfig, logger)
}

// displayStartupInfo prints startup information
func displayStartupInfo(cfg *config.Config, logger *zap.Logger) {
	logger.Info("Ingestion service started successfully",
		zap.String("address", cfg.GetServerAddress()),
		zap.String("health_endpoint", "/health"),
		zap.String("events_endpoint", "/api/v1/events/track"),
		zap.String("stats_endpoint", "/api/v1/stats"),
	)
}
