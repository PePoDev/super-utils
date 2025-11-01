package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/metrics"
)

// RequestCounter middleware increments the total request count
func RequestCounter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		metrics.TotalRequests++
		return c.Next()
	}
}
