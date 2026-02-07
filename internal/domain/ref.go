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
