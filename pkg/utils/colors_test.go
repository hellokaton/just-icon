package utils

import (
	"strings"
	"testing"
)

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		expected string
	}{
		{
			name:     "normal API key",
			apiKey:   "sk-1234567890123456789012345678901234567890",
			expected: "sk-...7890",
		},
		{
			name:     "short API key",
			apiKey:   "sk-123",
			expected: "sk-...-123",
		},
		{
			name:     "very short API key",
			apiKey:   "sk",
			expected: "****",
		},
		{
			name:     "empty API key",
			apiKey:   "",
			expected: "****",
		},
		{
			name:     "exactly 4 characters",
			apiKey:   "1234",
			expected: "****",
		},
		{
			name:     "5 characters",
			apiKey:   "12345",
			expected: "sk-...2345",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskAPIKey(tt.apiKey)
			if result != tt.expected {
				t.Errorf("MaskAPIKey(%s) = %s, want %s", tt.apiKey, result, tt.expected)
			}
		})
	}
}

func TestColorFunctions(t *testing.T) {
	// Test that color functions don't panic and return non-empty strings
	testText := "test"

	// Test individual color functions
	tests := []struct {
		name string
		fn   func(a ...interface{}) string
	}{
		{"Green", Green},
		{"Red", Red},
		{"Blue", Blue},
		{"Cyan", Cyan},
		{"Yellow", Yellow},
		{"Gray", Gray},
		{"Dim", Dim},
		{"Bold", Bold},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(testText)
			if result == "" {
				t.Errorf("Color function returned empty string")
			}
			// The result should contain the original text
			if !strings.Contains(result, testText) {
				t.Errorf("Color function result doesn't contain original text")
			}
		})
	}
}
