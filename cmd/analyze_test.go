package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFormatOutputPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "Simple filename",
			input:    "report.json",
			expected: time.Now().Format("060102") + "_report.json",
		},
		{
			name:     "Path with directory",
			input:    "/path/to/report.json",
			expected: filepath.Join("/path/to", time.Now().Format("060102")+"_report.json"),
		},
		{
			name:     "Filename without extension",
			input:    "report",
			expected: time.Now().Format("060102") + "_report",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatOutputPath(tt.input)
			if tt.input == "" {
				if got != "" {
					t.Errorf("formatOutputPath() = %v, want empty string", got)
				}
				return
			}

			// For non-empty paths, verify the format
			if !strings.Contains(got, time.Now().Format("060102")) {
				t.Errorf("formatOutputPath() output %v does not contain timestamp", got)
			}

			// Verify the base filename is preserved
			baseGot := filepath.Base(got)
			baseExpected := filepath.Base(tt.expected)
			if baseGot != baseExpected {
				t.Errorf("formatOutputPath() base filename = %v, want %v", baseGot, baseExpected)
			}

			// Verify the directory is preserved
			dirGot := filepath.Dir(got)
			dirExpected := filepath.Dir(tt.input)
			if dirGot != dirExpected {
				t.Errorf("formatOutputPath() directory = %v, want %v", dirGot, dirExpected)
			}
		})
	}
}

func TestAnalyzeCommand(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "loganalyzer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config file
	configContent := `[
		{
			"id": "log1",
			"path": "test.log",
			"type": "nginx"
		}
	]`
	configFilePath := filepath.Join(tempDir, "config.json")
	if err := os.WriteFile(configFilePath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Create a test log file
	logContent := `192.168.1.1 - - [01/Jan/2024:00:00:00 +0000] "GET / HTTP/1.1" 200 1234 "-" "Mozilla/5.0"`
	logPath := filepath.Join(tempDir, "test.log")
	if err := os.WriteFile(logPath, []byte(logContent), 0644); err != nil {
		t.Fatalf("Failed to write test log: %v", err)
	}

	// Test cases
	tests := []struct {
		name        string
		setConfig   bool
		setOutput   bool
		expectError bool
	}{
		{
			name:        "Missing config file",
			setConfig:   false,
			setOutput:   false,
			expectError: true,
		},
		{
			name:        "Valid config file",
			setConfig:   true,
			setOutput:   false,
			expectError: false,
		},
		{
			name:        "Valid config and output files",
			setConfig:   true,
			setOutput:   true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setConfig {
				configPath = configFilePath
			} else {
				configPath = ""
			}
			if tt.setOutput {
				outputPath = filepath.Join(tempDir, "output.json")
			} else {
				outputPath = ""
			}

			err := runAnalyze(nil, nil)
			if (err != nil) != tt.expectError {
				t.Errorf("runAnalyze() error = %v, expectError %v", err, tt.expectError)
			}

			// If output was specified, verify the timestamped file was created
			if tt.setOutput {
				expectedOutput := formatOutputPath(outputPath)
				if _, err := os.Stat(expectedOutput); os.IsNotExist(err) {
					t.Errorf("Expected output file %s was not created", expectedOutput)
				}
			}
		})
	}
}
