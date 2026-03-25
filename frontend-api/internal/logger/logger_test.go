package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogger_New(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &Config{
		Level:    "debug",
		Path:     tmpDir,
		Filename: "test.log",
		Console:  true,
	}

	logger := New(cfg)
	if logger == nil {
		t.Fatal("Expected logger to be created")
	}
	defer logger.Close()

	// Test logging at different levels
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")

	// Check log file was created
	logPath := filepath.Join(tmpDir, "test.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Expected log file to be created")
	}
}

func TestLogger_Level(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		minLevel Level
	}{
		{"debug level", "debug", DEBUG},
		{"info level", "info", INFO},
		{"warn level", "warn", WARN},
		{"error level", "error", ERROR},
		{"unknown level defaults to info", "unknown", INFO},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &Logger{
				config: &Config{Level: tt.level},
				stdout: os.Stdout,
			}

			if !logger.shouldLog(tt.minLevel) {
				t.Errorf("Expected shouldLog(%v) to be true", tt.minLevel)
			}
		})
	}
}

func TestLogger_shouldLog(t *testing.T) {
	logger := &Logger{
		config: &Config{Level: "info"},
		stdout: os.Stdout,
	}

	// DEBUG should not log when level is INFO
	if logger.shouldLog(DEBUG) {
		t.Error("DEBUG should not log when level is INFO")
	}

	// INFO should log when level is INFO
	if !logger.shouldLog(INFO) {
		t.Error("INFO should log when level is INFO")
	}

	// WARN should log when level is INFO
	if !logger.shouldLog(WARN) {
		t.Error("WARN should log when level is INFO")
	}
}

func TestLogger_write(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		config:  &Config{Level: "debug", Console: true}, // Console must be true to write to stdout
		stdout:  &buf,
		file:    nil,
	}

	logger.write(INFO, "test message %s", "arg")

	output := buf.String()
	if !strings.Contains(output, "test message arg") {
		t.Errorf("Expected log output to contain 'test message arg', got '%s'", output)
	}
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected log output to contain 'INFO', got '%s'", output)
	}
}

func TestLogger_Default(t *testing.T) {
	// Reset defaultLogger for testing
	defaultLogger = nil

	// Create a default logger
	defaultLogger = &Logger{
		config: &Config{
			Level:   "debug",
			Console: true,
		},
		stdout: os.Stdout,
	}

	logger := Default()
	if logger == nil {
		t.Fatal("Expected default logger to be created")
	}

	// Calling Default again should return the same logger
	logger2 := Default()
	if logger != logger2 {
		t.Error("Default should return the same logger instance")
	}
}

func TestLogger_Close_NilFile(t *testing.T) {
	logger := &Logger{
		config: &Config{Level: "debug"},
		stdout: os.Stdout,
		file:   nil,
	}

	// Close should not error when file is nil
	if err := logger.Close(); err != nil {
		t.Errorf("Close should not return error: %v", err)
	}
}

func TestLogger_Fatal(t *testing.T) {
	// Fatal calls os.Exit, so we can't test it directly
	// Just verify it doesn't panic
	logger := &Logger{
		config: &Config{Level: "debug"},
		stdout: os.Stdout,
	}
	t.Log("Fatal test: os.Exit cannot be tested directly")
	_ = logger
}

func TestEntry_FormatMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		config: &Config{Level: "debug"},
		stdout: &buf,
	}

	entry := &Entry{
		logger: logger,
		fields: map[string]any{"key1": "value1", "key2": "value2"},
	}

	msg := entry.formatMessage("test %s", []any{"arg"})
	if !strings.Contains(msg, "key1") || !strings.Contains(msg, "value1") {
		t.Errorf("Expected formatted message to contain fields, got '%s'", msg)
	}
	if !strings.Contains(msg, "test arg") {
		t.Errorf("Expected formatted message to contain 'test arg', got '%s'", msg)
	}
}

func TestEntry_FormatMessage_NoArgs(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		config: &Config{Level: "debug"},
		stdout: &buf,
	}

	entry := &Entry{
		logger: logger,
		fields: map[string]any{"field": "value"},
	}

	msg := entry.formatMessage("simple message", nil)
	if !strings.Contains(msg, "simple message") {
		t.Errorf("Expected formatted message to contain 'simple message', got '%s'", msg)
	}
}

func TestEntry_LogMethods(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		config:  &Config{Level: "debug", Console: true}, // Console must be true
		stdout:  &buf,
		file:    nil,
	}

	entry := &Entry{
		logger: logger,
		fields: map[string]any{"test": "value"},
	}

	entry.Debug("debug message")
	entry.Info("info message")
	entry.Warn("warn message")
	entry.Error("error message")

	// Verify log entries were written
	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Error("Expected debug message in output")
	}
}

func TestLevelNames(t *testing.T) {
	expected := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	if len(levelNames) != len(expected) {
		t.Errorf("Expected %d level names, got %d", len(expected), len(levelNames))
	}

	for i, name := range levelNames {
		if name != expected[i] {
			t.Errorf("Expected level name '%s' at index %d, got '%s'", expected[i], i, name)
		}
	}
}

func TestLevel_Constants(t *testing.T) {
	if DEBUG != 0 {
		t.Errorf("Expected DEBUG to be 0, got %d", DEBUG)
	}
	if INFO != 1 {
		t.Errorf("Expected INFO to be 1, got %d", INFO)
	}
	if WARN != 2 {
		t.Errorf("Expected WARN to be 2, got %d", WARN)
	}
	if ERROR != 3 {
		t.Errorf("Expected ERROR to be 3, got %d", ERROR)
	}
	if FATAL != 4 {
		t.Errorf("Expected FATAL to be 4, got %d", FATAL)
	}
}

// Test convenience functions
func TestDebug(t *testing.T) {
	// Reset default logger
	defaultLogger = nil

	var buf bytes.Buffer
	defaultLogger = &Logger{
		config:  &Config{Level: "debug", Console: true}, // Console must be true
		stdout:  &buf,
		file:    nil,
	}

	Debug("test debug message")
	output := buf.String()
	if !strings.Contains(output, "test debug message") {
		t.Errorf("Expected debug message in output, got '%s'", output)
	}
}

func TestInfo(t *testing.T) {
	// Reset default logger
	defaultLogger = nil

	var buf bytes.Buffer
	defaultLogger = &Logger{
		config:  &Config{Level: "debug", Console: true}, // Console must be true
		stdout:  &buf,
		file:    nil,
	}

	Info("test info message")
	output := buf.String()
	if !strings.Contains(output, "test info message") {
		t.Errorf("Expected info message in output, got '%s'", output)
	}
}

func TestWarn(t *testing.T) {
	// Reset default logger
	defaultLogger = nil

	var buf bytes.Buffer
	defaultLogger = &Logger{
		config:  &Config{Level: "debug", Console: true}, // Console must be true
		stdout:  &buf,
		file:    nil,
	}

	Warn("test warn message")
	output := buf.String()
	if !strings.Contains(output, "test warn message") {
		t.Errorf("Expected warn message in output, got '%s'", output)
	}
}

func TestError(t *testing.T) {
	// Reset default logger
	defaultLogger = nil

	var buf bytes.Buffer
	defaultLogger = &Logger{
		config:  &Config{Level: "debug", Console: true}, // Console must be true
		stdout:  &buf,
		file:    nil,
	}

	Error("test error message")
	output := buf.String()
	if !strings.Contains(output, "test error message") {
		t.Errorf("Expected error message in output, got '%s'", output)
	}
}
