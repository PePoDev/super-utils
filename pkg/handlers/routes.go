package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pepodev/super-utils/pkg/logger"
	"github.com/pepodev/super-utils/pkg/services"
	"go.uber.org/zap"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(app *fiber.App) {
	logger.Info("Registering application routes")

	// Home page
	app.Get("/", HomeHandler)

	// Book shop frontend (customer interface)
	app.Get("/shop", BookShopEmulatorHandler)

	// API routes
	app.Get("/api/sysinfo", SystemInfoHandler)
	app.Post("/api/shell", ShellHandler)
	app.Post("/api/curl", CurlHandler)
	app.Post("/api/network", NetworkHandler)

	// Service emulator API routes
	app.Get("/api/services", GetServicesHandler)
	app.Get("/api/services/status", GetServiceStatusHandler)
	app.Post("/api/services/start", StartServiceHandler)
	app.Post("/api/services/stop", StopServiceHandler)
	app.Post("/api/services/reset", ResetServicesHandler)

	// Speed test endpoint
	app.Post("/api/speedtest/upload", SpeedTestUploadHandler)

	logger.Info("Successfully registered all routes")
}

// HomeHandler serves the main HTML page
func HomeHandler(c *fiber.Ctx) error {
	logger.Debug("Serving home page")

	// Read the HTML file at runtime
	htmlPath := filepath.Join("web", "index.html")
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		htmlPath = filepath.Join("..", "web", "index.html")
	}

	htmlContent, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		logger.Error("Failed to load HTML page",
			zap.String("path", htmlPath),
			zap.Error(err),
		)
		return c.Status(500).SendString("Error loading page: " + err.Error())
	}

	logger.Info("Home page served successfully")
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Send(htmlContent)
}

// SystemInfoHandler returns system information
func SystemInfoHandler(c *fiber.Ctx) error {
	logger.Info("Fetching system information")

	info := services.GetSystemInfo()

	logger.Info("System information retrieved successfully",
		zap.Int("info_fields", len(info)),
	)

	return c.JSON(info)
}

// ShellHandler executes shell commands
func ShellHandler(c *fiber.Ctx) error {
	// Support both new command format and legacy tool+args format
	command := c.FormValue("command")
	tool := c.FormValue("tool")
	args := c.FormValue("args")

	logger.Info("Shell command request",
		zap.String("command", command),
		zap.String("tool", tool),
		zap.Bool("has_args", args != ""),
		zap.String("client_ip", c.IP()),
	)

	var output string
	var err error

	if command != "" {
		// New format: full command line
		output, err = services.ExecuteCommand(command)
	} else if tool != "" {
		// Legacy format: tool + args
		output, err = services.ExecuteToolCommand(tool, args)
	} else {
		logger.Warn("Shell command request missing parameters",
			zap.String("client_ip", c.IP()),
		)
		return c.Status(400).SendString("Command or tool is required")
	}

	if err != nil {
		logger.Error("Shell command execution failed",
			zap.String("command", command),
			zap.String("tool", tool),
			zap.Error(err),
		)
		return c.SendString(output)
	}

	logger.Info("Shell command executed successfully",
		zap.String("command", command),
		zap.String("tool", tool),
		zap.Int("output_length", len(output)),
	)

	return c.SendString(output)
}

// CurlHandler handles curl and GraphQL requests
func CurlHandler(c *fiber.Ctx) error {
	url := c.FormValue("url")
	graphql := c.FormValue("graphql")

	logger.Info("HTTP request",
		zap.String("url", url),
		zap.Bool("is_graphql", graphql != ""),
		zap.String("client_ip", c.IP()),
	)

	var resp *http.Response
	var err error

	if graphql != "" {
		// GraphQL POST
		payload := map[string]string{"query": graphql}
		body, _ := json.Marshal(payload)
		resp, err = http.Post(url, "application/json", strings.NewReader(string(body)))
	} else {
		resp, err = http.Get(url)
	}

	if err != nil {
		logger.Error("HTTP request failed",
			zap.String("url", url),
			zap.Bool("is_graphql", graphql != ""),
			zap.Error(err),
		)
		return c.Status(500).SendString("Request error: " + err.Error())
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	logger.Info("HTTP request completed",
		zap.String("url", url),
		zap.Int("status_code", resp.StatusCode),
		zap.Int("response_size", len(b)),
	)

	return c.SendString(string(b))
}

// NetworkHandler provides network information lookups
func NetworkHandler(c *fiber.Ctx) error {
	host := c.FormValue("host")

	logger.Info("Network lookup request",
		zap.String("host", host),
		zap.String("client_ip", c.IP()),
	)

	ips, err := net.LookupIP(host)
	if err != nil {
		logger.Error("Network lookup failed",
			zap.String("host", host),
			zap.Error(err),
		)
		return c.Status(500).SendString("Lookup error: " + err.Error())
	}

	var ipStrs []string
	for _, ip := range ips {
		ipStrs = append(ipStrs, ip.String())
	}

	addrs, err := net.LookupAddr(host)
	if err != nil {
		addrs = []string{"No PTR record found"}
	}

	result := map[string]interface{}{
		"ips":   ipStrs,
		"addrs": addrs,
	}

	logger.Info("Network lookup completed",
		zap.String("host", host),
		zap.Int("ip_count", len(ipStrs)),
		zap.Int("addr_count", len(addrs)),
	)

	res, _ := json.MarshalIndent(result, "", "  ")
	return c.SendString(string(res))
}

// SpeedTestUploadHandler handles upload speed test requests
func SpeedTestUploadHandler(c *fiber.Ctx) error {
	// Receive and discard the upload data
	body := c.Body()

	logger.Info("Speed test upload",
		zap.Int("bytes_received", len(body)),
		zap.String("client_ip", c.IP()),
	)

	return c.JSON(fiber.Map{
		"received": len(body),
		"status":   "ok",
	})
}
