package metrics

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	StartTime     = time.Now()
	TotalRequests = 0
)

// RegisterHealthMetrics registers health and metrics endpoints
func RegisterHealthMetrics(app *fiber.App) {
	// Health check endpoint
	app.Get("/health", HealthHandler)

	// Readiness probe
	app.Get("/ready", ReadinessHandler)

	// Liveness probe
	app.Get("/live", LivenessHandler)

	// Metrics endpoint (JSON format)
	app.Get("/metrics", MetricsHandler)

	// Prometheus-style metrics endpoint
	app.Get("/metrics/prometheus", PrometheusMetricsHandler)
}

// HealthHandler returns the health status of the application
func HealthHandler(c *fiber.Ctx) error {
	uptime := time.Since(StartTime)

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    uptime.String(),
		"version":   "1.0.0",
	}

	return c.JSON(health)
}

// ReadinessHandler returns the readiness status
func ReadinessHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"ready":     true,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// LivenessHandler returns the liveness status
func LivenessHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"alive":     true,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// MetricsHandler returns application metrics in JSON format
func MetricsHandler(c *fiber.Ctx) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(StartTime)

	metrics := map[string]interface{}{
		"uptime_seconds": uptime.Seconds(),
		"uptime":         uptime.String(),
		"goroutines":     runtime.NumGoroutine(),
		"total_requests": TotalRequests,
		"memory": map[string]interface{}{
			"alloc_bytes":       m.Alloc,
			"alloc_mb":          float64(m.Alloc) / 1024 / 1024,
			"total_alloc_bytes": m.TotalAlloc,
			"total_alloc_mb":    float64(m.TotalAlloc) / 1024 / 1024,
			"sys_bytes":         m.Sys,
			"sys_mb":            float64(m.Sys) / 1024 / 1024,
			"num_gc":            m.NumGC,
			"heap_alloc_bytes":  m.HeapAlloc,
			"heap_alloc_mb":     float64(m.HeapAlloc) / 1024 / 1024,
			"heap_sys_bytes":    m.HeapSys,
			"heap_sys_mb":       float64(m.HeapSys) / 1024 / 1024,
			"heap_inuse_bytes":  m.HeapInuse,
			"heap_inuse_mb":     float64(m.HeapInuse) / 1024 / 1024,
		},
		"runtime": map[string]interface{}{
			"go_version": runtime.Version(),
			"go_os":      runtime.GOOS,
			"go_arch":    runtime.GOARCH,
			"num_cpu":    runtime.NumCPU(),
		},
	}

	return c.JSON(metrics)
}

// PrometheusMetricsHandler returns metrics in Prometheus format
func PrometheusMetricsHandler(c *fiber.Ctx) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(StartTime)

	prometheusMetrics := `# HELP super_utils_uptime_seconds Application uptime in seconds
# TYPE super_utils_uptime_seconds gauge
super_utils_uptime_seconds ` + fmt.Sprintf("%.0f", uptime.Seconds()) + `

# HELP super_utils_goroutines Number of goroutines
# TYPE super_utils_goroutines gauge
super_utils_goroutines ` + fmt.Sprintf("%d", runtime.NumGoroutine()) + `

# HELP super_utils_memory_alloc_bytes Bytes of allocated heap objects
# TYPE super_utils_memory_alloc_bytes gauge
super_utils_memory_alloc_bytes ` + fmt.Sprintf("%d", m.Alloc) + `

# HELP super_utils_memory_sys_bytes Total bytes of memory obtained from the OS
# TYPE super_utils_memory_sys_bytes gauge
super_utils_memory_sys_bytes ` + fmt.Sprintf("%d", m.Sys) + `

# HELP super_utils_gc_total Total number of GC runs
# TYPE super_utils_gc_total counter
super_utils_gc_total ` + fmt.Sprintf("%d", m.NumGC) + `

# HELP super_utils_requests_total Total number of HTTP requests
# TYPE super_utils_requests_total counter
super_utils_requests_total ` + fmt.Sprintf("%d", TotalRequests) + `
`

	c.Set("Content-Type", "text/plain; version=0.0.4")
	return c.SendString(prometheusMetrics)
}
