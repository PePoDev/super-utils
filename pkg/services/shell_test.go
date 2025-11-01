package services

import (
	"strings"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	tests := []struct {
		name           string
		command        string
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "Simple echo command",
			command:        "echo hello",
			expectError:    false,
			expectedOutput: "hello",
		},
		{
			name:           "Command with output",
			command:        "echo 'test output'",
			expectError:    false,
			expectedOutput: "test output",
		},
		{
			name:        "Invalid command",
			command:     "nonexistentcommand123",
			expectError: true,
		},
		{
			name:           "Command with newline",
			command:        "printf 'line1\\nline2'",
			expectError:    false,
			expectedOutput: "line1\nline2",
		},
		{
			name:           "ls command",
			command:        "ls /tmp > /dev/null && echo success",
			expectError:    false,
			expectedOutput: "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ExecuteCommand(tt.command)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v, output: %s", err, output)
				}
				if tt.expectedOutput != "" && !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
				}
			}
		})
	}
}

func TestExecuteToolCommand(t *testing.T) {
	tests := []struct {
		name           string
		tool           string
		args           string
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "echo with single arg",
			tool:           "echo",
			args:           "hello",
			expectError:    false,
			expectedOutput: "hello",
		},
		{
			name:           "echo with multiple args",
			tool:           "echo",
			args:           "hello world",
			expectError:    false,
			expectedOutput: "hello world",
		},
		{
			name:           "echo with quoted args",
			tool:           "echo",
			args:           `"hello world"`,
			expectError:    false,
			expectedOutput: "hello world",
		},
		{
			name:        "Invalid tool",
			tool:        "nonexistenttool123",
			args:        "",
			expectError: true,
		},
		{
			name:           "printf with format",
			tool:           "printf",
			args:           "test",
			expectError:    false,
			expectedOutput: "test",
		},
		{
			name:           "ls with directory",
			tool:           "ls",
			args:           "/tmp",
			expectError:    false,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ExecuteToolCommand(tt.tool, tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v, output: %s", err, output)
				}
				if tt.expectedOutput != "" && !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
				}
			}
		})
	}
}

func TestExecuteCommandWithError(t *testing.T) {
	// Test command that exits with non-zero status
	output, err := ExecuteCommand("sh -c 'exit 1'")

	if err == nil {
		t.Error("Expected error for command with non-zero exit status")
	}

	if !strings.Contains(output, "Error:") {
		t.Error("Expected error message in output")
	}
}

func TestExecuteToolCommandEmptyArgs(t *testing.T) {
	// Test tool with empty args
	output, err := ExecuteToolCommand("echo", "")

	if err != nil {
		t.Errorf("Unexpected error with empty args: %v", err)
	}

	// echo with no args should output empty line or newline
	if output != "\n" && output != "" {
		t.Logf("Output from echo with no args: %q (acceptable)", output)
	}
}

func TestExecuteCommandOutputCapture(t *testing.T) {
	// Test that both stdout and stderr are captured
	output, err := ExecuteCommand("sh -c 'echo stdout; echo stderr >&2'")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !strings.Contains(output, "stdout") {
		t.Error("Expected stdout in output")
	}

	if !strings.Contains(output, "stderr") {
		t.Error("Expected stderr in output")
	}
}

func BenchmarkExecuteCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExecuteCommand("echo test")
	}
}

func BenchmarkExecuteToolCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExecuteToolCommand("echo", "test")
	}
}
