package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			fmt.Printf("✅ [CORS] Preflight request handled - allowing actual request\n")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// For actual requests, continue to the handler
		fmt.Printf("➡️ [CORS] Actual request proceeding to handler\n")
		c.Next()
	}
}
