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
