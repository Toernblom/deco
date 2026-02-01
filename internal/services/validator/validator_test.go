package validator_test

import (
	"testing"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
	"github.com/Toernblom/deco/internal/services/validator"
)

// ===== SCHEMA VALIDATOR TESTS =====

// Test validating a node with all required fields
func TestSchemaValidator_ValidNode(t *testing.T) {
	sv := validator.NewSchemaValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
	}

	collector := errors.NewCollectorWithLimit(100)
	sv.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors, got %d", collector.Count())
	}
}

// Test missing ID field
func TestSchemaValidator_MissingID(t *testing.T) {
	sv := validator.NewSchemaValidator()

	node := domain.Node{
		// ID missing
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
	}

	collector := errors.NewCollectorWithLimit(100)
	sv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for missing ID")
	}

	errs := collector.Errors()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}

	if errs[0].Code != "E008" {
		t.Errorf("expected error code E008, got %s", errs[0].Code)
	}
}

// Test missing Kind field
func TestSchemaValidator_MissingKind(t *testing.T) {
	sv := validator.NewSchemaValidator()

	node := domain.Node{
		ID: "test-node",
		// Kind missing
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
	}

	collector := errors.NewCollectorWithLimit(100)
	sv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for missing Kind")
	}

	errs := collector.Errors()
	if errs[0].Code != "E008" {
		t.Errorf("expected error code E008, got %s", errs[0].Code)
	}
}

// Test missing Version (zero value)
func TestSchemaValidator_MissingVersion(t *testing.T) {
	sv := validator.NewSchemaValidator()

	node := domain.Node{
		ID:   "test-node",
		Kind: "system",
		// Version is 0 (missing)
		Status: "draft",
		Title:  "Test Node",
	}

	collector := errors.NewCollectorWithLimit(100)
	sv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for missing Version")
	}

	errs := collector.Errors()
	if errs[0].Code != "E008" {
		t.Errorf("expected error code E008, got %s", errs[0].Code)
	}
}

// Test missing Status field
func TestSchemaValidator_MissingStatus(t *testing.T) {
	sv := validator.NewSchemaValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		// Status missing
		Title: "Test Node",
	}

	collector := errors.NewCollectorWithLimit(100)
	sv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for missing Status")
	}

	errs := collector.Errors()
	if errs[0].Code != "E008" {
		t.Errorf("expected error code E008, got %s", errs[0].Code)
	}
}

// Test missing Title field
func TestSchemaValidator_MissingTitle(t *testing.T) {
	sv := validator.NewSchemaValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		// Title missing
	}

	collector := errors.NewCollectorWithLimit(100)
	sv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for missing Title")
	}

	errs := collector.Errors()
	if errs[0].Code != "E008" {
		t.Errorf("expected error code E008, got %s", errs[0].Code)
	}
}

// Test multiple missing fields
func TestSchemaValidator_MultipleMissingFields(t *testing.T) {
	sv := validator.NewSchemaValidator()

	node := domain.Node{
		// ID, Kind, Title missing
		Version: 1,
		Status:  "draft",
	}

	collector := errors.NewCollectorWithLimit(100)
	sv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected errors for multiple missing fields")
	}

	// Should have at least 3 errors (ID, Kind, Title)
	if collector.Count() < 3 {
		t.Errorf("expected at least 3 errors, got %d", collector.Count())
	}
}

// Test nil node
func TestSchemaValidator_NilNode(t *testing.T) {
	sv := validator.NewSchemaValidator()

	collector := errors.NewCollectorWithLimit(100)
	sv.Validate(nil, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for nil node")
	}
}

// ===== REFERENCE VALIDATOR TESTS =====

// Test valid references
func TestReferenceValidator_ValidReferences(t *testing.T) {
	rv := validator.NewReferenceValidator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "node1"}}}},
	}

	collector := errors.NewCollectorWithLimit(100)
	rv.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid references, got %d", collector.Count())
	}
}

// Test broken reference (not found)
func TestReferenceValidator_BrokenReference(t *testing.T) {
	rv := validator.NewReferenceValidator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "nonexistent"}}}},
	}

	collector := errors.NewCollectorWithLimit(100)
	rv.Validate(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for broken reference")
	}

	errs := collector.Errors()
	if errs[0].Code != "E020" {
		t.Errorf("expected error code E020, got %s", errs[0].Code)
	}
}

// Test suggestion generation for typos
func TestReferenceValidator_SuggestionForTypo(t *testing.T) {
	rv := validator.NewReferenceValidator()

	nodes := []domain.Node{
		{ID: "combat-system", Kind: "system", Version: 1, Status: "draft", Title: "Combat System"},
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "combat-systm"}}}}, // typo: missing 'e'
	}

	collector := errors.NewCollectorWithLimit(100)
	rv.Validate(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for broken reference")
	}

	errs := collector.Errors()
	if errs[0].Suggestion == "" {
		t.Error("expected suggestion for similar ID, got none")
	}
}

// Test multiple broken references
func TestReferenceValidator_MultipleBrokenReferences(t *testing.T) {
	rv := validator.NewReferenceValidator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "missing1"}, {Target: "missing2"}}}},
	}

	collector := errors.NewCollectorWithLimit(100)
	rv.Validate(nodes, collector)

	if collector.Count() < 2 {
		t.Errorf("expected at least 2 errors, got %d", collector.Count())
	}
}

// Test Related references
func TestReferenceValidator_RelatedReferences(t *testing.T) {
	rv := validator.NewReferenceValidator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2",
			Refs: domain.Ref{Related: []domain.RefLink{{Target: "node1"}}}},
	}

	collector := errors.NewCollectorWithLimit(100)
	rv.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid Related references, got %d", collector.Count())
	}
}

// Test broken Related reference
func TestReferenceValidator_BrokenRelatedReference(t *testing.T) {
	rv := validator.NewReferenceValidator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1",
			Refs: domain.Ref{Related: []domain.RefLink{{Target: "missing"}}}},
	}

	collector := errors.NewCollectorWithLimit(100)
	rv.Validate(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for broken Related reference")
	}
}

// Test empty references (should pass)
func TestReferenceValidator_EmptyReferences(t *testing.T) {
	rv := validator.NewReferenceValidator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1",
			Refs: domain.Ref{}},
	}

	collector := errors.NewCollectorWithLimit(100)
	rv.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for empty references, got %d", collector.Count())
	}
}

// ===== CONSTRAINT VALIDATOR TESTS =====

// Test passing constraint
func TestConstraintValidator_PassingConstraint(t *testing.T) {
	cv := validator.NewConstraintValidator()

	node := domain.Node{
		ID:      "node1",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Constraints: []domain.Constraint{
			{
				Expr:    "version > 0",
				Message: "Version must be positive",
				Scope:   "all",
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, []domain.Node{node}, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for passing constraint, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test failing constraint
func TestConstraintValidator_FailingConstraint(t *testing.T) {
	cv := validator.NewConstraintValidator()

	node := domain.Node{
		ID:      "node1",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Constraints: []domain.Constraint{
			{
				Expr:    "version > 5",
				Message: "Version must be greater than 5",
				Scope:   "all",
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, []domain.Node{node}, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for failing constraint")
	}

	errs := collector.Errors()
	if errs[0].Code != "E041" {
		t.Errorf("expected error code E041, got %s", errs[0].Code)
	}
}

// Test constraint with invalid CEL expression
func TestConstraintValidator_InvalidCELExpression(t *testing.T) {
	cv := validator.NewConstraintValidator()

	node := domain.Node{
		ID:      "node1",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Constraints: []domain.Constraint{
			{
				Expr:    "invalid syntax !!!",
				Message: "Invalid expression",
				Scope:   "all",
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, []domain.Node{node}, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for invalid CEL expression")
	}

	errs := collector.Errors()
	if errs[0].Code != "E042" {
		t.Errorf("expected error code E042, got %s", errs[0].Code)
	}
}

// Test multiple constraints
func TestConstraintValidator_MultipleConstraints(t *testing.T) {
	cv := validator.NewConstraintValidator()

	node := domain.Node{
		ID:      "node1",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Tags:    []string{"gameplay"},
		Constraints: []domain.Constraint{
			{
				Expr:    "version > 0",
				Message: "Version must be positive",
				Scope:   "all",
			},
			{
				Expr:    `"gameplay" in tags`,
				Message: "Must have gameplay tag",
				Scope:   "all",
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, []domain.Node{node}, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test node with no constraints (should pass)
func TestConstraintValidator_NoConstraints(t *testing.T) {
	cv := validator.NewConstraintValidator()

	node := domain.Node{
		ID:          "node1",
		Kind:        "system",
		Version:     1,
		Status:      "draft",
		Title:       "Test Node",
		Constraints: []domain.Constraint{},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, []domain.Node{node}, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for node without constraints, got %d", collector.Count())
	}
}

func TestConstraintValidator_ScopeMatching(t *testing.T) {
	tests := []struct {
		name          string
		nodeID        string
		nodeKind      string
		scope         string
		shouldApply   bool
		expectedErrs  int // 0 if scope doesn't match, 1 if it does and constraint fails
	}{
		{
			name:         "all scope matches any node",
			nodeID:       "test-node",
			nodeKind:     "mechanic",
			scope:        "all",
			shouldApply:  true,
			expectedErrs: 1,
		},
		{
			name:         "empty scope treated as all",
			nodeID:       "test-node",
			nodeKind:     "system",
			scope:        "",
			shouldApply:  true,
			expectedErrs: 1,
		},
		{
			name:         "exact kind match",
			nodeID:       "test-node",
			nodeKind:     "mechanic",
			scope:        "mechanic",
			shouldApply:  true,
			expectedErrs: 1,
		},
		{
			name:         "kind mismatch skips constraint",
			nodeID:       "test-node",
			nodeKind:     "system",
			scope:        "mechanic",
			shouldApply:  false,
			expectedErrs: 0,
		},
		{
			name:         "glob pattern matches ID",
			nodeID:       "systems/combat/damage",
			nodeKind:     "system",
			scope:        "systems/combat/*",
			shouldApply:  true,
			expectedErrs: 1,
		},
		{
			name:         "glob pattern no match",
			nodeID:       "systems/ui/menu",
			nodeKind:     "system",
			scope:        "systems/combat/*",
			shouldApply:  false,
			expectedErrs: 0,
		},
		{
			name:         "question mark glob",
			nodeID:       "item-a",
			nodeKind:     "item",
			scope:        "item-?",
			shouldApply:  true,
			expectedErrs: 1,
		},
		{
			name:         "pattern with multiple slashes",
			nodeID:       "systems/combat/damage",
			nodeKind:     "system",
			scope:        "systems/*/*",
			shouldApply:  true, // filepath.Match supports single * per path segment
			expectedErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := validator.NewConstraintValidator()

			// Create a constraint that will always fail (version < 0 is always false for positive version)
			constraint := domain.Constraint{
				Expr:    "version < 0", // Always fails for valid nodes
				Message: "This constraint always fails",
				Scope:   tt.scope,
			}

			node := domain.Node{
				ID:          tt.nodeID,
				Kind:        tt.nodeKind,
				Version:     1,
				Status:      "draft",
				Title:       "Test Node",
				Constraints: []domain.Constraint{constraint},
			}

			collector := errors.NewCollectorWithLimit(100)
			cv.Validate(&node, []domain.Node{node}, collector)

			if collector.Count() != tt.expectedErrs {
				t.Errorf("expected %d errors, got %d (shouldApply=%v)", tt.expectedErrs, collector.Count(), tt.shouldApply)
			}
		})
	}
}

func TestConstraintValidator_ScopeWithMultipleConstraints(t *testing.T) {
	cv := validator.NewConstraintValidator()

	// Create constraints with different scopes and different expressions
	// (to avoid deduplication based on code+summary)
	constraints := []domain.Constraint{
		{
			Expr:    "version < 0", // Always fails
			Message: "Mechanic constraint failed",
			Scope:   "mechanic",
		},
		{
			Expr:    "version < -1", // Always fails - different expression
			Message: "System constraint failed",
			Scope:   "system",
		},
		{
			Expr:    "version < -2", // Always fails - different expression
			Message: "All constraint failed",
			Scope:   "all",
		},
	}

	// System node should only fail on "system" and "all" scoped constraints
	node := domain.Node{
		ID:          "test-system",
		Kind:        "system",
		Version:     1,
		Status:      "draft",
		Title:       "Test System",
		Constraints: constraints,
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, []domain.Node{node}, collector)

	// Should have 2 errors: one from "system" scope, one from "all" scope
	if collector.Count() != 2 {
		t.Errorf("expected 2 errors (system + all scopes), got %d", collector.Count())
	}
}

// ===== VALIDATOR ORCHESTRATOR TESTS =====

// Test orchestrator runs all validators
func TestOrchestrator_RunsAllValidators(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
	}

	collector := orch.ValidateAll(nodes)

	// Should not have errors for valid nodes
	if collector.HasErrors() {
		t.Errorf("expected no errors for valid nodes, got %d", collector.Count())
	}
}

// Test orchestrator aggregates errors from multiple validators
func TestOrchestrator_AggregatesErrors(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{
		{
			// Missing required fields (schema error)
			ID:      "node1",
			Kind:    "system",
			Version: 1,
			// Missing Status and Title
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "nonexistent"}}}, // Broken ref (reference error)
		},
	}

	collector := orch.ValidateAll(nodes)

	// Should have errors from both schema and reference validators
	if !collector.HasErrors() {
		t.Fatal("expected errors from multiple validators")
	}

	if collector.Count() < 2 {
		t.Errorf("expected at least 2 errors (schema + reference), got %d", collector.Count())
	}
}

// Test orchestrator with valid nodes
func TestOrchestrator_ValidNodes(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node2", Kind: "feature", Version: 1, Status: "draft", Title: "Node 2",
			Refs: domain.Ref{Uses: []domain.RefLink{{Target: "node1"}}}},
	}

	collector := orch.ValidateAll(nodes)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid nodes, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test orchestrator with empty node list
func TestOrchestrator_EmptyNodes(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{}

	collector := orch.ValidateAll(nodes)

	if collector.HasErrors() {
		t.Errorf("expected no errors for empty node list, got %d", collector.Count())
	}
}

// Test orchestrator validates all nodes
func TestOrchestrator_ValidatesAllNodes(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node2", Kind: "feature", Version: 1}, // Missing Status and Title
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Node 3"},
	}

	collector := orch.ValidateAll(nodes)

	if !collector.HasErrors() {
		t.Fatal("expected errors for invalid node")
	}

	// Should detect errors in node2
	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E008" { // Missing required field
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find missing field error for node2")
	}
}

// ===== DUPLICATE ID VALIDATOR TESTS =====

// Test no duplicates (should pass)
func TestDuplicateIDValidator_NoDuplicates(t *testing.T) {
	dv := validator.NewDuplicateIDValidator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2"},
		{ID: "node3", Kind: "system", Version: 1, Status: "draft", Title: "Node 3"},
	}

	collector := errors.NewCollectorWithLimit(100)
	dv.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for unique IDs, got %d", collector.Count())
	}
}

// Test duplicate IDs detected
func TestDuplicateIDValidator_DuplicatesDetected(t *testing.T) {
	dv := validator.NewDuplicateIDValidator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2"},
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1 Duplicate"}, // Duplicate ID
	}

	collector := errors.NewCollectorWithLimit(100)
	dv.Validate(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for duplicate ID")
	}

	errs := collector.Errors()
	if errs[0].Code != "E009" {
		t.Errorf("expected error code E009, got %s", errs[0].Code)
	}
}

// Test multiple duplicates
func TestDuplicateIDValidator_MultipleDuplicates(t *testing.T) {
	dv := validator.NewDuplicateIDValidator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1 Dup 1"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2"},
		{ID: "node2", Kind: "system", Version: 1, Status: "draft", Title: "Node 2 Dup 1"},
	}

	collector := errors.NewCollectorWithLimit(100)
	dv.Validate(nodes, collector)

	// Should have 2 errors (1 for node1 duplicate, 1 for node2 duplicate)
	if collector.Count() != 2 {
		t.Errorf("expected 2 duplicate errors, got %d", collector.Count())
	}
}

// Test empty ID is skipped (handled by schema validator)
func TestDuplicateIDValidator_EmptyIDSkipped(t *testing.T) {
	dv := validator.NewDuplicateIDValidator()

	nodes := []domain.Node{
		{ID: "", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"}, // Empty ID
		{ID: "", Kind: "system", Version: 1, Status: "draft", Title: "Node 2"}, // Another empty ID
	}

	collector := errors.NewCollectorWithLimit(100)
	dv.Validate(nodes, collector)

	// Empty IDs should not trigger duplicate error (they're handled by schema validator)
	if collector.HasErrors() {
		t.Errorf("expected no duplicate error for empty IDs, got %d", collector.Count())
	}
}

// Test orchestrator includes duplicate ID validation
func TestOrchestrator_DetectsDuplicateIDs(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "node1", Kind: "feature", Version: 1, Status: "draft", Title: "Node 1 Duplicate"},
	}

	collector := orch.ValidateAll(nodes)

	if !collector.HasErrors() {
		t.Fatal("expected error for duplicate IDs")
	}

	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E009" { // Duplicate ID
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find duplicate ID error (E009)")
	}
}

// ===== UNKNOWN FIELD VALIDATOR TESTS =====

// Test no unknown fields
func TestUnknownFieldValidator_AllKnownFields(t *testing.T) {
	uf := validator.NewUnknownFieldValidator()

	keys := []string{"id", "kind", "version", "status", "title", "tags", "content", "custom"}

	collector := errors.NewCollectorWithLimit(100)
	uf.ValidateYAML("test-node", "", keys, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for known fields, got %d", collector.Count())
	}
}

// Test unknown field detected
func TestUnknownFieldValidator_UnknownFieldDetected(t *testing.T) {
	uf := validator.NewUnknownFieldValidator()

	keys := []string{"id", "kind", "version", "status", "title", "unknown_field"}

	collector := errors.NewCollectorWithLimit(100)
	uf.ValidateYAML("test-node", "", keys, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for unknown field")
	}

	errs := collector.Errors()
	if errs[0].Code != "E010" {
		t.Errorf("expected error code E010, got %s", errs[0].Code)
	}
	if errs[0].Suggestion == "" {
		t.Log("Note: no suggestion generated (expected for dissimilar field names)")
	}
}

// Test typo in field name suggests correct field
func TestUnknownFieldValidator_SuggestsCorrection(t *testing.T) {
	uf := validator.NewUnknownFieldValidator()

	// "contnt" is a typo for "content"
	keys := []string{"id", "kind", "version", "status", "title", "contnt"}

	collector := errors.NewCollectorWithLimit(100)
	uf.ValidateYAML("test-node", "", keys, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for typo field")
	}

	errs := collector.Errors()
	if errs[0].Suggestion == "" {
		t.Error("expected suggestion for typo")
	}
}

// Test custom field is allowed (extension namespace)
func TestUnknownFieldValidator_CustomAllowed(t *testing.T) {
	uf := validator.NewUnknownFieldValidator()

	keys := []string{"id", "kind", "version", "status", "title", "custom"}

	collector := errors.NewCollectorWithLimit(100)
	uf.ValidateYAML("test-node", "", keys, collector)

	if collector.HasErrors() {
		t.Errorf("expected custom field to be allowed, got %d errors", collector.Count())
	}
}

// Test all known fields are accepted
func TestUnknownFieldValidator_AllSchemaFields(t *testing.T) {
	uf := validator.NewUnknownFieldValidator()

	// All valid top-level fields from the Node struct
	keys := []string{
		"id", "kind", "version", "status", "title", "tags",
		"refs", "content", "issues", "summary", "glossary",
		"contracts", "llm_context", "constraints", "custom",
	}

	collector := errors.NewCollectorWithLimit(100)
	uf.ValidateYAML("test-node", "", keys, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for all schema fields, got %d", collector.Count())
	}
}

// Test multiple unknown fields
func TestUnknownFieldValidator_MultipleUnknownFields(t *testing.T) {
	uf := validator.NewUnknownFieldValidator()

	keys := []string{"id", "kind", "foo", "bar", "baz"}

	collector := errors.NewCollectorWithLimit(100)
	uf.ValidateYAML("test-node", "", keys, collector)

	// Should have 3 errors (foo, bar, baz)
	if collector.Count() != 3 {
		t.Errorf("expected 3 unknown field errors, got %d", collector.Count())
	}
}

// ===== CONTRACT VALIDATOR TESTS =====

// Test valid contract with all steps
func TestContractValidator_ValidContract(t *testing.T) {
	cv := validator.NewContractValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Contracts: []domain.Contract{
			{
				Name:     "Basic Flow",
				Scenario: "Test scenario description",
				Given:    []string{"the player is in the game"},
				When:     []string{"the player performs an action"},
				Then:     []string{"the expected outcome occurs"},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid contract, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test contract missing name (E100)
func TestContractValidator_MissingName(t *testing.T) {
	cv := validator.NewContractValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Contracts: []domain.Contract{
			{
				Name:  "", // Missing name
				Given: []string{"some precondition"},
				When:  []string{"some action"},
				Then:  []string{"some result"},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for missing contract name")
	}

	errs := collector.Errors()
	if errs[0].Code != "E100" {
		t.Errorf("expected error code E100, got %s", errs[0].Code)
	}
}

// Test duplicate contract names within a node (E103)
func TestContractValidator_DuplicateNames(t *testing.T) {
	cv := validator.NewContractValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Contracts: []domain.Contract{
			{
				Name:  "Same Name",
				Given: []string{"precondition 1"},
				When:  []string{"action 1"},
				Then:  []string{"result 1"},
			},
			{
				Name:  "Same Name", // Duplicate
				Given: []string{"precondition 2"},
				When:  []string{"action 2"},
				Then:  []string{"result 2"},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for duplicate contract name")
	}

	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E103" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find duplicate name error (E103)")
	}
}

// Test empty step (E101)
func TestContractValidator_EmptyStep(t *testing.T) {
	cv := validator.NewContractValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Contracts: []domain.Contract{
			{
				Name:  "Test Contract",
				Given: []string{"valid step", ""}, // Empty step
				When:  []string{"action"},
				Then:  []string{"result"},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for empty step")
	}

	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E101" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find empty step error (E101)")
	}
}

// Test contract with no steps (E104)
func TestContractValidator_NoSteps(t *testing.T) {
	cv := validator.NewContractValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Contracts: []domain.Contract{
			{
				Name:     "Empty Contract",
				Scenario: "Has no steps",
				Given:    []string{},
				When:     []string{},
				Then:     []string{},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for contract with no steps")
	}

	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E104" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find no steps error (E104)")
	}
}

// Test contract with only given steps (valid)
func TestContractValidator_OnlyGivenSteps(t *testing.T) {
	cv := validator.NewContractValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Contracts: []domain.Contract{
			{
				Name:  "Setup Only",
				Given: []string{"the system is initialized"},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for contract with only given steps, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test node with no contracts (should pass)
func TestContractValidator_NoContracts(t *testing.T) {
	cv := validator.NewContractValidator()

	node := domain.Node{
		ID:        "test-node",
		Kind:      "system",
		Version:   1,
		Status:    "draft",
		Title:     "Test Node",
		Contracts: []domain.Contract{},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for node without contracts, got %d", collector.Count())
	}
}

// Test nil node (should pass)
func TestContractValidator_NilNode(t *testing.T) {
	cv := validator.NewContractValidator()

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(nil, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for nil node, got %d", collector.Count())
	}
}

// Test multiple contracts - one valid, one invalid
func TestContractValidator_MultipleContracts(t *testing.T) {
	cv := validator.NewContractValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Contracts: []domain.Contract{
			{
				Name:  "Valid Contract",
				Given: []string{"precondition"},
				When:  []string{"action"},
				Then:  []string{"result"},
			},
			{
				Name:  "", // Invalid: missing name
				Given: []string{"precondition"},
				When:  []string{"action"},
				Then:  []string{"result"},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for invalid contract")
	}

	// Should only have 1 error (missing name in second contract)
	if collector.Count() != 1 {
		t.Errorf("expected 1 error, got %d", collector.Count())
	}
}

// Test ValidateAll across multiple nodes
func TestContractValidator_ValidateAll(t *testing.T) {
	cv := validator.NewContractValidator()

	nodes := []domain.Node{
		{
			ID:      "node1",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Node 1",
			Contracts: []domain.Contract{
				{Name: "Contract 1", Given: []string{"step"}},
			},
		},
		{
			ID:      "node2",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Node 2",
			Contracts: []domain.Contract{
				{Name: "", Given: []string{"step"}}, // Invalid: missing name
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.ValidateAll(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for invalid contract in node2")
	}

	if collector.Count() != 1 {
		t.Errorf("expected 1 error, got %d", collector.Count())
	}
}

// Test orchestrator includes contract validation
func TestOrchestrator_ValidatesContracts(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{
		{
			ID:      "node1",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Node 1",
			Contracts: []domain.Contract{
				{
					Name:  "Duplicate",
					Given: []string{"step"},
				},
				{
					Name:  "Duplicate", // Same name - should trigger E103
					Given: []string{"step"},
				},
			},
		},
	}

	collector := orch.ValidateAll(nodes)

	if !collector.HasErrors() {
		t.Fatal("expected error for duplicate contract names")
	}

	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E103" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected orchestrator to catch duplicate contract name error (E103)")
	}
}

// ===== CONTRACT NODE REFERENCE VALIDATION TESTS =====

// Test valid node references in contract
func TestContractValidator_ValidNodeRefs(t *testing.T) {
	cv := validator.NewContractValidator()

	nodes := []domain.Node{
		{ID: "systems/combat", Kind: "system", Version: 1, Status: "draft", Title: "Combat System"},
		{ID: "mechanics/damage", Kind: "mechanic", Version: 1, Status: "draft", Title: "Damage Mechanic"},
		{
			ID:      "features/attack",
			Kind:    "feature",
			Version: 1,
			Status:  "draft",
			Title:   "Attack Feature",
			Contracts: []domain.Contract{
				{
					Name:  "Attack Flow",
					Given: []string{"@systems/combat is active"},
					When:  []string{"player attacks using @mechanics/damage"},
					Then:  []string{"damage is applied"},
				},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.ValidateAll(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid node references, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test invalid node reference (E102)
func TestContractValidator_InvalidNodeRef(t *testing.T) {
	cv := validator.NewContractValidator()

	nodes := []domain.Node{
		{ID: "systems/combat", Kind: "system", Version: 1, Status: "draft", Title: "Combat System"},
		{
			ID:      "features/attack",
			Kind:    "feature",
			Version: 1,
			Status:  "draft",
			Title:   "Attack Feature",
			Contracts: []domain.Contract{
				{
					Name:  "Attack Flow",
					Given: []string{"@systems/combat is active"},
					When:  []string{"player attacks using @mechanics/nonexistent"}, // Invalid ref
					Then:  []string{"damage is applied"},
				},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.ValidateAll(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for invalid node reference")
	}

	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E102" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find invalid node reference error (E102)")
	}
}

// Test suggestion for typo in node reference
func TestContractValidator_NodeRefSuggestion(t *testing.T) {
	cv := validator.NewContractValidator()

	nodes := []domain.Node{
		{ID: "systems/combat", Kind: "system", Version: 1, Status: "draft", Title: "Combat System"},
		{
			ID:      "features/attack",
			Kind:    "feature",
			Version: 1,
			Status:  "draft",
			Title:   "Attack Feature",
			Contracts: []domain.Contract{
				{
					Name:  "Attack Flow",
					Given: []string{"@systems/combt is active"}, // Typo: combt instead of combat
					When:  []string{"player attacks"},
					Then:  []string{"damage is applied"},
				},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.ValidateAll(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for typo in node reference")
	}

	errs := collector.Errors()
	if errs[0].Suggestion == "" {
		t.Error("expected suggestion for similar node ID")
	}
}

// Test multiple invalid node references
func TestContractValidator_MultipleInvalidNodeRefs(t *testing.T) {
	cv := validator.NewContractValidator()

	nodes := []domain.Node{
		{ID: "systems/combat", Kind: "system", Version: 1, Status: "draft", Title: "Combat System"},
		{
			ID:      "features/attack",
			Kind:    "feature",
			Version: 1,
			Status:  "draft",
			Title:   "Attack Feature",
			Contracts: []domain.Contract{
				{
					Name:  "Attack Flow",
					Given: []string{"@missing/node1 is active"},
					When:  []string{"player attacks using @missing/node2"},
					Then:  []string{"@missing/node3 is updated"},
				},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.ValidateAll(nodes, collector)

	// Should have 3 errors (one for each missing reference)
	if collector.Count() != 3 {
		t.Errorf("expected 3 errors for invalid references, got %d", collector.Count())
	}
}

// Test contract without node references (should pass)
func TestContractValidator_NoNodeRefs(t *testing.T) {
	cv := validator.NewContractValidator()

	nodes := []domain.Node{
		{
			ID:      "features/simple",
			Kind:    "feature",
			Version: 1,
			Status:  "draft",
			Title:   "Simple Feature",
			Contracts: []domain.Contract{
				{
					Name:  "Simple Flow",
					Given: []string{"the system is ready"},
					When:  []string{"user performs action"},
					Then:  []string{"expected result occurs"},
				},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.ValidateAll(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for contract without node refs, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test orchestrator validates contract node references
func TestOrchestrator_ValidatesContractNodeRefs(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{
		{ID: "systems/combat", Kind: "system", Version: 1, Status: "draft", Title: "Combat System"},
		{
			ID:      "features/attack",
			Kind:    "feature",
			Version: 1,
			Status:  "draft",
			Title:   "Attack Feature",
			Contracts: []domain.Contract{
				{
					Name:  "Attack Flow",
					Given: []string{"@systems/nonexistent is active"}, // Invalid ref
					When:  []string{"player attacks"},
					Then:  []string{"damage is applied"},
				},
			},
		},
	}

	collector := orch.ValidateAll(nodes)

	if !collector.HasErrors() {
		t.Fatal("expected error for invalid node reference")
	}

	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E102" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected orchestrator to catch invalid node reference error (E102)")
	}
}

// ===== CONTENT VALIDATOR TESTS =====

// Test draft node without content (should pass)
func TestContentValidator_DraftWithoutContent(t *testing.T) {
	cv := validator.NewContentValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		// No content
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for draft node without content, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test approved node without content (should fail)
func TestContentValidator_ApprovedWithoutContent(t *testing.T) {
	cv := validator.NewContentValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "approved",
		Title:   "Test Node",
		// No content
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for approved node without content")
	}

	errs := collector.Errors()
	if errs[0].Code != "E046" {
		t.Errorf("expected error code E046, got %s", errs[0].Code)
	}
}

// Test published node without content (should fail)
func TestContentValidator_PublishedWithoutContent(t *testing.T) {
	cv := validator.NewContentValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "published",
		Title:   "Test Node",
		// No content
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for published node without content")
	}

	errs := collector.Errors()
	if errs[0].Code != "E046" {
		t.Errorf("expected error code E046, got %s", errs[0].Code)
	}
}

// Test approved node with empty content (no sections, should fail)
func TestContentValidator_ApprovedWithEmptyContent(t *testing.T) {
	cv := validator.NewContentValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "approved",
		Title:   "Test Node",
		Content: &domain.Content{
			Sections: []domain.Section{}, // Empty sections
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for approved node with empty content")
	}

	errs := collector.Errors()
	if errs[0].Code != "E046" {
		t.Errorf("expected error code E046, got %s", errs[0].Code)
	}
}

// Test approved node with content (should pass)
func TestContentValidator_ApprovedWithContent(t *testing.T) {
	cv := validator.NewContentValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "approved",
		Title:   "Test Node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Overview",
					Blocks: []domain.Block{
						{Type: "rule", Data: map[string]interface{}{"text": "A rule"}},
					},
				},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for approved node with content, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test published node with content (should pass)
func TestContentValidator_PublishedWithContent(t *testing.T) {
	cv := validator.NewContentValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "published",
		Title:   "Test Node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Overview",
					Blocks: []domain.Block{
						{Type: "rule", Data: map[string]interface{}{"text": "A rule"}},
					},
				},
			},
		},
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for published node with content, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test deprecated node without content (should pass - not approved/published)
func TestContentValidator_DeprecatedWithoutContent(t *testing.T) {
	cv := validator.NewContentValidator()

	node := domain.Node{
		ID:      "test-node",
		Kind:    "system",
		Version: 1,
		Status:  "deprecated",
		Title:   "Test Node",
		// No content
	}

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for deprecated node without content, got %d: %v", collector.Count(), collector.Errors())
	}
}

// Test nil node (should pass)
func TestContentValidator_NilNode(t *testing.T) {
	cv := validator.NewContentValidator()

	collector := errors.NewCollectorWithLimit(100)
	cv.Validate(nil, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for nil node, got %d", collector.Count())
	}
}

// Test orchestrator validates content requirements
func TestOrchestrator_ValidatesContentRequirements(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{
		{
			ID:      "node1",
			Kind:    "system",
			Version: 1,
			Status:  "approved",
			Title:   "Node 1",
			// No content - should fail
		},
	}

	collector := orch.ValidateAll(nodes)

	if !collector.HasErrors() {
		t.Fatal("expected error for approved node without content")
	}

	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E046" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected orchestrator to catch content requirement error (E046)")
	}
}

// Test orchestrator validates block requirements
func TestOrchestrator_ValidatesBlockRequirements(t *testing.T) {
	orch := validator.NewOrchestrator()

	nodes := []domain.Node{
		{
			ID:      "node1",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Node 1",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Rules",
						Blocks: []domain.Block{
							{
								Type: "rule",
								Data: map[string]interface{}{
									"id": "test_rule",
									// missing "text" - should fail
								},
							},
						},
					},
				},
			},
		},
	}

	collector := orch.ValidateAll(nodes)

	if !collector.HasErrors() {
		t.Fatal("expected error for rule block missing text")
	}

	errs := collector.Errors()
	found := false
	for _, err := range errs {
		if err.Code == "E047" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected orchestrator to catch block requirement error (E047)")
	}
}

// ===== APPROVAL VALIDATOR TESTS =====

func TestApprovalValidator(t *testing.T) {
	t.Run("approved node without enough approvals fails", func(t *testing.T) {
		av := validator.NewApprovalValidator(2) // require 2 approvals
		collector := errors.NewCollector()

		node := &domain.Node{
			ID:         "test/node",
			Kind:       "mechanic",
			Version:    1,
			Status:     "approved",
			Title:      "Test",
			SourceFile: "test.yaml",
			Reviewers: []domain.Reviewer{
				{Name: "alice", Timestamp: time.Now(), Version: 1},
			},
		}

		av.Validate(node, collector)

		if !collector.HasErrors() {
			t.Error("Expected validation error for insufficient approvals")
		}
	})

	t.Run("approved node with enough approvals passes", func(t *testing.T) {
		av := validator.NewApprovalValidator(2)
		collector := errors.NewCollector()

		node := &domain.Node{
			ID:         "test/node",
			Kind:       "mechanic",
			Version:    1,
			Status:     "approved",
			Title:      "Test",
			SourceFile: "test.yaml",
			Reviewers: []domain.Reviewer{
				{Name: "alice", Timestamp: time.Now(), Version: 1},
				{Name: "bob", Timestamp: time.Now(), Version: 1},
			},
		}

		av.Validate(node, collector)

		if collector.HasErrors() {
			t.Errorf("Expected no errors, got: %v", collector.Errors())
		}
	})

	t.Run("draft node skips approval check", func(t *testing.T) {
		av := validator.NewApprovalValidator(2)
		collector := errors.NewCollector()

		node := &domain.Node{
			ID:         "test/node",
			Kind:       "mechanic",
			Version:    1,
			Status:     "draft",
			Title:      "Test",
			SourceFile: "test.yaml",
		}

		av.Validate(node, collector)

		if collector.HasErrors() {
			t.Errorf("Expected no errors for draft node, got: %v", collector.Errors())
		}
	})

	t.Run("approvals must match current version", func(t *testing.T) {
		av := validator.NewApprovalValidator(1)
		collector := errors.NewCollector()

		node := &domain.Node{
			ID:         "test/node",
			Kind:       "mechanic",
			Version:    2, // current version is 2
			Status:     "approved",
			Title:      "Test",
			SourceFile: "test.yaml",
			Reviewers: []domain.Reviewer{
				{Name: "alice", Timestamp: time.Now(), Version: 1}, // approved version 1
			},
		}

		av.Validate(node, collector)

		if !collector.HasErrors() {
			t.Error("Expected error for stale approval")
		}
	})
}
