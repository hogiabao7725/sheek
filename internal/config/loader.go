package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDirName  = ".config"
	configAppName  = "sheek"
	configFileName = "config.json"
)

// GetConfigDir returns the path to the sheek config directory (~/.config/sheek)
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, configDirName, configAppName)
	return configDir, nil
}

// GetConfigPath returns the full path to the config file (~/.config/sheek/config.json)
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, configFileName), nil
}

// LoadConfig loads the configuration from the config file.
// If the config file doesn't exist, it creates it with default values.
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, create it with default values
		defaultConfig := DefaultConfig()
		if err := SaveConfig(defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return defaultConfig, nil
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfgWithDefaults := DefaultConfig()
	if err := json.Unmarshal(data, cfgWithDefaults); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and apply defaults for missing fields
	cfg := validateAndMergeDefaults(*cfgWithDefaults)

	return &cfg, nil
}

// SaveConfig saves the configuration to the config file.
// Creates the config directory if it doesn't exist.
func SaveConfig(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// validateAndMergeDefaults validates the config and merges missing fields with defaults
func validateAndMergeDefaults(cfg Config) Config {
	defaults := DefaultConfig()

	// Validate and apply defaults for numeric fields
	if cfg.MaxItems <= 0 {
		cfg.MaxItems = defaults.MaxItems
	}
	if cfg.Height <= 0 {
		cfg.Height = defaults.Height
	}
	if cfg.Margin < 0 {
		cfg.Margin = defaults.Margin
	}

	// Validate and apply defaults for string fields
	if cfg.Mode != "exact" && cfg.Mode != "fuzzy" {
		cfg.Mode = defaults.Mode
	}
	if cfg.Placeholder == "" {
		cfg.Placeholder = defaults.Placeholder
	}
	if cfg.Title == "" {
		cfg.Title = defaults.Title
	}

	// Validate input char limit
	if cfg.Limit <= 0 {
		cfg.Limit = defaults.Limit
	}

	// Validate and merge color config
	cfg.Colors = validateColors(cfg.Colors, defaults.Colors)

	return cfg
}

// validateColors validates color values and applies defaults for missing/invalid colors
func validateColors(colors, defaults ColorConfig) ColorConfig {
	if colors.Primary == "" {
		colors.Primary = defaults.Primary
	}
	if colors.Secondary == "" {
		colors.Secondary = defaults.Secondary
	}
	if colors.Text == "" {
		colors.Text = defaults.Text
	}
	if colors.Border == "" {
		colors.Border = defaults.Border
	}
	if colors.Muted == "" {
		colors.Muted = defaults.Muted
	}
	if colors.Selected == "" {
		colors.Selected = defaults.Selected
	}
	if colors.Highlight == "" {
		colors.Highlight = defaults.Highlight
	}
	if colors.Background == "" {
		colors.Background = defaults.Background
	}

	return colors
}
