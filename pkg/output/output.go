package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var DebugEnabled bool

var (
	infoColor    = color.New(color.FgCyan)
	warnColor    = color.New(color.FgYellow)
	errorColor   = color.New(color.FgRed)
	successColor = color.New(color.FgGreen)
	dimColor     = color.New(color.FgHiBlack)
)

// Info prints an informational message.
func Info(format string, args ...any) {
	infoColor.Fprintf(os.Stderr, "ℹ "+format+"\n", args...)
}

// Warning prints a warning message.
func Warning(format string, args ...any) {
	warnColor.Fprintf(os.Stderr, "⚠ "+format+"\n", args...)
}

// Error prints an error message.
func Error(format string, args ...any) {
	errorColor.Fprintf(os.Stderr, "✗ "+format+"\n", args...)
}

// Success prints a success message.
func Success(format string, args ...any) {
	successColor.Fprintf(os.Stderr, "✓ "+format+"\n", args...)
}

// Debug prints a debug message when debug mode is on.
func Debug(format string, args ...any) {
	if DebugEnabled {
		dimColor.Fprintf(os.Stderr, "⋯ "+format+"\n", args...)
	}
}

// Dim prints dimmed text.
func Dim(format string, args ...any) {
	dimColor.Fprintf(os.Stdout, format+"\n", args...)
}

// Print prints to stdout.
func Print(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}
