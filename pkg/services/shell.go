package services

import (
	"os/exec"

	"github.com/pepodev/super-utils/pkg/utils"
)

// ExecuteCommand executes a shell command and returns the output
func ExecuteCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out) + "\nError: " + err.Error(), err
	}

	return string(out), nil
}

// ExecuteToolCommand executes a tool with arguments (legacy format)
func ExecuteToolCommand(tool string, args string) (string, error) {
	argList := []string{}
	if args != "" {
		argList = utils.ParseArgs(args)
	}

	cmd := exec.Command(tool, argList...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out) + "\nError: " + err.Error(), err
	}

	return string(out), nil
}
