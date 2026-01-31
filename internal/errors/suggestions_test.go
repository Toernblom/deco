package errors_test

import (
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/errors"
)

// Test basic suggestion for close match
func TestSuggester_CloseMatch(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"draft", "approved", "archived"}
	suggestions := suggester.Suggest("draf", candidates)

	if len(suggestions) == 0 {
		t.Fatal("Expected at least one suggestion for 'draf'")
	}

	// Should suggest "draft" as it's very close (1 character difference)
	if suggestions[0] != "draft" {
		t.Errorf("Expected 'draft' as top suggestion, got %q", suggestions[0])
	}
}

// Test multiple close matches
func TestSuggester_MultipleMatches(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"food", "good", "hood", "wood", "mood"}
	suggestions := suggester.Suggest("fod", candidates)

	if len(suggestions) == 0 {
		t.Fatal("Expected suggestions for 'fod'")
	}

	// Should suggest "food" as the closest (1 character difference)
	if suggestions[0] != "food" {
		t.Errorf("Expected 'food' as top suggestion, got %q", suggestions[0])
	}

	// Should include other close matches too
	if len(suggestions) < 2 {
		t.Error("Expected multiple suggestions for close matches")
	}
}

// Test no suggestion when too different
func TestSuggester_TooDifferent(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"draft", "approved", "archived"}
	suggestions := suggester.Suggest("xyz", candidates)

	// Should not suggest anything as "xyz" is too different
	if len(suggestions) > 0 {
		t.Errorf("Expected no suggestions for 'xyz', got %v", suggestions)
	}
}

// Test exact match should not suggest
func TestSuggester_ExactMatch(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"draft", "approved", "archived"}
	suggestions := suggester.Suggest("draft", candidates)

	// Should not suggest when exact match exists
	if len(suggestions) > 0 {
		t.Errorf("Expected no suggestions for exact match, got %v", suggestions)
	}
}

// Test case insensitive matching
func TestSuggester_CaseInsensitive(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"Draft", "Approved", "Archived"}
	suggestions := suggester.Suggest("draf", candidates)

	if len(suggestions) == 0 {
		t.Fatal("Expected suggestions for 'draf'")
	}

	// Should suggest "Draft" even though case differs
	if !strings.EqualFold(suggestions[0], "draft") {
		t.Errorf("Expected 'Draft' as suggestion, got %q", suggestions[0])
	}
}

// Test empty input
func TestSuggester_EmptyInput(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"draft", "approved", "archived"}
	suggestions := suggester.Suggest("", candidates)

	// Empty input should not produce suggestions
	if len(suggestions) > 0 {
		t.Errorf("Expected no suggestions for empty input, got %v", suggestions)
	}
}

// Test empty candidates
func TestSuggester_EmptyCandidates(t *testing.T) {
	suggester := errors.NewSuggester()

	suggestions := suggester.Suggest("draft", []string{})

	// No candidates means no suggestions
	if len(suggestions) > 0 {
		t.Errorf("Expected no suggestions with empty candidates, got %v", suggestions)
	}
}

// Test suggestions are sorted by distance
func TestSuggester_SortedByDistance(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"food", "wood", "hood", "good"}
	suggestions := suggester.Suggest("foo", candidates)

	if len(suggestions) < 2 {
		t.Fatal("Expected multiple suggestions")
	}

	// "food" should be first (1 character difference)
	if suggestions[0] != "food" {
		t.Errorf("Expected 'food' first, got %q", suggestions[0])
	}

	// "good" should be second (2 character differences)
	if suggestions[1] != "good" {
		t.Errorf("Expected 'good' second, got %q", suggestions[1])
	}
}

// Test limit on number of suggestions
func TestSuggester_LimitSuggestions(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	suggestions := suggester.Suggest("x", candidates)

	// Should limit to reasonable number (e.g., 3-5)
	if len(suggestions) > 5 {
		t.Errorf("Expected at most 5 suggestions, got %d", len(suggestions))
	}
}

// Test common typos
func TestSuggester_CommonTypos(t *testing.T) {
	suggester := errors.NewSuggester()

	tests := []struct {
		input      string
		candidates []string
		expected   string
	}{
		{"staus", []string{"status", "state", "stats"}, "status"}, // Could be "stats" or "status" (both distance 1)
		{"sytem", []string{"system", "item", "steam"}, "system"},
		{"refrence", []string{"reference", "preference", "deference"}, "reference"},
		{"mechainc", []string{"mechanic", "mechanical", "mechanism"}, "mechanic"},
	}

	for _, tt := range tests {
		suggestions := suggester.Suggest(tt.input, tt.candidates)
		if len(suggestions) == 0 {
			t.Errorf("Expected suggestion for %q, got none", tt.input)
			continue
		}

		// Check if expected suggestion is in the results (not necessarily first)
		found := false
		for _, sug := range suggestions {
			if sug == tt.expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("For %q: expected %q in suggestions, got %v", tt.input, tt.expected, suggestions)
		}
	}
}

// Test single character transposition
func TestSuggester_Transposition(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"draft", "craft", "graft"}
	suggestions := suggester.Suggest("darft", candidates) // "a" and "r" swapped

	if len(suggestions) == 0 {
		t.Fatal("Expected suggestion for transposition")
	}

	if suggestions[0] != "draft" {
		t.Errorf("Expected 'draft' for transposition, got %q", suggestions[0])
	}
}

// Test with threshold
func TestSuggester_WithThreshold(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"draft", "approved", "archived"}

	// Set a strict threshold
	suggestions := suggester.SuggestWithThreshold("drft", candidates, 1)

	// "drft" is 1 edit from "draft" (missing 'a'), should be included
	if len(suggestions) == 0 {
		t.Error("Expected suggestion within threshold of 1")
	}

	// Try with threshold 0 (exact match only)
	suggestions = suggester.SuggestWithThreshold("drft", candidates, 0)
	if len(suggestions) > 0 {
		t.Error("Expected no suggestions with threshold 0 for non-exact match")
	}
}

// Test prefix matching gets bonus
func TestSuggester_PrefixBonus(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"draft", "craft", "graft", "aft"}
	suggestions := suggester.Suggest("dra", candidates)

	if len(suggestions) == 0 {
		t.Fatal("Expected suggestions")
	}

	// "draft" has matching prefix "dra", should be preferred
	if suggestions[0] != "draft" {
		t.Errorf("Expected 'draft' with prefix bonus, got %q", suggestions[0])
	}
}

// Test suffix matching
func TestSuggester_SuffixMatch(t *testing.T) {
	suggester := errors.NewSuggester()

	candidates := []string{"mechanic", "ceramic", "dynamic", "panic"}
	suggestions := suggester.Suggest("manic", candidates)

	if len(suggestions) == 0 {
		t.Fatal("Expected suggestions")
	}

	// "mechanic" ends with "anic", close to "manic"
	// "panic" ends with "anic", even closer
	// Should suggest one of these
	found := false
	for _, s := range suggestions {
		if s == "mechanic" || s == "panic" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'mechanic' or 'panic' in suggestions, got %v", suggestions)
	}
}

// Test Levenshtein distance calculation directly
func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"a", "a", 0},
		{"draft", "draft", 0},
		{"draft", "draf", 1},
		{"draft", "drft", 1},
		{"draft", "craft", 1},
		{"kitten", "sitting", 3},
		{"saturday", "sunday", 3},
	}

	for _, tt := range tests {
		dist := errors.LevenshteinDistance(tt.s1, tt.s2)
		if dist != tt.expected {
			t.Errorf("LevenshteinDistance(%q, %q) = %d, expected %d",
				tt.s1, tt.s2, dist, tt.expected)
		}
	}
}

// Test generating error message with suggestion
func TestSuggester_FormatSuggestion(t *testing.T) {
	suggester := errors.NewSuggester()

	suggestion := suggester.FormatSuggestion("staus", []string{"status", "state", "stats"})

	// Should format as a helpful message
	if suggestion == "" {
		t.Error("Expected non-empty formatted suggestion")
	}

	// Should mention the correct value
	if !strings.Contains(suggestion, "status") {
		t.Errorf("Expected suggestion to mention 'status', got: %s", suggestion)
	}

	// Should be phrased as a question or help
	lower := strings.ToLower(suggestion)
	if !strings.Contains(lower, "did you mean") && !strings.Contains(lower, "try") {
		t.Errorf("Expected helpful phrasing in suggestion, got: %s", suggestion)
	}
}

// Test no formatted suggestion when no match
func TestSuggester_FormatSuggestion_NoMatch(t *testing.T) {
	suggester := errors.NewSuggester()

	suggestion := suggester.FormatSuggestion("xyz", []string{"draft", "approved"})

	// Should return empty string when no good suggestions
	if suggestion != "" {
		t.Errorf("Expected empty suggestion for no match, got: %s", suggestion)
	}
}
