package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"just-icon/internal/i18n"
	justcli "just-icon/pkg/cli"
	"just-icon/pkg/interactive"
	"just-icon/pkg/banner"
)

func main() {
	// Initialize default language (English)
	i18n.InitLocalizer(i18n.English)

	// Show ASCII art banner
	banner.ShowBanner()

	cmd := &cli.Command{
		Name:        i18n.T("app_name"),
		Usage:       i18n.T("app_usage"),
		Description: i18n.T("app_description"),
		Version:     "1.0.0",
		Commands: []*cli.Command{
			justcli.NewConfigCommand(),
			justcli.NewResetCommand(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Check for language argument
			args := cmd.Args()
			if args.Len() > 0 {
				lang := args.Get(0)
				i18n.SwitchLanguage(lang)
			}

			// Run interactive mode
			err := interactive.RunInteractiveMode()
			if err != nil {
				// Check if it's a user quit error, exit silently
				if errors.Is(err, interactive.ErrUserQuit) {
					os.Exit(0)
				}
			}
			return err
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// Handle language switching for subcommands
			args := cmd.Args()
			if args.Len() > 0 {
				firstArg := args.Get(0)
				i18n.SwitchLanguage(firstArg)
			}
			return ctx, nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		// Check if it's a user quit error, exit silently
		if errors.Is(err, interactive.ErrUserQuit) {
			os.Exit(0)
		}
		// For other errors, show the error message
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
