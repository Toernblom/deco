// Copyright (C) 2026 Anton Törnblom
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

// Package style provides consistent terminal styling for deco CLI output.
package style

import (
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// ColorMode controls when colors are used
type ColorMode int

const (
	ColorAuto   ColorMode = iota // Auto-detect TTY
	ColorAlways                  // Always use colors
	ColorNever                   // Never use colors
)

var (
	// Global color mode, defaults to auto
	mode = ColorAuto

	// Semantic styles
	Error   = color.New(color.FgRed, color.Bold)
	Warning = color.New(color.FgYellow)
	Success = color.New(color.FgGreen)
	Info    = color.New(color.FgCyan)
	Header  = color.New(color.Bold)
	Muted   = color.New(color.FgHiBlack)
	Code    = color.New(color.FgMagenta)

	// Severity colors
	Critical = color.New(color.FgRed, color.Bold)
	High     = color.New(color.FgRed)
	Medium   = color.New(color.FgYellow)
	Low      = color.New(color.FgBlue)

	// Status colors
	StatusDraft     = color.New(color.FgYellow)
	StatusPublished = color.New(color.FgGreen)
	StatusArchived  = color.New(color.FgHiBlack)
)

// SetMode sets the global color mode
func SetMode(m ColorMode) {
	mode = m
	updateColors()
}

// GetMode returns the current color mode
func GetMode() ColorMode {
	return mode
}

// IsEnabled returns whether colors are currently enabled
func IsEnabled() bool {
	switch mode {
	case ColorAlways:
		return true
	case ColorNever:
		return false
	default: // ColorAuto
		return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	}
}

// updateColors enables/disables all color objects based on mode
func updateColors() {
	enabled := IsEnabled()

	if enabled {
		color.NoColor = false
	} else {
		color.NoColor = true
	}
}

// Init initializes the style system. Call this early in main.
func Init() {
	updateColors()
}

// ParseColorMode parses a color mode string (auto, always, never)
func ParseColorMode(s string) ColorMode {
	switch s {
	case "always", "true", "1":
		return ColorAlways
	case "never", "false", "0":
		return ColorNever
	default:
		return ColorAuto
	}
}

// Symbols for consistent iconography
const (
	SymbolSuccess = "✓"
	SymbolError   = "✗"
	SymbolWarning = "!"
	SymbolInfo    = "•"
	SymbolArrow   = "→"
	SymbolBullet  = "•"
)

// Fmt returns a styled string using the given color
func Fmt(c *color.Color, format string, args ...interface{}) string {
	return c.Sprintf(format, args...)
}

// ErrorIcon returns a styled error symbol
func ErrorIcon() string {
	return Error.Sprint(SymbolError)
}

// SuccessIcon returns a styled success symbol
func SuccessIcon() string {
	return Success.Sprint(SymbolSuccess)
}

// WarningIcon returns a styled warning symbol
func WarningIcon() string {
	return Warning.Sprint(SymbolWarning)
}

// SeverityColor returns the color for a severity level
func SeverityColor(severity string) *color.Color {
	switch severity {
	case "critical":
		return Critical
	case "high":
		return High
	case "medium":
		return Medium
	case "low":
		return Low
	default:
		return Muted
	}
}

// StatusColor returns the color for a status
func StatusColor(status string) *color.Color {
	switch status {
	case "draft":
		return StatusDraft
	case "published":
		return StatusPublished
	case "archived":
		return StatusArchived
	default:
		return nil // No color
	}
}
