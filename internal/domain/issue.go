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
