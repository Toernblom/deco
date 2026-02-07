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

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
)

// ValidStatuses is the canonical set of valid node status values.
var ValidStatuses = []string{"draft", "review", "approved", "deprecated", "archived"}

// ValidSeverities is the canonical set of valid issue severity values.
var ValidSeverities = []string{"low", "medium", "high", "critical"}

// validateStatus checks if the given status is valid.
func validateStatus(status string) error {
	if status == "" {
		return nil
	}
	for _, s := range ValidStatuses {
		if status == s {
			return nil
		}
	}
	return newFilterError("status", status, ValidStatuses)
}

// validateSeverity checks if the given severity is valid.
func validateSeverity(severity string) error {
	if severity == "" {
		return nil
	}
	for _, s := range ValidSeverities {
		if severity == s {
			return nil
		}
	}
	return newFilterError("severity", severity, ValidSeverities)
}

// validateKind checks if the given kind exists among the loaded nodes.
func validateKind(kind string, nodes []domain.Node) error {
	if kind == "" {
		return nil
	}
	validKinds := collectKinds(nodes)
	for _, k := range validKinds {
		if kind == k {
			return nil
		}
	}
	return newFilterError("kind", kind, validKinds)
}

// validBuiltInBlockTypes is the canonical set of built-in block types.
var validBuiltInBlockTypes = []string{"doc", "list", "mechanic", "param", "rule", "table"}

// validateBlockType checks if the given block type is valid (built-in or custom).
func validateBlockType(blockType string, customBlockTypes map[string]config.BlockTypeConfig) error {
	if blockType == "" {
		return nil
	}
	// Check built-in types
	for _, bt := range validBuiltInBlockTypes {
		if blockType == bt {
			return nil
		}
	}
	// Check custom types
	if customBlockTypes != nil {
		if _, ok := customBlockTypes[blockType]; ok {
			return nil
		}
	}
	// Build full list for error message
	all := make([]string, len(validBuiltInBlockTypes))
	copy(all, validBuiltInBlockTypes)
	if customBlockTypes != nil {
		for name := range customBlockTypes {
			all = append(all, name)
		}
		sort.Strings(all)
	}
	return newFilterError("block-type", blockType, all)
}

// validateFieldFilter checks that a field filter has key=value format.
func validateFieldFilter(field string) error {
	if !strings.Contains(field, "=") {
		return fmt.Errorf("invalid --field value %q: expected key=value format (e.g. --field age=bronze)", field)
	}
	return nil
}

// collectKinds returns sorted unique kinds from a set of nodes.
func collectKinds(nodes []domain.Node) []string {
	seen := make(map[string]bool)
	for _, n := range nodes {
		seen[n.Kind] = true
	}
	kinds := make([]string, 0, len(seen))
	for k := range seen {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	return kinds
}

// newFilterError creates a user-friendly error for an invalid filter value.
func newFilterError(filterName, value string, validValues []string) error {
	msg := fmt.Sprintf("unknown %s %q\nValid values: %s", filterName, value, strings.Join(validValues, ", "))
	if suggestion := closestMatch(value, validValues); suggestion != "" {
		msg += fmt.Sprintf("\nDid you mean %q?", suggestion)
	}
	return fmt.Errorf("%s", msg)
}

// closestMatch finds the closest string match using Levenshtein distance.
// Returns empty string if no close match found (distance > half the target length).
func closestMatch(input string, candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}
	best := ""
	bestDist := len(input)/2 + 2 // max acceptable distance
	for _, c := range candidates {
		d := levenshtein(strings.ToLower(input), strings.ToLower(c))
		if d < bestDist {
			bestDist = d
			best = c
		}
	}
	return best
}

// levenshtein computes the edit distance between two strings.
func levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	// Use single-row optimization
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr := make([]int, len(b)+1)
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(curr[j-1]+1, min(prev[j]+1, prev[j-1]+cost))
		}
		prev = curr
	}
	return prev[len(b)]
}
