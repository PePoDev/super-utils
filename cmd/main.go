package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/handlers"
	"github.com/pepodev/super-utils/pkg/logger"
	"github.com/pepodev/super-utils/pkg/metrics"
	"github.com/pepodev/super-utils/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting Super-Utils application",
		zap.String("version", "1.0.0"),
		zap.Int("body_limit_mb", 10),
	)

	app := fiber.New(fiber.Config{
		BodyLimit:             10 * 1024 * 1024, // 10MB limit for speed test uploads
		DisableStartupMessage: true,             // Disable Fiber banner, we use structured logging
	})

	// Apply middleware
	app.Use(middleware.RequestCounter())
	logger.Info("Middleware configured",
		zap.String("middleware", "request_counter"),
	)

	// Register health and metrics endpoints
	metrics.RegisterHealthMetrics(app)
	logger.Info("Health and metrics endpoints registered")

	// Register application routes
	handlers.RegisterRoutes(app)
	logger.Info("Application routes registered")

	// Setup graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		logger.Info("Received shutdown signal, gracefully shutting down")
		app.Shutdown()
	}()

	// Start server
	logger.Info("Server starting", zap.String("address", ":8000"))
	if err := app.Listen(":8000"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
	logger.Info("Server stopped")
}
