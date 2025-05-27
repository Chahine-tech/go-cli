package reporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewReporter(t *testing.T) {
	reporter := NewReporter()
	if reporter == nil {
		t.Fatal("NewReporter() returned nil")
	}

	if reporter.results == nil {
		t.Error("NewReporter() did not initialize results slice")
	}

	if len(reporter.results) != 0 {
		t.Error("NewReporter() did not initialize empty results slice")
	}
}

func TestAddResult(t *testing.T) {
	reporter := NewReporter()

	// Test adding a success result
	successResult := CreateSuccessResult("log1", "/var/log/test1.log")
	reporter.AddResult(successResult)

	if len(reporter.results) != 1 {
		t.Errorf("AddResult() did not add result, got %d results", len(reporter.results))
	}

	if reporter.results[0].LogID != "log1" {
		t.Errorf("AddResult() wrong LogID, got %s, want %s", reporter.results[0].LogID, "log1")
	}

	// Test adding a failure result
	failureResult := CreateFailureResult("log2", "/var/log/test2.log", "Test error", "Error details")
	reporter.AddResult(failureResult)

	if len(reporter.results) != 2 {
		t.Errorf("AddResult() did not add second result, got %d results", len(reporter.results))
	}

	if reporter.results[1].Status != "FAILURE" {
		t.Errorf("AddResult() wrong Status, got %s, want %s", reporter.results[1].Status, "FAILURE")
	}
}

func TestGetResults(t *testing.T) {
	reporter := NewReporter()

	// Add some test results
	results := []AnalysisResult{
		CreateSuccessResult("log1", "/var/log/test1.log"),
		CreateFailureResult("log2", "/var/log/test2.log", "Test error", "Error details"),
	}

	for _, result := range results {
		reporter.AddResult(result)
	}

	gotResults := reporter.GetResults()
	if len(gotResults) != len(results) {
		t.Errorf("GetResults() returned wrong number of results, got %d, want %d", len(gotResults), len(results))
	}

	for i, result := range gotResults {
		if result.LogID != results[i].LogID {
			t.Errorf("GetResults() result %d wrong LogID, got %s, want %s", i, result.LogID, results[i].LogID)
		}
	}
}

func TestSaveToFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "loganalyzer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	reporter := NewReporter()

	// Add test results
	results := []AnalysisResult{
		CreateSuccessResult("log1", "/var/log/test1.log"),
		CreateFailureResult("log2", "/var/log/test2.log", "Test error", "Error details"),
	}

	for _, result := range results {
		reporter.AddResult(result)
	}

	// Test saving to file
	outputPath := filepath.Join(tempDir, "results.json")
	err = reporter.SaveToFile(outputPath)
	if err != nil {
		t.Errorf("SaveToFile() error = %v", err)
	}

	// Verify file contents
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Errorf("Failed to read saved file: %v", err)
	}

	var savedResults []AnalysisResult
	if err := json.Unmarshal(data, &savedResults); err != nil {
		t.Errorf("Failed to unmarshal saved results: %v", err)
	}

	if len(savedResults) != len(results) {
		t.Errorf("Saved results count mismatch, got %d, want %d", len(savedResults), len(results))
	}

	// Test saving to invalid path
	err = reporter.SaveToFile("/nonexistent/path/results.json")
	if err == nil {
		t.Error("SaveToFile() expected error for invalid path, got nil")
	}
}

func TestCreateSuccessResult(t *testing.T) {
	result := CreateSuccessResult("log1", "/var/log/test1.log")

	if result.LogID != "log1" {
		t.Errorf("CreateSuccessResult() wrong LogID, got %s, want %s", result.LogID, "log1")
	}

	if result.FilePath != "/var/log/test1.log" {
		t.Errorf("CreateSuccessResult() wrong FilePath, got %s, want %s", result.FilePath, "/var/log/test1.log")
	}

	if result.Status != "OK" {
		t.Errorf("CreateSuccessResult() wrong Status, got %s, want %s", result.Status, "OK")
	}

	if result.Message != "Analysis completed successfully." {
		t.Errorf("CreateSuccessResult() wrong Message, got %s, want %s", result.Message, "Analysis completed successfully.")
	}

	if result.ErrorDetails != "" {
		t.Errorf("CreateSuccessResult() wrong ErrorDetails, got %s, want empty string", result.ErrorDetails)
	}
}

func TestCreateFailureResult(t *testing.T) {
	result := CreateFailureResult("log1", "/var/log/test1.log", "Test error", "Error details")

	if result.LogID != "log1" {
		t.Errorf("CreateFailureResult() wrong LogID, got %s, want %s", result.LogID, "log1")
	}

	if result.FilePath != "/var/log/test1.log" {
		t.Errorf("CreateFailureResult() wrong FilePath, got %s, want %s", result.FilePath, "/var/log/test1.log")
	}

	if result.Status != "FAILURE" {
		t.Errorf("CreateFailureResult() wrong Status, got %s, want %s", result.Status, "FAILURE")
	}

	if result.Message != "Test error" {
		t.Errorf("CreateFailureResult() wrong Message, got %s, want %s", result.Message, "Test error")
	}

	if result.ErrorDetails != "Error details" {
		t.Errorf("CreateFailureResult() wrong ErrorDetails, got %s, want %s", result.ErrorDetails, "Error details")
	}
}
