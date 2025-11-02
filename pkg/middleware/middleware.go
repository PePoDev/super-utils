package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/logger"
	"github.com/pepodev/super-utils/pkg/metrics"
	"go.uber.org/zap"
)

// RequestCounter middleware increments the total request count
func RequestCounter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		method := c.Method()

		// Increment counter
		metrics.TotalRequests++

		// Process request
		err := c.Next()

		// Log request details
		duration := time.Since(start)
		statusCode := c.Response().StatusCode()

		// Determine log level based on status code
		if statusCode >= 500 {
			logger.Error("Request completed with error",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.Duration("duration", duration),
				zap.String("client_ip", c.IP()),
			)
		} else if statusCode >= 400 {
			logger.Warn("Request completed with client error",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.Duration("duration", duration),
				zap.String("client_ip", c.IP()),
			)
		} else {
			logger.Info("Request completed",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.Duration("duration", duration),
			)
		}

		return err
	}
}
