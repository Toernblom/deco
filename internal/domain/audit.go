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
	"time"
)

// AuditEntry represents a single entry in the audit log.
// It tracks changes to nodes over time (who, what, when).
type AuditEntry struct {
	Timestamp   time.Time              `json:"timestamp" yaml:"timestamp"`
	NodeID      string                 `json:"node_id" yaml:"node_id"`
	Operation   string                 `json:"operation" yaml:"operation"` // create, update, delete, set, append, unset, move, baseline
	User        string                 `json:"user" yaml:"user"`
	ContentHash string                 `json:"content_hash,omitempty" yaml:"content_hash,omitempty"`
	Before      map[string]interface{} `json:"before,omitempty" yaml:"before,omitempty"`
	After       map[string]interface{} `json:"after,omitempty" yaml:"after,omitempty"`
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
		"create":   true,
		"update":   true,
		"delete":   true,
		"set":      true,
		"append":   true,
		"unset":    true,
		"move":     true,
		"submit":   true, // draft -> review
		"approve":  true, // add approval
		"reject":   true, // review -> draft
		"sync":     true, // auto-fix unversioned edits
		"baseline": true, // record current state without modification
		"migrate":  true, // schema migration
		"rewrite":  true, // full node replacement
	}
	if !validOperations[a.Operation] {
		return fmt.Errorf("audit entry Operation must be one of: create, update, delete, set, append, unset, move, submit, approve, reject, sync, baseline, migrate, rewrite")
	}

	return nil
}
