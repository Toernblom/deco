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
