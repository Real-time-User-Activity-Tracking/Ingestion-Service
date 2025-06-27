package models

import (
	"time"

	"github.com/google/uuid"
)

// EventPayload represents the event structure from frontend
type EventPayload struct {
	EventType  string                 `json:"event_type"`
	Timestamp  string                 `json:"timestamp"`
	UserID     string                 `json:"user_id"`
	SessionID  string                 `json:"session_id"`
	PageURL    string                 `json:"page_url"`
	EventData  map[string]interface{} `json:"event_data"`
	ClientInfo ClientInfo             `json:"client_info"`
}

// ClientInfo represents client information
type ClientInfo struct {
	UserAgent        string `json:"user_agent"`
	ScreenResolution string `json:"screen_resolution"`
	Language         string `json:"language"`
}

// EnrichedEvent represents the enriched event sent to Kafka
type EnrichedEvent struct {
	EventID        string                 `json:"event_id"`
	RequestID      string                 `json:"request_id"`
	EventType      string                 `json:"event_type"`
	Timestamp      time.Time              `json:"timestamp"`
	UserID         string                 `json:"user_id"`
	SessionID      string                 `json:"session_id"`
	PageURL        string                 `json:"page_url"`
	EventData      map[string]interface{} `json:"event_data"`
	ClientInfo     ClientInfo             `json:"client_info"`
	ServiceInfo    ServiceInfo            `json:"service_info"`
	ProcessingInfo ProcessingInfo         `json:"processing_info"`
}

// ServiceInfo represents service metadata
type ServiceInfo struct {
	ServiceName    string `json:"service_name"`
	ServiceVersion string `json:"service_version"`
	Environment    string `json:"environment"`
}

// ProcessingInfo represents processing metadata
type ProcessingInfo struct {
	ReceivedAt   time.Time `json:"received_at"`
	ProcessedAt  time.Time `json:"processed_at"`
	ProcessingMs int64     `json:"processing_ms"`
}

// EventResponse represents the response sent back to the client
type EventResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	EventID   string    `json:"event_id"`
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		RequestID string `json:"request_id"`
	} `json:"error"`
}

// NewEventResponse creates a new event response
func NewEventResponse(requestID string) EventResponse {
	return EventResponse{
		Status:    "success",
		Message:   "Event received and processed successfully",
		EventID:   uuid.New().String(),
		RequestID: requestID,
		Timestamp: time.Now().UTC(),
	}
}

// NewHealthResponse creates a new health response
func NewHealthResponse() HealthResponse {
	return HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "ingestion-service",
		Version:   "1.0.0",
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message, requestID string) ErrorResponse {
	return ErrorResponse{
		Error: struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"request_id"`
		}{
			Code:      code,
			Message:   message,
			RequestID: requestID,
		},
	}
}

// EnrichEvent creates an enriched event from the payload
func EnrichEvent(payload EventPayload, requestID string) EnrichedEvent {
	now := time.Now().UTC()

	// Parse timestamp from payload or use current time
	timestamp := now
	if payload.Timestamp != "" {
		if parsed, err := time.Parse(time.RFC3339, payload.Timestamp); err == nil {
			timestamp = parsed
		}
	}

	return EnrichedEvent{
		EventID:    uuid.New().String(),
		RequestID:  requestID,
		EventType:  payload.EventType,
		Timestamp:  timestamp,
		UserID:     payload.UserID,
		SessionID:  payload.SessionID,
		PageURL:    payload.PageURL,
		EventData:  payload.EventData,
		ClientInfo: payload.ClientInfo,
		ServiceInfo: ServiceInfo{
			ServiceName:    "ingestion-service",
			ServiceVersion: "1.0.0",
			Environment:    getEnvironment(),
		},
		ProcessingInfo: ProcessingInfo{
			ReceivedAt:   now,
			ProcessedAt:  now,
			ProcessingMs: 0, // Will be calculated in handler
		},
	}
}

// getEnvironment returns the current environment
func getEnvironment() string {
	// This could be enhanced to read from environment variables
	return "development"
}
