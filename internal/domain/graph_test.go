package domain_test

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestGraph_EmptyGraph(t *testing.T) {
	g := domain.NewGraph()

	if g.Count() != 0 {
		t.Errorf("expected empty graph to have count 0, got %d", g.Count())
	}

	nodes := g.All()
	if len(nodes) != 0 {
		t.Errorf("expected empty graph to return empty slice, got %d nodes", len(nodes))
	}
}

func TestGraph_AddNode(t *testing.T) {
	g := domain.NewGraph()

	node := domain.Node{
		ID:      "test/node1",
		Kind:    "mechanic",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node 1",
	}

	err := g.Add(node)
	if err != nil {
		t.Fatalf("failed to add node: %v", err)
	}

	if g.Count() != 1 {
		t.Errorf("expected count 1 after adding node, got %d", g.Count())
	}
}

func TestGraph_LookupByID(t *testing.T) {
	g := domain.NewGraph()

	node := domain.Node{
		ID:      "test/node1",
		Kind:    "mechanic",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node 1",
	}

	if err := g.Add(node); err != nil {
		t.Fatalf("failed to add node: %v", err)
	}

	// Lookup existing node
	found, exists := g.Get("test/node1")
	if !exists {
		t.Errorf("expected to find node 'test/node1'")
	}
	if found.ID != "test/node1" {
		t.Errorf("expected found node ID 'test/node1', got %q", found.ID)
	}
	if found.Title != "Test Node 1" {
		t.Errorf("expected found node Title 'Test Node 1', got %q", found.Title)
	}

	// Lookup non-existing node
	_, exists = g.Get("test/nonexistent")
	if exists {
		t.Errorf("expected not to find node 'test/nonexistent'")
	}
}

func TestGraph_DuplicateID(t *testing.T) {
	g := domain.NewGraph()

	node1 := domain.Node{
		ID:      "test/node1",
		Kind:    "mechanic",
		Version: 1,
		Status:  "draft",
		Title:   "First Node",
	}

	node2 := domain.Node{
		ID:      "test/node1",
		Kind:    "system",
		Version: 2,
		Status:  "approved",
		Title:   "Duplicate ID",
	}

	// Add first node
	if err := g.Add(node1); err != nil {
		t.Fatalf("failed to add first node: %v", err)
	}

	// Attempt to add duplicate ID
	err := g.Add(node2)
	if err == nil {
		t.Errorf("expected error when adding duplicate ID, got nil")
	}

	// Verify original node is unchanged
	found, exists := g.Get("test/node1")
	if !exists {
		t.Fatalf("node should still exist")
	}
	if found.Title != "First Node" {
		t.Errorf("expected original node to be unchanged, got Title %q", found.Title)
	}
}

func TestGraph_Iteration(t *testing.T) {
	g := domain.NewGraph()

	nodes := []domain.Node{
		{ID: "test/node1", Kind: "mechanic", Version: 1, Status: "draft", Title: "Node 1"},
		{ID: "test/node2", Kind: "system", Version: 1, Status: "approved", Title: "Node 2"},
		{ID: "test/node3", Kind: "feature", Version: 2, Status: "draft", Title: "Node 3"},
	}

	for _, node := range nodes {
		if err := g.Add(node); err != nil {
			t.Fatalf("failed to add node %q: %v", node.ID, err)
		}
	}

	// Get all nodes
	all := g.All()
	if len(all) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(all))
	}

	// Verify all added nodes are present
	seen := make(map[string]bool)
	for _, node := range all {
		seen[node.ID] = true
	}

	for _, expected := range nodes {
		if !seen[expected.ID] {
			t.Errorf("expected to find node %q in iteration", expected.ID)
		}
	}
}

func TestGraph_RemoveNode(t *testing.T) {
	g := domain.NewGraph()

	node := domain.Node{
		ID:      "test/node1",
		Kind:    "mechanic",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
	}

	if err := g.Add(node); err != nil {
		t.Fatalf("failed to add node: %v", err)
	}

	if g.Count() != 1 {
		t.Errorf("expected count 1, got %d", g.Count())
	}

	// Remove the node
	removed := g.Remove("test/node1")
	if !removed {
		t.Errorf("expected Remove to return true for existing node")
	}

	if g.Count() != 0 {
		t.Errorf("expected count 0 after removal, got %d", g.Count())
	}

	// Try to remove non-existent node
	removed = g.Remove("test/nonexistent")
	if removed {
		t.Errorf("expected Remove to return false for non-existent node")
	}
}

func TestGraph_Update(t *testing.T) {
	g := domain.NewGraph()

	original := domain.Node{
		ID:      "test/node1",
		Kind:    "mechanic",
		Version: 1,
		Status:  "draft",
		Title:   "Original Title",
	}

	if err := g.Add(original); err != nil {
		t.Fatalf("failed to add node: %v", err)
	}

	// Update the node
	updated := domain.Node{
		ID:      "test/node1",
		Kind:    "mechanic",
		Version: 2,
		Status:  "approved",
		Title:   "Updated Title",
	}

	if err := g.Update(updated); err != nil {
		t.Fatalf("failed to update node: %v", err)
	}

	// Verify update
	found, exists := g.Get("test/node1")
	if !exists {
		t.Fatalf("node should exist after update")
	}
	if found.Version != 2 {
		t.Errorf("expected Version 2, got %d", found.Version)
	}
	if found.Status != "approved" {
		t.Errorf("expected Status 'approved', got %q", found.Status)
	}
	if found.Title != "Updated Title" {
		t.Errorf("expected Title 'Updated Title', got %q", found.Title)
	}
}
