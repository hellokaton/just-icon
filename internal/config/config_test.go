package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Fatal("NewService() returned nil")
	}
	
	if service.configPath == "" {
		t.Fatal("configPath is empty")
	}
}

func TestValidateAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid API key",
			apiKey:  "sk-1234567890123456789012345678901234567890",
			wantErr: false,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			wantErr: true,
		},
		{
			name:    "too short API key",
			apiKey:  "sk-123",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			apiKey:  "invalid-1234567890123456789012345678901234567890",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIKey(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAPIKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigOperations(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "just-icon-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a service with custom config path
	service := &Service{
		configPath: filepath.Join(tempDir, "config.json"),
	}
	
	// Test getting empty config
	config, err := service.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() failed: %v", err)
	}
	
	if config.OpenAIAPIKey != "" {
		t.Errorf("Expected empty API key, got %s", config.OpenAIAPIKey)
	}
	
	// Test setting API key
	testAPIKey := "sk-test1234567890123456789012345678901234567890"
	err = service.SetAPIKey(testAPIKey)
	if err != nil {
		t.Fatalf("SetAPIKey() failed: %v", err)
	}
	
	// Test getting API key
	apiKey, err := service.GetAPIKey()
	if err != nil {
		t.Fatalf("GetAPIKey() failed: %v", err)
	}
	
	if apiKey != testAPIKey {
		t.Errorf("Expected API key %s, got %s", testAPIKey, apiKey)
	}
	
	// Test setting output path
	testOutputPath := "/test/path"
	err = service.SetDefaultOutputPath(testOutputPath)
	if err != nil {
		t.Fatalf("SetDefaultOutputPath() failed: %v", err)
	}
	
	// Test getting output path
	outputPath, err := service.GetDefaultOutputPath()
	if err != nil {
		t.Fatalf("GetDefaultOutputPath() failed: %v", err)
	}
	
	if outputPath != testOutputPath {
		t.Errorf("Expected output path %s, got %s", testOutputPath, outputPath)
	}
	
	// Test full config
	config, err = service.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() failed: %v", err)
	}
	
	if config.OpenAIAPIKey != testAPIKey {
		t.Errorf("Expected API key %s, got %s", testAPIKey, config.OpenAIAPIKey)
	}
	
	if config.DefaultOutputPath != testOutputPath {
		t.Errorf("Expected output path %s, got %s", testOutputPath, config.DefaultOutputPath)
	}
}

func TestUpdateConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "just-icon-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a service with custom config path
	service := &Service{
		configPath: filepath.Join(tempDir, "config.json"),
	}
	
	// Test updating multiple fields
	updates := map[string]interface{}{
		"openai_api_key":       "sk-test1234567890123456789012345678901234567890",
		"default_output_path":  "/test/path",
	}
	
	err = service.UpdateConfig(updates)
	if err != nil {
		t.Fatalf("UpdateConfig() failed: %v", err)
	}
	
	// Verify updates
	config, err := service.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() failed: %v", err)
	}
	
	if config.OpenAIAPIKey != "sk-test1234567890123456789012345678901234567890" {
		t.Errorf("API key not updated correctly")
	}
	
	if config.DefaultOutputPath != "/test/path" {
		t.Errorf("Output path not updated correctly")
	}
}
