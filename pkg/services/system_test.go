package services

import (
	"testing"

	"github.com/pepodev/super-utils/pkg/logger"
)

func init() {
	// Initialize test logger (no-op) for all tests
	logger.InitTestLogger()
}

func TestGetSystemInfo(t *testing.T) {
	info := GetSystemInfo()

	if info == nil {
		t.Fatal("GetSystemInfo() returned nil")
	}

	// Test that we get some basic fields
	// Note: Some fields might not be available in all environments
	t.Run("Contains expected fields", func(t *testing.T) {
		expectedFields := []string{
			"hostname",
			"kernel",
			"uptime",
			"memory",
			"disk",
			"network",
		}

		for _, field := range expectedFields {
			if _, exists := info[field]; !exists {
				t.Logf("Warning: Field %s not found in system info (might be environment-specific)", field)
			}
		}
	})

	// Test hostname is not empty if it exists
	if hostname, exists := info["hostname"]; exists {
		if hostname == "" {
			t.Error("hostname should not be empty")
		}
		t.Logf("Hostname: %s", hostname)
	}

	// Test kernel is not empty if it exists
	if kernel, exists := info["kernel"]; exists {
		if kernel == "" {
			t.Error("kernel should not be empty")
		}
		t.Logf("Kernel: %s", kernel)
	}

	// Test memory info if it exists
	if memory, exists := info["memory"]; exists {
		if memory == "" {
			t.Error("memory should not be empty")
		}
		t.Logf("Memory info length: %d bytes", len(memory))
	}

	// Test disk info if it exists
	if disk, exists := info["disk"]; exists {
		if disk == "" {
			t.Error("disk should not be empty")
		}
		t.Logf("Disk info length: %d bytes", len(disk))
	}
}

func TestGetSystemInfoConsistency(t *testing.T) {
	// Call GetSystemInfo multiple times and ensure consistency
	info1 := GetSystemInfo()
	info2 := GetSystemInfo()

	if info1["hostname"] != info2["hostname"] {
		t.Error("Hostname should be consistent across calls")
	}

	if info1["kernel"] != info2["kernel"] {
		t.Error("Kernel version should be consistent across calls")
	}
}

func TestGetSystemInfoTypes(t *testing.T) {
	info := GetSystemInfo()

	// Verify map is properly structured (all values should be strings by type definition)
	for key, value := range info {
		if value == "" {
			t.Logf("Warning: Empty value for key %s", key)
		}
	}
}

func TestGetSystemInfoNotEmpty(t *testing.T) {
	info := GetSystemInfo()

	if len(info) == 0 {
		t.Error("GetSystemInfo() returned empty map")
	}

	t.Logf("GetSystemInfo() returned %d fields", len(info))
}

func TestGetSystemInfoNetworkField(t *testing.T) {
	info := GetSystemInfo()

	if network, exists := info["network"]; exists {
		if network == "" {
			t.Error("network should not be empty")
		}
		t.Logf("Network info: %s", network[:min(100, len(network))])
	}
}

func TestGetSystemInfoUptimeField(t *testing.T) {
	info := GetSystemInfo()

	if uptime, exists := info["uptime"]; exists {
		if uptime == "" {
			t.Error("uptime should not be empty")
		}
		t.Logf("Uptime: %s", uptime)
	}
}

func TestGetSystemInfoDockerField(t *testing.T) {
	info := GetSystemInfo()

	// Docker info is optional - it's only present if Docker is installed
	if docker, exists := info["docker"]; exists {
		t.Logf("Docker info found: %s", docker[:min(50, len(docker))])
	} else {
		t.Log("Docker info not available (Docker might not be installed)")
	}
}

func BenchmarkGetSystemInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetSystemInfo()
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
