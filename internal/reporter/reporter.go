package reporter

import (
	"encoding/json"
	"fmt"
	"os"
)

type AnalysisResult struct {
	LogID        string `json:"log_id"`
	FilePath     string `json:"file_path"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details"`
}

type Reporter struct {
	results []AnalysisResult
}

func NewReporter() *Reporter {
	return &Reporter{
		results: make([]AnalysisResult, 0),
	}
}

func (r *Reporter) AddResult(result AnalysisResult) {
	r.results = append(r.results, result)
}

func (r *Reporter) GetResults() []AnalysisResult {
	return r.results
}

func (r *Reporter) PrintSummary() {
	fmt.Println("\n=== Analysis Summary ===")
	
	successCount := 0
	failureCount := 0
	
	for _, result := range r.results {
		status := "✓"
		if result.Status == "FAILURE" {
			status = "✗"
			failureCount++
		} else {
			successCount++
		}
		
		fmt.Printf("%s [%s] %s: %s\n", status, result.LogID, result.FilePath, result.Message)
		if result.ErrorDetails != "" {
			fmt.Printf("   Error: %s\n", result.ErrorDetails)
		}
	}
	
	fmt.Printf("\nTotal: %d logs analyzed (%d successful, %d failed)\n", 
		len(r.results), successCount, failureCount)
}

func (r *Reporter) SaveToFile(outputPath string) error {
	data, err := json.MarshalIndent(r.results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write results to file %s: %w", outputPath, err)
	}

	fmt.Printf("Analysis results saved to: %s\n", outputPath)
	return nil
}

func CreateSuccessResult(logID, filePath string) AnalysisResult {
	return AnalysisResult{
		LogID:        logID,
		FilePath:     filePath,
		Status:       "OK",
		Message:      "Analysis completed successfully.",
		ErrorDetails: "",
	}
}

func CreateFailureResult(logID, filePath, message, errorDetails string) AnalysisResult {
	return AnalysisResult{
		LogID:        logID,
		FilePath:     filePath,
		Status:       "FAILURE",
		Message:      message,
		ErrorDetails: errorDetails,
	}
} 