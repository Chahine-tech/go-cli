package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "loganalyzer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		configJSON  string
		expectError bool
	}{
		{
			name: "Valid config",
			configJSON: `[
				{
					"id": "log1",
					"path": "/var/log/app1.log",
					"type": "nginx"
				},
				{
					"id": "log2",
					"path": "/var/log/app2.log",
					"type": "apache"
				}
			]`,
			expectError: false,
		},
		{
			name:        "Empty config",
			configJSON:  `[]`,
			expectError: true,
		},
		{
			name: "Missing ID",
			configJSON: `[
				{
					"path": "/var/log/app1.log",
					"type": "nginx"
				}
			]`,
			expectError: true,
		},
		{
			name: "Missing path",
			configJSON: `[
				{
					"id": "log1",
					"type": "nginx"
				}
			]`,
			expectError: true,
		},
		{
			name: "Missing type",
			configJSON: `[
				{
					"id": "log1",
					"path": "/var/log/app1.log"
				}
			]`,
			expectError: true,
		},
		{
			name: "Duplicate ID",
			configJSON: `[
				{
					"id": "log1",
					"path": "/var/log/app1.log",
					"type": "nginx"
				},
				{
					"id": "log1",
					"path": "/var/log/app2.log",
					"type": "apache"
				}
			]`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test config file
			configPath := filepath.Join(tempDir, "config.json")
			if err := os.WriteFile(configPath, []byte(tt.configJSON), 0644); err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			cfg, err := LoadConfig(configPath)
			if (err != nil) != tt.expectError {
				t.Errorf("LoadConfig() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && cfg == nil {
				t.Error("LoadConfig() returned nil config when no error was expected")
			}
		})
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.json")
	if err == nil {
		t.Error("LoadConfig() expected error for nonexistent file, got nil")
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "loganalyzer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test config file with invalid JSON
	configPath := filepath.Join(tempDir, "config.json")
	invalidJSON := `{invalid json}`
	if err := os.WriteFile(configPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("LoadConfig() expected error for invalid JSON, got nil")
	}
}
