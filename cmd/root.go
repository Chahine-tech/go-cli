package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loganalyzer",
	Short: "A distributed log analysis tool",
	Long: `LogAnalyzer is a command-line tool designed to help system administrators 
analyze log files from various sources (servers, applications) in parallel.

It provides concurrent processing of multiple log files with robust error handling
and JSON-based configuration and reporting capabilities.`,
	Version: "1.0.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
} 