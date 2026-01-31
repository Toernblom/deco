package domain_test

import (
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestDecoError_Creation(t *testing.T) {
	err := domain.DecoError{
		Code:       "E001",
		Summary:    "Node not found",
		Detail:     "The node 'systems/food' does not exist",
		Location:   &domain.Location{File: "nodes/systems/food.yaml", Line: 1},
		Context:    []string{"while loading node references"},
		Suggestion: "Check that the file exists in .deco/nodes/systems/",
		Related:    []domain.Related{{NodeID: "systems/health", Reason: "references this node"}},
	}

	if err.Code != "E001" {
		t.Errorf("expected Code 'E001', got %q", err.Code)
	}
	if err.Summary != "Node not found" {
		t.Errorf("expected Summary 'Node not found', got %q", err.Summary)
	}
	if err.Detail != "The node 'systems/food' does not exist" {
		t.Errorf("expected Detail to match, got %q", err.Detail)
	}
	if err.Location == nil || err.Location.File != "nodes/systems/food.yaml" {
		t.Errorf("expected Location.File 'nodes/systems/food.yaml', got %v", err.Location)
	}
	if len(err.Context) != 1 {
		t.Errorf("expected 1 context item, got %d", len(err.Context))
	}
	if err.Suggestion != "Check that the file exists in .deco/nodes/systems/" {
		t.Errorf("expected Suggestion to match, got %q", err.Suggestion)
	}
	if len(err.Related) != 1 {
		t.Errorf("expected 1 related item, got %d", len(err.Related))
	}
}

func TestDecoError_AllFieldsPopulated(t *testing.T) {
	err := domain.DecoError{
		Code:    "E002",
		Summary: "Invalid reference",
		Detail:  "Reference to 'unknown/node' cannot be resolved",
		Location: &domain.Location{
			File:   ".deco/nodes/mechanics/combat.yaml",
			Line:   15,
			Column: 8,
		},
		Context: []string{
			"in node 'mechanics/combat'",
			"in field 'refs.uses[0]'",
		},
		Suggestion: "Ensure the referenced node exists and is spelled correctly",
		Related: []domain.Related{
			{NodeID: "mechanics/combat", Reason: "contains invalid reference"},
		},
	}

	if err.Code == "" {
		t.Error("Code should be populated")
	}
	if err.Summary == "" {
		t.Error("Summary should be populated")
	}
	if err.Detail == "" {
		t.Error("Detail should be populated")
	}
	if err.Location == nil {
		t.Error("Location should be populated")
	}
	if err.Location.Line != 15 {
		t.Errorf("expected Location.Line 15, got %d", err.Location.Line)
	}
	if err.Location.Column != 8 {
		t.Errorf("expected Location.Column 8, got %d", err.Location.Column)
	}
	if len(err.Context) != 2 {
		t.Errorf("expected 2 context items, got %d", len(err.Context))
	}
	if err.Suggestion == "" {
		t.Error("Suggestion should be populated")
	}
	if len(err.Related) != 1 {
		t.Errorf("expected 1 related item, got %d", len(err.Related))
	}
}

func TestLocation_Type(t *testing.T) {
	tests := []struct {
		name     string
		location domain.Location
	}{
		{
			name: "file with line",
			location: domain.Location{
				File: "nodes/systems/food.yaml",
				Line: 42,
			},
		},
		{
			name: "file with line and column",
			location: domain.Location{
				File:   "nodes/mechanics/combat.yaml",
				Line:   15,
				Column: 8,
			},
		},
		{
			name: "file only",
			location: domain.Location{
				File: "config.yaml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.location.File == "" {
				t.Error("Location.File should not be empty")
			}
		})
	}
}

func TestRelated_Type(t *testing.T) {
	related := domain.Related{
		NodeID: "systems/health",
		Reason: "depends on this system",
	}

	if related.NodeID != "systems/health" {
		t.Errorf("expected NodeID 'systems/health', got %q", related.NodeID)
	}
	if related.Reason != "depends on this system" {
		t.Errorf("expected Reason 'depends on this system', got %q", related.Reason)
	}
}

func TestDecoError_ErrorMethod(t *testing.T) {
	err := domain.DecoError{
		Code:    "E001",
		Summary: "Node not found",
		Detail:  "The node 'systems/food' does not exist",
	}

	errMsg := err.Error()

	// Error() should return a string containing key information
	if errMsg == "" {
		t.Error("Error() should return non-empty string")
	}

	// Should contain the code
	if !strings.Contains(errMsg, "E001") {
		t.Errorf("Error() should contain code 'E001', got %q", errMsg)
	}

	// Should contain the summary
	if !strings.Contains(errMsg, "Node not found") {
		t.Errorf("Error() should contain summary, got %q", errMsg)
	}
}

func TestDecoError_ErrorMethodWithLocation(t *testing.T) {
	err := domain.DecoError{
		Code:     "E002",
		Summary:  "Invalid reference",
		Detail:   "Reference to 'unknown/node' cannot be resolved",
		Location: &domain.Location{File: "test.yaml", Line: 10, Column: 5},
	}

	errMsg := err.Error()

	// Should contain file location
	if !strings.Contains(errMsg, "test.yaml") {
		t.Errorf("Error() should contain file location, got %q", errMsg)
	}

	// Should contain line number
	if !strings.Contains(errMsg, "10") {
		t.Errorf("Error() should contain line number, got %q", errMsg)
	}
}

func TestDecoError_ErrorMethodWithContext(t *testing.T) {
	err := domain.DecoError{
		Code:    "E003",
		Summary: "Validation failed",
		Detail:  "Field 'status' has invalid value",
		Context: []string{"in node 'systems/food'", "during validation"},
	}

	errMsg := err.Error()

	// Context should be included in error message
	if !strings.Contains(errMsg, "in node 'systems/food'") {
		t.Errorf("Error() should include context, got %q", errMsg)
	}
}

func TestDecoError_MultipleRelated(t *testing.T) {
	err := domain.DecoError{
		Code:    "E004",
		Summary: "Circular dependency",
		Detail:  "Nodes form a dependency cycle",
		Related: []domain.Related{
			{NodeID: "systems/food", Reason: "depends on systems/water"},
			{NodeID: "systems/water", Reason: "depends on systems/food"},
		},
	}

	if len(err.Related) != 2 {
		t.Errorf("expected 2 related items, got %d", len(err.Related))
	}

	// Verify both related items
	if err.Related[0].NodeID != "systems/food" {
		t.Errorf("expected first related NodeID 'systems/food', got %q", err.Related[0].NodeID)
	}
	if err.Related[1].NodeID != "systems/water" {
		t.Errorf("expected second related NodeID 'systems/water', got %q", err.Related[1].NodeID)
	}
}

func TestDecoError_OptionalFields(t *testing.T) {
	// Error with minimal fields
	err := domain.DecoError{
		Code:    "E005",
		Summary: "Simple error",
	}

	if err.Detail != "" {
		t.Error("Detail should be empty when not set")
	}
	if err.Location != nil {
		t.Error("Location should be nil when not set")
	}
	if len(err.Context) != 0 {
		t.Error("Context should be empty when not set")
	}
	if err.Suggestion != "" {
		t.Error("Suggestion should be empty when not set")
	}
	if len(err.Related) != 0 {
		t.Error("Related should be empty when not set")
	}
}

func TestLocation_String(t *testing.T) {
	tests := []struct {
		name     string
		location domain.Location
		want     string
	}{
		{
			name:     "file with line and column",
			location: domain.Location{File: "test.yaml", Line: 10, Column: 5},
			want:     "test.yaml:10:5",
		},
		{
			name:     "file with line only",
			location: domain.Location{File: "test.yaml", Line: 10},
			want:     "test.yaml:10",
		},
		{
			name:     "file only",
			location: domain.Location{File: "test.yaml"},
			want:     "test.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.location.String()
			if got != tt.want {
				t.Errorf("Location.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
