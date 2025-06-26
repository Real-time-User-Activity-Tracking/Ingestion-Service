package main

import (
	"fmt"
	"ingestion-service/config"
	"ingestion-service/router"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode to debug for development
	gin.SetMode(gin.DebugMode)

	// Setup router
	router := router.SetupRouter()

	// Display startup information
	displayStartupInfo(cfg)

	// Start server
	if err := router.Run(cfg.GetServerAddress()); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// displayStartupInfo prints startup information
func displayStartupInfo(cfg *config.Config) {
	fmt.Printf("🚀 Starting Ingestion Service on %s...\n", cfg.GetServerAddress())
	fmt.Printf("📡 Health check: http://%s/health\n", cfg.GetServerAddress())
	fmt.Printf("📊 Status: http://%s/api/v1/status\n", cfg.GetServerAddress())
	fmt.Printf("📝 Event tracking: http://%s/api/v1/events/track\n", cfg.GetServerAddress())
}
