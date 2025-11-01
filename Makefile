.PHONY: help test test-verbose test-coverage test-unit test-integration clean build run docker-build docker-run lint fmt

# Default target
help:
	@echo "Available targets:"
	@echo "  make test              - Run all tests"
	@echo "  make test-verbose      - Run tests with verbose output"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-unit         - Run only unit tests"
	@echo "  make test-integration  - Run only integration tests"
	@echo "  make test-short        - Run tests in short mode (skip slow tests)"
	@echo "  make bench             - Run benchmarks"
	@echo "  make lint              - Run linter"
	@echo "  make fmt               - Format code"
	@echo "  make build             - Build the application"
	@echo "  make run               - Run the application"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make docker-build      - Build Docker image"
	@echo "  make docker-run        - Run Docker container"

# Run all tests
test:
	@echo "Running all tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	go tool cover -func=coverage.out

# Run only unit tests (exclude integration tests)
test-unit:
	@echo "Running unit tests..."
	go test -v ./pkg/...

# Run only integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v ./cmd/...

# Run tests in short mode
test-short:
	@echo "Running tests in short mode..."
	go test -short ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Run benchmarks with CPU profiling
bench-cpu:
	@echo "Running benchmarks with CPU profiling..."
	go test -bench=. -benchmem -cpuprofile=cpu.prof ./...
	@echo "CPU profile saved to cpu.prof"
	@echo "View with: go tool pprof cpu.prof"

# Run benchmarks with memory profiling
bench-mem:
	@echo "Running benchmarks with memory profiling..."
	go test -bench=. -benchmem -memprofile=mem.prof ./...
	@echo "Memory profile saved to mem.prof"
	@echo "View with: go tool pprof mem.prof"

# Run specific test by name
# Usage: make test-run TEST=TestName
test-run:
	@echo "Running test: $(TEST)"
	go test -v -run $(TEST) ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/"; \
		echo "Or run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "All checks passed!"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/super-utils ./cmd/

# Build with version info
build-version:
	@echo "Building application with version info..."
	@VERSION=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	BUILD_TIME=$$(date -u '+%Y-%m-%d_%H:%M:%S'); \
	go build -ldflags "-X main.Version=$$VERSION -X main.BuildTime=$$BUILD_TIME" -o bin/super-utils ./cmd/

# Run the application
run:
	@echo "Running application..."
	go run ./cmd/

# Run with race detector
run-race:
	@echo "Running application with race detector..."
	go run -race ./cmd/

# Clean build artifacts and test cache
clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/super-utils
	rm -f coverage.out coverage.html
	rm -f cpu.prof mem.prof
	go clean -testcache
	go clean -cache

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t super-utils:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8000:8000 super-utils:latest

# Show test coverage by package
coverage-by-package:
	@echo "Test coverage by package:"
	go test -coverprofile=coverage.out ./... > /dev/null 2>&1
	go tool cover -func=coverage.out | grep -E '^github.com' | column -t

# Generate test report
test-report:
	@echo "Generating test report..."
	go test -v -json ./... > test-report.json
	@echo "Test report saved to test-report.json"

# Watch and run tests on file changes (requires entr or similar)
test-watch:
	@echo "Watching for changes and running tests..."
	@if command -v find >/dev/null 2>&1 && command -v entr >/dev/null 2>&1; then \
		find . -name '*.go' | entr -c make test; \
	else \
		echo "This target requires 'entr'. Install it with your package manager."; \
		echo "  Ubuntu/Debian: apt-get install entr"; \
		echo "  macOS: brew install entr"; \
	fi

# Initialize project (install tools)
init:
	@echo "Initializing project..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod download
	@echo "Project initialized!"

# Show project statistics
stats:
	@echo "Project Statistics:"
	@echo "===================="
	@echo "Total Go files: $$(find . -name '*.go' | wc -l)"
	@echo "Total lines of code: $$(find . -name '*.go' -exec cat {} \; | wc -l)"
	@echo "Total test files: $$(find . -name '*_test.go' | wc -l)"
	@echo "Test coverage: $$(go test -coverprofile=coverage.out ./... > /dev/null 2>&1 && go tool cover -func=coverage.out | grep total | awk '{print $$3}')"
