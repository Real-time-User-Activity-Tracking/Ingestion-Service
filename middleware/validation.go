package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidationMiddleware creates a validation middleware for requests
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate Content-Type for POST requests
		if c.Request.Method == "POST" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": gin.H{
						"code":    "INVALID_CONTENT_TYPE",
						"message": "Content-Type must be application/json",
					},
				})
				c.Abort()
				return
			}
		}

		// Validate request size (optional)
		if c.Request.ContentLength > 10*1024*1024 { // 10MB limit
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": gin.H{
					"code":    "REQUEST_TOO_LARGE",
					"message": "Request body too large",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
