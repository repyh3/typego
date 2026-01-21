package internal

import (
	"fmt"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// Info prints an informational message with blue prefix
func Info(msg string) {
	fmt.Printf("%s[INFO]%s %s\n", colorBlue, colorReset, msg)
}

// Success prints a success message with green checkmark
func Success(msg string) {
	fmt.Printf("%s✅%s %s\n", colorGreen, colorReset, msg)
}

// Warn prints a warning message with yellow prefix
func Warn(msg string) {
	fmt.Printf("%s⚠️  Warning:%s %s\n", colorYellow, colorReset, msg)
}

// Error prints an error message with red prefix
func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s❌ Error:%s %s\n", colorRed, colorReset, msg)
}

// Step prints a step indicator (emoji + message)
func Step(emoji, msg string) {
	fmt.Printf("%s %s\n", emoji, msg)
}

// Verbose prints only if verbose mode is enabled
func Verbose(msg string) {
	if os.Getenv("TYPEGO_VERBOSE") == "true" {
		fmt.Printf("%s[DEBUG]%s %s\n", colorCyan, colorReset, msg)
	}
}
