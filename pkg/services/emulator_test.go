package services

import (
	"testing"
	"time"

	"github.com/pepodev/super-utils/pkg/logger"
)

func init() {
	// Initialize test logger (no-op) for all tests
	logger.InitTestLogger()
}

func TestInitEmulator(t *testing.T) {
	emulator := InitEmulator()

	if emulator == nil {
		t.Fatal("InitEmulator() returned nil")
	}

	if len(emulator.services) == 0 {
		t.Error("InitEmulator() created emulator with no services")
	}

	expectedServices := []string{
		"api-gateway",
		"service-registry",
		"user-service",
		"catalog-service",
		"cart-service",
		"order-service",
		"payment-service",
		"review-service",
		"notification-service",
		"inventory-service",
	}

	for _, serviceID := range expectedServices {
		if _, exists := emulator.services[serviceID]; !exists {
			t.Errorf("Expected service %s not found in emulator", serviceID)
		}
	}
}

func TestGetEmulator(t *testing.T) {
	// Reset global emulator
	emulator = nil

	e1 := GetEmulator()
	if e1 == nil {
		t.Fatal("GetEmulator() returned nil")
	}

	e2 := GetEmulator()
	if e1 != e2 {
		t.Error("GetEmulator() should return the same instance (singleton)")
	}
}

func TestGetAllServices(t *testing.T) {
	em := InitEmulator()

	services := em.GetAllServices()

	if len(services) == 0 {
		t.Error("GetAllServices() returned no services")
	}

	if len(services) != 10 {
		t.Errorf("Expected 10 services, got %d", len(services))
	}

	// Verify all services have required fields
	for _, svc := range services {
		if svc.ID == "" {
			t.Error("Service with empty ID found")
		}
		if svc.Name == "" {
			t.Error("Service with empty Name found")
		}
		if svc.Status != "stopped" {
			t.Errorf("Expected initial status to be 'stopped', got %s", svc.Status)
		}
	}
}

func TestGetService(t *testing.T) {
	em := InitEmulator()

	tests := []struct {
		name        string
		serviceID   string
		expectError bool
	}{
		{
			name:        "Valid service - api-gateway",
			serviceID:   "api-gateway",
			expectError: false,
		},
		{
			name:        "Valid service - user-service",
			serviceID:   "user-service",
			expectError: false,
		},
		{
			name:        "Invalid service",
			serviceID:   "non-existent-service",
			expectError: true,
		},
		{
			name:        "Empty service ID",
			serviceID:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := em.GetService(tt.serviceID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if svc != nil {
					t.Error("Expected nil service but got non-nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if svc == nil {
					t.Error("Expected service but got nil")
				}
				if svc.ID != tt.serviceID {
					t.Errorf("Expected service ID %s, got %s", tt.serviceID, svc.ID)
				}
			}
		})
	}
}

func TestStartService(t *testing.T) {
	em := InitEmulator()

	// Test starting a valid service
	err := em.StartService("api-gateway")
	if err != nil {
		t.Errorf("Unexpected error starting service: %v", err)
	}

	svc, _ := em.GetService("api-gateway")
	if svc.Status != "running" {
		t.Errorf("Expected status 'running', got %s", svc.Status)
	}

	if svc.StartedAt.IsZero() {
		t.Error("StartedAt should be set after starting service")
	}

	// Test starting an already running service
	err = em.StartService("api-gateway")
	if err == nil {
		t.Error("Expected error when starting already running service")
	}

	// Test starting non-existent service
	err = em.StartService("non-existent")
	if err == nil {
		t.Error("Expected error when starting non-existent service")
	}
}

func TestStopService(t *testing.T) {
	em := InitEmulator()

	// First start a service
	em.StartService("user-service")

	// Test stopping a running service
	err := em.StopService("user-service")
	if err != nil {
		t.Errorf("Unexpected error stopping service: %v", err)
	}

	svc, _ := em.GetService("user-service")
	if svc.Status != "stopped" {
		t.Errorf("Expected status 'stopped', got %s", svc.Status)
	}

	if !svc.StartedAt.IsZero() {
		t.Error("StartedAt should be reset after stopping service")
	}

	// Test stopping an already stopped service
	err = em.StopService("user-service")
	if err == nil {
		t.Error("Expected error when stopping already stopped service")
	}

	// Test stopping non-existent service
	err = em.StopService("non-existent")
	if err == nil {
		t.Error("Expected error when stopping non-existent service")
	}
}

func TestResetAllServices(t *testing.T) {
	em := InitEmulator()

	// Start several services
	services := []string{"api-gateway", "user-service", "catalog-service"}
	for _, svcID := range services {
		em.StartService(svcID)
	}

	// Verify services are running
	for _, svcID := range services {
		svc, _ := em.GetService(svcID)
		if svc.Status != "running" {
			t.Errorf("Service %s should be running before reset", svcID)
		}
	}

	// Reset all services
	em.ResetAllServices()

	// Verify all services are stopped
	allServices := em.GetAllServices()
	for _, svc := range allServices {
		if svc.Status != "stopped" {
			t.Errorf("Service %s should be stopped after reset, got %s", svc.ID, svc.Status)
		}
		if !svc.StartedAt.IsZero() {
			t.Errorf("Service %s StartedAt should be reset", svc.ID)
		}
	}
}

func TestGetServiceStatus(t *testing.T) {
	em := InitEmulator()

	// Initially all services should be stopped
	status := em.GetServiceStatus()

	total, ok := status["total"].(int)
	if !ok || total != 10 {
		t.Errorf("Expected total to be 10, got %v", status["total"])
	}

	stopped, ok := status["stopped"].(int)
	if !ok || stopped != 10 {
		t.Errorf("Expected stopped to be 10, got %v", status["stopped"])
	}

	running, ok := status["running"].(int)
	if !ok || running != 0 {
		t.Errorf("Expected running to be 0, got %v", status["running"])
	}

	// Start some services
	em.StartService("api-gateway")
	em.StartService("user-service")
	em.StartService("catalog-service")

	status = em.GetServiceStatus()

	running = status["running"].(int)
	stopped = status["stopped"].(int)

	if running != 3 {
		t.Errorf("Expected 3 running services, got %d", running)
	}

	if stopped != 7 {
		t.Errorf("Expected 7 stopped services, got %d", stopped)
	}
}

func TestConcurrentAccess(t *testing.T) {
	em := InitEmulator()

	// Test concurrent reads and writes
	done := make(chan bool)

	// Start multiple goroutines trying to start/stop services
	for i := 0; i < 10; i++ {
		go func() {
			em.StartService("api-gateway")
			em.GetService("api-gateway")
			em.StopService("api-gateway")
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
	close(done)

	// Verify emulator is still in a valid state
	svc, err := em.GetService("api-gateway")
	if err != nil {
		t.Errorf("Error getting service after concurrent access: %v", err)
	}
	if svc == nil {
		t.Error("Service should still exist after concurrent access")
	}
}

func TestServiceStartedAtTimestamp(t *testing.T) {
	em := InitEmulator()

	before := time.Now()
	time.Sleep(10 * time.Millisecond)

	err := em.StartService("payment-service")
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	after := time.Now()

	svc, _ := em.GetService("payment-service")

	if svc.StartedAt.Before(before) {
		t.Error("StartedAt timestamp is before service was started")
	}

	if svc.StartedAt.After(after) {
		t.Error("StartedAt timestamp is after service was started")
	}
}

func BenchmarkStartService(b *testing.B) {
	em := InitEmulator()
	serviceIDs := []string{"api-gateway", "user-service", "catalog-service"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svcID := serviceIDs[i%len(serviceIDs)]
		em.StopService(svcID) // Ensure it's stopped
		em.StartService(svcID)
	}
}

func BenchmarkGetAllServices(b *testing.B) {
	em := InitEmulator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		em.GetAllServices()
	}
}
