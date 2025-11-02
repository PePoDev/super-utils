package services

import (
	"os/exec"

	"github.com/pepodev/super-utils/pkg/logger"
	"github.com/pepodev/super-utils/pkg/utils"
	"go.uber.org/zap"
)

// ExecuteCommand executes a shell command and returns the output
func ExecuteCommand(command string) (string, error) {
	// Sanitize command for logging (truncate if too long)
	logCommand := command
	if len(logCommand) > 100 {
		logCommand = logCommand[:97] + "..."
	}

	logger.Info("Executing shell command",
		zap.String("command", logCommand),
	)

	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()

	if err != nil {
		logger.Error("Shell command failed",
			zap.String("command", logCommand),
			zap.String("output", string(out)),
			zap.Error(err),
		)
		return string(out) + "\nError: " + err.Error(), err
	}

	logger.Info("Shell command executed successfully",
		zap.String("command", logCommand),
		zap.Int("output_length", len(out)),
	)

	return string(out), nil
}

// ExecuteToolCommand executes a tool with arguments (legacy format)
func ExecuteToolCommand(tool string, args string) (string, error) {
	argList := []string{}
	if args != "" {
		argList = utils.ParseArgs(args)
	}

	logger.Info("Executing tool command",
		zap.String("tool", tool),
		zap.String("args", args),
		zap.Int("arg_count", len(argList)),
	)

	cmd := exec.Command(tool, argList...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		logger.Error("Tool command failed",
			zap.String("tool", tool),
			zap.Strings("parsed_args", argList),
			zap.String("output", string(out)),
			zap.Error(err),
		)
		return string(out) + "\nError: " + err.Error(), err
	}

	logger.Info("Tool command executed successfully",
		zap.String("tool", tool),
		zap.Int("output_length", len(out)),
	)

	return string(out), nil
}
