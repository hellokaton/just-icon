package utils

import (
	"fmt"

	"github.com/pterm/pterm"
)

// Color functions for consistent output styling using pterm
var (
	// Success colors
	Green = pterm.FgGreen.Sprint
	Bold  = pterm.Bold.Sprint

	// Error colors
	Red = pterm.FgRed.Sprint

	// Info colors
	Blue = pterm.FgBlue.Sprint
	Cyan = pterm.FgCyan.Sprint

	// Warning colors
	Yellow = pterm.FgYellow.Sprint

	// Dim colors
	Gray = pterm.FgGray.Sprint
	Dim  = pterm.FgGray.Sprint
)

// PrintSuccess prints a success message with green checkmark
func PrintSuccess(message string) {
	pterm.Printf("%s %s\n", Green("✅"), message)
}

// PrintError prints an error message with red X
func PrintError(message string) {
	pterm.Printf("%s %s\n", Red("❌"), message)
}

// PrintInfo prints an info message with blue info icon
func PrintInfo(message string) {
	pterm.Printf("%s %s\n", Blue("ℹ️"), message)
}

// PrintWarning prints a warning message with yellow warning icon
func PrintWarning(message string) {
	pterm.Printf("%s %s\n", Yellow("⚠️"), message)
}

// PrintHeader prints a header with emoji using pterm's enhanced styling
func PrintHeader(message string) {
	// Create a styled header with pterm
	headerStyle := pterm.NewStyle(pterm.FgLightCyan, pterm.Bold)
	pterm.Printf("%s %s\n", "🎨", headerStyle.Sprint(message))
}

// PrintSubHeader prints a sub-header with enhanced styling
func PrintSubHeader(message string) {
	subHeaderStyle := pterm.NewStyle(pterm.FgLightGreen, pterm.Bold)
	pterm.Printf("\n%s\n", subHeaderStyle.Sprint(message))
}

// PrintKeyValue prints a key-value pair with consistent formatting
func PrintKeyValue(key, value string) {
	pterm.Printf("%s %s\n", Bold(key+":"), value)
}

// PrintDim prints dimmed text
func PrintDim(message string) {
	pterm.Println(Dim(message))
}

// MaskAPIKey masks an API key for display, showing only the last 4 characters
func MaskAPIKey(apiKey string) string {
	if len(apiKey) <= 4 {
		return "****"
	}
	return fmt.Sprintf("sk-...%s", apiKey[len(apiKey)-4:])
}
