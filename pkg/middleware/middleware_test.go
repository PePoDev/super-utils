package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/logger"
	"github.com/pepodev/super-utils/pkg/metrics"
)

func init() {
	// Initialize test logger (no-op) for all tests
	logger.InitTestLogger()
}

func TestRequestCounter(t *testing.T) {
	app := fiber.New()

	// Reset the counter before test
	initialCount := metrics.TotalRequests

	// Apply middleware
	app.Use(RequestCounter())

	// Add a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Make a request
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check that counter was incremented
	if metrics.TotalRequests != initialCount+1 {
		t.Errorf("Expected TotalRequests to be %d, got %d", initialCount+1, metrics.TotalRequests)
	}

	// Make another request
	req2 := httptest.NewRequest("GET", "/test", nil)
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp2.Body.Close()

	// Check that counter was incremented again
	if metrics.TotalRequests != initialCount+2 {
		t.Errorf("Expected TotalRequests to be %d, got %d", initialCount+2, metrics.TotalRequests)
	}
}

func TestRequestCounterMultipleRoutes(t *testing.T) {
	app := fiber.New()

	// Reset the counter before test
	initialCount := metrics.TotalRequests

	// Apply middleware
	app.Use(RequestCounter())

	// Add multiple test routes
	app.Get("/route1", func(c *fiber.Ctx) error {
		return c.SendString("Route 1")
	})
	app.Get("/route2", func(c *fiber.Ctx) error {
		return c.SendString("Route 2")
	})
	app.Post("/route3", func(c *fiber.Ctx) error {
		return c.SendString("Route 3")
	})

	// Make requests to different routes
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/route1"},
		{"GET", "/route2"},
		{"POST", "/route3"},
		{"GET", "/route1"},
	}

	for i, route := range routes {
		req := httptest.NewRequest(route.method, route.path, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request %d: %v", i, err)
		}
		resp.Body.Close()

		expectedCount := initialCount + i + 1
		if metrics.TotalRequests != expectedCount {
			t.Errorf("After request %d: Expected TotalRequests to be %d, got %d", i+1, expectedCount, metrics.TotalRequests)
		}
	}
}

func TestRequestCounterWith404(t *testing.T) {
	app := fiber.New()

	// Reset the counter before test
	initialCount := metrics.TotalRequests

	// Apply middleware
	app.Use(RequestCounter())

	// Add a test route
	app.Get("/exists", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Make a request to non-existent route (should still count)
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check that counter was incremented even for 404
	if metrics.TotalRequests != initialCount+1 {
		t.Errorf("Expected TotalRequests to be %d even for 404, got %d", initialCount+1, metrics.TotalRequests)
	}
}

func TestRequestCounterPassesRequestThrough(t *testing.T) {
	app := fiber.New()

	// Apply middleware
	app.Use(RequestCounter())

	// Add a test route that returns specific content
	testContent := "Test Response Content"
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString(testContent)
	})

	// Make a request
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read response body
	body := make([]byte, len(testContent))
	resp.Body.Read(body)
	if string(body) != testContent {
		t.Errorf("Expected body %q, got %q", testContent, string(body))
	}
}

func BenchmarkRequestCounter(b *testing.B) {
	app := fiber.New()
	app.Use(RequestCounter())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}
