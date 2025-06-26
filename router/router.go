package router

import (
	"ingestion-service/handlers"
	"ingestion-service/middleware"
	"ingestion-service/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router
func SetupRouter() *gin.Engine {
	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", healthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Initialize handlers
		eventHandler := handlers.NewEventHandler()

		// Event tracking endpoint
		api.POST("/events/track", eventHandler.TrackEvent)

		// Status endpoint
		api.GET("/status", statusCheck)
	}

	return router
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	response := models.NewHealthResponse()
	c.JSON(http.StatusOK, response)
}

// statusCheck handles status check requests
func statusCheck(c *gin.Context) {
	response := gin.H{
		"status":    "running",
		"service":   "ingestion-service",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
		"message":   "Service is ready to receive events",
	}
	c.JSON(http.StatusOK, response)
}
