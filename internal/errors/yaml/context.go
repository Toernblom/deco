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

package yaml

import (
	"strings"

	"github.com/Toernblom/deco/internal/domain"
)

// ExtractContext extracts source lines around a location for error display.
// before and after specify how many lines to include before and after the target line.
// Returns empty slice if the location is invalid.
func ExtractContext(content string, loc domain.Location, before, after int) []string {
	return ExtractContextBytes([]byte(content), loc, before, after)
}

// ExtractContextBytes is like ExtractContext but takes []byte input.
func ExtractContextBytes(content []byte, loc domain.Location, before, after int) []string {
	if loc.Line <= 0 {
		return []string{}
	}

	// Handle empty content
	if len(content) == 0 {
		return []string{}
	}

	// Split into lines, handling both LF and CRLF
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "\r\n", "\n")
	lines := strings.Split(contentStr, "\n")

	// Check if line number is valid
	if loc.Line > len(lines) {
		return []string{}
	}

	// Calculate range
	startLine := loc.Line - before
	if startLine < 1 {
		startLine = 1
	}

	endLine := loc.Line + after
	if endLine > len(lines) {
		endLine = len(lines)
	}

	// Extract the range (convert to 0-based indexing)
	result := make([]string, 0, endLine-startLine+1)
	for i := startLine - 1; i < endLine; i++ {
		result = append(result, lines[i])
	}

	return result
}

// HighlightColumn creates a visual pointer at the specified column.
// Returns a string with the source line and a pointer line below it.
func HighlightColumn(content string, loc domain.Location) string {
	return HighlightColumnWithLength(content, loc, 1)
}

// HighlightColumnWithLength creates a visual pointer spanning multiple characters.
// Returns a string with the source line and a pointer line below it.
func HighlightColumnWithLength(content string, loc domain.Location, length int) string {
	if loc.Line <= 0 || loc.Column <= 0 {
		return ""
	}

	// Split into lines
	contentStr := strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(contentStr, "\n")

	// Check if line number is valid
	if loc.Line > len(lines) {
		return ""
	}

	// Get the target line (1-based to 0-based)
	sourceLine := lines[loc.Line-1]

	// Build pointer line
	var pointer strings.Builder

	// Add spaces before the pointer (column is 1-based)
	for i := 0; i < loc.Column-1; i++ {
		pointer.WriteByte(' ')
	}

	// Add pointer characters
	if length == 1 {
		pointer.WriteByte('^')
	} else {
		// Use '^' for first character, '~' for the rest
		pointer.WriteByte('^')
		for i := 1; i < length; i++ {
			pointer.WriteByte('~')
		}
	}

	// Combine source line and pointer line
	return sourceLine + "\n" + pointer.String()
}
