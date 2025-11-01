package metrics

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestHealthHandler(t *testing.T) {
	app := fiber.New()
	RegisterHealthMetrics(app)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var health map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check required fields
	if status, ok := health["status"]; !ok || status != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", health["status"])
	}

	if _, ok := health["timestamp"]; !ok {
		t.Error("Expected timestamp field")
	}

	if _, ok := health["uptime"]; !ok {
		t.Error("Expected uptime field")
	}

	if version, ok := health["version"]; !ok || version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %v", health["version"])
	}
}

func TestReadinessHandler(t *testing.T) {
	app := fiber.New()
	RegisterHealthMetrics(app)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var readiness map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&readiness); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if ready, ok := readiness["ready"]; !ok || ready != true {
		t.Errorf("Expected ready to be true, got %v", readiness["ready"])
	}

	if _, ok := readiness["timestamp"]; !ok {
		t.Error("Expected timestamp field")
	}
}

func TestLivenessHandler(t *testing.T) {
	app := fiber.New()
	RegisterHealthMetrics(app)

	req := httptest.NewRequest("GET", "/live", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var liveness map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&liveness); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if alive, ok := liveness["alive"]; !ok || alive != true {
		t.Errorf("Expected alive to be true, got %v", liveness["alive"])
	}

	if _, ok := liveness["timestamp"]; !ok {
		t.Error("Expected timestamp field")
	}
}

func TestMetricsHandler(t *testing.T) {
	app := fiber.New()
	RegisterHealthMetrics(app)

	req := httptest.NewRequest("GET", "/metrics", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var metrics map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check required fields
	requiredFields := []string{"uptime_seconds", "uptime", "goroutines", "total_requests", "memory", "runtime"}
	for _, field := range requiredFields {
		if _, ok := metrics[field]; !ok {
			t.Errorf("Expected field %s in metrics", field)
		}
	}

	// Check memory metrics
	if memory, ok := metrics["memory"].(map[string]interface{}); ok {
		memoryFields := []string{"alloc_bytes", "alloc_mb", "sys_bytes", "sys_mb", "num_gc"}
		for _, field := range memoryFields {
			if _, ok := memory[field]; !ok {
				t.Errorf("Expected field %s in memory metrics", field)
			}
		}
	} else {
		t.Error("Expected memory to be a map")
	}

	// Check runtime metrics
	if runtime, ok := metrics["runtime"].(map[string]interface{}); ok {
		runtimeFields := []string{"go_version", "go_os", "go_arch", "num_cpu"}
		for _, field := range runtimeFields {
			if _, ok := runtime[field]; !ok {
				t.Errorf("Expected field %s in runtime metrics", field)
			}
		}
	} else {
		t.Error("Expected runtime to be a map")
	}
}

func TestPrometheusMetricsHandler(t *testing.T) {
	app := fiber.New()
	RegisterHealthMetrics(app)

	req := httptest.NewRequest("GET", "/metrics/prometheus", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Errorf("Expected Content-Type to contain text/plain, got %s", contentType)
	}

	// Read body
	body := make([]byte, 2048)
	n, _ := resp.Body.Read(body)
	bodyStr := string(body[:n])

	// Check for Prometheus format metrics
	expectedMetrics := []string{
		"super_utils_uptime_seconds",
		"super_utils_goroutines",
		"super_utils_memory_alloc_bytes",
		"super_utils_memory_sys_bytes",
		"super_utils_gc_total",
		"super_utils_requests_total",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(bodyStr, metric) {
			t.Errorf("Expected metric %s in Prometheus output", metric)
		}
	}

	// Check for HELP and TYPE comments
	if !strings.Contains(bodyStr, "# HELP") {
		t.Error("Expected HELP comments in Prometheus output")
	}

	if !strings.Contains(bodyStr, "# TYPE") {
		t.Error("Expected TYPE comments in Prometheus output")
	}
}

func TestStartTime(t *testing.T) {
	// Verify StartTime was initialized
	if StartTime.IsZero() {
		t.Error("StartTime should be initialized")
	}

	// Verify StartTime is in the past
	if StartTime.After(time.Now()) {
		t.Error("StartTime should be in the past")
	}
}

func TestMetricsHandlerUptime(t *testing.T) {
	// Reset StartTime to a known value
	originalStartTime := StartTime
	StartTime = time.Now().Add(-1 * time.Hour)
	defer func() { StartTime = originalStartTime }()

	app := fiber.New()
	RegisterHealthMetrics(app)

	req := httptest.NewRequest("GET", "/metrics", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var metrics map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&metrics)

	// Check that uptime_seconds is reasonable
	uptimeSeconds, ok := metrics["uptime_seconds"].(float64)
	if !ok {
		t.Error("Expected uptime_seconds to be a number")
	}

	// Should be approximately 1 hour (3600 seconds)
	if uptimeSeconds < 3500 || uptimeSeconds > 3700 {
		t.Errorf("Expected uptime_seconds to be around 3600, got %f", uptimeSeconds)
	}
}

func TestRegisterHealthMetrics(t *testing.T) {
	app := fiber.New()
	RegisterHealthMetrics(app)

	// Test all endpoints are registered
	endpoints := []string{
		"/health",
		"/ready",
		"/live",
		"/metrics",
		"/metrics/prometheus",
	}

	for _, endpoint := range endpoints {
		req := httptest.NewRequest("GET", endpoint, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Errorf("Failed to make request to %s: %v", endpoint, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == 404 {
			t.Errorf("Endpoint %s not registered (got 404)", endpoint)
		}
	}
}

func BenchmarkHealthHandler(b *testing.B) {
	app := fiber.New()
	RegisterHealthMetrics(app)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}

func BenchmarkMetricsHandler(b *testing.B) {
	app := fiber.New()
	RegisterHealthMetrics(app)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/metrics", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}

func BenchmarkPrometheusMetricsHandler(b *testing.B) {
	app := fiber.New()
	RegisterHealthMetrics(app)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/metrics/prometheus", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}
