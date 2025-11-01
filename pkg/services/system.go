package services

import (
	"os/exec"
	"strings"
)

// SystemInfo represents system information
type SystemInfo map[string]string

// GetSystemInfo collects and returns system information
func GetSystemInfo() SystemInfo {
	info := make(SystemInfo)

	// Hostname
	if out, err := exec.Command("hostname").Output(); err == nil {
		info["hostname"] = strings.TrimSpace(string(out))
	}

	// OS Info
	if out, err := exec.Command("sh", "-c", "cat /etc/os-release 2>/dev/null | grep PRETTY_NAME | cut -d'=' -f2 | tr -d '\"'").Output(); err == nil {
		info["os"] = strings.TrimSpace(string(out))
	}

	// Kernel version
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		info["kernel"] = strings.TrimSpace(string(out))
	}

	// Uptime
	if out, err := exec.Command("uptime", "-p").Output(); err == nil {
		info["uptime"] = strings.TrimSpace(string(out))
	}

	// CPU Info
	if out, err := exec.Command("sh", "-c", "lscpu | grep -E 'Model name|Architecture|CPU\\(s\\)|Thread|Core|Socket|MHz'").Output(); err == nil {
		info["cpu"] = string(out)
	}

	// Memory Info
	if out, err := exec.Command("free", "-h").Output(); err == nil {
		info["memory"] = string(out)
	}

	// Disk Usage
	if out, err := exec.Command("df", "-h").Output(); err == nil {
		info["disk"] = string(out)
	}

	// Network Interfaces
	if out, err := exec.Command("sh", "-c", "ip -br addr show || ifconfig").Output(); err == nil {
		info["network"] = string(out)
	}

	// Load Average
	if out, err := exec.Command("sh", "-c", "uptime | awk -F'load average:' '{print $2}'").Output(); err == nil {
		info["load"] = "Load Average:" + string(out)
	}

	// Top Processes
	if out, err := exec.Command("sh", "-c", "ps aux --sort=-%mem | head -11").Output(); err == nil {
		info["processes"] = string(out)
	}

	// Docker info (if available)
	if out, err := exec.Command("docker", "info", "--format", "{{.ServerVersion}}\nContainers: {{.Containers}}\nImages: {{.Images}}").Output(); err == nil {
		info["docker"] = string(out)
	}

	return info
}
