package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"just-icon/internal/types"
)

const (
	ConfigFileName = "just-icon.json"
)

// Service handles configuration management
type Service struct {
	configPath string
}

// NewService creates a new configuration service
func NewService() *Service {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory is not available
		homeDir = "."
	}

	configPath := filepath.Join(homeDir, ConfigFileName)
	return &Service{
		configPath: configPath,
	}
}

// GetConfig reads the configuration from file
func (s *Service) GetConfig() (*types.Config, error) {
	config := &types.Config{}
	
	// Check if config file exists
	if _, err := os.Stat(s.configPath); os.IsNotExist(err) {
		// Return empty config if file doesn't exist
		return config, nil
	}
	
	// Read config file
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Parse JSON
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return config, nil
}

// SetConfig writes the configuration to file
func (s *Service) SetConfig(config *types.Config) error {
	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file (directly to home directory)
	if err := os.WriteFile(s.configPath, data, types.ConfigFilePerm); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// UpdateConfig updates specific fields in the configuration
func (s *Service) UpdateConfig(updates map[string]interface{}) error {
	config, err := s.GetConfig()
	if err != nil {
		return err
	}
	
	// Apply updates
	for key, value := range updates {
		switch key {
		case "openai_api_key":
			if v, ok := value.(string); ok {
				config.OpenAIAPIKey = v
			}
		case "base_url":
			if v, ok := value.(string); ok {
				config.BaseURL = v
			}
		case "default_output_path":
			if v, ok := value.(string); ok {
				config.DefaultOutputPath = v
			}
		case "language":
			if v, ok := value.(string); ok {
				config.Language = v
			}
		case "initialized":
			if v, ok := value.(bool); ok {
				config.Initialized = v
			}
		}
	}
	
	return s.SetConfig(config)
}

// getConfigField is a generic method to get a config field with fallback
func (s *Service) getConfigField(getter func(*types.Config) string, defaultValue string) (string, error) {
	config, err := s.GetConfig()
	if err != nil {
		return "", err
	}
	
	value := getter(config)
	if value == "" {
		return defaultValue, nil
	}
	
	return value, nil
}

// setConfigField is a generic method to set a config field
func (s *Service) setConfigField(key string, value interface{}) error {
	return s.UpdateConfig(map[string]interface{}{
		key: value,
	})
}

// SetConfigField is a generic method to set a config field
func (s *Service) SetConfigField(key string, value interface{}) error {
	return s.setConfigField(key, value)
}

// GetAPIKey returns the OpenAI API key
func (s *Service) GetAPIKey() (string, error) {
	return s.getConfigField(func(c *types.Config) string { return c.OpenAIAPIKey }, "")
}

// GetBaseURL returns the base URL
func (s *Service) GetBaseURL() (string, error) {
	return s.getConfigField(func(c *types.Config) string { return c.BaseURL }, types.DefaultValues.BaseURL)
}

// GetDefaultOutputPath returns the default output path
func (s *Service) GetDefaultOutputPath() (string, error) {
	return s.getConfigField(func(c *types.Config) string { return c.DefaultOutputPath }, types.DefaultValues.OutputPath)
}

// GetLanguage returns the configured language
func (s *Service) GetLanguage() (string, error) {
	return s.getConfigField(func(c *types.Config) string { return c.Language }, types.DefaultValues.Language)
}

// IsInitialized returns whether the configuration has been initialized
func (s *Service) IsInitialized() (bool, error) {
	config, err := s.GetConfig()
	if err != nil {
		return false, err
	}
	return config.Initialized, nil
}

// SetInitialized marks the configuration as initialized
func (s *Service) SetInitialized() error {
	return s.setConfigField("initialized", true)
}

// ResetConfig resets the configuration to default values
func (s *Service) ResetConfig() error {
	// Remove existing config file first
	if err := os.Remove(s.configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing config file: %w", err)
	}

	defaultConfig := map[string]interface{}{
		"language":             types.DefaultValues.Language,
		"base_url":            types.DefaultValues.BaseURL,
		"default_output_path": types.DefaultValues.OutputPath,
		"initialized":         false,
	}

	return s.UpdateConfig(defaultConfig)
}

// SetAPIKey sets the OpenAI API key
func (s *Service) SetAPIKey(apiKey string) error {
	return s.setConfigField("openai_api_key", apiKey)
}

// SetDefaultOutputPath sets the default output path
func (s *Service) SetDefaultOutputPath(path string) error {
	return s.setConfigField("default_output_path", path)
}

// SetLanguage sets the interface language
func (s *Service) SetLanguage(language string) error {
	return s.setConfigField("language", language)
}

// GetConfigPath returns the path to the config file
func (s *Service) GetConfigPath() string {
	return s.configPath
}

// ValidateAPIKey validates the format of an API key
func ValidateAPIKey(apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	return nil
}

// DefaultService returns a default configuration service instance
var DefaultService = NewService()
