package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware creates a logging middleware for HTTP requests
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log structured information
		logger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
			zap.Int("status_code", param.StatusCode),
			zap.Int("body_size", param.BodySize),
			zap.Duration("latency", param.Latency),
			zap.String("error", param.ErrorMessage),
		)

		// Return empty string as we're handling logging ourselves
		return ""
	})
}
