package yaml_test

import (
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	yaml_errors "github.com/Toernblom/deco/internal/errors/yaml"
)

// Test basic context extraction with surrounding lines
func TestContextExtractor_BasicExtraction(t *testing.T) {
	content := `line 1
line 2
line 3
line 4
line 5
line 6
line 7`

	loc := domain.Location{Line: 4, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 2, 2)

	expected := []string{
		"line 2",
		"line 3",
		"line 4",
		"line 5",
		"line 6",
	}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

// Test extraction at start of file
func TestContextExtractor_StartOfFile(t *testing.T) {
	content := `line 1
line 2
line 3
line 4
line 5`

	loc := domain.Location{Line: 1, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 2, 2)

	// Should only get lines 1-3 (can't go before line 1)
	expected := []string{
		"line 1",
		"line 2",
		"line 3",
	}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

// Test extraction at end of file
func TestContextExtractor_EndOfFile(t *testing.T) {
	content := `line 1
line 2
line 3
line 4
line 5`

	loc := domain.Location{Line: 5, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 2, 2)

	// Should only get lines 3-5 (can't go past line 5)
	expected := []string{
		"line 3",
		"line 4",
		"line 5",
	}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

// Test with zero context
func TestContextExtractor_NoContext(t *testing.T) {
	content := `line 1
line 2
line 3`

	loc := domain.Location{Line: 2, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 0, 0)

	// Should only get the target line
	expected := []string{"line 2"}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}

	if lines[0] != expected[0] {
		t.Errorf("Expected %q, got %q", expected[0], lines[0])
	}
}

// Test with asymmetric context
func TestContextExtractor_AsymmetricContext(t *testing.T) {
	content := `line 1
line 2
line 3
line 4
line 5
line 6
line 7`

	loc := domain.Location{Line: 4, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 1, 3)

	expected := []string{
		"line 3",
		"line 4",
		"line 5",
		"line 6",
		"line 7",
	}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

// Test with invalid line number
func TestContextExtractor_InvalidLine(t *testing.T) {
	content := `line 1
line 2
line 3`

	// Line 0 (invalid)
	loc := domain.Location{Line: 0, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 2, 2)
	if len(lines) != 0 {
		t.Errorf("Expected empty result for line 0, got %d lines", len(lines))
	}

	// Line beyond end
	loc = domain.Location{Line: 100, Column: 1}
	lines = yaml_errors.ExtractContext(content, loc, 2, 2)
	if len(lines) != 0 {
		t.Errorf("Expected empty result for line beyond end, got %d lines", len(lines))
	}
}

// Test with empty content
func TestContextExtractor_EmptyContent(t *testing.T) {
	content := ""

	loc := domain.Location{Line: 1, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 2, 2)

	if len(lines) != 0 {
		t.Errorf("Expected empty result for empty content, got %d lines", len(lines))
	}
}

// Test with single line
func TestContextExtractor_SingleLine(t *testing.T) {
	content := "only one line"

	loc := domain.Location{Line: 1, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 5, 5)

	expected := []string{"only one line"}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}

	if lines[0] != expected[0] {
		t.Errorf("Expected %q, got %q", expected[0], lines[0])
	}
}

// Test preserving indentation
func TestContextExtractor_PreservesIndentation(t *testing.T) {
	content := `root:
  nested:
    deep:
      value: 123
  other: value`

	loc := domain.Location{Line: 4, Column: 7}
	lines := yaml_errors.ExtractContext(content, loc, 1, 1)

	expected := []string{
		"    deep:",
		"      value: 123",
		"  other: value",
	}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

// Test with trailing newlines
func TestContextExtractor_TrailingNewlines(t *testing.T) {
	content := `line 1
line 2
line 3

`

	loc := domain.Location{Line: 2, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 1, 2)

	// Should handle trailing empty lines gracefully
	if len(lines) == 0 {
		t.Error("Expected non-empty result")
	}

	// Should include line 2
	found := false
	for _, line := range lines {
		if line == "line 2" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'line 2' in context")
	}
}

// Test column highlighting
func TestContextExtractor_ColumnHighlight(t *testing.T) {
	content := `id: systems/food
kind: system
status: draft`

	loc := domain.Location{Line: 2, Column: 7}

	highlighted := yaml_errors.HighlightColumn(content, loc)

	// Should contain the original line
	if !strings.Contains(highlighted, "kind: system") {
		t.Error("Highlighted output should contain the original line")
	}

	// Should contain a pointer at column 7
	lines := strings.Split(highlighted, "\n")
	if len(lines) < 2 {
		t.Fatal("Expected at least 2 lines in highlighted output")
	}

	// The second line should be the pointer
	pointerLine := lines[1]
	if !strings.Contains(pointerLine, "^") {
		t.Errorf("Expected pointer line to contain '^', got: %q", pointerLine)
	}

	// Count spaces before the pointer
	spaces := 0
	for _, ch := range pointerLine {
		if ch == ' ' {
			spaces++
		} else {
			break
		}
	}

	// Should have column-1 spaces (column is 1-based)
	expectedSpaces := loc.Column - 1
	if spaces != expectedSpaces {
		t.Errorf("Expected %d spaces before pointer, got %d", expectedSpaces, spaces)
	}
}

// Test column highlighting at start of line
func TestContextExtractor_ColumnHighlightStart(t *testing.T) {
	content := `id: systems/food`

	loc := domain.Location{Line: 1, Column: 1}

	highlighted := yaml_errors.HighlightColumn(content, loc)
	lines := strings.Split(highlighted, "\n")

	if len(lines) < 2 {
		t.Fatal("Expected at least 2 lines in highlighted output")
	}

	// Pointer should be at the start (no leading spaces)
	pointerLine := lines[1]
	if !strings.HasPrefix(pointerLine, "^") {
		t.Errorf("Expected pointer at start, got: %q", pointerLine)
	}
}

// Test column highlighting with multi-character pointer
func TestContextExtractor_ColumnHighlightLength(t *testing.T) {
	content := `id: systems/food`

	loc := domain.Location{Line: 1, Column: 5}

	// Highlight 7 characters (the word "systems")
	highlighted := yaml_errors.HighlightColumnWithLength(content, loc, 7)
	lines := strings.Split(highlighted, "\n")

	if len(lines) < 2 {
		t.Fatal("Expected at least 2 lines in highlighted output")
	}

	pointerLine := lines[1]

	// Should have 7 pointer characters
	carets := strings.Count(pointerLine, "^") + strings.Count(pointerLine, "~")
	if carets < 5 { // Allow some flexibility in implementation
		t.Errorf("Expected multiple pointer characters for length=7, got: %q", pointerLine)
	}
}

// Test with bytes input
func TestContextExtractor_BytesInput(t *testing.T) {
	content := []byte(`line 1
line 2
line 3`)

	loc := domain.Location{Line: 2, Column: 1}
	lines := yaml_errors.ExtractContextBytes(content, loc, 1, 1)

	expected := []string{
		"line 1",
		"line 2",
		"line 3",
	}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

// Test with CRLF line endings (Windows)
func TestContextExtractor_CRLFEndings(t *testing.T) {
	content := "line 1\r\nline 2\r\nline 3\r\nline 4\r\nline 5"

	loc := domain.Location{Line: 3, Column: 1}
	lines := yaml_errors.ExtractContext(content, loc, 1, 1)

	expected := []string{
		"line 2",
		"line 3",
		"line 4",
	}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}
