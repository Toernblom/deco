package errors_test

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
)

// Test basic error collection
func TestCollector_AddSingleError(t *testing.T) {
	collector := errors.NewCollector()

	err := domain.DecoError{
		Code:    "E001",
		Summary: "Node not found",
		Detail:  "The node 'systems/food' does not exist",
		Location: &domain.Location{
			File: "nodes/systems/food.yaml",
			Line: 10,
		},
	}

	collector.Add(err)

	collected := collector.Errors()
	if len(collected) != 1 {
		t.Fatalf("expected 1 error, got %d", len(collected))
	}

	if collected[0].Code != "E001" {
		t.Errorf("expected Code 'E001', got %q", collected[0].Code)
	}
}

// Test multiple errors are collected
func TestCollector_AddMultipleErrors(t *testing.T) {
	collector := errors.NewCollector()

	err1 := domain.DecoError{
		Code:     "E001",
		Summary:  "First error",
		Location: &domain.Location{File: "file1.yaml", Line: 1},
	}
	err2 := domain.DecoError{
		Code:     "E002",
		Summary:  "Second error",
		Location: &domain.Location{File: "file2.yaml", Line: 5},
	}
	err3 := domain.DecoError{
		Code:     "E003",
		Summary:  "Third error",
		Location: &domain.Location{File: "file3.yaml", Line: 10},
	}

	collector.Add(err1)
	collector.Add(err2)
	collector.Add(err3)

	collected := collector.Errors()
	if len(collected) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(collected))
	}
}

// Test errors are sorted by file path
func TestCollector_SortByFile(t *testing.T) {
	collector := errors.NewCollector()

	// Add errors in reverse alphabetical order
	collector.Add(domain.DecoError{
		Code:     "E003",
		Summary:  "Error in file3",
		Location: &domain.Location{File: "nodes/z.yaml", Line: 1},
	})
	collector.Add(domain.DecoError{
		Code:     "E002",
		Summary:  "Error in file2",
		Location: &domain.Location{File: "nodes/m.yaml", Line: 1},
	})
	collector.Add(domain.DecoError{
		Code:     "E001",
		Summary:  "Error in file1",
		Location: &domain.Location{File: "nodes/a.yaml", Line: 1},
	})

	collected := collector.Errors()
	if len(collected) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(collected))
	}

	// Should be sorted alphabetically by file
	if collected[0].Location.File != "nodes/a.yaml" {
		t.Errorf("expected first error in 'nodes/a.yaml', got %q", collected[0].Location.File)
	}
	if collected[1].Location.File != "nodes/m.yaml" {
		t.Errorf("expected second error in 'nodes/m.yaml', got %q", collected[1].Location.File)
	}
	if collected[2].Location.File != "nodes/z.yaml" {
		t.Errorf("expected third error in 'nodes/z.yaml', got %q", collected[2].Location.File)
	}
}

// Test errors in the same file are sorted by line number
func TestCollector_SortByLine(t *testing.T) {
	collector := errors.NewCollector()

	// Add errors in reverse line number order
	collector.Add(domain.DecoError{
		Code:     "E003",
		Summary:  "Error at line 30",
		Location: &domain.Location{File: "nodes/test.yaml", Line: 30},
	})
	collector.Add(domain.DecoError{
		Code:     "E002",
		Summary:  "Error at line 20",
		Location: &domain.Location{File: "nodes/test.yaml", Line: 20},
	})
	collector.Add(domain.DecoError{
		Code:     "E001",
		Summary:  "Error at line 10",
		Location: &domain.Location{File: "nodes/test.yaml", Line: 10},
	})

	collected := collector.Errors()
	if len(collected) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(collected))
	}

	// Should be sorted by line number
	if collected[0].Location.Line != 10 {
		t.Errorf("expected first error at line 10, got %d", collected[0].Location.Line)
	}
	if collected[1].Location.Line != 20 {
		t.Errorf("expected second error at line 20, got %d", collected[1].Location.Line)
	}
	if collected[2].Location.Line != 30 {
		t.Errorf("expected third error at line 30, got %d", collected[2].Location.Line)
	}
}

// Test sorting by file first, then line
func TestCollector_SortByFileAndLine(t *testing.T) {
	collector := errors.NewCollector()

	// Add errors in mixed order
	collector.Add(domain.DecoError{
		Summary:  "b.yaml:20",
		Location: &domain.Location{File: "b.yaml", Line: 20},
	})
	collector.Add(domain.DecoError{
		Summary:  "a.yaml:30",
		Location: &domain.Location{File: "a.yaml", Line: 30},
	})
	collector.Add(domain.DecoError{
		Summary:  "b.yaml:10",
		Location: &domain.Location{File: "b.yaml", Line: 10},
	})
	collector.Add(domain.DecoError{
		Summary:  "a.yaml:5",
		Location: &domain.Location{File: "a.yaml", Line: 5},
	})

	collected := collector.Errors()
	if len(collected) != 4 {
		t.Fatalf("expected 4 errors, got %d", len(collected))
	}

	// Should be sorted by file first, then line
	if collected[0].Summary != "a.yaml:5" {
		t.Errorf("expected 'a.yaml:5', got %q", collected[0].Summary)
	}
	if collected[1].Summary != "a.yaml:30" {
		t.Errorf("expected 'a.yaml:30', got %q", collected[1].Summary)
	}
	if collected[2].Summary != "b.yaml:10" {
		t.Errorf("expected 'b.yaml:10', got %q", collected[2].Summary)
	}
	if collected[3].Summary != "b.yaml:20" {
		t.Errorf("expected 'b.yaml:20', got %q", collected[3].Summary)
	}
}

// Test errors without location come last
func TestCollector_ErrorsWithoutLocation(t *testing.T) {
	collector := errors.NewCollector()

	collector.Add(domain.DecoError{
		Summary:  "With location",
		Location: &domain.Location{File: "test.yaml", Line: 10},
	})
	collector.Add(domain.DecoError{
		Summary:  "Without location",
		Location: nil,
	})

	collected := collector.Errors()
	if len(collected) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(collected))
	}

	// Error with location should come first
	if collected[0].Location == nil {
		t.Error("expected first error to have location")
	}
	if collected[1].Location != nil {
		t.Error("expected second error to have no location")
	}
}

// Test deduplication of exact duplicate errors
func TestCollector_DeduplicateExact(t *testing.T) {
	collector := errors.NewCollector()

	err := domain.DecoError{
		Code:     "E001",
		Summary:  "Duplicate error",
		Detail:   "This is a duplicate",
		Location: &domain.Location{File: "test.yaml", Line: 10, Column: 5},
	}

	// Add the same error three times
	collector.Add(err)
	collector.Add(err)
	collector.Add(err)

	collected := collector.Errors()
	// Should only keep one copy
	if len(collected) != 1 {
		t.Errorf("expected deduplication to leave 1 error, got %d", len(collected))
	}
}

// Test similar errors with same location are deduplicated
func TestCollector_DeduplicateSameLocation(t *testing.T) {
	collector := errors.NewCollector()

	// Same code and location, slightly different detail
	collector.Add(domain.DecoError{
		Code:     "E001",
		Summary:  "Node not found",
		Detail:   "First detail",
		Location: &domain.Location{File: "test.yaml", Line: 10},
	})
	collector.Add(domain.DecoError{
		Code:     "E001",
		Summary:  "Node not found",
		Detail:   "Second detail",
		Location: &domain.Location{File: "test.yaml", Line: 10},
	})

	collected := collector.Errors()
	// Should keep only one (same code + location = likely same root cause)
	if len(collected) != 1 {
		t.Errorf("expected deduplication to leave 1 error, got %d", len(collected))
	}
}

// Test different errors at same location are kept
func TestCollector_DifferentErrorsSameLocation(t *testing.T) {
	collector := errors.NewCollector()

	// Different error codes at same location
	collector.Add(domain.DecoError{
		Code:     "E001",
		Summary:  "First type of error",
		Location: &domain.Location{File: "test.yaml", Line: 10},
	})
	collector.Add(domain.DecoError{
		Code:     "E002",
		Summary:  "Second type of error",
		Location: &domain.Location{File: "test.yaml", Line: 10},
	})

	collected := collector.Errors()
	// Should keep both (different error codes)
	if len(collected) != 2 {
		t.Errorf("expected 2 different errors, got %d", len(collected))
	}
}

// Test max errors limit
func TestCollector_MaxErrorsLimit(t *testing.T) {
	collector := errors.NewCollectorWithLimit(3)

	// Add 5 errors, but limit is 3
	for i := 1; i <= 5; i++ {
		collector.Add(domain.DecoError{
			Code:     "E001",
			Summary:  "Error",
			Location: &domain.Location{File: "test.yaml", Line: i},
		})
	}

	collected := collector.Errors()
	// Should only return 3 errors
	if len(collected) != 3 {
		t.Errorf("expected max 3 errors, got %d", len(collected))
	}

	// Should still track total count
	if collector.Count() != 5 {
		t.Errorf("expected Count() to return 5, got %d", collector.Count())
	}

	// Should know it was truncated
	if !collector.Truncated() {
		t.Error("expected Truncated() to return true")
	}
}

// Test default limit (no limit)
func TestCollector_NoLimit(t *testing.T) {
	collector := errors.NewCollector()

	// Add many errors
	for i := 1; i <= 100; i++ {
		collector.Add(domain.DecoError{
			Code:     "E001",
			Summary:  "Error",
			Location: &domain.Location{File: "test.yaml", Line: i},
		})
	}

	collected := collector.Errors()
	// Should keep all of them
	if len(collected) != 100 {
		t.Errorf("expected 100 errors, got %d", len(collected))
	}

	if collector.Truncated() {
		t.Error("expected Truncated() to return false when no limit")
	}
}

// Test HasErrors method
func TestCollector_HasErrors(t *testing.T) {
	collector := errors.NewCollector()

	if collector.HasErrors() {
		t.Error("expected HasErrors() to return false for empty collector")
	}

	collector.Add(domain.DecoError{
		Code:    "E001",
		Summary: "Test error",
	})

	if !collector.HasErrors() {
		t.Error("expected HasErrors() to return true after adding error")
	}
}

// Test Count method
func TestCollector_Count(t *testing.T) {
	collector := errors.NewCollector()

	if collector.Count() != 0 {
		t.Errorf("expected Count() to return 0 for empty collector, got %d", collector.Count())
	}

	collector.Add(domain.DecoError{Code: "E001"})
	collector.Add(domain.DecoError{Code: "E002"})

	if collector.Count() != 2 {
		t.Errorf("expected Count() to return 2, got %d", collector.Count())
	}
}

// Test Reset method
func TestCollector_Reset(t *testing.T) {
	collector := errors.NewCollector()

	collector.Add(domain.DecoError{Code: "E001"})
	collector.Add(domain.DecoError{Code: "E002"})

	if collector.Count() != 2 {
		t.Fatalf("expected 2 errors before reset, got %d", collector.Count())
	}

	collector.Reset()

	if collector.HasErrors() {
		t.Error("expected no errors after reset")
	}
	if collector.Count() != 0 {
		t.Errorf("expected Count() to return 0 after reset, got %d", collector.Count())
	}
}

// Test sorting stability with column numbers
func TestCollector_SortByColumn(t *testing.T) {
	collector := errors.NewCollector()

	// Same file and line, different columns
	collector.Add(domain.DecoError{
		Summary:  "col:30",
		Location: &domain.Location{File: "test.yaml", Line: 10, Column: 30},
	})
	collector.Add(domain.DecoError{
		Summary:  "col:10",
		Location: &domain.Location{File: "test.yaml", Line: 10, Column: 10},
	})
	collector.Add(domain.DecoError{
		Summary:  "col:20",
		Location: &domain.Location{File: "test.yaml", Line: 10, Column: 20},
	})

	collected := collector.Errors()
	if len(collected) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(collected))
	}

	// Should be sorted by column within same file and line
	if collected[0].Location.Column != 10 {
		t.Errorf("expected first error at column 10, got %d", collected[0].Location.Column)
	}
	if collected[1].Location.Column != 20 {
		t.Errorf("expected second error at column 20, got %d", collected[1].Location.Column)
	}
	if collected[2].Location.Column != 30 {
		t.Errorf("expected third error at column 30, got %d", collected[2].Location.Column)
	}
}

// Test deduplication with AddBatch
func TestCollector_AddBatch(t *testing.T) {
	collector := errors.NewCollector()

	errors := []domain.DecoError{
		{Code: "E001", Location: &domain.Location{File: "a.yaml", Line: 1}},
		{Code: "E002", Location: &domain.Location{File: "b.yaml", Line: 2}},
		{Code: "E003", Location: &domain.Location{File: "c.yaml", Line: 3}},
	}

	collector.AddBatch(errors)

	collected := collector.Errors()
	if len(collected) != 3 {
		t.Errorf("expected 3 errors from batch, got %d", len(collected))
	}
}

// Test empty collector
func TestCollector_Empty(t *testing.T) {
	collector := errors.NewCollector()

	if collector.HasErrors() {
		t.Error("expected empty collector to have no errors")
	}

	errors := collector.Errors()
	if errors == nil {
		t.Error("expected Errors() to return non-nil slice")
	}
	if len(errors) != 0 {
		t.Errorf("expected empty slice, got %d errors", len(errors))
	}
}
