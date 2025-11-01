package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/handlers"
	"github.com/pepodev/super-utils/pkg/metrics"
	"github.com/pepodev/super-utils/pkg/middleware"
)

// setupTestApp creates a test Fiber app with all routes registered
func setupTestApp() *fiber.App {
	app := fiber.New()

	// Apply middleware
	app.Use(middleware.RequestCounter())

	// Register health and metrics endpoints
	metrics.RegisterHealthMetrics(app)

	// Register application routes
	handlers.RegisterRoutes(app)

	return app
}

func TestIntegrationHealthCheck(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var health map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if health["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", health["status"])
	}
}

func TestIntegrationReadinessAndLiveness(t *testing.T) {
	app := setupTestApp()

	tests := []struct {
		name     string
		endpoint string
		field    string
	}{
		{"Readiness", "/ready", "ready"},
		{"Liveness", "/live", "alive"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.endpoint, nil)
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

			if result[tt.field] != true {
				t.Errorf("Expected %s to be true", tt.field)
			}
		})
	}
}

func TestIntegrationMetricsEndpoint(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/metrics", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var metricsData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&metricsData); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify key metrics are present
	if _, ok := metricsData["uptime"]; !ok {
		t.Error("Expected uptime in metrics")
	}
	if _, ok := metricsData["goroutines"]; !ok {
		t.Error("Expected goroutines in metrics")
	}
	if _, ok := metricsData["total_requests"]; !ok {
		t.Error("Expected total_requests in metrics")
	}
}

func TestIntegrationPrometheusMetrics(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/metrics/prometheus", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/plain; version=0.0.4" {
		t.Errorf("Expected Content-Type 'text/plain; version=0.0.4', got %s", contentType)
	}
}

func TestIntegrationRequestCounter(t *testing.T) {
	app := setupTestApp()

	// Record initial count
	initialCount := metrics.TotalRequests

	// Make several requests
	endpoints := []string{"/health", "/ready", "/live", "/metrics"}
	for _, endpoint := range endpoints {
		req := httptest.NewRequest("GET", endpoint, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Errorf("Failed to make request to %s: %v", endpoint, err)
			continue
		}
		resp.Body.Close()
	}

	// Verify counter was incremented
	if metrics.TotalRequests != initialCount+len(endpoints) {
		t.Errorf("Expected TotalRequests to be %d, got %d", initialCount+len(endpoints), metrics.TotalRequests)
	}
}

func TestIntegrationSystemInfo(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/api/sysinfo", nil)
	resp, err := app.Test(req, 5000) // 5 second timeout for system info
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var sysinfo map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&sysinfo); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(sysinfo) == 0 {
		t.Error("Expected system info to not be empty")
	}

	t.Logf("System info has %d fields", len(sysinfo))
}

func TestIntegrationServiceEmulator(t *testing.T) {
	app := setupTestApp()

	// Reset all services first
	resetReq := httptest.NewRequest("POST", "/api/services/reset", nil)
	resetResp, _ := app.Test(resetReq)
	resetResp.Body.Close()

	// Get all services
	req := httptest.NewRequest("GET", "/api/services", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to get services: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var services []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&services)

	if len(services) == 0 {
		t.Fatal("Expected services list to not be empty")
	}

	t.Logf("Found %d services", len(services))
}

func TestIntegrationServiceLifecycle(t *testing.T) {
	app := setupTestApp()

	// Reset all services
	resetReq := httptest.NewRequest("POST", "/api/services/reset", nil)
	resetResp, _ := app.Test(resetReq)
	resetResp.Body.Close()

	// Check initial status - all should be stopped
	statusReq := httptest.NewRequest("GET", "/api/services/status", nil)
	statusResp, _ := app.Test(statusReq)
	var status map[string]interface{}
	json.NewDecoder(statusResp.Body).Decode(&status)
	statusResp.Body.Close()

	if status["running"].(float64) != 0 {
		t.Error("Expected all services to be stopped initially")
	}

	// Start a service using JSON payload
	startPayload := map[string]string{"service_id": "api-gateway"}
	startBody, _ := json.Marshal(startPayload)
	startReq := httptest.NewRequest("POST", "/api/services/start", bytes.NewReader(startBody))
	startReq.Header.Set("Content-Type", "application/json")
	startResp, err := app.Test(startReq)
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}
	defer startResp.Body.Close()

	if startResp.StatusCode != 200 {
		// Read response for debugging
		body := make([]byte, 1024)
		n, _ := startResp.Body.Read(body)
		t.Errorf("Expected status 200 when starting service, got %d. Response: %s", startResp.StatusCode, string(body[:n]))
	}

	// Verify service is running
	statusReq2 := httptest.NewRequest("GET", "/api/services/status", nil)
	statusResp2, _ := app.Test(statusReq2)
	var status2 map[string]interface{}
	json.NewDecoder(statusResp2.Body).Decode(&status2)
	statusResp2.Body.Close()

	if status2["running"].(float64) != 1 {
		t.Errorf("Expected 1 running service, got %v", status2["running"])
	}

	// Stop the service
	stopPayload := map[string]string{"service_id": "api-gateway"}
	stopBody, _ := json.Marshal(stopPayload)
	stopReq := httptest.NewRequest("POST", "/api/services/stop", bytes.NewReader(stopBody))
	stopReq.Header.Set("Content-Type", "application/json")
	stopResp, err := app.Test(stopReq)
	if err != nil {
		t.Fatalf("Failed to stop service: %v", err)
	}
	defer stopResp.Body.Close()

	// Verify service is stopped
	statusReq3 := httptest.NewRequest("GET", "/api/services/status", nil)
	statusResp3, _ := app.Test(statusReq3)
	var status3 map[string]interface{}
	json.NewDecoder(statusResp3.Body).Decode(&status3)
	statusResp3.Body.Close()

	if status3["running"].(float64) != 0 {
		t.Errorf("Expected 0 running services, got %v", status3["running"])
	}
}

func TestIntegrationMultipleRequests(t *testing.T) {
	app := setupTestApp()

	// Simulate multiple concurrent requests
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/health", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Errorf("Request failed: %v", err)
			} else {
				resp.Body.Close()
			}
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	t.Log("All concurrent requests completed successfully")
}

func TestIntegrationAllEndpointsRespond(t *testing.T) {
	app := setupTestApp()

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/health"},
		{"GET", "/ready"},
		{"GET", "/live"},
		{"GET", "/metrics"},
		{"GET", "/metrics/prometheus"},
		{"GET", "/api/sysinfo"},
		{"GET", "/api/services"},
		{"GET", "/api/services/status"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			req := httptest.NewRequest(endpoint.method, endpoint.path, nil)
			resp, err := app.Test(req, 5000) // 5 second timeout
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 500 {
				t.Errorf("Endpoint returned server error: %d", resp.StatusCode)
			}

			if resp.StatusCode == 404 {
				t.Errorf("Endpoint not found: %s %s", endpoint.method, endpoint.path)
			}
		})
	}
}

func TestIntegrationMetricsStartTime(t *testing.T) {
	// Verify that StartTime was initialized on import
	if metrics.StartTime.IsZero() {
		t.Error("StartTime should be initialized")
	}

	if metrics.StartTime.After(time.Now()) {
		t.Error("StartTime should be in the past")
	}

	t.Logf("Application start time: %s", metrics.StartTime.Format(time.RFC3339))
}

func BenchmarkIntegrationHealthCheck(b *testing.B) {
	app := setupTestApp()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}

func BenchmarkIntegrationMetrics(b *testing.B) {
	app := setupTestApp()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/metrics", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}

func BenchmarkIntegrationServicesList(b *testing.B) {
	app := setupTestApp()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/services", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}
