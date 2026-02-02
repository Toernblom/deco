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
	ExitCodeSuccess         = 0
	ExitCodeError           = 1
	ExitCodeSchemaMismatch  = 2 // Schema version mismatch, migration needed
)
