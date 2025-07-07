package banner

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

// ASCII art for "JUST ICON" with gradient effect
var asciiArt = []string{
	"     ██╗██╗   ██╗███████╗████████╗    ██╗ ██████╗ ██████╗ ███╗   ██╗",
	"     ██║██║   ██║██╔════╝╚══██╔══╝    ██║██╔════╝██╔═══██╗████╗  ██║",
	"     ██║██║   ██║███████╗   ██║       ██║██║     ██║   ██║██╔██╗ ██║",
	"██   ██║██║   ██║╚════██║   ██║       ██║██║     ██║   ██║██║╚██╗██║",
	"╚█████╔╝╚██████╔╝███████║   ██║       ██║╚██████╗╚██████╔╝██║ ╚████║",
	" ╚════╝  ╚═════╝ ╚══════╝   ╚═╝       ╚═╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝",
}

// ShowBanner displays the ASCII art banner with gradient colors
func ShowBanner() {
	fmt.Println() // Add some space at the top

	// Define RGB gradient colors with blue theme
	startColor := pterm.NewRGB(100, 200, 255)  // Light Blue
	firstPoint := pterm.NewRGB(70, 130, 255)   // Sky Blue
	secondPoint := pterm.NewRGB(30, 100, 255)  // Royal Blue
	endColor := pterm.NewRGB(0, 50, 200)       // Deep Blue

	// Print each line with gradient effect
	for _, line := range asciiArt {
		// Create gradient effect for each line
		var gradientLine string
		lineChars := strings.Split(line, "")

		for j, char := range lineChars {
			// Calculate the fade position for this character
			fadePos := float32(j)
			totalChars := float32(len(lineChars))

			// Apply gradient color to each character
			coloredChar := startColor.Fade(0, totalChars, fadePos, firstPoint, secondPoint, endColor).Sprint(char)
			gradientLine += coloredChar
		}

		fmt.Println(gradientLine)
	}

	// Add some space and a subtitle
	fmt.Println()
	subtitle := "🎨 AI-Powered Icon Generator"

	// Apply gradient to subtitle as well
	var gradientSubtitle string
	subtitleChars := strings.Split(subtitle, "")
	for j, char := range subtitleChars {
		fadePos := float32(j)
		totalChars := float32(len(subtitleChars))
		coloredChar := pterm.NewRGB(150, 220, 255).Fade(0, totalChars, fadePos, pterm.NewRGB(50, 150, 255)).Sprint(char)
		gradientSubtitle += coloredChar
	}

	fmt.Println(centerText(gradientSubtitle, 70))

	// Add creator information
	createdBy := "created by hellokaton"
	centeredCreatedBy := centerText(createdBy, 70)
	pterm.NewStyle(pterm.FgGray).Println(centeredCreatedBy)

	fmt.Println() // Add some space at the bottom
}

// centerText centers text within a given width
func centerText(text string, width int) string {
	// Remove ANSI color codes for length calculation
	cleanText := stripAnsiCodes(text)
	textLen := len(cleanText)
	
	if textLen >= width {
		return text
	}
	
	padding := (width - textLen) / 2
	return strings.Repeat(" ", padding) + text
}

// stripAnsiCodes removes ANSI escape sequences for accurate length calculation
func stripAnsiCodes(text string) string {
	// Simple implementation - for more complex cases, use a proper ANSI parser
	result := ""
	inEscape := false
	
	for _, char := range text {
		if char == '\033' { // ESC character
			inEscape = true
			continue
		}
		if inEscape {
			if char == 'm' {
				inEscape = false
			}
			continue
		}
		result += string(char)
	}
	
	return result
}
