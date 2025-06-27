package handlers

import (
	"context"
	"fmt"
	"ingestion-service/models"
	"ingestion-service/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EventHandler handles event-related HTTP requests
type EventHandler struct {
	kafkaService *services.KafkaService
	logger       *zap.Logger
}

// NewEventHandler creates a new event handler
func NewEventHandler(kafkaService *services.KafkaService, logger *zap.Logger) *EventHandler {
	return &EventHandler{
		kafkaService: kafkaService,
		logger:       logger,
	}
}

// TrackEvent handles event tracking requests
func (h *EventHandler) TrackEvent(c *gin.Context) {
	startTime := time.Now()
	requestID := uuid.New().String()

	// Add request ID to context for logging
	ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
	c.Request = c.Request.WithContext(ctx)

	h.logger.Info("Processing event tracking request",
		zap.String("request_id", requestID),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	// Parse and validate the event payload
	event, err := h.parseAndValidateEvent(c, requestID)
	if err != nil {
		h.logger.Error("Event validation failed",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return
	}

	// Enrich the event with metadata
	enrichedEvent := models.EnrichEvent(event, requestID)

	// Publish event to Kafka
	if err := h.publishEventToKafka(ctx, enrichedEvent, requestID); err != nil {
		h.logger.Error("Failed to publish event to Kafka",
			zap.String("request_id", requestID),
			zap.String("event_id", enrichedEvent.EventID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"KAFKA_ERROR",
			"Failed to process event",
			requestID,
		))
		return
	}

	// Calculate processing time
	processingTime := time.Since(startTime)
	enrichedEvent.ProcessingInfo.ProcessingMs = processingTime.Milliseconds()

	// Create response
	response := models.NewEventResponse(requestID)
	response.EventID = enrichedEvent.EventID

	h.logger.Info("Event processed successfully",
		zap.String("request_id", requestID),
		zap.String("event_id", enrichedEvent.EventID),
		zap.String("event_type", enrichedEvent.EventType),
		zap.String("user_id", enrichedEvent.UserID),
		zap.Duration("processing_time", processingTime),
	)

	// Send response back to frontend
	c.JSON(http.StatusOK, response)
}

// parseAndValidateEvent parses and validates the incoming event
func (h *EventHandler) parseAndValidateEvent(c *gin.Context, requestID string) (models.EventPayload, error) {
	var event models.EventPayload

	// Parse JSON payload
	if err := c.ShouldBindJSON(&event); err != nil {
		h.logger.Error("Failed to parse JSON payload",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_JSON",
			"Invalid JSON payload",
			requestID,
		))
		return models.EventPayload{}, err
	}

	// Validate required fields
	if err := h.validateEvent(event); err != nil {
		h.logger.Error("Event validation failed",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"VALIDATION_ERROR",
			err.Error(),
			requestID,
		))
		return models.EventPayload{}, err
	}

	return event, nil
}

// validateEvent validates the event payload
func (h *EventHandler) validateEvent(event models.EventPayload) error {
	if event.EventType == "" {
		return fmt.Errorf("event_type is required")
	}

	if event.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if event.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}

	if event.PageURL == "" {
		return fmt.Errorf("page_url is required")
	}

	// Validate timestamp format if provided
	if event.Timestamp != "" {
		if _, err := time.Parse(time.RFC3339, event.Timestamp); err != nil {
			return fmt.Errorf("invalid timestamp format, expected RFC3339")
		}
	}

	return nil
}

// publishEventToKafka publishes the enriched event to Kafka
func (h *EventHandler) publishEventToKafka(ctx context.Context, event models.EnrichedEvent, requestID string) error {
	h.logger.Debug("Publishing event to Kafka",
		zap.String("request_id", requestID),
		zap.String("event_id", event.EventID),
		zap.String("topic", "user-activity-events"),
	)

	// Publish to Kafka with retry logic
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := h.kafkaService.PublishEvent(ctx, event, event.EventID)
		if err == nil {
			return nil
		}

		h.logger.Warn("Kafka publish attempt failed",
			zap.String("request_id", requestID),
			zap.String("event_id", event.EventID),
			zap.Int("attempt", attempt),
			zap.Int("max_retries", maxRetries),
			zap.Error(err),
		)

		if attempt < maxRetries {
			// Exponential backoff
			backoffTime := time.Duration(attempt*attempt) * time.Millisecond * 100
			time.Sleep(backoffTime)
		}
	}

	return fmt.Errorf("failed to publish event to Kafka after %d attempts", maxRetries)
}

// HealthCheck handles health check requests
func (h *EventHandler) HealthCheck(c *gin.Context) {
	requestID := uuid.New().String()

	// Check Kafka service health
	kafkaHealth := "healthy"
	if err := h.kafkaService.HealthCheck(); err != nil {
		kafkaHealth = "unhealthy"
		h.logger.Error("Kafka health check failed",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
	}

	response := models.NewHealthResponse()
	response.Status = kafkaHealth

	c.JSON(http.StatusOK, response)
}

// GetStats returns service statistics
func (h *EventHandler) GetStats(c *gin.Context) {
	stats := map[string]interface{}{
		"service":   "ingestion-service",
		"version":   "1.0.0",
		"kafka":     h.kafkaService.GetStats(),
		"timestamp": time.Now().UTC(),
	}

	c.JSON(http.StatusOK, stats)
}
