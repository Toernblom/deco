package graph_test

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/graph"
)

// Test building reverse index with no references
func TestReverseIndex_NoReferences(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "a", Kind: "system", Version: 1, Status: "draft", Title: "A"},
		{ID: "b", Kind: "system", Version: 1, Status: "draft", Title: "B"},
	}

	g, _ := builder.Build(nodes)
	reverseIndex := builder.BuildReverseIndex(g)

	// Both nodes should have empty reverse ref lists
	if len(reverseIndex["a"]) != 0 {
		t.Errorf("expected no reverse refs for 'a', got %d", len(reverseIndex["a"]))
	}
	if len(reverseIndex["b"]) != 0 {
		t.Errorf("expected no reverse refs for 'b', got %d", len(reverseIndex["b"]))
	}
}

// Test simple reverse reference (A uses B -> B is used by A)
func TestReverseIndex_SimpleReference(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "a", Kind: "system", Version: 1, Status: "draft", Title: "A"},
		{
			ID:      "b",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "B",
			Refs: domain.Ref{
				Uses: []domain.RefLink{{Target: "a"}},
			},
		},
	}

	g, _ := builder.Build(nodes)
	reverseIndex := builder.BuildReverseIndex(g)

	// 'a' should be used by 'b'
	usedBy := reverseIndex["a"]
	if len(usedBy) != 1 {
		t.Fatalf("expected 1 reverse ref for 'a', got %d", len(usedBy))
	}
	if usedBy[0] != "b" {
		t.Errorf("expected 'a' to be used by 'b', got %q", usedBy[0])
	}

	// 'b' should have no reverse refs
	if len(reverseIndex["b"]) != 0 {
		t.Errorf("expected no reverse refs for 'b', got %d", len(reverseIndex["b"]))
	}
}

// Test multiple nodes referencing the same node
func TestReverseIndex_MultipleReferences(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "a", Kind: "system", Version: 1, Status: "draft", Title: "A"},
		{
			ID:      "b",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "B",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "a"}}},
		},
		{
			ID:      "c",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "C",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "a"}}},
		},
		{
			ID:      "d",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "D",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "a"}}},
		},
	}

	g, _ := builder.Build(nodes)
	reverseIndex := builder.BuildReverseIndex(g)

	// 'a' should be used by b, c, and d
	usedBy := reverseIndex["a"]
	if len(usedBy) != 3 {
		t.Fatalf("expected 3 reverse refs for 'a', got %d", len(usedBy))
	}

	// Check all three are present
	foundB, foundC, foundD := false, false, false
	for _, nodeID := range usedBy {
		switch nodeID {
		case "b":
			foundB = true
		case "c":
			foundC = true
		case "d":
			foundD = true
		}
	}

	if !foundB || !foundC || !foundD {
		t.Errorf("expected 'a' to be used by b, c, and d, got %v", usedBy)
	}
}

// Test circular references are handled
func TestReverseIndex_CircularReferences(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{
			ID:      "a",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "A",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "b"}}},
		},
		{
			ID:      "b",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "B",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "a"}}},
		},
	}

	g, _ := builder.Build(nodes)
	reverseIndex := builder.BuildReverseIndex(g)

	// 'a' should be used by 'b'
	if len(reverseIndex["a"]) != 1 || reverseIndex["a"][0] != "b" {
		t.Errorf("expected 'a' to be used by 'b', got %v", reverseIndex["a"])
	}

	// 'b' should be used by 'a'
	if len(reverseIndex["b"]) != 1 || reverseIndex["b"][0] != "a" {
		t.Errorf("expected 'b' to be used by 'a', got %v", reverseIndex["b"])
	}
}

// Test orphan nodes (nodes not referenced by anyone)
func TestReverseIndex_OrphanNodes(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "orphan1", Kind: "system", Version: 1, Status: "draft", Title: "Orphan 1"},
		{ID: "orphan2", Kind: "system", Version: 1, Status: "draft", Title: "Orphan 2"},
		{
			ID:      "connected",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Connected",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "orphan1"}}},
		},
	}

	g, _ := builder.Build(nodes)
	reverseIndex := builder.BuildReverseIndex(g)

	// orphan1 should be used by 'connected'
	if len(reverseIndex["orphan1"]) != 1 {
		t.Errorf("expected 1 reverse ref for orphan1, got %d", len(reverseIndex["orphan1"]))
	}

	// orphan2 should have no reverse refs (true orphan)
	if len(reverseIndex["orphan2"]) != 0 {
		t.Errorf("expected no reverse refs for orphan2, got %d", len(reverseIndex["orphan2"]))
	}

	// connected should have no reverse refs
	if len(reverseIndex["connected"]) != 0 {
		t.Errorf("expected no reverse refs for connected, got %d", len(reverseIndex["connected"]))
	}
}

// Test both Uses and Related references are tracked
func TestReverseIndex_UsesAndRelated(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{
			ID:      "user",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "User",
			Refs: domain.Ref{
				Uses: []domain.RefLink{{Target: "target"}},
			},
		},
		{
			ID:      "related",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Related",
			Refs: domain.Ref{
				Related: []domain.RefLink{{Target: "target"}},
			},
		},
	}

	g, _ := builder.Build(nodes)
	reverseIndex := builder.BuildReverseIndex(g)

	// 'target' should be referenced by both 'user' and 'related'
	refs := reverseIndex["target"]
	if len(refs) != 2 {
		t.Fatalf("expected 2 reverse refs for target, got %d", len(refs))
	}

	foundUser, foundRelated := false, false
	for _, nodeID := range refs {
		if nodeID == "user" {
			foundUser = true
		}
		if nodeID == "related" {
			foundRelated = true
		}
	}

	if !foundUser || !foundRelated {
		t.Errorf("expected target to be referenced by both user and related, got %v", refs)
	}
}

// Test node referencing itself
func TestReverseIndex_SelfReference(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{
			ID:      "self",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Self",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "self"}}},
		},
	}

	g, _ := builder.Build(nodes)
	reverseIndex := builder.BuildReverseIndex(g)

	// 'self' should reference itself
	refs := reverseIndex["self"]
	if len(refs) != 1 {
		t.Fatalf("expected 1 reverse ref for self, got %d", len(refs))
	}
	if refs[0] != "self" {
		t.Errorf("expected self to reference itself, got %q", refs[0])
	}
}

// Test node with multiple reference types to same target
func TestReverseIndex_MultipleRefTypes(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "target", Kind: "system", Version: 1, Status: "draft", Title: "Target"},
		{
			ID:      "source",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Source",
			Refs: domain.Ref{
				Uses:    []domain.RefLink{{Target: "target"}},
				Related: []domain.RefLink{{Target: "target"}},
			},
		},
	}

	g, _ := builder.Build(nodes)
	reverseIndex := builder.BuildReverseIndex(g)

	// 'target' should be referenced by 'source' (deduplicated)
	refs := reverseIndex["target"]
	if len(refs) != 1 {
		t.Errorf("expected 1 reverse ref for target (deduplicated), got %d", len(refs))
	}
	if len(refs) > 0 && refs[0] != "source" {
		t.Errorf("expected target to be referenced by source, got %q", refs[0])
	}
}

// Test dangling references (references to non-existent nodes)
func TestReverseIndex_DanglingReferences(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{
			ID:      "source",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Source",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "nonexistent"}}},
		},
	}

	g, _ := builder.Build(nodes)
	reverseIndex := builder.BuildReverseIndex(g)

	// 'nonexistent' should not be in the index (it's not a node)
	if _, exists := reverseIndex["nonexistent"]; exists {
		t.Error("expected dangling reference target to not be in index")
	}

	// 'source' should have no reverse refs
	if len(reverseIndex["source"]) != 0 {
		t.Errorf("expected no reverse refs for source, got %d", len(reverseIndex["source"]))
	}
}

// Test reverse index is consistent with dependency map
func TestReverseIndex_ConsistentWithDeps(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "a", Kind: "system", Version: 1, Status: "draft", Title: "A"},
		{
			ID:      "b",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "B",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "a"}}},
		},
		{
			ID:      "c",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "C",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "a"}, {Target: "b"}}},
		},
	}

	g, _ := builder.Build(nodes)
	depMap := builder.BuildDependencyMap(g)
	reverseIndex := builder.BuildReverseIndex(g)

	// For each dependency, check reverse index consistency
	for nodeID, deps := range depMap {
		for _, depID := range deps {
			// depID should have nodeID in its reverse refs
			found := false
			for _, refBy := range reverseIndex[depID] {
				if refBy == nodeID {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("inconsistency: %q depends on %q, but %q doesn't list %q in reverse refs",
					nodeID, depID, depID, nodeID)
			}
		}
	}
}
