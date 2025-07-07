package interactive

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/sveltinio/prompti/choose"
	"github.com/sveltinio/prompti/input"

	"just-icon/internal/config"
	"just-icon/internal/i18n"
	"just-icon/internal/types"
	"just-icon/pkg/utils"
)

var ErrEmptyAPIKey = errors.New("API key cannot be empty")
var ErrSkipAPIKey = errors.New("user skipped API key setup")

// validateAPIKey validates that the API key is not empty and has correct format
func validateAPIKey(apiKey string) error {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return ErrEmptyAPIKey
	}

	// Use the more comprehensive validation from config package
	if err := config.ValidateAPIKey(apiKey); err != nil {
		return err
	}

	return nil
}

// SetupWizard runs the initial setup wizard for first-time users
func SetupWizard() error {
	fmt.Println()
	fmt.Printf("🎨 %s\n", "Welcome to Just Icon!")
	fmt.Println()

	configService := config.DefaultService

	// Create default config file first
	if err := createDefaultConfig(configService); err != nil {
		return err
	}

	// Step 1: Language selection
	if err := setupLanguage(configService); err != nil {
		if errors.Is(err, ErrUserQuit) {
			return ErrUserQuit
		}
		return err
	}

	// Clear screen and show welcome message in selected language
	fmt.Print("\033[2J\033[H") // Clear screen and move cursor to top
	fmt.Printf("🎨 %s\n", "Welcome to Just Icon!")
	fmt.Println("Let's continue setting up your configuration.")
	fmt.Println()

	// Step 2: Base URL configuration
	if err := setupBaseURL(configService); err != nil {
		if errors.Is(err, ErrUserQuit) {
			return ErrUserQuit
		}
		return err
	}

	// Step 2: API Key configuration
	if err := setupAPIKey(configService); err != nil {
		if errors.Is(err, ErrSkipAPIKey) {
			// User pressed Ctrl+C to quit, return special error
			return ErrUserQuit
		}
		return err
	}

	// Step 3: Output directory configuration
	if err := setupOutputDirectory(configService); err != nil {
		if errors.Is(err, ErrUserQuit) {
			return ErrUserQuit
		}
		return err
	}

	fmt.Println()
	fmt.Printf("%s\n", i18n.T("interactive_setup_complete"))
	fmt.Println()

	return nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configService *config.Service) error {
	// Create default config with English language
	defaultConfig := map[string]interface{}{
		"language":             "en",
		"base_url":            types.DefaultValues.BaseURL,
		"default_output_path": types.DefaultValues.OutputPath,
	}

	return configService.UpdateConfig(defaultConfig)
}

// setupLanguage handles language selection
func setupLanguage(configService *config.Service) error {
	cfg := &choose.Config{
		Title:    "Select language",
		ErrorMsg: "Please select a language.",
	}

	entries := []list.Item{
		choose.Item{Name: "English", Desc: "English"},
		choose.Item{Name: "中文", Desc: "中文"},
	}

	result, err := choose.Run(cfg, entries)
	if err != nil {
		// Check for user quit (Ctrl+C)
		if strings.Contains(err.Error(), "interrupt") ||
			strings.Contains(err.Error(), "cancelled") ||
			strings.Contains(err.Error(), "user quit") {
			return ErrUserQuit
		}
		return err
	}

	var language string
	switch result {
	case "English":
		language = "en"
	case "中文":
		language = "zh"
	default:
		language = "en"
	}

	// Save language setting immediately
	if err := configService.UpdateConfig(map[string]interface{}{
		"language": language,
	}); err != nil {
		return err
	}

	// Update localizer immediately to make the change take effect
	var lang i18n.Language
	if language == "zh" {
		lang = i18n.Chinese
	} else {
		lang = i18n.English
	}

	// Re-initialize the entire i18n system with the new language
	i18n.InitLocalizer(lang)

	return nil
}

// setupBaseURL handles base URL configuration
func setupBaseURL(configService *config.Service) error {
	defaultURL := types.DefaultValues.BaseURL

	cfg := &input.Config{
		Message:     i18n.T("interactive_base_url_prompt"),
		Initial:     defaultURL,
		ShowResult:  false,
		Styles:      input.DefaultStyles(),
	}

	result, err := input.Run(cfg)
	if err != nil {
		// Check for user quit (Ctrl+C)
		if strings.Contains(err.Error(), "interrupt") ||
			strings.Contains(err.Error(), "cancelled") ||
			strings.Contains(err.Error(), "user quit") {
			return ErrUserQuit
		}
		return err
	}

	if result == "" {
		result = defaultURL
	}

	// Save base URL
	if err := configService.UpdateConfig(map[string]interface{}{
		"base_url": result,
	}); err != nil {
		return err
	}

	return nil
}

// setupAPIKey handles API key configuration
func setupAPIKey(configService *config.Service) error {
	fmt.Println(utils.Gray(i18n.T("interactive_api_key_help")))

	// Loop until valid API key is entered
	for {
		cfg := &input.Config{
			Message:      i18n.T("interactive_api_key_prompt"),
			Password:     true,
			ValidateFunc: validateAPIKey,
			ShowResult:   false,
			Styles:       input.DefaultStyles(),
		}

		result, err := input.Run(cfg)
		if err != nil {
			// Check if user quit (Ctrl+C or ESC)
			if strings.Contains(err.Error(), "interrupt") ||
				strings.Contains(err.Error(), "cancelled") ||
				strings.Contains(err.Error(), "user quit") {
				return ErrSkipAPIKey
			}

			// Check for validation errors that indicate invalid API key
			if errors.Is(err, ErrEmptyAPIKey) ||
				strings.Contains(err.Error(), "API key cannot be empty") {
				fmt.Printf("❌ %s\n", err.Error())
				continue // Loop back to ask for input again
			}
			return err
		}

		// Additional validation check (in case prompt validation was bypassed)
		apiKey := strings.TrimSpace(result)
		if err := config.ValidateAPIKey(apiKey); err != nil {
			fmt.Printf("❌ %s\n", err.Error())
			continue // Loop back to ask for input again
		}

		// Save API key
		if err := configService.UpdateConfig(map[string]interface{}{
			"openai_api_key": apiKey,
		}); err != nil {
			return err
		}

		// Successfully saved, break out of loop
		break
	}

	return nil
}

// setupOutputDirectory handles output directory configuration
func setupOutputDirectory(configService *config.Service) error {
	defaultPath := types.DefaultValues.OutputPath

	cfg := &input.Config{
		Message:     i18n.T("interactive_output_dir_prompt"),
		Initial:     defaultPath,
		ShowResult:  false,
		Styles:      input.DefaultStyles(),
	}

	result, err := input.Run(cfg)
	if err != nil {
		// Check for user quit (Ctrl+C)
		if strings.Contains(err.Error(), "interrupt") ||
			strings.Contains(err.Error(), "cancelled") ||
			strings.Contains(err.Error(), "user quit") {
			return ErrUserQuit
		}
		return err
	}

	if result == "" {
		result = defaultPath
	}

	// Expand relative path to absolute path
	if !filepath.IsAbs(result) {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		result = filepath.Join(cwd, result)
	}

	// Save output directory and mark as initialized
	if err := configService.UpdateConfig(map[string]interface{}{
		"default_output_path": result,
		"initialized":         true,
	}); err != nil {
		return err
	}

	return nil
}

// IsFirstRun checks if this is the first time the user runs the application
func IsFirstRun() bool {
	configService := config.DefaultService
	initialized, err := configService.IsInitialized()
	if err != nil {
		return true
	}

	// Consider it first run if not initialized
	return !initialized
}
