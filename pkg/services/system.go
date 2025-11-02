package services

import (
	"os/exec"
	"strings"

	"github.com/pepodev/super-utils/pkg/logger"
	"go.uber.org/zap"
)

// SystemInfo represents system information
type SystemInfo map[string]string

// GetSystemInfo collects and returns system information
func GetSystemInfo() SystemInfo {
	logger.Debug("Starting system information collection")
	info := make(SystemInfo)
	fieldsCollected := 0

	// Hostname
	if out, err := exec.Command("hostname").Output(); err == nil {
		info["hostname"] = strings.TrimSpace(string(out))
		fieldsCollected++
	} else {
		logger.Warn("Failed to get hostname", zap.Error(err))
	}

	// OS Info
	if out, err := exec.Command("sh", "-c", "cat /etc/os-release 2>/dev/null | grep PRETTY_NAME | cut -d'=' -f2 | tr -d '\"'").Output(); err == nil {
		info["os"] = strings.TrimSpace(string(out))
		fieldsCollected++
	} else {
		logger.Debug("Failed to get OS info", zap.Error(err))
	}

	// Kernel version
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		info["kernel"] = strings.TrimSpace(string(out))
		fieldsCollected++
	} else {
		logger.Warn("Failed to get kernel version", zap.Error(err))
	}

	// Uptime
	if out, err := exec.Command("uptime", "-p").Output(); err == nil {
		info["uptime"] = strings.TrimSpace(string(out))
		fieldsCollected++
	} else {
		logger.Debug("Failed to get uptime", zap.Error(err))
	}

	// CPU Info
	if out, err := exec.Command("sh", "-c", "lscpu | grep -E 'Model name|Architecture|CPU\\(s\\)|Thread|Core|Socket|MHz'").Output(); err == nil {
		info["cpu"] = string(out)
		fieldsCollected++
	} else {
		logger.Debug("Failed to get CPU info", zap.Error(err))
	}

	// Memory Info
	if out, err := exec.Command("free", "-h").Output(); err == nil {
		info["memory"] = string(out)
		fieldsCollected++
	} else {
		logger.Warn("Failed to get memory info", zap.Error(err))
	}

	// Disk Usage
	if out, err := exec.Command("df", "-h").Output(); err == nil {
		info["disk"] = string(out)
		fieldsCollected++
	} else {
		logger.Warn("Failed to get disk info", zap.Error(err))
	}

	// Network Interfaces
	if out, err := exec.Command("sh", "-c", "ip -br addr show || ifconfig").Output(); err == nil {
		info["network"] = string(out)
		fieldsCollected++
	} else {
		logger.Debug("Failed to get network info", zap.Error(err))
	}

	// Load Average
	if out, err := exec.Command("sh", "-c", "uptime | awk -F'load average:' '{print $2}'").Output(); err == nil {
		info["load"] = "Load Average:" + string(out)
		fieldsCollected++
	} else {
		logger.Debug("Failed to get load average", zap.Error(err))
	}

	// Top Processes
	if out, err := exec.Command("sh", "-c", "ps aux --sort=-%mem | head -11").Output(); err == nil {
		info["processes"] = string(out)
		fieldsCollected++
	} else {
		logger.Debug("Failed to get process list", zap.Error(err))
	}

	// Docker info (if available)
	if out, err := exec.Command("docker", "info", "--format", "{{.ServerVersion}}\nContainers: {{.Containers}}\nImages: {{.Images}}").Output(); err == nil {
		info["docker"] = string(out)
		fieldsCollected++
		logger.Debug("Docker info collected")
	} else {
		logger.Debug("Docker not available or failed to collect info", zap.Error(err))
	}

	logger.Info("System information collection completed",
		zap.Int("fields_collected", fieldsCollected),
		zap.Int("total_fields", len(info)),
	)

	return info
}
