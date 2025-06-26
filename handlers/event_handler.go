package handlers

import (
	"encoding/json"
	"fmt"
	"ingestion-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EventHandler handles event-related HTTP requests
type EventHandler struct{}

// NewEventHandler creates a new event handler
func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

// TrackEvent handles event tracking requests
func (h *EventHandler) TrackEvent(c *gin.Context) {
	// Generate request ID for tracking
	requestID := uuid.New().String()

	// Parse the event payload
	var event models.EventPayload
	if err := c.ShouldBindJSON(&event); err != nil {
		fmt.Printf("❌ [%s] Failed to parse JSON: %v\n", requestID, err)
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("INVALID_JSON", "Invalid JSON payload", requestID))
		return
	}

	// Display the received event data
	h.displayEvent(event, requestID)

	// Create response
	response := models.NewEventResponse(requestID)

	fmt.Printf("✅ [%s] Event processed successfully\n", requestID)

	// Send response back to frontend
	c.JSON(http.StatusOK, response)
}

// displayEvent prints the event details to console
func (h *EventHandler) displayEvent(event models.EventPayload, requestID string) {
	fmt.Printf("\n📥 [%s] Received Event:\n", requestID)
	fmt.Printf("   Event Type: %s\n", event.EventType)
	fmt.Printf("   User ID: %s\n", event.UserID)
	fmt.Printf("   Session ID: %s\n", event.SessionID)
	fmt.Printf("   Page URL: %s\n", event.PageURL)
	fmt.Printf("   Timestamp: %s\n", event.Timestamp)

	// Display event data
	if len(event.EventData) > 0 {
		fmt.Printf("   Event Data:\n")
		for key, value := range event.EventData {
			fmt.Printf("     %s: %v\n", key, value)
		}
	}

	// Display client info
	fmt.Printf("   Client Info:\n")
	fmt.Printf("     User Agent: %s\n", event.ClientInfo.UserAgent)
	fmt.Printf("     Screen Resolution: %s\n", event.ClientInfo.ScreenResolution)
	fmt.Printf("     Language: %s\n", event.ClientInfo.Language)

	// Pretty print the full event for debugging
	if eventJSON, err := json.MarshalIndent(event, "", "  "); err == nil {
		fmt.Printf("   Full Event JSON:\n%s\n", string(eventJSON))
	}
}
