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

package domain

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorGray   = "\033[90m"
)

// ErrorFormatter formats DecoError objects in a Rust-like style
type ErrorFormatter struct {
	useColor bool
}

// NewErrorFormatter creates a new error formatter with color enabled by default
func NewErrorFormatter() *ErrorFormatter {
	return &ErrorFormatter{
		useColor: false, // Default to no color for tests
	}
}

// SetColor enables or disables color output
func (f *ErrorFormatter) SetColor(enabled bool) {
	f.useColor = enabled
}

// Format formats a single DecoError into a human-readable string
func (f *ErrorFormatter) Format(err DecoError) string {
	var b strings.Builder

	// Error header with code and summary
	if err.Code != "" {
		if f.useColor {
			b.WriteString(f.colorize(colorBold+colorRed, fmt.Sprintf("error[%s]", err.Code)))
		} else {
			b.WriteString(fmt.Sprintf("error[%s]", err.Code))
		}
		b.WriteString(": ")
	} else {
		if f.useColor {
			b.WriteString(f.colorize(colorBold+colorRed, "error"))
		} else {
			b.WriteString("error")
		}
		b.WriteString(": ")
	}

	if err.Summary != "" {
		if f.useColor {
			b.WriteString(f.colorize(colorBold, err.Summary))
		} else {
			b.WriteString(err.Summary)
		}
	}
	b.WriteString("\n")

	// Location
	if err.Location != nil {
		if f.useColor {
			b.WriteString(f.colorize(colorBlue, fmt.Sprintf("  --> %s\n", err.Location.String())))
		} else {
			b.WriteString(fmt.Sprintf("  --> %s\n", err.Location.String()))
		}
	}

	// Detail
	if err.Detail != "" {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("  %s\n", err.Detail))
	}

	// Context
	if len(err.Context) > 0 {
		b.WriteString("\n")
		if f.useColor {
			b.WriteString(f.colorize(colorCyan, "  Context:\n"))
		} else {
			b.WriteString("  Context:\n")
		}
		for _, ctx := range err.Context {
			b.WriteString(fmt.Sprintf("    • %s\n", ctx))
		}
	}

	// Suggestion
	if err.Suggestion != "" {
		b.WriteString("\n")
		if f.useColor {
			b.WriteString(f.colorize(colorYellow, "  Suggestion: "))
		} else {
			b.WriteString("  Suggestion: ")
		}
		b.WriteString(err.Suggestion)
		b.WriteString("\n")
	}

	// Related
	if len(err.Related) > 0 {
		b.WriteString("\n")
		if f.useColor {
			b.WriteString(f.colorize(colorCyan, "  Related:\n"))
		} else {
			b.WriteString("  Related:\n")
		}
		for _, rel := range err.Related {
			b.WriteString(fmt.Sprintf("    • %s (%s)\n", rel.NodeID, rel.Reason))
		}
	}

	return b.String()
}

// FormatWithSource formats a DecoError with source code context
func (f *ErrorFormatter) FormatWithSource(err DecoError, sourceLines []string) string {
	var b strings.Builder

	// Write the basic error format
	b.WriteString(f.Format(err))

	// Add source context if location is available
	if err.Location != nil && err.Location.Line > 0 && len(sourceLines) > 0 {
		b.WriteString("\n")

		lineNum := err.Location.Line
		contextBefore := 2
		contextAfter := 2

		startLine := lineNum - contextBefore
		if startLine < 1 {
			startLine = 1
		}

		endLine := lineNum + contextAfter
		if endLine > len(sourceLines) {
			endLine = len(sourceLines)
		}

		// Calculate max line number width for alignment
		maxLineNumWidth := len(fmt.Sprintf("%d", endLine))

		// Print context lines
		for i := startLine; i <= endLine; i++ {
			linePrefix := fmt.Sprintf("%*d | ", maxLineNumWidth, i)

			if i == lineNum {
				// Error line
				if f.useColor {
					b.WriteString(f.colorize(colorBlue, linePrefix))
					b.WriteString(f.colorize(colorBold, sourceLines[i-1]))
				} else {
					b.WriteString(linePrefix)
					b.WriteString(sourceLines[i-1])
				}
				b.WriteString("\n")

				// Add column pointer if column is specified
				if err.Location.Column > 0 {
					padding := strings.Repeat(" ", maxLineNumWidth+3+err.Location.Column-1)
					if f.useColor {
						b.WriteString(f.colorize(colorRed, padding+"^"))
					} else {
						b.WriteString(padding + "^")
					}
					b.WriteString("\n")
				}
			} else {
				// Context line
				if f.useColor {
					b.WriteString(f.colorize(colorGray, linePrefix))
					b.WriteString(f.colorize(colorGray, sourceLines[i-1]))
				} else {
					b.WriteString(linePrefix)
					b.WriteString(sourceLines[i-1])
				}
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

// FormatMultiple formats multiple errors
func (f *ErrorFormatter) FormatMultiple(errors []DecoError) string {
	if len(errors) == 0 {
		return ""
	}

	var b strings.Builder

	// Error count summary
	if f.useColor {
		b.WriteString(f.colorize(colorBold+colorRed, fmt.Sprintf("%d error(s) found:\n\n", len(errors))))
	} else {
		b.WriteString(fmt.Sprintf("%d error(s) found:\n\n", len(errors)))
	}

	// Format each error
	for i, err := range errors {
		if i > 0 {
			b.WriteString("\n")
			b.WriteString(strings.Repeat("-", 60))
			b.WriteString("\n\n")
		}
		b.WriteString(f.Format(err))
	}

	return b.String()
}

// colorize wraps text with ANSI color codes if color is enabled
func (f *ErrorFormatter) colorize(color, text string) string {
	if !f.useColor {
		return text
	}
	return color + text + colorReset
}
