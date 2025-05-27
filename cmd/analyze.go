package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	parser "loganalyzer/internal/analyzer"
	"loganalyzer/internal/config"

	"github.com/spf13/cobra"
)

var (
	configPath string
	outputPath string
)

func formatOutputPath(path string) string {
	if path == "" {
		return ""
	}

	// Get the directory and filename
	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	// Get current date in YYMMDD format
	timestamp := time.Now().Format("060102") // Go's reference time is 2006-01-02

	// Insert timestamp before the extension
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]
	newFilename := fmt.Sprintf("%s_%s%s", timestamp, nameWithoutExt, ext)

	// Combine directory and new filename
	return filepath.Join(dir, newFilename)
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze log files based on JSON configuration",
	Long: `The analyze command processes multiple log files concurrently based on a JSON 
configuration file. It performs parallel analysis using goroutines and outputs 
results both to the console and optionally to a JSON report file.

Features:
- Concurrent processing of multiple log files using goroutines
- Custom error handling for file access and parsing errors
- JSON configuration input and optional JSON report output
- Real-time progress updates and detailed error reporting
- Automatic timestamp in output filenames (YYMMDD format)

Example usage:
  loganalyzer analyze --config config.json --output report.json
  loganalyzer analyze -c config.json -o report.json`,
	RunE: runAnalyze,
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	if configPath == "" {
		return fmt.Errorf("config file path is required (use --config or -c flag)")
	}

	fmt.Printf("Loading configuration from: %s\n", configPath)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Printf("Loaded configuration with %d log files\n", len(cfg.Logs))

	analyzer := parser.NewAnalyzer(cfg)

	if err := analyzer.AnalyzeAllLogs(); err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	reporter := analyzer.GetReporter()
	reporter.PrintSummary()

	if outputPath != "" {
		timestampedPath := formatOutputPath(outputPath)
		if err := reporter.SaveToFile(timestampedPath); err != nil {
			return fmt.Errorf("failed to save results: %w", err)
		}
	}

	fmt.Println("\nAnalysis completed successfully!")
	return nil
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to JSON configuration file (required)")
	analyzeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to JSON output file (optional)")

	if err := analyzeCmd.MarkFlagRequired("config"); err != nil {
		panic(fmt.Sprintf("Failed to mark config flag as required: %v", err))
	}

	analyzeCmd.Example = `  # Analyze logs with config file only
  loganalyzer analyze --config config.json

  # Analyze logs and save results to file (with timestamp)
  loganalyzer analyze --config config.json --output report.json
  # Output will be saved as: YYMMDD_report.json (e.g., 240524_report.json)

  # Using short flags
  loganalyzer analyze -c config.json -o report.json`
}
