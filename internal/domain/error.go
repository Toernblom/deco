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

// Location represents a position in a file
type Location struct {
	File   string
	Line   int
	Column int
}

// String formats the location as "file:line:column", "file:line", or "file"
func (l Location) String() string {
	if l.Column > 0 {
		return fmt.Sprintf("%s:%d:%d", l.File, l.Line, l.Column)
	}
	if l.Line > 0 {
		return fmt.Sprintf("%s:%d", l.File, l.Line)
	}
	return l.File
}

// Related represents a related node in an error context
type Related struct {
	NodeID string
	Reason string
}

// DecoError represents a structured error following Rust-like error patterns
type DecoError struct {
	Code       string
	Summary    string
	Detail     string
	Location   *Location
	Context    []string
	Suggestion string
	Related    []Related
}

// Error implements the error interface
func (e DecoError) Error() string {
	var parts []string

	// Start with code and summary
	if e.Code != "" {
		parts = append(parts, fmt.Sprintf("[%s]", e.Code))
	}
	if e.Summary != "" {
		parts = append(parts, e.Summary)
	}

	// Add location if present
	if e.Location != nil {
		parts = append(parts, fmt.Sprintf("at %s", e.Location.String()))
	}

	// Join the main parts
	result := strings.Join(parts, " ")

	// Add detail if present
	if e.Detail != "" {
		result += fmt.Sprintf(": %s", e.Detail)
	}

	// Add context if present
	if len(e.Context) > 0 {
		result += fmt.Sprintf(" (%s)", strings.Join(e.Context, ", "))
	}

	// Add suggestion if present
	if e.Suggestion != "" {
		result += fmt.Sprintf(" [Hint: %s]", e.Suggestion)
	}

	return result
}
