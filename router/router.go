package router

import (
	"ingestion-service/handlers"
	"ingestion-service/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRouter configures and returns the Gin router with dependencies
func SetupRouter(eventHandler *handlers.EventHandler, logger *zap.Logger) *gin.Engine {
	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.ValidationMiddleware())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", eventHandler.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Event tracking endpoint
		api.POST("/events/track", eventHandler.TrackEvent)

		// Stats endpoint
		api.GET("/stats", eventHandler.GetStats)

		// Status endpoint (legacy)
		api.GET("/status", statusCheck)
	}

	logger.Info("Router configured successfully",
		zap.String("health_endpoint", "/health"),
		zap.String("events_endpoint", "/api/v1/events/track"),
		zap.String("stats_endpoint", "/api/v1/stats"),
	)

	return router
}

// statusCheck handles status check requests (legacy endpoint)
func statusCheck(c *gin.Context) {
	response := gin.H{
		"status":    "running",
		"service":   "ingestion-service",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
		"message":   "Service is ready to receive events",
		"endpoints": gin.H{
			"health": "/health",
			"events": "/api/v1/events/track",
			"stats":  "/api/v1/stats",
		},
	}
	c.JSON(http.StatusOK, response)
}
