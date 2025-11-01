package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/handlers"
	"github.com/pepodev/super-utils/pkg/metrics"
	"github.com/pepodev/super-utils/pkg/middleware"
)

func main() {
	app := fiber.New()

	// Apply middleware
	app.Use(middleware.RequestCounter())

	// Register health and metrics endpoints
	metrics.RegisterHealthMetrics(app)

	// Register application routes
	handlers.RegisterRoutes(app)

	// Start server
	app.Listen(":8000")
}
