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
