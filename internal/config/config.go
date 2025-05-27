package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type LogConfig struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type Config struct {
	Logs []LogConfig
}

func LoadConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var logs []LogConfig
	if err := json.Unmarshal(data, &logs); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	if err := validateConfig(logs); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &Config{Logs: logs}, nil
}

func validateConfig(logs []LogConfig) error {
	if len(logs) == 0 {
		return fmt.Errorf("no logs configured")
	}

	ids := make(map[string]bool)
	for _, log := range logs {
		if log.ID == "" {
			return fmt.Errorf("log entry missing ID")
		}
		if log.Path == "" {
			return fmt.Errorf("log entry %s missing path", log.ID)
		}
		if log.Type == "" {
			return fmt.Errorf("log entry %s missing type", log.ID)
		}
		if ids[log.ID] {
			return fmt.Errorf("duplicate log ID: %s", log.ID)
		}
		ids[log.ID] = true
	}

	return nil
} 