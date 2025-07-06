package utils

import (
	"regexp"
	"strings"
	"testing"
)

func TestGenerateTimestampFileName(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{
			name:   "png format",
			format: "png",
		},
		{
			name:   "jpeg format",
			format: "jpeg",
		},
		{
			name:   "webp format",
			format: "webp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := GenerateTimestampFileName(tt.format)
			
			// Check that filename has correct format: icon_YYYYMMDDHHMMSS###.format
			// Pattern: icon_ + 14 digits (YYYYMMDDHHMMSS) + 3 digits (random) + .format
			expectedPattern := `^icon_\d{14}\d{3}\.\w+$`
			matched, err := regexp.MatchString(expectedPattern, filename)
			if err != nil {
				t.Fatalf("Failed to compile regex: %v", err)
			}
			if !matched {
				t.Errorf("GenerateTimestampFileName(%s) = %s, doesn't match expected pattern %s", tt.format, filename, expectedPattern)
			}

			// Check that it ends with the correct format
			if !strings.HasSuffix(filename, "."+tt.format) {
				t.Errorf("GenerateTimestampFileName(%s) = %s, doesn't end with .%s", tt.format, filename, tt.format)
			}

			// Check that it starts with "icon_"
			if !strings.HasPrefix(filename, "icon_") {
				t.Errorf("GenerateTimestampFileName(%s) = %s, doesn't start with 'icon_'", tt.format, filename)
			}

			// Check total length: icon_ (5) + timestamp (14) + random (3) + . (1) + format (len)
			expectedLength := 5 + 14 + 3 + 1 + len(tt.format)
			if len(filename) != expectedLength {
				t.Errorf("GenerateTimestampFileName(%s) = %s, length %d, expected %d", tt.format, filename, len(filename), expectedLength)
			}
		})
	}
}

func TestGenerateTimestampFileNameUniqueness(t *testing.T) {
	// Generate multiple filenames and check they are different
	// (due to timestamp and random number)
	format := "png"
	filenames := make(map[string]bool)

	for i := 0; i < 10; i++ {
		filename := GenerateTimestampFileName(format)
		if filenames[filename] {
			t.Errorf("GenerateTimestampFileName generated duplicate filename: %s", filename)
		}
		filenames[filename] = true
	}
}

func TestGenerateFileNameWithFormat(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		format string
	}{
		{
			name:   "icon prefix with png",
			prefix: "icon",
			format: "png",
		},
		{
			name:   "test prefix with jpeg",
			prefix: "test",
			format: "jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := GenerateFileNameWithFormat(tt.prefix, tt.format)
			
			// Check that it starts with prefix
			if !strings.HasPrefix(filename, tt.prefix+"-") {
				t.Errorf("GenerateFileNameWithFormat(%s, %s) = %s, doesn't start with '%s-'", tt.prefix, tt.format, filename, tt.prefix)
			}
			
			// Check that it ends with format
			if !strings.HasSuffix(filename, "."+tt.format) {
				t.Errorf("GenerateFileNameWithFormat(%s, %s) = %s, doesn't end with '.%s'", tt.prefix, tt.format, filename, tt.format)
			}
		})
	}
}
