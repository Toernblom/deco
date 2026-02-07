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

package errors

import (
	"fmt"
	"sort"
	"strings"
)

// Suggester generates "did you mean?" suggestions for typos.
type Suggester struct {
	defaultThreshold int
	maxSuggestions   int
}

// NewSuggester creates a new suggestion engine with default settings.
func NewSuggester() *Suggester {
	return &Suggester{
		defaultThreshold: 2,
		maxSuggestions:   3,
	}
}

// candidate represents a suggestion candidate with its distance
type candidate struct {
	value    string
	distance int
}

// Suggest generates suggestions for the input based on candidates.
// Returns up to maxSuggestions closest matches within the default threshold.
func (s *Suggester) Suggest(input string, candidates []string) []string {
	return s.SuggestWithThreshold(input, candidates, s.defaultThreshold)
}

// SuggestWithThreshold generates suggestions using a custom distance threshold.
// Only candidates within the threshold are returned.
func (s *Suggester) SuggestWithThreshold(input string, candidates []string, threshold int) []string {
	if input == "" || len(candidates) == 0 {
		return []string{}
	}

	inputLower := strings.ToLower(input)
	var matches []candidate

	for _, cand := range candidates {
		candLower := strings.ToLower(cand)

		// Skip exact matches
		if inputLower == candLower {
			return []string{}
		}

		// Calculate distance
		dist := LevenshteinDistance(inputLower, candLower)

		// Apply prefix bonus (reduce distance for matching prefix)
		if strings.HasPrefix(candLower, inputLower) || strings.HasPrefix(inputLower, candLower) {
			dist = max(0, dist-1)
		}

		// Only include if within threshold
		if dist <= threshold {
			matches = append(matches, candidate{
				value:    cand,
				distance: dist,
			})
		}
	}

	// Sort by distance (closest first)
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].distance != matches[j].distance {
			return matches[i].distance < matches[j].distance
		}
		// If equal distance, prefer similar length
		lenDiffI := abs(len(matches[i].value) - len(input))
		lenDiffJ := abs(len(matches[j].value) - len(input))
		if lenDiffI != lenDiffJ {
			return lenDiffI < lenDiffJ
		}
		// If still equal, prefer alphabetically
		return matches[i].value < matches[j].value
	})

	// Limit results
	limit := min(len(matches), s.maxSuggestions)
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = matches[i].value
	}

	return result
}

// FormatSuggestion creates a formatted error message with suggestions.
// Returns empty string if no good suggestions are found.
func (s *Suggester) FormatSuggestion(input string, candidates []string) string {
	suggestions := s.Suggest(input, candidates)

	if len(suggestions) == 0 {
		return ""
	}

	if len(suggestions) == 1 {
		return fmt.Sprintf("did you mean '%s'?", suggestions[0])
	}

	// Multiple suggestions
	quoted := make([]string, len(suggestions))
	for i, sug := range suggestions {
		quoted[i] = "'" + sug + "'"
	}

	if len(suggestions) == 2 {
		return fmt.Sprintf("did you mean %s or %s?", quoted[0], quoted[1])
	}

	// More than 2
	last := quoted[len(quoted)-1]
	rest := strings.Join(quoted[:len(quoted)-1], ", ")
	return fmt.Sprintf("did you mean %s, or %s?", rest, last)
}

// LevenshteinDistance calculates the edit distance between two strings.
// This is the minimum number of single-character edits (insertions, deletions, or substitutions)
// required to change one string into the other.
func LevenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}

	if len(s1) == 0 {
		return len(s2)
	}

	if len(s2) == 0 {
		return len(s1)
	}

	// Create a matrix for dynamic programming
	// matrix[i][j] represents the distance between s1[0:i] and s2[0:j]
	rows := len(s1) + 1
	cols := len(s2) + 1
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}

	// Initialize first row and column
	for i := 0; i < rows; i++ {
		matrix[i][0] = i
	}
	for j := 0; j < cols; j++ {
		matrix[0][j] = j
	}

	// Fill in the rest of the matrix
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1, // deletion
				min(
					matrix[i][j-1]+1,      // insertion
					matrix[i-1][j-1]+cost, // substitution
				),
			)
		}
	}

	return matrix[rows-1][cols-1]
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
