package cmd

import (
	"fmt"

	parser "loganalyzer/internal/analyzer"
	"loganalyzer/internal/config"

	"github.com/spf13/cobra"
)

var (
	configPath string
	outputPath string
)

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
		if err := reporter.SaveToFile(outputPath); err != nil {
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

	analyzeCmd.MarkFlagRequired("config")

	analyzeCmd.Example = `  # Analyze logs with config file only
  loganalyzer analyze --config config.json

  # Analyze logs and save results to file
  loganalyzer analyze --config config.json --output report.json

  # Using short flags
  loganalyzer analyze -c config.json -o report.json`
} 