package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// DownloadFile downloads a file from URL and saves it to the specified path
func DownloadFile(url, filepath string) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	
	// Make request
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()
	
	// Copy data
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	
	return nil
}

// EnsureDir ensures that a directory exists, creating it if necessary
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// GenerateFileName generates a unique filename with timestamp
func GenerateFileName(prefix, extension string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d.%s", prefix, timestamp, extension)
}

// GetFileExtension returns the file extension from a filename
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

// FileExists checks if a file exists
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// SaveBase64Image decodes base64 image data and saves it to the specified path
func SaveBase64Image(base64Data, filepath string) error {
	// Decode base64 data
	imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Write image to file
	err = os.WriteFile(filepath, imgBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write image file: %w", err)
	}

	return nil
}

// GenerateFileNameWithFormat generates a unique filename with the specified format
func GenerateFileNameWithFormat(prefix, format string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d.%s", prefix, timestamp, format)
}

// GenerateTimestampFileName generates a filename with format: icon_2025070622012323
// Format: icon_YYYYMMDDHHMMSS + 3-digit random number
func GenerateTimestampFileName(format string) string {
	now := time.Now()
	timestamp := now.Format("20060102150405") // YYYYMMDDHHMMSS

	// Generate 3-digit random number (000-999)
	rand.Seed(now.UnixNano())
	randomNum := rand.Intn(1000)

	return fmt.Sprintf("icon_%s%03d.%s", timestamp, randomNum, format)
}
