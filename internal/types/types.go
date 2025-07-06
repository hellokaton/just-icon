package types

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
	"gpt-image-1": {"1024x1024", "1536x1024", "1024x1536"},
}

// SupportedQualities contains the mapping of models to their supported qualities
var SupportedQualities = map[string][]string{
	"gpt-image-1": {"auto", "high", "medium", "low"},
}

// SupportedBackgrounds contains the mapping of models to their supported background types
var SupportedBackgrounds = map[string][]string{
	"gpt-image-1": {"auto", "transparent", "opaque"},
}

// SupportedOutputFormats contains the mapping of models to their supported output formats
var SupportedOutputFormats = map[string][]string{
	"gpt-image-1": {"png", "jpeg", "webp"},
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
	Model:        "gpt-image-1",
	Size:         "1024x1024",
	Quality:      "auto",
	OutputPath:   "./output",
	NumImages:    1,
	Background:   "auto",
	OutputFormat: "png",
	Moderation:   "auto",
	Language:     "en",
	BaseURL:      "https://api.katonai.dev",
}
