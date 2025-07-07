package types

// Constants for the application
const (
	// Model constants
	ModelGPTImage1 = "gpt-image-1"

	// Quality constants
	QualityAuto   = "auto"
	QualityHigh   = "high"
	QualityMedium = "medium"
	QualityLow    = "low"
	
	// Size constants
	SizeSmall  = "1024x1024"
	SizeMedium = "1536x1024"
	SizeLarge  = "1024x1536"
	
	// API constants
	DefaultBaseURL = "https://api.katonai.dev"
	APIVersion     = "/v1"
	
	// File constants
	ConfigDirPerm  = 0755
	ConfigFilePerm = 0600
	
	// Validation constants
	MinImages = 1
	MaxImages = 10

	// API Key constants
	APIKeyPrefix = "sk-"
	APIKeyMaskLength = 4

	// UI constants
	DefaultPromptPlaceholder = "e.g., minimalist weather app with sun and cloud"
	
	// Error message constants
	ErrorInterrupt = "interrupt"
	ErrorCancelled = "cancelled"
	ErrorUserQuit  = "user quit"
)

// Config represents the application configuration
type Config struct {
	OpenAIAPIKey       string `json:"openai_api_key,omitempty"`
	BaseURL            string `json:"base_url,omitempty"`
	DefaultOutputPath  string `json:"default_output_path,omitempty"`
	Language           string `json:"language,omitempty"`
	Initialized        bool   `json:"initialized"`
}

// IconGenerationOptions represents options for icon generation
type IconGenerationOptions struct {
	Prompt       string `json:"prompt"`
	Output       string `json:"output,omitempty"`
	Size         string `json:"size,omitempty"`
	Quality      string `json:"quality,omitempty"`
	Background   string `json:"background,omitempty"`
	OutputFormat string `json:"output_format,omitempty"`
	Model        string `json:"model,omitempty"`
	NumImages    int    `json:"num_images,omitempty"`
	Moderation   string `json:"moderation,omitempty"`
	RawPrompt    bool   `json:"raw_prompt,omitempty"`
}

// OpenAIImageResponse represents the response from OpenAI image generation API
type OpenAIImageResponse struct {
	Data []OpenAIImageData `json:"data"`
}

// OpenAIImageData represents individual image data from OpenAI response
type OpenAIImageData struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// OpenAIImageRequest represents the request to OpenAI image generation API
type OpenAIImageRequest struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model,omitempty"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Style          string `json:"style,omitempty"`
}

// SupportedModels contains the list of supported AI models
var SupportedModels = []string{
	"gpt-image-1",
}

// SupportedSizes contains the mapping of models to their supported sizes
var SupportedSizes = map[string][]string{
	ModelGPTImage1: {SizeSmall, SizeMedium, SizeLarge},
}

// SupportedQualities contains the mapping of models to their supported qualities
var SupportedQualities = map[string][]string{
	ModelGPTImage1: {QualityAuto, QualityHigh, QualityMedium, QualityLow},
}

// SupportedBackgrounds contains the mapping of models to their supported background types
var SupportedBackgrounds = map[string][]string{
	ModelGPTImage1: {"auto", "transparent", "opaque"},
}

// SupportedOutputFormats contains the mapping of models to their supported output formats
var SupportedOutputFormats = map[string][]string{
	ModelGPTImage1: {"png", "jpeg", "webp"},
}

// DefaultValues contains default configuration values
var DefaultValues = struct {
	Model        string
	Size         string
	Quality      string
	OutputPath   string
	NumImages    int
	Background   string
	OutputFormat string
	Moderation   string
	Language     string
	BaseURL      string
}{
	Model:        ModelGPTImage1,
	Size:         SizeSmall,
	Quality:      QualityAuto,
	OutputPath:   "./output",
	NumImages:    MinImages,
	Background:   "auto",
	OutputFormat: "png",
	Moderation:   "auto",
	Language:     "en",
	BaseURL:      DefaultBaseURL,
}
