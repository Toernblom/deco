package domain

import (
	"fmt"
	"time"
)

// AuditEntry represents a single entry in the audit log.
// It tracks changes to nodes over time (who, what, when).
type AuditEntry struct {
	Timestamp time.Time              `json:"timestamp" yaml:"timestamp"`
	NodeID    string                 `json:"node_id" yaml:"node_id"`
	Operation string                 `json:"operation" yaml:"operation"` // create, update, delete, set, append, unset, move
	User      string                 `json:"user" yaml:"user"`
	Before    map[string]interface{} `json:"before,omitempty" yaml:"before,omitempty"`
	After     map[string]interface{} `json:"after,omitempty" yaml:"after,omitempty"`
}

// Validate checks that all required fields are present and valid.
func (a *AuditEntry) Validate() error {
	if a.Timestamp.IsZero() {
		return fmt.Errorf("audit entry Timestamp is required")
	}
	if a.NodeID == "" {
		return fmt.Errorf("audit entry NodeID is required")
	}
	if a.User == "" {
		return fmt.Errorf("audit entry User is required")
	}

	// Validate operation type
	validOperations := map[string]bool{
		"create":  true,
		"update":  true,
		"delete":  true,
		"set":     true,
		"append":  true,
		"unset":   true,
		"move":    true,
		"submit":  true, // draft -> review
		"approve": true, // add approval
		"reject":  true, // review -> draft
	}
	if !validOperations[a.Operation] {
		return fmt.Errorf("audit entry Operation must be one of: create, update, delete, set, append, unset, move, submit, approve, reject")
	}

	return nil
}
