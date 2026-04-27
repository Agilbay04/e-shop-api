package middlewares

import (
	"e-shop-api/internal/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		requestID, _ := c.Get(RequestIDKey)

		// Process request
		c.Next()

		// If request is done, log it
		latency := time.Since(start)
		status := c.Writer.Status()

		// Logging HTTP request
		logger.L.Info("HTTP Request",
			zap.String("request_id", requestID.(string)),
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	}
}