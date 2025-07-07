package interactive

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/pterm/pterm"
	"github.com/sveltinio/prompti/choose"
	"github.com/sveltinio/prompti/input"

	"just-icon/internal/config"
	"just-icon/internal/i18n"
	"just-icon/internal/openai"
	"just-icon/internal/types"
	"just-icon/pkg/utils"
)

var ErrEmptyPrompt = errors.New("empty")
var ErrPlaceholderPrompt = errors.New("placeholder")
var ErrUserQuit = errors.New("user quit")

// validatePrompt validates that the prompt is not empty and not the placeholder text
func validatePrompt(prompt string) error {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return ErrEmptyPrompt
	}
	if prompt == types.DefaultPromptPlaceholder {
		return ErrPlaceholderPrompt
	}
	return nil
}

// loadLanguageSetting loads and applies the language setting from config
func loadLanguageSetting(configService *config.Service) error {
	config, err := configService.GetConfig()
	if err != nil {
		return err
	}

	language := config.Language
	if language == "" {
		language = "en" // Default to English
	}

	// Apply language setting
	var lang i18n.Language
	if language == "zh" {
		lang = i18n.Chinese
	} else {
		lang = i18n.English
	}
	i18n.GetLocalizer().SetLanguage(lang)

	return nil
}

// RunInteractiveMode starts the interactive mode for icon generation
func RunInteractiveMode() error {
	configService := config.DefaultService

	// Load and apply language setting before starting
	if err := loadLanguageSetting(configService); err != nil {
		// If error loading language, continue with default
	}

	// Check if this is first run
	if IsFirstRun() {
		if err := SetupWizard(); err != nil {
			if errors.Is(err, ErrUserQuit) {
				// User pressed Ctrl+C during setup, exit silently
				os.Exit(0)
			}
			return err
		}
		// Reload language setting after setup
		if err := loadLanguageSetting(configService); err != nil {
			// If error loading language, continue with default
		}
	}

	// Check if API key is configured before starting main loop
	apiKey, err := configService.GetAPIKey()
	if err != nil {
		return err
	}
	if apiKey == "" {
		fmt.Printf("❌ %s\n", i18n.T("interactive_api_key_required"))
		fmt.Printf("💡 %s: just-icon config --api-key YOUR_KEY\n", i18n.T("interactive_api_key_set_hint"))
		return nil
	}

	// Check if output directory is configured
	outputDir, err := configService.GetDefaultOutputPath()
	if err != nil {
		return err
	}
	if outputDir == "" {
		fmt.Printf("❌ %s\n", i18n.T("interactive_output_dir_required"))
		fmt.Printf("💡 %s: just-icon config --output-path YOUR_PATH\n", i18n.T("interactive_output_dir_set_hint"))
		return nil
	}

	// Main interactive loop
	for {

		// Get icon prompt from user
		prompt_text, err := getIconPrompt()
		if err != nil {
			if errors.Is(err, ErrUserQuit) {
				// Exit silently without goodbye message when user presses Ctrl+C
				return ErrUserQuit
			}
			if errors.Is(err, ErrEmptyPrompt) {
				fmt.Printf("❌ %s\n", i18n.T("validation_prompt_empty"))
				continue
			}
			if errors.Is(err, ErrPlaceholderPrompt) {
				fmt.Printf("❌ %s\n", i18n.T("validation_prompt_placeholder"))
				continue
			}
			return err
		}

		// Get quantity selection
		quantity, err := getQuantitySelection()
		if err != nil {
			if errors.Is(err, ErrUserQuit) {
				// Exit silently without goodbye message when user presses Ctrl+C
				return ErrUserQuit
			}
			return err
		}

		// Get quality selection
		quality_level, err := getQualitySelection()
		if err != nil {
			if errors.Is(err, ErrUserQuit) {
				// Exit silently without goodbye message when user presses Ctrl+C
				return ErrUserQuit
			}
			return err
		}

		// Get output directory from config
		config, err := configService.GetConfig()
		if err != nil {
			return err
		}

		outputDir := config.DefaultOutputPath
		if outputDir == "" {
			outputDir = types.DefaultValues.OutputPath
		}

		// Generate icon
		if err := generateIcon(prompt_text, quantity, quality_level, outputDir); err != nil {
			utils.PrintError(fmt.Sprintf("Failed to generate icon: %v", err))
			continue
		}

		// Ask if user wants to generate another icon
		if !askForAnother() {
			// Exit silently when user chooses not to continue
			break
		}
	}

	return nil
}

// getIconPrompt gets the icon description from user
func getIconPrompt() (string, error) {
	cfg := &input.Config{
		Message:      i18n.T("interactive_prompt_input"),
		Placeholder:  types.DefaultPromptPlaceholder,
		ValidateFunc: validatePrompt,
		ShowResult:   false,
		Styles:       input.DefaultStyles(),
	}

	result, err := input.Run(cfg)
	if err != nil {
		// Check for user quit (Ctrl+C)
		if strings.Contains(err.Error(), "interrupt") ||
			strings.Contains(err.Error(), "cancelled") ||
			strings.Contains(err.Error(), "user quit") ||
			err.Error() == "" {
			return "", ErrUserQuit
		}
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// getQuantitySelection gets the quantity selection from user
func getQuantitySelection() (int, error) {
	cfg := &input.Config{
		Message:      i18n.T("interactive_quantity_prompt"),
		Initial:      "1",
		ValidateFunc: input.ValidateInteger,
		ShowResult:   false,
		Styles:       input.DefaultStyles(),
	}

	result, err := input.Run(cfg)
	if err != nil {
		// Check for user quit (Ctrl+C)
		if strings.Contains(err.Error(), "interrupt") ||
			strings.Contains(err.Error(), "cancelled") ||
			strings.Contains(err.Error(), "user quit") ||
			err.Error() == "" {
			return 0, ErrUserQuit
		}
		return 0, err
	}

	// Convert string to int
	quantity := 1 // default value
	if result != "" {
		if q, parseErr := strconv.Atoi(result); parseErr == nil && q > 0 {
			quantity = q
		}
	}

	return quantity, nil
}

// getQualitySelection gets the quality selection from user
func getQualitySelection() (string, error) {
	cfg := &choose.Config{
		Title:    i18n.T("interactive_quality_prompt"),
		ErrorMsg: i18n.T("interactive_quality_error"),
	}

	entries := []list.Item{
		choose.Item{Name: i18n.T("interactive_quality_auto"), Desc: i18n.T("interactive_quality_auto")},
		choose.Item{Name: i18n.T("interactive_quality_high"), Desc: i18n.T("interactive_quality_high")},
		choose.Item{Name: i18n.T("interactive_quality_medium"), Desc: i18n.T("interactive_quality_medium")},
		choose.Item{Name: i18n.T("interactive_quality_low"), Desc: i18n.T("interactive_quality_low")},
	}

	result, err := choose.Run(cfg, entries)
	if err != nil {
		// Check for user quit (Ctrl+C)
		if strings.Contains(err.Error(), "interrupt") ||
			strings.Contains(err.Error(), "cancelled") ||
			strings.Contains(err.Error(), "user quit") ||
			err.Error() == "" {
			return "", ErrUserQuit
		}
		return "", err
	}

	// Map display text to actual quality values
	switch result {
	case i18n.T("interactive_quality_auto"):
		return types.QualityAuto, nil
	case i18n.T("interactive_quality_high"):
		return types.QualityHigh, nil
	case i18n.T("interactive_quality_medium"):
		return types.QualityMedium, nil
	case i18n.T("interactive_quality_low"):
		return types.QualityLow, nil
	default:
		return types.QualityAuto, nil
	}
}

// generateIcon generates the icon with given parameters
func generateIcon(prompt_text string, numImages int, quality, outputDir string) error {
	// Build options
	options := &types.IconGenerationOptions{
		Prompt:       prompt_text,
		Output:       outputDir,
		Model:        types.ModelGPTImage1, // Fixed model
		Size:         types.DefaultValues.Size,
		Quality:      quality,
		NumImages:    numImages,
		Background:   types.DefaultValues.Background,
		OutputFormat: types.DefaultValues.OutputFormat,
		Moderation:   types.DefaultValues.Moderation,
		RawPrompt:    false,
	}

	// Create OpenAI client
	client, err := openai.NewClientFromConfig()
	if err != nil {
		return err
	}

	// Show generation info
	fmt.Println()
	fmt.Printf("%s %s\n", i18n.T("prompt_label"), utils.Bold(options.Prompt))
	fmt.Printf("%s %s\n", i18n.T("model_label"), utils.Blue(types.ModelGPTImage1))
	fmt.Printf("%s %s\n", i18n.T("size_label"), utils.Cyan(options.Size))
	fmt.Printf("%s %s\n", i18n.T("quality_label"), utils.Cyan(options.Quality))
	fmt.Printf("%s %s\n", i18n.T("quantity_label"), utils.Cyan(fmt.Sprintf("%d", options.NumImages)))
	fmt.Println()

	// Create and start spinner for API call
	spinner, _ := pterm.DefaultSpinner.Start(i18n.T("interactive_generating_spinner"))

	// Generate icons with spinner
	var imageBase64Data []string
	var genErr error

	// Run the API call in a goroutine to allow spinner animation
	done := make(chan bool)
	go func() {
		imageBase64Data, genErr = client.GenerateIcon(options)
		done <- true
	}()

	// Wait for completion
	<-done

	// Stop spinner
	if genErr != nil {
		spinner.Fail("Generation failed")
		return genErr
	} else {
		spinner.Success(i18n.T("interactive_generating_success"))
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(options.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save images
	var savedFiles []string
	for i, base64Data := range imageBase64Data {
		// Use timestamp format for all files: icon_YYYYMMDDHHMMSS + 3-digit random number
		filename := utils.GenerateTimestampFileName(options.OutputFormat)
		filePath := filepath.Join(options.Output, filename)

		if err := utils.SaveBase64Image(base64Data, filePath); err != nil {
			// Generate alternative filename with timestamp and random number
			errorFilename := utils.GenerateTimestampFileName(options.OutputFormat)
			errorFilePath := filepath.Join(options.Output, errorFilename)

			// Try to save with error filename
			if saveErr := utils.SaveBase64Image(base64Data, errorFilePath); saveErr != nil {
				utils.PrintWarning(fmt.Sprintf("Failed to save image %d (both normal and error filename): %v", i+1, err))
				continue
			}

			utils.PrintWarning(fmt.Sprintf("Failed to save image %d with normal filename, saved as error file: %s", i+1, errorFilename))
			savedFiles = append(savedFiles, errorFilePath)
			continue
		}

		savedFiles = append(savedFiles, filePath)
		fmt.Println()
		utils.PrintSuccess(fmt.Sprintf("Saved: %s", filePath))
	}

	if len(savedFiles) == 0 {
		return fmt.Errorf("no images were saved successfully")
	}

	// Show summary
	fmt.Println()
	fmt.Printf(i18n.T("icon_generation_summary")+"\n",
		utils.Green(fmt.Sprintf("%d", len(savedFiles))),
		utils.Blue(options.Output))
	fmt.Println()
	return nil
}

// askForAnother asks if user wants to generate another icon
func askForAnother() bool {
	cfg := &choose.Config{
		Title:    i18n.T("interactive_another"),
		ErrorMsg: i18n.T("interactive_another_error"),
	}

	entries := []list.Item{
		choose.Item{Name: i18n.T("interactive_yes"), Desc: i18n.T("interactive_yes")},
		choose.Item{Name: i18n.T("interactive_no"), Desc: i18n.T("interactive_no")},
	}

	result, err := choose.Run(cfg, entries)
	if err != nil {
		// If user presses Ctrl+C or any error occurs, return false to exit
		return false
	}

	return result == i18n.T("interactive_yes")
}
