package domain

import "fmt"

// Ref holds references from one node to other nodes.
// It tracks dependencies, relationships, events, and shared vocabulary.
type Ref struct {
	Uses        []RefLink `json:"uses,omitempty" yaml:"uses,omitempty"`
	Related     []RefLink `json:"related,omitempty" yaml:"related,omitempty"`
	EmitsEvents []string  `json:"emits_events,omitempty" yaml:"emits_events,omitempty"`
	Vocabulary  []string  `json:"vocabulary,omitempty" yaml:"vocabulary,omitempty"`
}

// RefLink represents a single reference to another node with optional context.
type RefLink struct {
	Target   string `json:"target" yaml:"target"`
	Context  string `json:"context,omitempty" yaml:"context,omitempty"`
	Resolved bool   `json:"resolved,omitempty" yaml:"resolved,omitempty"`
}

// Validate checks that the RefLink has a valid target.
func (r *RefLink) Validate() error {
	if r.Target == "" {
		return fmt.Errorf("reflink Target is required")
	}
	return nil
}
