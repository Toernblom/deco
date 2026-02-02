package patcher_test

import (
	"fmt"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/patcher"
)

// ===== SET OPERATION TESTS =====

// Test setting a simple top-level field
func TestPatcher_SetSimpleField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Old Title",
	}

	err := p.Set(&node, "title", "New Title")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Title != "New Title" {
		t.Errorf("expected title 'New Title', got %q", node.Title)
	}
}

// Test setting a nested field with dot notation
func TestPatcher_SetNestedField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Summary: "Old summary",
	}

	err := p.Set(&node, "summary", "New summary")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Summary != "New summary" {
		t.Errorf("expected summary 'New summary', got %q", node.Summary)
	}
}

// Test setting a new field that didn't exist before
func TestPatcher_SetNewField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	// Summary doesn't exist yet
	err := p.Set(&node, "summary", "New summary")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Summary != "New summary" {
		t.Errorf("expected summary to be set, got %q", node.Summary)
	}
}

// Test setting with invalid path
func TestPatcher_SetInvalidPath(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	// Try to set a field that doesn't exist and can't be created
	err := p.Set(&node, "invalid.deeply.nested.field", "value")

	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

// Test setting an array element
func TestPatcher_SetArrayElement(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"old-tag", "another-tag"},
	}

	err := p.Set(&node, "tags[0]", "new-tag")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(node.Tags))
	}

	if node.Tags[0] != "new-tag" {
		t.Errorf("expected first tag to be 'new-tag', got %q", node.Tags[0])
	}
}

// Test setting with type conversion
func TestPatcher_SetWithTypeConversion(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	// Set version (int) from string
	err := p.Set(&node, "version", 2)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Version != 2 {
		t.Errorf("expected version 2, got %d", node.Version)
	}
}

// ===== APPEND OPERATION TESTS =====

// Test appending to an existing array
func TestPatcher_AppendToArray(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1", "tag2"},
	}

	err := p.Append(&node, "tags", "tag3")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(node.Tags))
	}

	if node.Tags[2] != "tag3" {
		t.Errorf("expected third tag to be 'tag3', got %q", node.Tags[2])
	}
}

// Test appending to an empty array
func TestPatcher_AppendToEmptyArray(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{},
	}

	err := p.Append(&node, "tags", "first-tag")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(node.Tags))
	}

	if node.Tags[0] != "first-tag" {
		t.Errorf("expected tag 'first-tag', got %q", node.Tags[0])
	}
}

// Test appending to nil array (should initialize it)
func TestPatcher_AppendToNilArray(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    nil,
	}

	err := p.Append(&node, "tags", "first-tag")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(node.Tags))
	}

	if node.Tags[0] != "first-tag" {
		t.Errorf("expected tag 'first-tag', got %q", node.Tags[0])
	}
}

// Test appending to a non-array field (should error)
func TestPatcher_AppendToNonArray(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	// Try to append to a string field
	err := p.Append(&node, "title", "something")

	if err == nil {
		t.Fatal("expected error when appending to non-array, got nil")
	}
}

// Test appending multiple values at once
func TestPatcher_AppendMultiple(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1"},
	}

	err := p.Append(&node, "tags", "tag2", "tag3", "tag4")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 4 {
		t.Fatalf("expected 4 tags, got %d", len(node.Tags))
	}

	expected := []string{"tag1", "tag2", "tag3", "tag4"}
	for i, tag := range expected {
		if node.Tags[i] != tag {
			t.Errorf("expected tag[%d] to be %q, got %q", i, tag, node.Tags[i])
		}
	}
}

// ===== UNSET OPERATION TESTS =====

// Test removing a simple field
func TestPatcher_UnsetSimpleField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Summary: "A summary",
	}

	err := p.Unset(&node, "summary")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Summary != "" {
		t.Errorf("expected summary to be empty, got %q", node.Summary)
	}
}

// Test removing a nested field
func TestPatcher_UnsetNestedField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:       "test",
		Kind:     "system",
		Version:  1,
		Status:   "draft",
		Title:    "Test",
		Glossary: map[string]string{"term1": "def1", "term2": "def2"},
	}

	err := p.Unset(&node, "glossary.term1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Glossary) != 1 {
		t.Fatalf("expected 1 glossary entry, got %d", len(node.Glossary))
	}

	if _, exists := node.Glossary["term1"]; exists {
		t.Error("expected term1 to be removed from glossary")
	}
}

// Test removing an array element by index
func TestPatcher_UnsetArrayElement(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1", "tag2", "tag3"},
	}

	err := p.Unset(&node, "tags[1]")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 2 {
		t.Fatalf("expected 2 tags after removal, got %d", len(node.Tags))
	}

	// Should have tag1 and tag3, tag2 removed
	if node.Tags[0] != "tag1" || node.Tags[1] != "tag3" {
		t.Errorf("expected tags [tag1, tag3], got %v", node.Tags)
	}
}

// Test removing entire array
func TestPatcher_UnsetArray(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1", "tag2"},
	}

	err := p.Unset(&node, "tags")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 0 {
		t.Errorf("expected tags to be empty, got %v", node.Tags)
	}
}

// Test unsetting a field that doesn't exist (should not error)
func TestPatcher_UnsetMissingField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	// Summary doesn't exist
	err := p.Unset(&node, "summary")

	if err != nil {
		t.Errorf("expected no error when unsetting missing field, got %v", err)
	}
}

// Test unsetting with invalid path
func TestPatcher_UnsetInvalidPath(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	err := p.Unset(&node, "invalid.deeply.nested.path")

	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

// Test unsetting required field (should error or validate)
func TestPatcher_UnsetRequiredField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	// ID is required - unsetting should error or be prevented
	err := p.Unset(&node, "id")

	if err == nil {
		t.Fatal("expected error when unsetting required field, got nil")
	}
}

// ===== APPLY OPERATION TESTS =====

// Test applying a single set operation
func TestPatcher_ApplySingleSet(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Old Title",
	}

	ops := []patcher.PatchOperation{
		{Op: "set", Path: "title", Value: "New Title"},
	}

	err := p.Apply(&node, ops)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Title != "New Title" {
		t.Errorf("expected title 'New Title', got %q", node.Title)
	}
}

// Test applying multiple operations in sequence
func TestPatcher_ApplyMultipleOperations(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Old Title",
		Tags:    []string{"tag1"},
	}

	ops := []patcher.PatchOperation{
		{Op: "set", Path: "title", Value: "New Title"},
		{Op: "append", Path: "tags", Value: "tag2"},
		{Op: "set", Path: "summary", Value: "A summary"},
	}

	err := p.Apply(&node, ops)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Title != "New Title" {
		t.Errorf("expected title 'New Title', got %q", node.Title)
	}

	if len(node.Tags) != 2 || node.Tags[1] != "tag2" {
		t.Errorf("expected tags [tag1, tag2], got %v", node.Tags)
	}

	if node.Summary != "A summary" {
		t.Errorf("expected summary 'A summary', got %q", node.Summary)
	}
}

// Test applying mixed set/append/unset operations
func TestPatcher_ApplyMixedOperations(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Summary: "Old summary",
		Tags:    []string{"tag1", "tag2"},
	}

	ops := []patcher.PatchOperation{
		{Op: "set", Path: "title", Value: "Updated Title"},
		{Op: "unset", Path: "summary"},
		{Op: "append", Path: "tags", Value: "tag3"},
		{Op: "set", Path: "tags[0]", Value: "new-tag1"},
	}

	err := p.Apply(&node, ops)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %q", node.Title)
	}

	if node.Summary != "" {
		t.Errorf("expected summary to be empty, got %q", node.Summary)
	}

	if len(node.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(node.Tags))
	}

	if node.Tags[0] != "new-tag1" {
		t.Errorf("expected first tag 'new-tag1', got %q", node.Tags[0])
	}
}

// Test applying empty operations list (should succeed with no changes)
func TestPatcher_ApplyEmptyOperations(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	ops := []patcher.PatchOperation{}

	err := p.Apply(&node, ops)

	if err != nil {
		t.Fatalf("expected no error for empty operations, got %v", err)
	}

	if node.Title != "Test" {
		t.Errorf("expected title unchanged, got %q", node.Title)
	}
}

// Test rollback on error - all operations should be reverted if one fails
func TestPatcher_ApplyRollbackOnError(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Original Title",
		Summary: "Original Summary",
		Tags:    []string{"tag1"},
	}

	// Second operation will fail (invalid field)
	ops := []patcher.PatchOperation{
		{Op: "set", Path: "title", Value: "New Title"},
		{Op: "set", Path: "invalid.field.path", Value: "value"}, // This should fail
		{Op: "set", Path: "summary", Value: "New Summary"},
	}

	err := p.Apply(&node, ops)

	if err == nil {
		t.Fatal("expected error for invalid operation, got nil")
	}

	// All changes should be rolled back
	if node.Title != "Original Title" {
		t.Errorf("expected title rolled back to 'Original Title', got %q", node.Title)
	}

	if node.Summary != "Original Summary" {
		t.Errorf("expected summary unchanged after rollback, got %q", node.Summary)
	}

	if len(node.Tags) != 1 || node.Tags[0] != "tag1" {
		t.Errorf("expected tags unchanged after rollback, got %v", node.Tags)
	}
}

// Test rollback with append operation failure
func TestPatcher_ApplyRollbackWithAppend(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Original Title",
		Tags:    []string{"tag1"},
	}

	// First operation succeeds, second fails (append to non-array)
	ops := []patcher.PatchOperation{
		{Op: "append", Path: "tags", Value: "tag2"},
		{Op: "append", Path: "title", Value: "invalid"}, // Can't append to string
	}

	err := p.Apply(&node, ops)

	if err == nil {
		t.Fatal("expected error when appending to non-array, got nil")
	}

	// Tags should be rolled back to original state
	if len(node.Tags) != 1 || node.Tags[0] != "tag1" {
		t.Errorf("expected tags rolled back to [tag1], got %v", node.Tags)
	}
}

// Test invalid operation type
func TestPatcher_ApplyInvalidOperation(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	ops := []patcher.PatchOperation{
		{Op: "invalid", Path: "title", Value: "value"},
	}

	err := p.Apply(&node, ops)

	if err == nil {
		t.Fatal("expected error for invalid operation type, got nil")
	}
}

// Test applying to nil node
func TestPatcher_ApplyToNilNode(t *testing.T) {
	p := patcher.New()

	ops := []patcher.PatchOperation{
		{Op: "set", Path: "title", Value: "Test"},
	}

	err := p.Apply(nil, ops)

	if err == nil {
		t.Fatal("expected error when applying to nil node, got nil")
	}
}

// Test operations are applied in order
func TestPatcher_ApplyOperationsInOrder(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    nil,
	}

	// Order matters: append to tags, then modify the appended element
	ops := []patcher.PatchOperation{
		{Op: "append", Path: "tags", Value: "tag1"},
		{Op: "append", Path: "tags", Value: "tag2"},
		{Op: "set", Path: "tags[1]", Value: "modified-tag2"},
	}

	err := p.Apply(&node, ops)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(node.Tags))
	}

	if node.Tags[0] != "tag1" {
		t.Errorf("expected first tag 'tag1', got %q", node.Tags[0])
	}

	if node.Tags[1] != "modified-tag2" {
		t.Errorf("expected second tag 'modified-tag2', got %q", node.Tags[1])
	}
}

// Test unset operation in apply
func TestPatcher_ApplyUnsetOperation(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Summary: "Summary to remove",
		Tags:    []string{"tag1", "tag2", "tag3"},
	}

	ops := []patcher.PatchOperation{
		{Op: "unset", Path: "summary"},
		{Op: "unset", Path: "tags[1]"}, // Remove "tag2"
	}

	err := p.Apply(&node, ops)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Summary != "" {
		t.Errorf("expected summary to be empty, got %q", node.Summary)
	}

	if len(node.Tags) != 2 {
		t.Fatalf("expected 2 tags after unset, got %d", len(node.Tags))
	}

	if node.Tags[0] != "tag1" || node.Tags[1] != "tag3" {
		t.Errorf("expected tags [tag1, tag3], got %v", node.Tags)
	}
}

// Test rollback on unset required field error
func TestPatcher_ApplyRollbackOnUnsetRequired(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Original Title",
	}

	// Try to unset a required field - should fail and rollback
	ops := []patcher.PatchOperation{
		{Op: "set", Path: "title", Value: "New Title"},
		{Op: "unset", Path: "id"}, // Required field - should fail
	}

	err := p.Apply(&node, ops)

	if err == nil {
		t.Fatal("expected error when unsetting required field, got nil")
	}

	// Title should be rolled back
	if node.Title != "Original Title" {
		t.Errorf("expected title rolled back to 'Original Title', got %q", node.Title)
	}

	// ID should still exist
	if node.ID != "test" {
		t.Errorf("expected ID to remain 'test', got %q", node.ID)
	}
}

// Test apply with append multiple values (variadic)
func TestPatcher_ApplyAppendMultipleValues(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1"},
	}

	// Note: For multiple values in append, we'd need to support array values
	// or handle it specially. This tests appending an array as a single operation.
	ops := []patcher.PatchOperation{
		{Op: "append", Path: "tags", Value: "tag2"},
		{Op: "append", Path: "tags", Value: "tag3"},
	}

	err := p.Apply(&node, ops)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(node.Tags))
	}

	expected := []string{"tag1", "tag2", "tag3"}
	for i, tag := range expected {
		if node.Tags[i] != tag {
			t.Errorf("expected tag[%d] to be %q, got %q", i, tag, node.Tags[i])
		}
	}
}

// ===== POINTER FIELD TESTS =====

// Test setting a field through a pointer (Content *Content)
func TestPatcher_SetThroughPointerField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Original Section",
					Blocks: []domain.Block{
						{Type: "text", Data: map[string]interface{}{"text": "original"}},
					},
				},
			},
		},
	}

	// This path traverses through Content (which is *Content pointer)
	err := p.Set(&node, "content.sections[0].name", "Updated Section")

	if err != nil {
		t.Fatalf("expected no error setting through pointer, got %v", err)
	}

	if node.Content.Sections[0].Name != "Updated Section" {
		t.Errorf("expected section name 'Updated Section', got %q", node.Content.Sections[0].Name)
	}
}

// Test setting a deeply nested field through pointer
func TestPatcher_SetDeeplyNestedThroughPointer(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Section",
					Blocks: []domain.Block{
						{Type: "text", Data: map[string]interface{}{"text": "original"}},
					},
				},
			},
		},
	}

	// Set block type through pointer path
	err := p.Set(&node, "content.sections[0].blocks[0].type", "updated-type")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Content.Sections[0].Blocks[0].Type != "updated-type" {
		t.Errorf("expected block type 'updated-type', got %q", node.Content.Sections[0].Blocks[0].Type)
	}
}

// Test setting through nil pointer should initialize it
func TestPatcher_SetThroughNilPointer(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: nil, // nil pointer
	}

	// Attempting to set through nil pointer should initialize it
	err := p.Set(&node, "content.sections", []domain.Section{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Content == nil {
		t.Fatal("expected Content to be initialized")
	}
}

// Test getField through pointer
func TestPatcher_GetFieldThroughPointer(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{Name: "Section1", Blocks: []domain.Block{}},
				{Name: "Section2", Blocks: []domain.Block{}},
			},
		},
	}

	// Append should work through pointer path
	err := p.Append(&node, "content.sections", domain.Section{Name: "Section3", Blocks: []domain.Block{}})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Content.Sections) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(node.Content.Sections))
	}

	if node.Content.Sections[2].Name != "Section3" {
		t.Errorf("expected third section name 'Section3', got %q", node.Content.Sections[2].Name)
	}
}

// Test unset through pointer
func TestPatcher_UnsetThroughPointer(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{Name: "Section1", Blocks: []domain.Block{}},
				{Name: "Section2", Blocks: []domain.Block{}},
			},
		},
	}

	// Unset a section through pointer path
	err := p.Unset(&node, "content.sections[0]")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Content.Sections) != 1 {
		t.Fatalf("expected 1 section after unset, got %d", len(node.Content.Sections))
	}

	if node.Content.Sections[0].Name != "Section2" {
		t.Errorf("expected remaining section to be 'Section2', got %q", node.Content.Sections[0].Name)
	}
}

// ===== NIL NODE ERROR TESTS =====

// Test Set with nil node returns error
func TestPatcher_SetNilNode(t *testing.T) {
	p := patcher.New()

	err := p.Set(nil, "title", "value")
	if err == nil {
		t.Fatal("expected error when setting on nil node, got nil")
	}
}

// Test Append with nil node returns error
func TestPatcher_AppendNilNode(t *testing.T) {
	p := patcher.New()

	err := p.Append(nil, "tags", "value")
	if err == nil {
		t.Fatal("expected error when appending to nil node, got nil")
	}
}

// Test Unset with nil node returns error
func TestPatcher_UnsetNilNode(t *testing.T) {
	p := patcher.New()

	err := p.Unset(nil, "tags")
	if err == nil {
		t.Fatal("expected error when unsetting on nil node, got nil")
	}
}

// ===== EMPTY PATH ERROR TESTS =====

// Test Set with empty path returns error
func TestPatcher_SetEmptyPath(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	err := p.Set(&node, "", "value")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

// Test Append with empty path returns error
func TestPatcher_AppendEmptyPath(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	err := p.Append(&node, "", "value")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

// Test Unset with empty path returns error
func TestPatcher_UnsetEmptyPath(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	err := p.Unset(&node, "")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

// ===== TYPE CONVERSION TESTS =====

// Test Set fails when type conversion is not possible
func TestPatcher_SetTypeConversionFailure(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	// Try to set a string field to a slice (incompatible)
	err := p.Set(&node, "title", []int{1, 2, 3})
	if err == nil {
		t.Fatal("expected error for incompatible type conversion, got nil")
	}
}

// ===== ARRAY INDEX EDGE CASES =====

// Test Set with array index out of bounds
func TestPatcher_SetArrayIndexOutOfBounds(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1"},
	}

	err := p.Set(&node, "tags[5]", "value")
	if err == nil {
		t.Fatal("expected error for out of bounds index, got nil")
	}
}

// Test Set with negative array index
func TestPatcher_SetNegativeArrayIndex(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1"},
	}

	// Negative index should be handled as invalid (parseArrayIndex returns false)
	err := p.Set(&node, "tags[-1]", "value")
	if err == nil {
		t.Fatal("expected error for negative index, got nil")
	}
}

// Test Set with malformed array notation
func TestPatcher_SetMalformedArrayNotation(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1"},
	}

	// Malformed bracket notation - no closing bracket
	err := p.Set(&node, "tags[0", "value")
	if err == nil {
		t.Fatal("expected error for malformed array notation, got nil")
	}
}

// Test Set with non-numeric array index
func TestPatcher_SetNonNumericArrayIndex(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1"},
	}

	// Non-numeric index
	err := p.Set(&node, "tags[abc]", "value")
	if err == nil {
		t.Fatal("expected error for non-numeric index, got nil")
	}
}

// Test Set with array notation on non-slice field
func TestPatcher_SetArrayNotationOnNonSlice(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	// Try to use array notation on a string field
	err := p.Set(&node, "title[0]", "value")
	if err == nil {
		t.Fatal("expected error for array notation on non-slice, got nil")
	}
}

// ===== MAP OPERATION TESTS =====

// Test Set value in a map (glossary)
func TestPatcher_SetMapValue(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:       "test",
		Kind:     "system",
		Version:  1,
		Status:   "draft",
		Title:    "Test",
		Glossary: map[string]string{"existing": "value"},
	}

	err := p.Set(&node, "glossary.newkey", "newvalue")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Glossary["newkey"] != "newvalue" {
		t.Errorf("expected glossary[newkey]='newvalue', got %q", node.Glossary["newkey"])
	}
}

// Test Set value in nil map initializes the map
func TestPatcher_SetMapValueNilMap(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:       "test",
		Kind:     "system",
		Version:  1,
		Status:   "draft",
		Title:    "Test",
		Glossary: nil,
	}

	err := p.Set(&node, "glossary.key", "value")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Glossary == nil {
		t.Fatal("expected glossary to be initialized")
	}

	if node.Glossary["key"] != "value" {
		t.Errorf("expected glossary[key]='value', got %q", node.Glossary["key"])
	}
}

// Test Unset map entry
func TestPatcher_UnsetMapEntry(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:       "test",
		Kind:     "system",
		Version:  1,
		Status:   "draft",
		Title:    "Test",
		Glossary: map[string]string{"key1": "val1", "key2": "val2"},
	}

	err := p.Unset(&node, "glossary.key1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Glossary) != 1 {
		t.Fatalf("expected 1 entry in glossary, got %d", len(node.Glossary))
	}

	if _, exists := node.Glossary["key1"]; exists {
		t.Error("expected key1 to be removed from glossary")
	}
}

// Test Unset from nil map (should not error)
func TestPatcher_UnsetFromNilMap(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:       "test",
		Kind:     "system",
		Version:  1,
		Status:   "draft",
		Title:    "Test",
		Glossary: nil,
	}

	err := p.Unset(&node, "glossary.key")
	if err != nil {
		t.Errorf("expected no error when unsetting from nil map, got %v", err)
	}
}

// ===== POINTER EDGE CASES =====

// Test Append to field through nil pointer fails gracefully
func TestPatcher_AppendThroughNilPointer(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: nil,
	}

	// Appending through nil pointer should fail
	err := p.Append(&node, "content.sections", domain.Section{})
	if err == nil {
		t.Fatal("expected error when appending through nil pointer, got nil")
	}
}

// Test getField through nil pointer returns error
func TestPatcher_GetFieldThroughNilPointerError(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: nil,
	}

	// Trying to access field through nil pointer for append should fail
	err := p.Append(&node, "content.sections", domain.Section{Name: "test"})
	if err == nil {
		t.Fatal("expected error for accessing through nil pointer, got nil")
	}
}

// ===== DEEPLY NESTED PATH TESTS =====

// Test Set deeply nested path in content blocks
func TestPatcher_SetDeeplyNestedBlockData(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Section1",
					Blocks: []domain.Block{
						{
							Type: "table",
							Data: map[string]interface{}{
								"headers": []string{"col1", "col2"},
							},
						},
					},
				},
			},
		},
	}

	// Set data field in nested block
	err := p.Set(&node, "content.sections[0].blocks[0].data.newfield", "newvalue")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Content.Sections[0].Blocks[0].Data["newfield"] != "newvalue" {
		t.Errorf("expected data[newfield]='newvalue', got %v",
			node.Content.Sections[0].Blocks[0].Data["newfield"])
	}
}

// ===== UNSET ARRAY ELEMENT EDGE CASES =====

// Test Unset array element at index 0 (first element)
func TestPatcher_UnsetFirstArrayElement(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"first", "second", "third"},
	}

	err := p.Unset(&node, "tags[0]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(node.Tags))
	}

	if node.Tags[0] != "second" || node.Tags[1] != "third" {
		t.Errorf("expected [second, third], got %v", node.Tags)
	}
}

// Test Unset last array element
func TestPatcher_UnsetLastArrayElement(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"first", "second", "third"},
	}

	err := p.Unset(&node, "tags[2]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(node.Tags))
	}

	if node.Tags[0] != "first" || node.Tags[1] != "second" {
		t.Errorf("expected [first, second], got %v", node.Tags)
	}
}

// Test Unset array element out of bounds (should not error)
func TestPatcher_UnsetArrayOutOfBounds(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1"},
	}

	// Out of bounds should not error (nothing to unset)
	err := p.Unset(&node, "tags[10]")
	if err != nil {
		t.Errorf("expected no error for out of bounds unset, got %v", err)
	}

	// Original should be unchanged
	if len(node.Tags) != 1 || node.Tags[0] != "tag1" {
		t.Errorf("expected tags unchanged, got %v", node.Tags)
	}
}

// ===== REQUIRED FIELDS TESTS =====

// Test Unset all required fields
func TestPatcher_UnsetAllRequiredFields(t *testing.T) {
	p := patcher.New()

	requiredFields := []string{"id", "kind", "version", "status", "title"}

	for _, field := range requiredFields {
		node := domain.Node{
			ID:      "test",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Test",
		}

		err := p.Unset(&node, field)
		if err == nil {
			t.Errorf("expected error when unsetting required field %q, got nil", field)
		}
	}
}

// ===== APPLY EDGE CASES =====

// Test Apply detects snapshot failure (unusual but possible)
func TestPatcher_ApplyWithComplexRollback(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Original",
		Tags:    []string{"tag1", "tag2"},
		Content: &domain.Content{
			Sections: []domain.Section{
				{Name: "Section1", Blocks: []domain.Block{}},
			},
		},
	}

	// First few operations succeed, then one fails
	ops := []patcher.PatchOperation{
		{Op: "set", Path: "title", Value: "Changed"},
		{Op: "append", Path: "tags", Value: "tag3"},
		{Op: "set", Path: "content.sections[0].name", Value: "Updated"},
		{Op: "set", Path: "nonexistent.deep.path", Value: "fail"}, // This will fail
	}

	err := p.Apply(&node, ops)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// All should be rolled back
	if node.Title != "Original" {
		t.Errorf("expected title rolled back to 'Original', got %q", node.Title)
	}
	if len(node.Tags) != 2 {
		t.Errorf("expected 2 tags after rollback, got %d", len(node.Tags))
	}
	if node.Content.Sections[0].Name != "Section1" {
		t.Errorf("expected section name rolled back to 'Section1', got %q",
			node.Content.Sections[0].Name)
	}
}

// ===== SPECIAL CHARACTER TESTS =====

// Test Set with various string values
func TestPatcher_SetSpecialStringValues(t *testing.T) {
	p := patcher.New()

	testCases := []struct {
		name  string
		value string
	}{
		{"empty string", ""},
		{"unicode", "日本語テスト"},
		{"emoji", "🎮🎯🎲"},
		{"special chars", "!@#$%^&*()"},
		{"newlines", "line1\nline2\nline3"},
		{"tabs", "col1\tcol2\tcol3"},
		{"quotes", `"quoted" and 'single'`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := domain.Node{
				ID:      "test",
				Kind:    "system",
				Version: 1,
				Status:  "draft",
				Title:   "Test",
			}

			err := p.Set(&node, "summary", tc.value)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if node.Summary != tc.value {
				t.Errorf("expected summary %q, got %q", tc.value, node.Summary)
			}
		})
	}
}

// ===== ISSUES FIELD TESTS (SLICE OF STRUCT) =====

// Test Append to Issues slice
func TestPatcher_AppendToIssuesSlice(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Issues:  []domain.Issue{},
	}

	newIssue := domain.Issue{
		ID:          "issue-1",
		Description: "Something to decide",
		Severity:    "medium",
		Location:    "content",
	}

	err := p.Append(&node, "issues", newIssue)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(node.Issues))
	}

	if node.Issues[0].ID != "issue-1" {
		t.Errorf("expected issue ID 'issue-1', got %q", node.Issues[0].ID)
	}
}

// Test Set issue field through array index
func TestPatcher_SetIssueFieldByIndex(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Issues: []domain.Issue{
			{ID: "issue-1", Description: "Original description", Severity: "low"},
		},
	}

	err := p.Set(&node, "issues[0].description", "Updated description")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Issues[0].Description != "Updated description" {
		t.Errorf("expected issue description 'Updated description', got %q", node.Issues[0].Description)
	}
}

// ===== CONTRACTS FIELD TESTS =====

// Test Set contract scenario field
func TestPatcher_SetContractScenarioField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Contracts: []domain.Contract{
			{Name: "Contract1", Scenario: "Original scenario"},
		},
	}

	err := p.Set(&node, "contracts[0].scenario", "Updated scenario")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Contracts[0].Scenario != "Updated scenario" {
		t.Errorf("expected 'Updated scenario', got %q", node.Contracts[0].Scenario)
	}
}

// Test Append to contract's Then slice
func TestPatcher_AppendToContractThenSlice(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Contracts: []domain.Contract{
			{Name: "Contract1", Scenario: "Test scenario", Then: []string{"result1"}},
		},
	}

	err := p.Append(&node, "contracts[0].then", "result2")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Contracts[0].Then) != 2 {
		t.Fatalf("expected 2 then items, got %d", len(node.Contracts[0].Then))
	}

	if node.Contracts[0].Then[1] != "result2" {
		t.Errorf("expected 'result2', got %q", node.Contracts[0].Then[1])
	}
}

// ===== REFS FIELD TESTS =====

// Test Append to refs.emitsEvents slice
func TestPatcher_AppendToRefsEmitsEvents(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Refs:    domain.Ref{EmitsEvents: []string{"event1"}},
	}

	err := p.Append(&node, "refs.emitsEvents", "event2")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Refs.EmitsEvents) != 2 {
		t.Fatalf("expected 2 events, got %d", len(node.Refs.EmitsEvents))
	}

	if node.Refs.EmitsEvents[1] != "event2" {
		t.Errorf("expected second event 'event2', got %q", node.Refs.EmitsEvents[1])
	}
}

// Test Append to refs.vocabulary slice
func TestPatcher_AppendToRefsVocabulary(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Refs:    domain.Ref{Vocabulary: []string{"term1"}},
	}

	err := p.Append(&node, "refs.vocabulary", "term2")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Refs.Vocabulary) != 2 {
		t.Fatalf("expected 2 terms, got %d", len(node.Refs.Vocabulary))
	}

	if node.Refs.Vocabulary[1] != "term2" {
		t.Errorf("expected second term 'term2', got %q", node.Refs.Vocabulary[1])
	}
}

// ===== FIELD NOT FOUND TESTS =====

// Test Set field not found at top level
func TestPatcher_SetFieldNotFoundTopLevel(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	err := p.Set(&node, "nonexistentfield", "value")
	if err == nil {
		t.Fatal("expected error for nonexistent field, got nil")
	}
}

// Test Set field not found in nested struct
func TestPatcher_SetFieldNotFoundNested(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Refs:    domain.Ref{},
	}

	err := p.Set(&node, "refs.nonexistent", "value")
	if err == nil {
		t.Fatal("expected error for nonexistent nested field, got nil")
	}
}

// ===== APPEND TYPE CONVERSION TESTS =====

// Test Append with type that needs conversion
func TestPatcher_AppendWithTypeConversion(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1"},
	}

	// String should convert fine
	err := p.Append(&node, "tags", "tag2")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 2 || node.Tags[1] != "tag2" {
		t.Errorf("expected [tag1, tag2], got %v", node.Tags)
	}
}

// ===== FIELD CAN'T BE SET TESTS =====

// Test getting field that exists for append operation
func TestPatcher_AppendToNestedSlice(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name:   "Section1",
					Blocks: []domain.Block{},
				},
			},
		},
	}

	newBlock := domain.Block{Type: "text", Data: map[string]interface{}{"text": "hello"}}
	err := p.Append(&node, "content.sections[0].blocks", newBlock)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Content.Sections[0].Blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(node.Content.Sections[0].Blocks))
	}
}

// ===== GETFIELD EDGE CASES =====

// Test Append to non-existent field returns error
func TestPatcher_AppendToNonExistentField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
	}

	err := p.Append(&node, "nonexistent", "value")
	if err == nil {
		t.Fatal("expected error for nonexistent field, got nil")
	}
}

// Test Append with array notation to non-existent array field
func TestPatcher_AppendToNonExistentArrayField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{},
		},
	}

	// Try to access non-existent index
	err := p.Append(&node, "content.sections[0].blocks", domain.Block{})
	if err == nil {
		t.Fatal("expected error for non-existent array element, got nil")
	}
}

// Test getField with map access
func TestPatcher_AppendToMapFieldError(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:       "test",
		Kind:     "system",
		Version:  1,
		Status:   "draft",
		Title:    "Test",
		Glossary: map[string]string{"key": "value"},
	}

	// Can't append to a map (map values aren't slices)
	err := p.Append(&node, "glossary.key", "value")
	if err == nil {
		t.Fatal("expected error when appending to map value, got nil")
	}
}

// ===== UNSETVALUE NESTED EDGE CASES =====

// Test Unset nested field inside array element
func TestPatcher_UnsetNestedInArrayElement(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Section1",
					Blocks: []domain.Block{
						{Type: "text", Data: map[string]interface{}{"text": "hello"}},
					},
				},
			},
		},
	}

	// Unset nested field within array element
	err := p.Unset(&node, "content.sections[0].blocks[0].data.text")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// The 'text' key should be removed from the data map
	if _, exists := node.Content.Sections[0].Blocks[0].Data["text"]; exists {
		t.Error("expected text key to be removed from data map")
	}
}

// Test Unset through nil pointer field (should not error)
func TestPatcher_UnsetThroughNilContentPointer(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: nil,
	}

	// Trying to unset through nil pointer should not error (nothing to unset)
	err := p.Unset(&node, "content.sections")
	if err != nil {
		t.Errorf("expected no error when unsetting through nil pointer, got %v", err)
	}
}

// ===== SETVALUE EDGE CASES =====

// Test Set with type conversion at array element
func TestPatcher_SetArrayElementWithTypeConversion(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"tag1", "tag2"},
	}

	// Setting string should work
	err := p.Set(&node, "tags[0]", "newtag")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Tags[0] != "newtag" {
		t.Errorf("expected 'newtag', got %q", node.Tags[0])
	}
}

// Test Set with nested path through array to map value
func TestPatcher_SetNestedMapThroughArray(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Section1",
					Blocks: []domain.Block{
						{Type: "text", Data: map[string]interface{}{}},
					},
				},
			},
		},
	}

	// Set a value inside the data map through array access
	err := p.Set(&node, "content.sections[0].blocks[0].data.newkey", "newvalue")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Content.Sections[0].Blocks[0].Data["newkey"] != "newvalue" {
		t.Errorf("expected data[newkey]='newvalue', got %v",
			node.Content.Sections[0].Blocks[0].Data["newkey"])
	}
}

// ===== ADDITIONAL ARRAY OPERATIONS =====

// Test Set all elements in an array using loop
func TestPatcher_SetMultipleArrayElements(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"a", "b", "c"},
	}

	// Set each element
	for i := 0; i < 3; i++ {
		err := p.Set(&node, fmt.Sprintf("tags[%d]", i), fmt.Sprintf("new%d", i))
		if err != nil {
			t.Fatalf("failed to set tags[%d]: %v", i, err)
		}
	}

	expected := []string{"new0", "new1", "new2"}
	for i, e := range expected {
		if node.Tags[i] != e {
			t.Errorf("expected tags[%d]=%q, got %q", i, e, node.Tags[i])
		}
	}
}

// Test nested Unset in array then parent array access
func TestPatcher_UnsetArrayElementThenParent(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Section1",
					Blocks: []domain.Block{
						{Type: "text"},
						{Type: "image"},
					},
				},
			},
		},
	}

	// Unset block from within section
	err := p.Unset(&node, "content.sections[0].blocks[1]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Content.Sections[0].Blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(node.Content.Sections[0].Blocks))
	}

	if node.Content.Sections[0].Blocks[0].Type != "text" {
		t.Errorf("expected remaining block type 'text', got %q",
			node.Content.Sections[0].Blocks[0].Type)
	}
}

// ===== CONSTRAINT TESTS =====

// Test operations on Constraints slice
func TestPatcher_AppendToConstraints(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:          "test",
		Kind:        "system",
		Version:     1,
		Status:      "draft",
		Title:       "Test",
		Constraints: []domain.Constraint{},
	}

	newConstraint := domain.Constraint{
		Expr:    "health > 0",
		Message: "Health must be positive",
		Scope:   "all",
	}

	err := p.Append(&node, "constraints", newConstraint)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Constraints) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(node.Constraints))
	}

	if node.Constraints[0].Expr != "health > 0" {
		t.Errorf("expected constraint expr 'health > 0', got %q", node.Constraints[0].Expr)
	}
}

// ===== REVIEWERS TESTS =====

// Test operations on Reviewers slice
func TestPatcher_AppendToReviewers(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:        "test",
		Kind:      "system",
		Version:   1,
		Status:    "draft",
		Title:     "Test",
		Reviewers: []domain.Reviewer{},
	}

	newReviewer := domain.Reviewer{
		Name:    "john@example.com",
		Version: 1,
	}

	err := p.Append(&node, "reviewers", newReviewer)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Reviewers) != 1 {
		t.Fatalf("expected 1 reviewer, got %d", len(node.Reviewers))
	}

	if node.Reviewers[0].Name != "john@example.com" {
		t.Errorf("expected reviewer name 'john@example.com', got %q", node.Reviewers[0].Name)
	}
}

// ===== CUSTOM MAP TESTS =====

// Test Set and Unset on Custom map
func TestPatcher_SetCustomMapValue(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Custom:  map[string]interface{}{},
	}

	err := p.Set(&node, "custom.mykey", "myvalue")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Custom["mykey"] != "myvalue" {
		t.Errorf("expected custom[mykey]='myvalue', got %v", node.Custom["mykey"])
	}
}

func TestPatcher_UnsetCustomMapKey(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Custom:  map[string]interface{}{"key1": "val1", "key2": "val2"},
	}

	err := p.Unset(&node, "custom.key1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, exists := node.Custom["key1"]; exists {
		t.Error("expected key1 to be removed from custom map")
	}
}

// ===== DEEP ARRAY TRAVERSAL TESTS =====

// Test Set through multiple nested arrays
func TestPatcher_SetThroughNestedArrays(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Section1",
					Blocks: []domain.Block{
						{Type: "text"},
						{Type: "image"},
					},
				},
				{
					Name: "Section2",
					Blocks: []domain.Block{
						{Type: "table"},
					},
				},
			},
		},
	}

	// Set through second section, first block
	err := p.Set(&node, "content.sections[1].blocks[0].type", "updated-table")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Content.Sections[1].Blocks[0].Type != "updated-table" {
		t.Errorf("expected 'updated-table', got %q", node.Content.Sections[1].Blocks[0].Type)
	}
}

// Test Unset from deeply nested structure
func TestPatcher_UnsetFromNestedStructure(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Section1",
					Blocks: []domain.Block{
						{Type: "text", Data: map[string]interface{}{"key": "value"}},
					},
				},
			},
		},
	}

	// Unset the entire blocks array of section
	err := p.Unset(&node, "content.sections[0].blocks")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Content.Sections[0].Blocks) != 0 {
		t.Errorf("expected empty blocks, got %d blocks", len(node.Content.Sections[0].Blocks))
	}
}

// ===== REFLINKS TESTS =====

// Test operations on RefLinks (Uses field in Refs)
func TestPatcher_AppendToRefLinks(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Refs:    domain.Ref{Uses: []domain.RefLink{}},
	}

	newLink := domain.RefLink{Target: "other-node", Context: "dependency"}
	err := p.Append(&node, "refs.uses", newLink)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Refs.Uses) != 1 {
		t.Fatalf("expected 1 ref link, got %d", len(node.Refs.Uses))
	}

	if node.Refs.Uses[0].Target != "other-node" {
		t.Errorf("expected target 'other-node', got %q", node.Refs.Uses[0].Target)
	}
}

// Test Set on RefLink field through array
func TestPatcher_SetRefLinkField(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "old-target", Context: "old-context"}}},
	}

	err := p.Set(&node, "refs.uses[0].context", "new-context")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Refs.Uses[0].Context != "new-context" {
		t.Errorf("expected 'new-context', got %q", node.Refs.Uses[0].Context)
	}
}

// ===== LLMCONTEXT TESTS =====

// Test Set on LLMContext field (uses exact Go field name casing)
func TestPatcher_SetLLMContext(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:         "test",
		Kind:       "system",
		Version:    1,
		Status:     "draft",
		Title:      "Test",
		LLMContext: "",
	}

	// The field is named LLMContext in Go, so we use lLMContext (first letter lowercased)
	// since capitalizeFirst will capitalize only the first letter
	err := p.Set(&node, "lLMContext", "This is context for AI")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.LLMContext != "This is context for AI" {
		t.Errorf("expected 'This is context for AI', got %q", node.LLMContext)
	}
}

// ===== SOURCEFILE TESTS =====

// Test Set on SourceFile (unexported in JSON but accessible)
func TestPatcher_SetSourceFile(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:         "test",
		Kind:       "system",
		Version:    1,
		Status:     "draft",
		Title:      "Test",
		SourceFile: "",
	}

	err := p.Set(&node, "sourceFile", "/path/to/file.yaml")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.SourceFile != "/path/to/file.yaml" {
		t.Errorf("expected '/path/to/file.yaml', got %q", node.SourceFile)
	}
}

// ===== BOUNDARY CONDITIONS =====

// Test very long path traversal
func TestPatcher_VeryLongPath(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Section",
					Blocks: []domain.Block{
						{Type: "text", Data: map[string]interface{}{"nested": "value"}},
					},
				},
			},
		},
	}

	// This is a moderately deep path
	err := p.Set(&node, "content.sections[0].blocks[0].data.nested", "updated")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Content.Sections[0].Blocks[0].Data["nested"] != "updated" {
		t.Errorf("expected 'updated', got %v", node.Content.Sections[0].Blocks[0].Data["nested"])
	}
}

// Test single element array operations
func TestPatcher_SingleElementArrayOperations(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{"only-tag"},
	}

	// Unset the only element
	err := p.Unset(&node, "tags[0]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Tags) != 0 {
		t.Errorf("expected empty tags, got %v", node.Tags)
	}
}

// ===== EMPTY ARRAY OPERATIONS =====

// Test Set on empty array returns error
func TestPatcher_SetOnEmptyArray(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{},
	}

	err := p.Set(&node, "tags[0]", "value")
	if err == nil {
		t.Fatal("expected error for index on empty array, got nil")
	}
}

// Test Unset on empty array does nothing
func TestPatcher_UnsetOnEmptyArray(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Tags:    []string{},
	}

	// Should not error
	err := p.Unset(&node, "tags[0]")
	if err != nil {
		t.Errorf("expected no error for unset on empty array, got %v", err)
	}
}

// ===== APPEND TO NIL SLICE TESTS =====

// Test Append to nil refs.uses slice
func TestPatcher_AppendToNilSlice(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:      "test",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Refs:    domain.Ref{}, // Uses is nil
	}

	newLink := domain.RefLink{Target: "target"}
	err := p.Append(&node, "refs.uses", newLink)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(node.Refs.Uses) != 1 {
		t.Fatalf("expected 1 ref link, got %d", len(node.Refs.Uses))
	}
}

// ===== APPLY WITH MIXED COMPLEX OPERATIONS =====

// Test Apply with operations that touch many different field types
func TestPatcher_ApplyComplexMixedOperations(t *testing.T) {
	p := patcher.New()

	node := domain.Node{
		ID:         "test",
		Kind:       "system",
		Version:    1,
		Status:     "draft",
		Title:      "Test",
		Tags:       []string{"tag1"},
		Glossary:   map[string]string{"term": "def"},
		LLMContext: "original context",
		Content: &domain.Content{
			Sections: []domain.Section{
				{Name: "Section1", Blocks: []domain.Block{}},
			},
		},
	}

	ops := []patcher.PatchOperation{
		{Op: "set", Path: "title", Value: "Updated Title"},
		{Op: "append", Path: "tags", Value: "tag2"},
		{Op: "set", Path: "glossary.newterm", Value: "newdef"},
		{Op: "set", Path: "lLMContext", Value: "updated context"},
		{Op: "set", Path: "content.sections[0].name", Value: "Updated Section"},
	}

	err := p.Apply(&node, ops)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %q", node.Title)
	}
	if len(node.Tags) != 2 || node.Tags[1] != "tag2" {
		t.Errorf("expected tags [tag1, tag2], got %v", node.Tags)
	}
	if node.Glossary["newterm"] != "newdef" {
		t.Errorf("expected glossary[newterm]='newdef', got %v", node.Glossary["newterm"])
	}
	if node.LLMContext != "updated context" {
		t.Errorf("expected LLMContext 'updated context', got %q", node.LLMContext)
	}
	if node.Content.Sections[0].Name != "Updated Section" {
		t.Errorf("expected section name 'Updated Section', got %q", node.Content.Sections[0].Name)
	}
}
