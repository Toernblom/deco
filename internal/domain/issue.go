package domain

import "fmt"

// Issue represents a tracked TBD or question within a node.
// Issues mark areas that need clarification or resolution.
type Issue struct {
	ID          string `json:"id" yaml:"id"`
	Description string `json:"description" yaml:"description"`
	Severity    string `json:"severity" yaml:"severity"` // low, medium, high, critical
	Location    string `json:"location" yaml:"location"` // path to field (e.g., "content.sections[0]")
	Resolved    bool   `json:"resolved" yaml:"resolved"`
}

// Validate checks that all required fields are present and valid.
func (i *Issue) Validate() error {
	if i.ID == "" {
		return fmt.Errorf("issue ID is required")
	}
	if i.Description == "" {
		return fmt.Errorf("issue Description is required")
	}
	if i.Location == "" {
		return fmt.Errorf("issue Location is required")
	}

	// Validate severity level
	validSeverities := map[string]bool{
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}
	if !validSeverities[i.Severity] {
		return fmt.Errorf("issue Severity must be one of: low, medium, high, critical")
	}

	return nil
}
