package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/logger"
)

func init() {
	// Initialize test logger (no-op) for all tests
	logger.InitTestLogger()
}

func TestHomeHandler(t *testing.T) {
	app := fiber.New()
	app.Get("/", HomeHandler)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Status could be 200 or 500 depending on whether file exists
	if resp.StatusCode != 200 && resp.StatusCode != 500 {
		t.Errorf("Expected status 200 or 500, got %d", resp.StatusCode)
	}

	// If successful, should have HTML content type
	if resp.StatusCode == 200 {
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("Expected Content-Type to contain text/html, got %s", contentType)
		}
	}
}

func TestSystemInfoHandler(t *testing.T) {
	app := fiber.New()
	app.Get("/api/sysinfo", SystemInfoHandler)

	req := httptest.NewRequest("GET", "/api/sysinfo", nil)
	resp, err := app.Test(req, 5000) // 5 second timeout
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var info map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(info) == 0 {
		t.Error("Expected system info to not be empty")
	}

	t.Logf("System info contains %d fields", len(info))
}

func TestShellHandler(t *testing.T) {
	app := fiber.New()
	app.Post("/api/shell", ShellHandler)

	tests := []struct {
		name           string
		method         string
		body           string
		contentType    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Execute simple command",
			method:         "POST",
			body:           "command=echo+hello",
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name:           "Execute with tool and args (legacy)",
			method:         "POST",
			body:           "tool=echo&args=hello",
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name:           "Missing command and tool",
			method:         "POST",
			body:           "",
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name:           "Invalid command",
			method:         "POST",
			body:           "command=nonexistentcommand123",
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: 200,
			expectError:    false, // Returns 200 but with error message
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/shell", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestCurlHandler(t *testing.T) {
	app := fiber.New()
	app.Post("/api/curl", CurlHandler)

	tests := []struct {
		name           string
		formData       string
		expectedStatus int
		skipTest       bool
	}{
		{
			name:           "GET request",
			formData:       "url=http://example.com",
			expectedStatus: 200,
			skipTest:       true, // Skip in CI/testing environment
		},
		{
			name:           "Missing URL",
			formData:       "",
			expectedStatus: 500,
			skipTest:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipTest {
				t.Skip("Skipping test that requires external network")
			}

			req := httptest.NewRequest("POST", "/api/curl?"+tt.formData, nil)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestNetworkHandler(t *testing.T) {
	app := fiber.New()
	app.Post("/api/network", NetworkHandler)

	tests := []struct {
		name           string
		formData       string
		expectedStatus int
		expectIPs      bool
	}{
		{
			name:           "Lookup localhost",
			formData:       "host=localhost",
			expectedStatus: 200,
			expectIPs:      true,
		},
		{
			name:           "Lookup 127.0.0.1",
			formData:       "host=127.0.0.1",
			expectedStatus: 200,
			expectIPs:      true,
		},
		{
			name:           "Empty host",
			formData:       "",
			expectedStatus: 500,
			expectIPs:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/network?"+tt.formData, nil)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectIPs && resp.StatusCode == 200 {
				body := make([]byte, 1024)
				n, _ := resp.Body.Read(body)
				bodyStr := string(body[:n])

				if !strings.Contains(bodyStr, "ips") {
					t.Error("Expected 'ips' field in response")
				}
			}
		})
	}
}

func TestSpeedTestUploadHandler(t *testing.T) {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB limit for speed test
	})
	app.Post("/api/speedtest/upload", SpeedTestUploadHandler)

	tests := []struct {
		name           string
		dataSize       int
		expectedStatus int
	}{
		{
			name:           "Upload 1KB data",
			dataSize:       1024,
			expectedStatus: 200,
		},
		{
			name:           "Upload 5MB data",
			dataSize:       5 * 1024 * 1024,
			expectedStatus: 200,
		},
		{
			name:           "Upload empty data",
			dataSize:       0,
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, tt.dataSize)
			req := httptest.NewRequest("POST", "/api/speedtest/upload", strings.NewReader(string(data)))
			req.Header.Set("Content-Type", "application/octet-stream")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Verify response
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if result["status"] != "ok" {
				t.Error("Expected status 'ok' in response")
			}

			if received, ok := result["received"].(float64); ok {
				if int(received) != tt.dataSize {
					t.Errorf("Expected received=%d, got %d", tt.dataSize, int(received))
				}
			}
		})
	}
}

func TestRegisterRoutes(t *testing.T) {
	app := fiber.New()
	RegisterRoutes(app)

	// Test that all routes are registered
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"GET", "/shop"},
		{"GET", "/api/sysinfo"},
		{"POST", "/api/shell"},
		{"POST", "/api/curl"},
		{"POST", "/api/network"},
		{"GET", "/api/services"},
		{"GET", "/api/services/status"},
		{"POST", "/api/services/start"},
		{"POST", "/api/services/stop"},
		{"POST", "/api/services/reset"},
	}

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			var req *http.Request
			if route.method == "GET" {
				req = httptest.NewRequest(route.method, route.path, nil)
			} else {
				req = httptest.NewRequest(route.method, route.path, strings.NewReader(""))
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Errorf("Failed to make request: %v", err)
				return
			}
			defer resp.Body.Close()

			// 404 means the route is not registered
			if resp.StatusCode == 404 {
				t.Errorf("Route %s %s not registered (got 404)", route.method, route.path)
			}
		})
	}
}

func TestHomeHandlerWithTestFile(t *testing.T) {
	// Create a temporary HTML file
	tmpDir := t.TempDir()
	htmlPath := tmpDir + "/index.html"
	htmlContent := "<html><body>Test Home Page</body></html>"
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Note: This test demonstrates the handler logic but won't actually use the temp file
	// because the handler looks for web/index.html
	app := fiber.New()
	app.Get("/", HomeHandler)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Just verify handler responds
	if resp.StatusCode != 200 && resp.StatusCode != 500 {
		t.Errorf("Expected status 200 or 500, got %d", resp.StatusCode)
	}
}

func TestShellHandlerCommandFormat(t *testing.T) {
	app := fiber.New()
	app.Post("/api/shell", ShellHandler)

	// Test new command format with form body
	body := "command=echo+test"
	req := httptest.NewRequest("POST", "/api/shell", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBody := make([]byte, 1024)
	n, _ := resp.Body.Read(respBody)
	output := string(respBody[:n])

	if !strings.Contains(output, "test") {
		t.Errorf("Expected output to contain 'test', got %q", output)
	}
}

func BenchmarkSystemInfoHandler(b *testing.B) {
	app := fiber.New()
	app.Get("/api/sysinfo", SystemInfoHandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/sysinfo", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}

func BenchmarkShellHandler(b *testing.B) {
	app := fiber.New()
	app.Post("/api/shell", ShellHandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/shell?command=echo test", nil)
		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}
