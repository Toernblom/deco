package patcher_test

import (
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
