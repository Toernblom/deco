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

package domain_test

import (
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestErrorFormatter_BasicError(t *testing.T) {
	formatter := domain.NewErrorFormatter()
	err := domain.DecoError{
		Code:    "E001",
		Summary: "Node not found",
		Detail:  "The node 'systems/food' does not exist",
	}

	output := formatter.Format(err)

	if output == "" {
		t.Error("Expected non-empty output")
	}

	// Should contain code
	if !strings.Contains(output, "E001") {
		t.Errorf("Output should contain error code, got:\n%s", output)
	}

	// Should contain summary
	if !strings.Contains(output, "Node not found") {
		t.Errorf("Output should contain summary, got:\n%s", output)
	}

	// Should contain detail
	if !strings.Contains(output, "The node 'systems/food' does not exist") {
		t.Errorf("Output should contain detail, got:\n%s", output)
	}
}

func TestErrorFormatter_WithLocation(t *testing.T) {
	formatter := domain.NewErrorFormatter()
	err := domain.DecoError{
		Code:     "E002",
		Summary:  "Invalid reference",
		Detail:   "Reference to 'unknown/node' cannot be resolved",
		Location: &domain.Location{File: "test.yaml", Line: 10, Column: 5},
	}

	output := formatter.Format(err)

	// Should contain file location
	if !strings.Contains(output, "test.yaml") {
		t.Errorf("Output should contain file location, got:\n%s", output)
	}

	// Should contain line number
	if !strings.Contains(output, "10") {
		t.Errorf("Output should contain line number, got:\n%s", output)
	}

	// Should contain column number
	if !strings.Contains(output, "5") {
		t.Errorf("Output should contain column number, got:\n%s", output)
	}
}

func TestErrorFormatter_WithContext(t *testing.T) {
	formatter := domain.NewErrorFormatter()
	err := domain.DecoError{
		Code:    "E003",
		Summary: "Validation failed",
		Detail:  "Field 'status' has invalid value",
		Context: []string{"in node 'systems/food'", "during validation"},
	}

	output := formatter.Format(err)

	// Should contain context items
	if !strings.Contains(output, "in node 'systems/food'") {
		t.Errorf("Output should contain first context item, got:\n%s", output)
	}
	if !strings.Contains(output, "during validation") {
		t.Errorf("Output should contain second context item, got:\n%s", output)
	}
}

func TestErrorFormatter_WithSuggestion(t *testing.T) {
	formatter := domain.NewErrorFormatter()
	err := domain.DecoError{
		Code:       "E004",
		Summary:    "Node not found",
		Suggestion: "Check that the file exists in .deco/nodes/systems/",
	}

	output := formatter.Format(err)

	// Should contain suggestion
	if !strings.Contains(output, "Check that the file exists") {
		t.Errorf("Output should contain suggestion, got:\n%s", output)
	}
}

func TestErrorFormatter_WithRelated(t *testing.T) {
	formatter := domain.NewErrorFormatter()
	err := domain.DecoError{
		Code:    "E005",
		Summary: "Circular dependency",
		Related: []domain.Related{
			{NodeID: "systems/food", Reason: "depends on systems/water"},
			{NodeID: "systems/water", Reason: "depends on systems/food"},
		},
	}

	output := formatter.Format(err)

	// Should contain related items
	if !strings.Contains(output, "systems/food") {
		t.Errorf("Output should contain first related NodeID, got:\n%s", output)
	}
	if !strings.Contains(output, "systems/water") {
		t.Errorf("Output should contain second related NodeID, got:\n%s", output)
	}
	if !strings.Contains(output, "depends on systems/water") {
		t.Errorf("Output should contain first related reason, got:\n%s", output)
	}
}

func TestErrorFormatter_ColorOutput(t *testing.T) {
	formatter := domain.NewErrorFormatter()
	formatter.SetColor(true)

	err := domain.DecoError{
		Code:    "E001",
		Summary: "Test error",
	}

	output := formatter.Format(err)

	// Should contain ANSI color codes when color is enabled
	// ANSI codes start with \x1b[ or \033[
	if !strings.Contains(output, "\x1b[") && !strings.Contains(output, "\033[") {
		t.Errorf("Output should contain ANSI color codes when color is enabled, got:\n%s", output)
	}
}

func TestErrorFormatter_NoColorOutput(t *testing.T) {
	formatter := domain.NewErrorFormatter()
	formatter.SetColor(false)

	err := domain.DecoError{
		Code:    "E001",
		Summary: "Test error",
	}

	output := formatter.Format(err)

	// Should not contain ANSI color codes when color is disabled
	if strings.Contains(output, "\x1b[") || strings.Contains(output, "\033[") {
		t.Errorf("Output should not contain ANSI color codes when color is disabled, got:\n%s", output)
	}
}

func TestErrorFormatter_WithSourceContext(t *testing.T) {
	formatter := domain.NewErrorFormatter()

	sourceLines := []string{
		"id: systems/food",
		"type: system",
		"refs:",
		"  uses: [unknown/node]",
		"metadata:",
		"  version: 1.0",
	}

	err := domain.DecoError{
		Code:     "E002",
		Summary:  "Invalid reference",
		Location: &domain.Location{File: "test.yaml", Line: 4, Column: 10},
	}

	output := formatter.FormatWithSource(err, sourceLines)

	// Should show lines around the error location
	if !strings.Contains(output, "uses: [unknown/node]") {
		t.Errorf("Output should contain the error line, got:\n%s", output)
	}

	// Should show context lines before and after
	if !strings.Contains(output, "type: system") {
		t.Errorf("Output should contain line before error, got:\n%s", output)
	}
}

func TestErrorFormatter_LineNumbersInSource(t *testing.T) {
	formatter := domain.NewErrorFormatter()

	sourceLines := []string{
		"id: systems/food",
		"type: system",
		"refs:",
		"  uses: [unknown/node]",
	}

	err := domain.DecoError{
		Code:     "E002",
		Summary:  "Invalid reference",
		Location: &domain.Location{File: "test.yaml", Line: 4},
	}

	output := formatter.FormatWithSource(err, sourceLines)

	// Should include line numbers in output
	if !strings.Contains(output, "4") {
		t.Errorf("Output should contain line number, got:\n%s", output)
	}
}

func TestErrorFormatter_ColumnPointer(t *testing.T) {
	formatter := domain.NewErrorFormatter()

	sourceLines := []string{
		"refs:",
		"  uses: [unknown/node]",
	}

	err := domain.DecoError{
		Code:     "E002",
		Summary:  "Invalid reference",
		Location: &domain.Location{File: "test.yaml", Line: 2, Column: 10},
	}

	output := formatter.FormatWithSource(err, sourceLines)

	// Should include a pointer to the column (like Rust's ^)
	// Look for some form of indicator at the column position
	lines := strings.Split(output, "\n")
	hasPointer := false
	for _, line := range lines {
		if strings.Contains(line, "^") || strings.Contains(line, "~") {
			hasPointer = true
			break
		}
	}

	if !hasPointer {
		t.Errorf("Output should contain a column pointer (^ or ~), got:\n%s", output)
	}
}

func TestErrorFormatter_MultipleErrors(t *testing.T) {
	formatter := domain.NewErrorFormatter()

	errors := []domain.DecoError{
		{Code: "E001", Summary: "First error"},
		{Code: "E002", Summary: "Second error"},
	}

	output := formatter.FormatMultiple(errors)

	// Should contain both errors
	if !strings.Contains(output, "First error") {
		t.Errorf("Output should contain first error, got:\n%s", output)
	}
	if !strings.Contains(output, "Second error") {
		t.Errorf("Output should contain second error, got:\n%s", output)
	}

	// Should contain error count
	if !strings.Contains(output, "2") {
		t.Errorf("Output should contain error count, got:\n%s", output)
	}
}

func TestErrorFormatter_EmptyError(t *testing.T) {
	formatter := domain.NewErrorFormatter()
	err := domain.DecoError{}

	output := formatter.Format(err)

	// Should handle empty error gracefully
	if output == "" {
		t.Error("Expected non-empty output even for empty error")
	}
}
