package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"

	_ "embed"

	"github.com/gofiber/fiber/v2"
)

// Only keep one main function below
// Only keep one parseArgs function below

// Curl/GraphQL endpoint and network info registration
func registerRoutes(app *fiber.App) {
	app.Post("/api/curl", func(c *fiber.Ctx) error {
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
	})

	// Network info endpoint
	app.Post("/api/network", func(c *fiber.Ctx) error {
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
	})
}

// Helper to split args respecting quotes
// e.g. parseArgs('foo "bar baz"') => [foo, bar baz]
func parseArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	for _, r := range input {
		switch r {
		case ' ':
			if inQuotes {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		case '"':
			inQuotes = !inQuotes
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

//go:embed web/index.html
var indexHTML string

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(indexHTML)
	})

	// Run selected tool with arguments
	app.Post("/api/shell", func(c *fiber.Ctx) error {
		tool := c.FormValue("tool")
		args := c.FormValue("args")
		if tool == "" {
			return c.Status(400).SendString("Tool is required")
		}
		argList := []string{}
		if args != "" {
			argList = parseArgs(args)
		}
		cmd := exec.Command(tool, argList...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return c.SendString(string(out) + "\nError: " + err.Error())
		}
		return c.SendString(string(out))
	})
	registerRoutes(app)
	app.Listen(":8080")
}
