package domain

import "fmt"

// Constraint defines a validation rule that must be satisfied.
// Constraints use CEL (Common Expression Language) for validation.
type Constraint struct {
	Expr    string `json:"expr" yaml:"expr"`       // CEL expression
	Message string `json:"message" yaml:"message"` // Error message if constraint fails
	Scope   string `json:"scope" yaml:"scope"`     // Which nodes this applies to (e.g., "all", "mechanic", "systems/*")
}

// Validate checks that all required fields are present.
func (c *Constraint) Validate() error {
	if c.Expr == "" {
		return fmt.Errorf("constraint Expr is required")
	}
	if c.Message == "" {
		return fmt.Errorf("constraint Message is required")
	}
	if c.Scope == "" {
		return fmt.Errorf("constraint Scope is required")
	}
	return nil
}
