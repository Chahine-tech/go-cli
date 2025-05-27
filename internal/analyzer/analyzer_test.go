package parser

import (
	"os"
	"path/filepath"
	"testing"

	"loganalyzer/internal/config"
)

func TestNewAnalyzer(t *testing.T) {
	cfg := &config.Config{
		Logs: []config.LogConfig{
			{
				ID:   "test1",
				Path: "/var/log/test1.log",
				Type: "nginx",
			},
		},
	}

	analyzer := NewAnalyzer(cfg)
	if analyzer == nil {
		t.Fatal("NewAnalyzer() returned nil")
	}

	if analyzer.config != cfg {
		t.Error("NewAnalyzer() did not set config correctly")
	}

	if analyzer.reporter == nil {
		t.Error("NewAnalyzer() did not initialize reporter")
	}
}

func TestAnalyzeAllLogs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "loganalyzer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	log1Path := filepath.Join(tempDir, "test1.log")
	log2Path := filepath.Join(tempDir, "test2.log")
	log3Path := filepath.Join(tempDir, "nonexistent.log")

	if err := os.WriteFile(log1Path, []byte("test log 1"), 0644); err != nil {
		t.Fatalf("Failed to write test log 1: %v", err)
	}
	if err := os.WriteFile(log2Path, []byte("test log 2"), 0644); err != nil {
		t.Fatalf("Failed to write test log 2: %v", err)
	}

	tests := []struct {
		name        string
		logs        []config.LogConfig
		expectError bool
	}{
		{
			name:        "Empty logs",
			logs:        []config.LogConfig{},
			expectError: true,
		},
		{
			name: "Valid logs",
			logs: []config.LogConfig{
				{
					ID:   "log1",
					Path: log1Path,
					Type: "nginx",
				},
				{
					ID:   "log2",
					Path: log2Path,
					Type: "apache",
				},
			},
			expectError: false,
		},
		{
			name: "Mix of valid and invalid logs",
			logs: []config.LogConfig{
				{
					ID:   "log1",
					Path: log1Path,
					Type: "nginx",
				},
				{
					ID:   "log3",
					Path: log3Path,
					Type: "nginx",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Logs: tt.logs}
			analyzer := NewAnalyzer(cfg)

			err := analyzer.AnalyzeAllLogs()
			if (err != nil) != tt.expectError {
				t.Errorf("AnalyzeAllLogs() error = %v, expectError %v", err, tt.expectError)
			}

			reporter := analyzer.GetReporter()
			if reporter == nil {
				t.Error("GetReporter() returned nil")
			}
		})
	}
}

func TestCheckFileAccess(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "loganalyzer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	validFile := filepath.Join(tempDir, "valid.log")
	if err := os.WriteFile(validFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	dirPath := filepath.Join(tempDir, "dir")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "Valid file",
			filePath:    validFile,
			expectError: false,
		},
		{
			name:        "Nonexistent file",
			filePath:    filepath.Join(tempDir, "nonexistent.log"),
			expectError: true,
		},
		{
			name:        "Directory instead of file",
			filePath:    dirPath,
			expectError: true,
		},
	}

	analyzer := NewAnalyzer(&config.Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := analyzer.checkFileAccess(tt.filePath)
			if (err != nil) != tt.expectError {
				t.Errorf("checkFileAccess() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
