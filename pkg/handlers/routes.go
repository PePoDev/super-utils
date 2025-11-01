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
	"github.com/pepodev/super-utils/pkg/services"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(app *fiber.App) {
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
}

// HomeHandler serves the main HTML page
func HomeHandler(c *fiber.Ctx) error {
	// Read the HTML file at runtime
	htmlPath := filepath.Join("web", "index.html")
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		htmlPath = filepath.Join("..", "web", "index.html")
	}

	htmlContent, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		return c.Status(500).SendString("Error loading page: " + err.Error())
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Send(htmlContent)
}

// SystemInfoHandler returns system information
func SystemInfoHandler(c *fiber.Ctx) error {
	info := services.GetSystemInfo()
	return c.JSON(info)
}

// ShellHandler executes shell commands
func ShellHandler(c *fiber.Ctx) error {
	// Support both new command format and legacy tool+args format
	command := c.FormValue("command")
	tool := c.FormValue("tool")
	args := c.FormValue("args")

	var output string
	var err error

	if command != "" {
		// New format: full command line
		output, err = services.ExecuteCommand(command)
	} else if tool != "" {
		// Legacy format: tool + args
		output, err = services.ExecuteToolCommand(tool, args)
	} else {
		return c.Status(400).SendString("Command or tool is required")
	}

	if err != nil {
		return c.SendString(output)
	}

	return c.SendString(output)
}

// CurlHandler handles curl and GraphQL requests
func CurlHandler(c *fiber.Ctx) error {
	url := c.FormValue("url")
	graphql := c.FormValue("graphql")

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
		return c.Status(500).SendString("Request error: " + err.Error())
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	return c.SendString(string(b))
}

// NetworkHandler provides network information lookups
func NetworkHandler(c *fiber.Ctx) error {
	host := c.FormValue("host")

	ips, err := net.LookupIP(host)
	if err != nil {
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

	res, _ := json.MarshalIndent(result, "", "  ")
	return c.SendString(string(res))
}
