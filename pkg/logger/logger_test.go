package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Fatalf("InitLogger() failed: %v", err)
	}

	if Log == nil {
		t.Error("InitLogger() did not initialize Log")
	}

	// Test logging
	Log.Info("Test log message",
		zap.String("test", "value"),
	)

	// Clean up
	Sync()
}

func TestInitTestLogger(t *testing.T) {
	InitTestLogger()

	if Log == nil {
		t.Error("InitTestLogger() did not initialize Log")
	}

	// Test that no-op logger doesn't panic
	Log.Info("This should not produce output")
	Log.Error("This should not produce output")
	Log.Debug("This should not produce output")

	// Clean up
	Sync()
}

func TestSync(t *testing.T) {
	// Test with nil logger
	Log = nil
	Sync() // Should not panic

	// Test with initialized logger
	InitTestLogger()
	Sync() // Should not panic
}

func TestWrapperFunctions(t *testing.T) {
	// Initialize test logger (no-op)
	InitTestLogger()

	// Test all wrapper functions - should not panic
	Debug("Debug message", zap.String("key", "value"))
	Info("Info message", zap.String("key", "value"))
	Warn("Warn message", zap.String("key", "value"))
	Error("Error message", zap.String("key", "value"))

	// Test with nil logger
	Log = nil
	Debug("Should not panic")
	Info("Should not panic")
	Warn("Should not panic")
	Error("Should not panic")
}

func TestWrapperFunctionsWithRealLogger(t *testing.T) {
	// Initialize real logger
	err := InitLogger()
	if err != nil {
		t.Fatalf("InitLogger() failed: %v", err)
	}

	// Test all wrapper functions
	Debug("Debug message from test", zap.String("test", "wrapper"))
	Info("Info message from test", zap.String("test", "wrapper"))
	Warn("Warn message from test", zap.String("test", "wrapper"))
	Error("Error message from test", zap.String("test", "wrapper"))

	Sync()
}
