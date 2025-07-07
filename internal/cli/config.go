package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"just-icon/internal/config"
	"just-icon/internal/i18n"
	"just-icon/internal/types"
	"just-icon/pkg/utils"
)

// NewConfigCommand creates the config command
func NewConfigCommand() *cli.Command {
	return &cli.Command{
		Name:        i18n.T("config_name"),
		Usage:       i18n.T("config_usage"),
		Description: i18n.T("config_description"),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "api-key",
				Usage:   i18n.T("config_flag_api_key"),
				Aliases: []string{"k"},
			},
			&cli.StringFlag{
				Name:    "base-url",
				Usage:   i18n.T("config_flag_base_url"),
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    "output-path",
				Usage:   i18n.T("config_flag_output_path"),
				Aliases: []string{"o"},
			},
			&cli.StringFlag{
				Name:    "language",
				Usage:   i18n.T("config_flag_language"),
				Aliases: []string{"l"},
			},
			&cli.BoolFlag{
				Name:    "show",
				Usage:   i18n.T("config_flag_show"),
				Aliases: []string{"s"},
			},
		},
		Action: configAction,
	}
}

func configAction(ctx context.Context, cmd *cli.Command) error {
	configService := config.DefaultService

	// Handle API key setting
	if apiKey := cmd.String("api-key"); apiKey != "" {
		if err := setAPIKey(configService, apiKey); err != nil {
			return err
		}
	}

	// Handle base URL setting
	if baseURL := cmd.String("base-url"); baseURL != "" {
		if err := setBaseURL(configService, baseURL); err != nil {
			return err
		}
	}

	// Handle output path setting
	if outputPath := cmd.String("output-path"); outputPath != "" {
		if err := setOutputPath(configService, outputPath); err != nil {
			return err
		}
	}

	// Handle language setting
	if language := cmd.String("language"); language != "" {
		if err := setLanguage(configService, language); err != nil {
			return err
		}
	}

	// Handle show configuration
	if cmd.Bool("show") {
		return showConfig(configService)
	}

	// If no flags provided, show configuration by default
	if !cmd.Bool("show") && cmd.String("api-key") == "" && cmd.String("output-path") == "" && cmd.String("language") == "" {
		return showConfig(configService)
	}

	return nil
}

// handleConfigError handles configuration errors consistently
func handleConfigError(operation string, err error) error {
	utils.PrintError(i18n.Tf("config_failed_to_save", err.Error()))
	return err
}

// setConfigValue is a generic function for setting configuration values
func setConfigValue(configService *config.Service, key, value string, validator func(string) error, setter func(string) error, successMsgKey string) error {
	// Validate if validator is provided
	if validator != nil {
		if err := validator(value); err != nil {
			utils.PrintError(i18n.Tf("config_invalid_"+key, err.Error()))
			return nil // Don't return error to avoid showing usage
		}
	}

	// Save the value
	if err := setter(value); err != nil {
		return handleConfigError("save "+key, err)
	}

	// Show success message
	if successMsgKey != "" {
		utils.PrintSuccess(i18n.Tf(successMsgKey, value))
	}

	return nil
}

func setAPIKey(configService *config.Service, apiKey string) error {
	err := setConfigValue(configService, "api_key", apiKey, 
		config.ValidateAPIKey, 
		configService.SetAPIKey, 
		"config_api_key_success")
	
	if err == nil {
		fmt.Println()
		utils.PrintDim(i18n.T("common_built_with"))
	}
	
	return err
}

func setBaseURL(configService *config.Service, baseURL string) error {
	return setConfigValue(configService, "base_url", baseURL, 
		nil, 
		func(url string) error {
			return configService.SetConfigField("base_url", url)
		}, 
		"config_base_url_success")
}

func setOutputPath(configService *config.Service, outputPath string) error {
	return setConfigValue(configService, "output_path", outputPath, 
		nil, 
		configService.SetDefaultOutputPath, 
		"config_output_path_success")
}

func setLanguage(configService *config.Service, language string) error {
	// Validate language
	validateLanguage := func(lang string) error {
		if lang != "en" && lang != "zh" {
			return fmt.Errorf("unsupported language: %s. Supported languages: en, zh", lang)
		}
		return nil
	}

	// Custom setter that also updates the localizer
	setLangWithLocalizer := func(lang string) error {
		if err := configService.SetLanguage(lang); err != nil {
			return err
		}
		// Update localizer using the new SwitchLanguage function
		i18n.SwitchLanguage(lang)
		return nil
	}

	return setConfigValue(configService, "language", language, 
		validateLanguage, 
		setLangWithLocalizer, 
		"config_language_success")
}

func showConfig(configService *config.Service) error {
	config, err := configService.GetConfig()
	if err != nil {
		utils.PrintError(i18n.Tf("config_failed_to_read", err.Error()))
		return err
	}

	utils.PrintSubHeader(i18n.T("config_current_title"))
	fmt.Println()

	// Show API key status
	if config.OpenAIAPIKey != "" {
		maskedKey := utils.MaskAPIKey(config.OpenAIAPIKey)
		utils.PrintKeyValue(i18n.T("config_api_key"), utils.Green(maskedKey))
	} else {
		utils.PrintKeyValue(i18n.T("config_api_key"), utils.Red(i18n.T("config_not_configured")))
		fmt.Printf("   %s\n", utils.Gray(i18n.T("config_get_api_key")))
		fmt.Printf("   %s\n", utils.Gray(i18n.T("config_set_api_key")))
	}

	// Show base URL
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = types.DefaultValues.BaseURL
	}
	utils.PrintKeyValue(i18n.T("config_base_url"), utils.Blue(baseURL))

	// Show default output path
	outputPath, err := configService.GetDefaultOutputPath()
	if err != nil {
		utils.PrintWarning(i18n.Tf("config_failed_to_read", err.Error()))
	} else {
		utils.PrintKeyValue(i18n.T("config_default_output"), utils.Blue(outputPath))
	}

	// Show language
	language, err := configService.GetLanguage()
	if err != nil {
		utils.PrintWarning(i18n.Tf("config_failed_to_read", err.Error()))
	} else {
		utils.PrintKeyValue(i18n.T("config_language"), utils.Cyan(language))
	}

	// Show config file location
	utils.PrintKeyValue(i18n.T("config_file"), utils.Gray(configService.GetConfigPath()))

	fmt.Println()
	utils.PrintDim(i18n.T("common_built_with"))

	return nil
}
