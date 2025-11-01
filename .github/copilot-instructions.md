# Super-Utils Development Guide

## Architecture Overview

Super-utils is a **dual-purpose Go/Fiber application** serving both as a web-based network toolkit and a microservices emulator:

1. **Network Toolkit**: Web UI for executing system/network commands (ping, curl, tcpdump, etc.)
2. **Service Emulator**: Simulates 10 microservices in a book shop architecture for testing

### Core Structure

```
cmd/main.go                    # Entry point: Fiber app setup on port 8000
pkg/
  handlers/                    # HTTP endpoint handlers
    routes.go                  # Main API routes + shell/curl/network handlers
    service-emulator.go        # Microservices emulation API
  services/                    # Business logic
    emulator.go               # In-memory service state management (singleton pattern)
    shell.go                  # Command execution via os/exec
    system.go                 # System info collection
  middleware/                  # Request counter middleware
  metrics/                     # Health/readiness/metrics endpoints
  utils/                       # Argument parser for shell commands
web/
  index.html                   # Main HTMX UI (embedded at runtime)
  book-shop.html              # Service emulator UI
```

## Key Patterns & Conventions

### 1. Service Emulator Singleton

The service emulator (`pkg/services/emulator.go`) uses a **global singleton** accessed via `GetEmulator()`. Always reset state in tests:

```go
em := services.GetEmulator()
em.ResetAllServices() // Required before each test
```

### 2. Handler Input Patterns

Handlers support **both JSON and form data** for backward compatibility:

```go
// Try JSON first, fallback to form value
if err := c.BodyParser(&req); err != nil {
    req.ServiceID = c.FormValue("service_id")
}
```

### 3. HTML File Loading

HTML files are read at **runtime** (not embedded) with fallback paths:

```go
htmlPath := filepath.Join("web", "index.html")
if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
    htmlPath = filepath.Join("..", "web", "index.html") // For tests
}
```

### 4. Command Execution

Two formats supported in `ShellHandler`:

- **New format**: `command=full command line` (preferred)
- **Legacy format**: `tool=cmd&args=arguments` (parsed via `utils.ParseArgs`)

## Development Workflows

### Running Locally

```bash
go run cmd/main.go              # Runs on http://localhost:8000
# OR
make run                        # Same as above
```

### Testing (93.6% coverage)

```bash
make test                       # Run all tests
make test-coverage              # Generate HTML coverage report
make test-unit                  # Only pkg/* tests
make test-integration           # Only cmd/* tests
make bench                      # Run benchmarks
```

**Critical**: Use `app.Test(req, timeout)` in tests to avoid default 1s timeout:

```go
resp, err := app.Test(req, 5000) // 5 second timeout for slow operations
```

### Building

```bash
make build                      # Builds to bin/super-utils
make docker-build               # Multi-arch Docker image
```

## API Endpoints Reference

### Core APIs

- `GET /` - Main HTMX UI
- `GET /shop` - Service emulator UI
- `POST /api/shell` - Execute shell commands
- `POST /api/curl` - HTTP/GraphQL requests
- `POST /api/network` - Network lookups

### Service Emulator APIs

- `GET /api/services` - List all 10 services
- `GET /api/services/status` - Service counts (running/stopped)
- `POST /api/services/start` - Start service (JSON: `{"service_id": "api-gateway"}`)
- `POST /api/services/stop` - Stop service
- `POST /api/services/reset` - Stop all services

### Health & Metrics

- `GET /health` - Health check with uptime
- `GET /ready` - Readiness probe
- `GET /live` - Liveness probe
- `GET /metrics` - JSON metrics
- `GET /metrics/prometheus` - Prometheus format

## Testing Conventions

### Table-Driven Tests

All tests use table-driven patterns with descriptive names:

```go
tests := []struct {
    name     string
    input    string
    expected []string
}{
    {"Empty string", "", []string{}},
    {"Single argument", "foo", []string{"foo"}},
}
```

### Fiber Test Pattern

```go
app := fiber.New()
app.Post("/endpoint", Handler)
req := httptest.New Request("POST", "/endpoint", body)
resp, err := app.Test(req, 5000) // Always specify timeout
```

### Service Emulator Test Setup

```go
em := services.InitEmulator()  // Creates fresh instance
em.ResetAllServices()           // Ensure clean state
// ... run tests
```

## Project-Specific Notes

1. **Port 8000** (not 8080): Main server runs on port 8000
2. **10 Microservices**: api-gateway, service-registry, user-service, catalog-service, cart-service, order-service, payment-service, review-service, notification-service, inventory-service
3. **No real service processes**: Emulator only tracks state (running/stopped) in memory
4. **HTMX-based UI**: Frontend uses HTMX for interactivity, not React/Vue
5. **Security caveat**: Shell execution is intentional but restricted to trusted environments

## Common Tasks

### Add New Shell Tool

1. Add to `web/index.html` dropdown
2. No backend changes needed (uses generic `ExecuteCommand`)

### Add New Service

1. Update `InitEmulator()` in `pkg/services/emulator.go`
2. Add to service array with ID, port, database

### Add New Endpoint

1. Define handler in `pkg/handlers/`
2. Register in `RegisterRoutes()`
3. Add test in corresponding `*_test.go`

### Debug Test Failures

```bash
make test-verbose               # See detailed output
go test -v -run TestName ./...  # Run specific test
go test -race ./...             # Check race conditions
```

## External Dependencies

- **Fiber v2**: HTTP framework (not Gin/Echo)
- **Standard library**: Uses `os/exec`, `net/http/httptest` (no testify/mockery)
- **Docker Alpine**: Base image for multi-tool container

See `AGENTS.md` for full tool list and `TESTING.md` for comprehensive test documentation.
