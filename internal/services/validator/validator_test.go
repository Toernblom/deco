package validator_test

import (
	"testing"

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
