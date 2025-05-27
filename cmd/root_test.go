package cmd

import (
	"testing"
)

func TestExecute(t *testing.T) {
	// Since Execute() calls os.Exit(1) on error, we can't directly test it
	// Instead, we'll test that the command is properly initialized
	if rootCmd.Use != "loganalyzer" {
		t.Errorf("Expected command name 'loganalyzer', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Expected non-empty short description")
	}

	if rootCmd.Long == "" {
		t.Error("Expected non-empty long description")
	}

	if rootCmd.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", rootCmd.Version)
	}
}
