package cmd

import (
	"testing"
)

func TestExecute(t *testing.T) {

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
