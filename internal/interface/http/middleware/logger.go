package middleware

import (
	"time"

	"github.com/Azzurriii/slythr/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// Log request details
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Log request details with formatted message
		switch {
		case status >= 500:
			logger.Default.Errorf("Server error - %s %s [%d] %v", method, path, status, latency)
		case status >= 400:
			logger.Default.Warnf("Client error - %s %s [%d] %v", method, path, status, latency)
		default:
			logger.Default.Infof("Request processed - %s %s [%d] %v", method, path, status, latency)
		}
	}
}
