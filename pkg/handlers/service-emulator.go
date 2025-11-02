package handlers

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/logger"
	"github.com/pepodev/super-utils/pkg/services"
	"go.uber.org/zap"
)

// BookShopEmulatorHandler serves the book shop emulator page
func BookShopEmulatorHandler(c *fiber.Ctx) error {
	logger.Debug("Serving book shop emulator page")

	// Read the HTML file at runtime
	htmlPath := filepath.Join("web", "book-shop.html")
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		htmlPath = filepath.Join("..", "web", "book-shop.html")
	}

	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		logger.Error("Failed to load book shop HTML",
			zap.String("path", htmlPath),
			zap.Error(err),
		)
		return c.Status(500).SendString("Error loading page: " + err.Error())
	}

	logger.Info("Book shop emulator page served")
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Send(htmlContent)
}

// GetServicesHandler returns all services
func GetServicesHandler(c *fiber.Ctx) error {
	logger.Debug("Fetching all services")

	emulator := services.GetEmulator()
	allServices := emulator.GetAllServices()

	logger.Info("Services list retrieved",
		zap.Int("service_count", len(allServices)),
	)

	return c.JSON(allServices)
}

// GetServiceStatusHandler returns service status summary
func GetServiceStatusHandler(c *fiber.Ctx) error {
	logger.Debug("Fetching service status summary")

	emulator := services.GetEmulator()
	status := emulator.GetServiceStatus()

	logger.Info("Service status retrieved",
		zap.Int("running", status["running"].(int)),
		zap.Int("stopped", status["stopped"].(int)),
	)

	return c.JSON(status)
}

// StartServiceRequest represents a start service request
type StartServiceRequest struct {
	ServiceID string `json:"service_id"`
}

// StartServiceHandler starts a service
func StartServiceHandler(c *fiber.Ctx) error {
	var req StartServiceRequest

	// Try to parse as JSON
	if err := c.BodyParser(&req); err != nil {
		// Fallback to form value
		req.ServiceID = c.FormValue("service_id")
	}

	if req.ServiceID == "" {
		logger.Warn("Start service request missing service_id",
			zap.String("client_ip", c.IP()),
		)
		return c.Status(400).JSON(fiber.Map{
			"error": "service_id is required",
		})
	}

	logger.Info("Starting service",
		zap.String("service_id", req.ServiceID),
		zap.String("client_ip", c.IP()),
	)

	emulator := services.GetEmulator()
	err := emulator.StartService(req.ServiceID)
	if err != nil {
		logger.Error("Failed to start service",
			zap.String("service_id", req.ServiceID),
			zap.Error(err),
		)
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	service, _ := emulator.GetService(req.ServiceID)
	logger.Info("Service started successfully",
		zap.String("service_id", req.ServiceID),
	)

	return c.JSON(fiber.Map{
		"success": true,
		"service": service,
	})
}

// StopServiceHandler stops a service
func StopServiceHandler(c *fiber.Ctx) error {
	var req StartServiceRequest

	// Try to parse as JSON
	if err := c.BodyParser(&req); err != nil {
		// Fallback to form value
		req.ServiceID = c.FormValue("service_id")
	}

	if req.ServiceID == "" {
		logger.Warn("Stop service request missing service_id",
			zap.String("client_ip", c.IP()),
		)
		return c.Status(400).JSON(fiber.Map{
			"error": "service_id is required",
		})
	}

	logger.Info("Stopping service",
		zap.String("service_id", req.ServiceID),
		zap.String("client_ip", c.IP()),
	)

	emulator := services.GetEmulator()
	err := emulator.StopService(req.ServiceID)
	if err != nil {
		logger.Error("Failed to stop service",
			zap.String("service_id", req.ServiceID),
			zap.Error(err),
		)
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	service, _ := emulator.GetService(req.ServiceID)
	logger.Info("Service stopped successfully",
		zap.String("service_id", req.ServiceID),
	)

	return c.JSON(fiber.Map{
		"success": true,
		"service": service,
	})
}

// ResetServicesHandler resets all services
func ResetServicesHandler(c *fiber.Ctx) error {
	logger.Info("Resetting all services",
		zap.String("client_ip", c.IP()),
	)

	emulator := services.GetEmulator()
	emulator.ResetAllServices()

	logger.Info("All services reset successfully")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "All services have been reset to stopped state",
	})
}
