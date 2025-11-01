package services

import (
	"fmt"
	"sync"
	"time"
)

// Service represents a microservice in the book shop architecture
type Service struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Port        int       `json:"port"`
	Database    string    `json:"database"`
	Status      string    `json:"status"` // "stopped" or "running"
	StartedAt   time.Time `json:"started_at,omitempty"`
}

// ServiceEmulator manages the state of all microservices
type ServiceEmulator struct {
	mu       sync.RWMutex
	services map[string]*Service
}

var emulator *ServiceEmulator

// InitEmulator initializes the service emulator with all book shop services
func InitEmulator() *ServiceEmulator {
	emulator = &ServiceEmulator{
		services: make(map[string]*Service),
	}

	// Define all book shop microservices
	services := []Service{
		{
			ID:          "api-gateway",
			Name:        "API Gateway",
			Description: "Routes requests to appropriate services, handles authentication and rate limiting",
			Port:        8000,
			Database:    "N/A",
			Status:      "stopped",
		},
		{
			ID:          "service-registry",
			Name:        "Service Registry",
			Description: "Service discovery and registration, health monitoring",
			Port:        8761,
			Database:    "N/A",
			Status:      "stopped",
		},
		{
			ID:          "user-service",
			Name:        "User Service",
			Description: "Manages user accounts, profiles, authentication, and authorization",
			Port:        8001,
			Database:    "PostgreSQL",
			Status:      "stopped",
		},
		{
			ID:          "catalog-service",
			Name:        "Catalog Service",
			Description: "Manages book inventory, product information, and search functionality",
			Port:        8002,
			Database:    "MongoDB",
			Status:      "stopped",
		},
		{
			ID:          "cart-service",
			Name:        "Shopping Cart Service",
			Description: "Manages user shopping carts and cart items",
			Port:        8003,
			Database:    "Redis",
			Status:      "stopped",
		},
		{
			ID:          "order-service",
			Name:        "Order Service",
			Description: "Handles order processing, management, and tracking",
			Port:        8004,
			Database:    "PostgreSQL",
			Status:      "stopped",
		},
		{
			ID:          "payment-service",
			Name:        "Payment Service",
			Description: "Processes payments and manages payment transactions",
			Port:        8005,
			Database:    "MySQL",
			Status:      "stopped",
		},
		{
			ID:          "review-service",
			Name:        "Review Service",
			Description: "Manages book reviews and ratings from customers",
			Port:        8006,
			Database:    "MongoDB",
			Status:      "stopped",
		},
		{
			ID:          "notification-service",
			Name:        "Notification Service",
			Description: "Sends notifications to users (email, SMS, in-app)",
			Port:        8007,
			Database:    "Redis",
			Status:      "stopped",
		},
		{
			ID:          "inventory-service",
			Name:        "Inventory Service",
			Description: "Manages real-time inventory levels and stock tracking",
			Port:        8008,
			Database:    "PostgreSQL",
			Status:      "stopped",
		},
	}

	for _, svc := range services {
		svcCopy := svc
		emulator.services[svc.ID] = &svcCopy
	}

	return emulator
}

// GetEmulator returns the global emulator instance
func GetEmulator() *ServiceEmulator {
	if emulator == nil {
		InitEmulator()
	}
	return emulator
}

// GetAllServices returns all services
func (e *ServiceEmulator) GetAllServices() []*Service {
	e.mu.RLock()
	defer e.mu.RUnlock()

	services := make([]*Service, 0, len(e.services))
	for _, svc := range e.services {
		services = append(services, svc)
	}
	return services
}

// GetService returns a specific service by ID
func (e *ServiceEmulator) GetService(id string) (*Service, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	svc, exists := e.services[id]
	if !exists {
		return nil, fmt.Errorf("service not found: %s", id)
	}
	return svc, nil
}

// StartService starts a service
func (e *ServiceEmulator) StartService(id string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	svc, exists := e.services[id]
	if !exists {
		return fmt.Errorf("service not found: %s", id)
	}

	if svc.Status == "running" {
		return fmt.Errorf("service is already running: %s", id)
	}

	svc.Status = "running"
	svc.StartedAt = time.Now()
	return nil
}

// StopService stops a service
func (e *ServiceEmulator) StopService(id string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	svc, exists := e.services[id]
	if !exists {
		return fmt.Errorf("service not found: %s", id)
	}

	if svc.Status == "stopped" {
		return fmt.Errorf("service is already stopped: %s", id)
	}

	svc.Status = "stopped"
	svc.StartedAt = time.Time{}
	return nil
}

// ResetAllServices stops all services
func (e *ServiceEmulator) ResetAllServices() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, svc := range e.services {
		svc.Status = "stopped"
		svc.StartedAt = time.Time{}
	}
}

// GetServiceStatus returns a summary of service statuses
func (e *ServiceEmulator) GetServiceStatus() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	running := 0
	stopped := 0

	for _, svc := range e.services {
		if svc.Status == "running" {
			running++
		} else {
			stopped++
		}
	}

	return map[string]interface{}{
		"total":   len(e.services),
		"running": running,
		"stopped": stopped,
	}
}
