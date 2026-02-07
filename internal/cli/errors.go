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

package cli

import "fmt"

// ExitError represents an error with a specific exit code.
type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string {
	return e.Message
}

// NewExitError creates an ExitError with the specified code and message.
func NewExitError(code int, message string) *ExitError {
	return &ExitError{Code: code, Message: message}
}

// NewExitErrorf creates an ExitError with a formatted message.
func NewExitErrorf(code int, format string, args ...interface{}) *ExitError {
	return &ExitError{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Exit code constants for CLI commands.
const (
	ExitCodeSuccess        = 0
	ExitCodeError          = 1
	ExitCodeSchemaMismatch = 2 // Schema version mismatch, migration needed
)
