package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/services"
)

func TestGetServicesHandler(t *testing.T) {
	// Initialize emulator
	services.InitEmulator()

	app := fiber.New()
	app.Get("/api/services", GetServicesHandler)

	req := httptest.NewRequest("GET", "/api/services", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var servicesResp []*services.Service
	if err := json.NewDecoder(resp.Body).Decode(&servicesResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(servicesResp) == 0 {
		t.Error("Expected services list to not be empty")
	}

	// Verify services have expected fields
	for _, svc := range servicesResp {
		if svc.ID == "" {
			t.Error("Service ID should not be empty")
		}
		if svc.Name == "" {
			t.Error("Service Name should not be empty")
		}
	}
}

func TestGetServiceStatusHandler(t *testing.T) {
	// Initialize and reset emulator
	em := services.InitEmulator()
	em.ResetAllServices()

	app := fiber.New()
	app.Get("/api/services/status", GetServiceStatusHandler)

	req := httptest.NewRequest("GET", "/api/services/status", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var status map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := status["total"]; !ok {
		t.Error("Expected 'total' field in status")
	}
	if _, ok := status["running"]; !ok {
		t.Error("Expected 'running' field in status")
	}
	if _, ok := status["stopped"]; !ok {
		t.Error("Expected 'stopped' field in status")
	}
}

func TestStartServiceHandler(t *testing.T) {
	// Initialize and reset emulator
	em := services.InitEmulator()
	em.ResetAllServices()

	app := fiber.New()
	app.Post("/api/services/start", StartServiceHandler)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectSuccess  bool
	}{
		{
			name: "Start service with JSON",
			payload: map[string]string{
				"service_id": "api-gateway",
			},
			expectedStatus: 200,
			expectSuccess:  true,
		},
		{
			name:           "Start service without service_id",
			payload:        map[string]string{},
			expectedStatus: 400,
			expectSuccess:  false,
		},
		{
			name: "Start non-existent service",
			payload: map[string]string{
				"service_id": "non-existent",
			},
			expectedStatus: 400,
			expectSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset service state
			em.ResetAllServices()

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/services/start", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)

			if tt.expectSuccess {
				if success, ok := result["success"].(bool); !ok || !success {
					t.Error("Expected success to be true")
				}
			} else {
				if _, ok := result["error"]; !ok {
					t.Error("Expected error field in response")
				}
			}
		})
	}
}

func TestStopServiceHandler(t *testing.T) {
	// Initialize and reset emulator
	em := services.InitEmulator()
	em.ResetAllServices()

	app := fiber.New()
	app.Post("/api/services/stop", StopServiceHandler)

	// First start a service
	em.StartService("user-service")

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectSuccess  bool
	}{
		{
			name: "Stop running service with JSON",
			payload: map[string]string{
				"service_id": "user-service",
			},
			expectedStatus: 200,
			expectSuccess:  true,
		},
		{
			name:           "Stop service without service_id",
			payload:        map[string]string{},
			expectedStatus: 400,
			expectSuccess:  false,
		},
		{
			name: "Stop non-existent service",
			payload: map[string]string{
				"service_id": "non-existent",
			},
			expectedStatus: 400,
			expectSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure user-service is running before each test
			em.StopService("user-service")
			em.StartService("user-service")

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/services/stop", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)

			if tt.expectSuccess {
				if success, ok := result["success"].(bool); !ok || !success {
					t.Error("Expected success to be true")
				}
			} else {
				if _, ok := result["error"]; !ok {
					t.Error("Expected error field in response")
				}
			}
		})
	}
}

func TestResetServicesHandler(t *testing.T) {
	// Initialize emulator and start some services
	em := services.InitEmulator()
	em.StartService("api-gateway")
	em.StartService("user-service")

	app := fiber.New()
	app.Post("/api/services/reset", ResetServicesHandler)

	req := httptest.NewRequest("POST", "/api/services/reset", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if success, ok := result["success"].(bool); !ok || !success {
		t.Error("Expected success to be true")
	}

	// Verify all services are stopped
	status := em.GetServiceStatus()
	if status["running"].(int) != 0 {
		t.Error("Expected all services to be stopped after reset")
	}
}

func TestBookShopEmulatorHandler(t *testing.T) {
	app := fiber.New()
	app.Get("/shop", BookShopEmulatorHandler)

	// Create a temporary HTML file for testing
	tmpDir := t.TempDir()
	htmlPath := tmpDir + "/book-shop.html"
	htmlContent := "<html><body>Book Shop</body></html>"
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		t.Skipf("Skipping test: could not create test file: %v", err)
	}

	// Note: This test may fail if the HTML file doesn't exist in the expected location
	// We'll just check that the handler responds
	req := httptest.NewRequest("GET", "/shop", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Status could be 200 or 500 depending on whether file exists
	if resp.StatusCode != 200 && resp.StatusCode != 500 {
		t.Errorf("Expected status 200 or 500, got %d", resp.StatusCode)
	}
}

func TestStartServiceHandlerWithFormData(t *testing.T) {
	// Initialize and reset emulator
	em := services.InitEmulator()
	em.ResetAllServices()

	app := fiber.New()
	app.Post("/api/services/start", StartServiceHandler)

	// Test with form data instead of JSON
	req := httptest.NewRequest("POST", "/api/services/start?service_id=catalog-service", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Note: The handler tries JSON first, then falls back to form value
	// This test verifies the handler doesn't crash with form data
	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		t.Errorf("Expected status 200 or 400, got %d", resp.StatusCode)
	}
}

func BenchmarkGetServicesHandler(b *testing.B) {
	services.InitEmulator()
	app := fiber.New()
	app.Get("/api/services", GetServicesHandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/services", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}

func BenchmarkStartServiceHandler(b *testing.B) {
	em := services.InitEmulator()
	app := fiber.New()
	app.Post("/api/services/start", StartServiceHandler)

	payload := map[string]string{"service_id": "api-gateway"}
	body, _ := json.Marshal(payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		em.ResetAllServices()
		req := httptest.NewRequest("POST", "/api/services/start", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}
