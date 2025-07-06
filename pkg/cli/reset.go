package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"just-icon/internal/config"
	"just-icon/internal/i18n"
	"just-icon/pkg/utils"
)

// NewResetCommand creates the reset command
func NewResetCommand() *cli.Command {
	return &cli.Command{
		Name:  "reset",
		Usage: "Reset configuration to default values",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Usage:   "Force reset without confirmation",
				Aliases: []string{"f"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configService := config.DefaultService

			// Load current language setting before showing messages
			config, err := configService.GetConfig()
			if err == nil && config.Language != "" {
				var lang i18n.Language
				if config.Language == "zh" {
					lang = i18n.Chinese
				} else {
					lang = i18n.English
				}
				i18n.GetLocalizer().SetLanguage(lang)
			}

			// Check if force flag is provided
			force := cmd.Bool("force")

			if !force {
				// Ask for confirmation using internationalized text
				fmt.Printf("⚠️  %s\n", i18n.T("reset_warning"))
				fmt.Printf("%s (y/N): ", i18n.T("reset_confirm"))

				var response string
				fmt.Scanln(&response)

				if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
					fmt.Println(i18n.T("reset_cancelled"))
					return nil
				}
			}

			// Reset configuration
			if err := configService.ResetConfig(); err != nil {
				utils.PrintError(i18n.Tf("reset_failed", err))
				return err
			}

			utils.PrintSuccess(i18n.T("reset_success"))
			fmt.Printf("%s\n", i18n.T("reset_next_run"))

			return nil
		},
	}
}
