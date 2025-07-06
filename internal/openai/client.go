package openai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sashabaranov/go-openai"

	"just-icon/internal/config"
	"just-icon/internal/types"
)

// Client handles OpenAI API interactions
type Client struct {
	client *openai.Client
}

// NewClient creates a new OpenAI client with custom base URL
func NewClient(apiKey, baseURL string) *Client {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL + "/v1"
	} else {
		config.BaseURL = "https://api.katonai.dev/v1"
	}

	return &Client{
		client: openai.NewClientWithConfig(config),
	}
}

// NewClientFromConfig creates a new OpenAI client using configuration
func NewClientFromConfig() (*Client, error) {
	configService := config.DefaultService
	apiKey, err := configService.GetAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get API key from config: %w", err)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key not configured. Get your key at https://api.katonai.dev and run: just-icon config --api-key YOUR_KEY")
	}

	baseURL, err := configService.GetBaseURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get base URL from config: %w", err)
	}

	return NewClient(apiKey, baseURL), nil
}

// GenerateIcon generates icons using OpenAI API and returns base64 encoded image data
func (c *Client) GenerateIcon(options *types.IconGenerationOptions) ([]string, error) {
	// Validate parameters
	if err := c.validateParameters(options); err != nil {
		return nil, err
	}

	// Build request
	request := c.buildRequest(options)

	// Log request details
	if err := c.logRequest(request); err != nil {
		fmt.Printf("Warning: Failed to log request: %v\n", err)
	}

	// Make API call
	ctx := context.Background()
	response, err := c.client.CreateImage(ctx, request)

	// Log response details (including errors)
	if err := c.logResponse(response, err); err != nil {
		fmt.Printf("Warning: Failed to log response: %v\n", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	// Extract base64 image data
	var imageBase64Data []string
	for _, data := range response.Data {
		if data.B64JSON != "" {
			imageBase64Data = append(imageBase64Data, data.B64JSON)
		} else if data.URL != "" {
			// Download image from URL and convert to base64
			base64Data, err := c.downloadImageAsBase64(data.URL)
			if err != nil {
				return nil, fmt.Errorf("failed to download image from URL %s: %w", data.URL, err)
			}
			imageBase64Data = append(imageBase64Data, base64Data)
		}
	}

	if len(imageBase64Data) == 0 {
		return nil, fmt.Errorf("no images generated")
	}

	return imageBase64Data, nil
}

// downloadImageAsBase64 downloads an image from URL and returns it as base64 encoded string
func (c *Client) downloadImageAsBase64(url string) (string, error) {
	// Create HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: HTTP %d", resp.StatusCode)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Convert to base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	return base64Data, nil
}

// validateParameters validates the generation options
func (c *Client) validateParameters(options *types.IconGenerationOptions) error {
	model := options.Model
	if model == "" {
		model = types.DefaultValues.Model
	}

	// Check if model is supported
	found := false
	for _, supportedModel := range types.SupportedModels {
		if model == supportedModel {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unsupported model: %s. Supported models: %v", model, types.SupportedModels)
	}

	// Validate size
	size := options.Size
	if size == "" {
		size = types.DefaultValues.Size
	}

	validSizes, exists := types.SupportedSizes[model]
	if !exists {
		return fmt.Errorf("no size configuration for model: %s", model)
	}

	found = false
	for _, validSize := range validSizes {
		if size == validSize {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("invalid size %s for model %s. Valid sizes: %v", size, model, validSizes)
	}

	// Validate quality
	quality := options.Quality
	if quality == "" {
		quality = types.DefaultValues.Quality
	}

	validQualities, exists := types.SupportedQualities[model]
	if !exists {
		return fmt.Errorf("no quality configuration for model: %s", model)
	}

	found = false
	for _, validQuality := range validQualities {
		if quality == validQuality {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("invalid quality %s for model %s. Valid qualities: %v", quality, model, validQualities)
	}

	// Validate number of images
	numImages := options.NumImages
	if numImages == 0 {
		numImages = types.DefaultValues.NumImages
	}

	if numImages < 1 || numImages > 10 {
		return fmt.Errorf("number of images must be between 1 and 10")
	}

	return nil
}

// buildRequest builds the API request
func (c *Client) buildRequest(options *types.IconGenerationOptions) openai.ImageRequest {
	model := options.Model
	if model == "" {
		model = types.DefaultValues.Model
	}

	size := options.Size
	if size == "" {
		size = types.DefaultValues.Size
	}

	quality := options.Quality
	if quality == "" {
		quality = types.DefaultValues.Quality
	}

	background := options.Background
	if background == "" {
		background = types.DefaultValues.Background
	}

	outputFormat := options.OutputFormat
	if outputFormat == "" {
		outputFormat = types.DefaultValues.OutputFormat
	}

	numImages := options.NumImages
	if numImages == 0 {
		numImages = types.DefaultValues.NumImages
	}

	// Build enhanced prompt for iOS app icons (unless raw prompt is requested)
	prompt := options.Prompt
	if !options.RawPrompt {
		prompt = fmt.Sprintf(
			"Create a full-bleed %s px iOS app icon: %s. Use crisp, minimal design with vibrant colors. Add a subtle inner bevel for gentle depth; no hard shadows or outlines. Center the design with comfortable breathing room from the edges. Solid, light-neutral background. IMPORTANT: Fill the entire canvas edge-to-edge with the design, no padding, no margins. Design elements should be centered with appropriate spacing from edges but the background must cover 100%% of the canvas. Add subtle depth with inner highlights, avoid hard shadows. Clean, minimal, Apple-style design. No borders, frames, or rounded corners.",
			size, prompt,
		)
	}

	// Map background values to OpenAI constants
	var backgroundValue string
	switch background {
	case "transparent":
		backgroundValue = "transparent"
	case "opaque":
		backgroundValue = "opaque"
	default: // "auto"
		backgroundValue = "auto" // Use auto as default
	}

	// Map output format values to OpenAI constants
	var outputFormatValue string
	switch outputFormat {
	case "jpeg":
		outputFormatValue = "jpeg"
	case "webp":
		outputFormatValue = "webp"
	default: // "png"
		outputFormatValue = "png"
	}

	// Map quality values to OpenAI constants
	var qualityValue string
	switch quality {
	case "high":
		qualityValue = "high"
	case "medium":
		qualityValue = "medium"
	case "low":
		qualityValue = "low"
	default: // "auto"
		qualityValue = "low" // Default to low for cost efficiency
	}

	request := openai.ImageRequest{
		Prompt:       prompt,
		Background:   backgroundValue,
		Model:        "gpt-image-1",
		Size:         size,
		N:            numImages,
		Quality:      qualityValue,
		OutputFormat: outputFormatValue,
	}

	return request
}

// logRequest logs the API request details to a file
func (c *Client) logRequest(request openai.ImageRequest) error {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(logDir, fmt.Sprintf("openai_request_%s.json", timestamp))

	// Create a sanitized version of the request for logging (remove sensitive data)
	logRequest := map[string]interface{}{
		"model":         request.Model,
		"prompt":        request.Prompt,
		"size":          request.Size,
		"quality":       request.Quality,
		"n":             request.N,
		"background":    request.Background,
		"output_format": request.OutputFormat,
	}

	data, err := json.MarshalIndent(logRequest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// logResponse logs the API response details to a file
func (c *Client) logResponse(response openai.ImageResponse, apiErr error) error {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(logDir, fmt.Sprintf("openai_response_%s.json", timestamp))

	logResponse := map[string]interface{}{}

	if apiErr != nil {
		logResponse["error"] = apiErr.Error()
		logResponse["success"] = false
	} else {
		logResponse["success"] = true
		logResponse["created"] = response.Created

		// Log data info without the actual base64 content (too large)
		var dataInfo []map[string]interface{}
		for i, data := range response.Data {
			info := map[string]interface{}{
				"index": i,
				"has_b64_json": data.B64JSON != "",
				"has_url": data.URL != "",
			}
			if data.B64JSON != "" {
				info["b64_json_length"] = len(data.B64JSON)
			}
			if data.URL != "" {
				info["url"] = data.URL
			}
			dataInfo = append(dataInfo, info)
		}
		logResponse["data"] = dataInfo
		logResponse["data_count"] = len(response.Data)
	}

	data, err := json.MarshalIndent(logResponse, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}