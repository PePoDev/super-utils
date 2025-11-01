# Testing Documentation

This document describes the comprehensive test suite for the super-utils project.

## Test Coverage

The project includes **93.6% test coverage** across all packages with the following test types:

- **Unit Tests**: Tests for individual functions and methods
- **Integration Tests**: End-to-end tests for the complete application
- **Benchmarks**: Performance tests for critical paths

## Test Structure

```
.
├── cmd/
│   └── main_test.go              # Integration tests
├── pkg/
│   ├── handlers/
│   │   ├── routes_test.go        # HTTP handler tests
│   │   └── service-emulator_test.go  # Service emulator handler tests
│   ├── metrics/
│   │   └── metrics_test.go       # Metrics and health check tests
│   ├── middleware/
│   │   └── middleware_test.go    # Middleware tests
│   ├── services/
│   │   ├── emulator_test.go      # Service emulator logic tests
│   │   ├── shell_test.go         # Shell command execution tests
│   │   └── system_test.go        # System information tests
│   └── utils/
│       └── parser_test.go        # Argument parser tests
```

## Running Tests

### Using Make (Recommended)

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests in short mode (skip slow tests)
make test-short

# Run benchmarks
make bench

# Run specific test by name
make test-run TEST=TestName
```

### Using Go CLI

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./pkg/services/...

# Run specific test
go test -run TestName ./...

# Run benchmarks
go test -bench=. ./...
```

## Test Categories

### Unit Tests

**Utils Package** (`pkg/utils/parser_test.go`)
- Tests the argument parser with various input formats
- Edge cases: empty strings, quoted arguments, special characters
- Coverage: 100%

**Services Package**
- **Emulator** (`pkg/services/emulator_test.go`): Tests service lifecycle management
- **Shell** (`pkg/services/shell_test.go`): Tests command execution
- **System** (`pkg/services/system_test.go`): Tests system information gathering
- Coverage: 100%

**Middleware Package** (`pkg/middleware/middleware_test.go`)
- Tests request counter middleware
- Coverage: 100%

**Metrics Package** (`pkg/metrics/metrics_test.go`)
- Tests health check endpoints
- Tests metrics collection and reporting
- Tests Prometheus format output
- Coverage: 100%

**Handlers Package**
- **Routes** (`pkg/handlers/routes_test.go`): Tests HTTP endpoints
- **Service Emulator** (`pkg/handlers/service-emulator_test.go`): Tests service management API
- Coverage: 93.3%

### Integration Tests

**Main Package** (`cmd/main_test.go`)
- Tests complete application with all middleware and routes
- Tests request flow through the full stack
- Tests concurrent requests
- Tests all endpoints respond correctly

### Benchmarks

All test files include benchmark functions for performance testing:

```bash
# Run all benchmarks
make bench

# Run with memory profiling
make bench-mem

# Run with CPU profiling
make bench-cpu
```

## Test Conventions

### Naming
- Test files: `*_test.go`
- Test functions: `TestFunctionName`
- Benchmark functions: `BenchmarkFunctionName`
- Table-driven tests with descriptive names

### Structure
```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
    }{
        {"description", input, expected},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Coverage Goals

- **Overall**: > 90% ✅ (Currently 93.6%)
- **Critical paths**: 100% ✅
- **Handlers**: > 90% ✅
- **Business logic**: 100% ✅

## CI/CD Integration

Tests can be easily integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run tests
  run: make test

- name: Generate coverage
  run: make test-coverage

- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    files: ./coverage.out
```

## Test Dependencies

- **testing**: Go standard library
- **github.com/gofiber/fiber/v2**: For HTTP testing
- **net/http/httptest**: For HTTP request/response testing

No external testing frameworks required - all tests use Go's built-in testing package.

## Writing New Tests

When adding new features, ensure:

1. **Unit tests** for all new functions
2. **Integration tests** for new endpoints
3. **Benchmarks** for performance-critical code
4. **Table-driven tests** for multiple scenarios
5. **Error cases** are tested
6. **Edge cases** are covered

Example:

```go
func TestNewFeature(t *testing.T) {
    // Arrange
    input := setupTestData()
    
    // Act
    result, err := NewFeature(input)
    
    // Assert
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

## Troubleshooting

### Tests timeout
Some tests (especially system info) may take longer on slower systems. Adjust timeouts:

```go
resp, err := app.Test(req, 5000) // 5 second timeout
```

### External dependencies
Tests that require external services (curl tests, docker info) are marked with `skipTest` flag and can be enabled locally.

### Flaky tests
If you encounter flaky tests:
1. Check for race conditions: `go test -race ./...`
2. Run multiple times: `go test -count=10 ./...`
3. Ensure proper test isolation

## Continuous Improvement

- Monitor coverage trends
- Add tests for bug fixes
- Update tests when refactoring
- Keep tests fast and reliable
- Document complex test scenarios

For questions or issues with tests, please open an issue on the repository.
