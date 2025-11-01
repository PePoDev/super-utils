package handlers

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/services"
)

// BookShopEmulatorHandler serves the book shop emulator page
func BookShopEmulatorHandler(c *fiber.Ctx) error {
	// Read the HTML file at runtime
	htmlPath := filepath.Join("web", "book-shop.html")
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		htmlPath = filepath.Join("..", "web", "book-shop.html")
	}

	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		return c.Status(500).SendString("Error loading page: " + err.Error())
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Send(htmlContent)
}

// GetServicesHandler returns all services
func GetServicesHandler(c *fiber.Ctx) error {
	emulator := services.GetEmulator()
	allServices := emulator.GetAllServices()
	return c.JSON(allServices)
}

// GetServiceStatusHandler returns service status summary
func GetServiceStatusHandler(c *fiber.Ctx) error {
	emulator := services.GetEmulator()
	status := emulator.GetServiceStatus()
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
		return c.Status(400).JSON(fiber.Map{
			"error": "service_id is required",
		})
	}

	emulator := services.GetEmulator()
	err := emulator.StartService(req.ServiceID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	service, _ := emulator.GetService(req.ServiceID)
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
		return c.Status(400).JSON(fiber.Map{
			"error": "service_id is required",
		})
	}

	emulator := services.GetEmulator()
	err := emulator.StopService(req.ServiceID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	service, _ := emulator.GetService(req.ServiceID)
	return c.JSON(fiber.Map{
		"success": true,
		"service": service,
	})
}

// ResetServicesHandler resets all services
func ResetServicesHandler(c *fiber.Ctx) error {
	emulator := services.GetEmulator()
	emulator.ResetAllServices()

	return c.JSON(fiber.Map{
		"success": true,
		"message": "All services have been reset to stopped state",
	})
}
