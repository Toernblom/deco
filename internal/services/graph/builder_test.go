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

package graph_test

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/graph"
)

// Test building empty graph from empty slice
func TestGraphBuilder_BuildFromEmpty(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{}
	g, err := builder.Build(nodes)

	if err != nil {
		t.Fatalf("expected no error building empty graph, got %v", err)
	}

	if g.Count() != 0 {
		t.Errorf("expected empty graph, got %d nodes", g.Count())
	}
}

// Test building graph from single node
func TestGraphBuilder_BuildFromSingle(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "systems/food", Kind: "system", Version: 1, Status: "draft", Title: "Food System"},
	}

	g, err := builder.Build(nodes)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Count() != 1 {
		t.Errorf("expected 1 node, got %d", g.Count())
	}

	node, exists := g.Get("systems/food")
	if !exists {
		t.Fatal("expected to find node 'systems/food'")
	}

	if node.ID != "systems/food" {
		t.Errorf("expected ID 'systems/food', got %q", node.ID)
	}
}

// Test building graph from multiple nodes
func TestGraphBuilder_BuildFromMultiple(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "systems/food", Kind: "system", Version: 1, Status: "draft", Title: "Food"},
		{ID: "systems/health", Kind: "system", Version: 1, Status: "draft", Title: "Health"},
		{ID: "mechanics/hunger", Kind: "mechanic", Version: 1, Status: "draft", Title: "Hunger"},
	}

	g, err := builder.Build(nodes)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if g.Count() != 3 {
		t.Errorf("expected 3 nodes, got %d", g.Count())
	}

	// Verify all nodes are present
	for _, expectedNode := range nodes {
		node, exists := g.Get(expectedNode.ID)
		if !exists {
			t.Errorf("expected to find node %q", expectedNode.ID)
			continue
		}
		if node.ID != expectedNode.ID {
			t.Errorf("node ID mismatch: expected %q, got %q", expectedNode.ID, node.ID)
		}
	}
}

// Test error when duplicate node IDs are provided
func TestGraphBuilder_DuplicateNodes(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "systems/food", Kind: "system", Version: 1, Status: "draft", Title: "Food v1"},
		{ID: "systems/food", Kind: "system", Version: 2, Status: "approved", Title: "Food v2"},
	}

	_, err := builder.Build(nodes)

	if err == nil {
		t.Fatal("expected error for duplicate node IDs, got nil")
	}
}

// Test building dependency map with no dependencies
func TestGraphBuilder_DependencyMapNoDeps(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "systems/food", Kind: "system", Version: 1, Status: "draft", Title: "Food"},
		{ID: "systems/health", Kind: "system", Version: 1, Status: "draft", Title: "Health"},
	}

	g, _ := builder.Build(nodes)
	depMap := builder.BuildDependencyMap(g)

	// Should have entries for all nodes
	if len(depMap) != 2 {
		t.Errorf("expected 2 entries in dependency map, got %d", len(depMap))
	}

	// Both should have empty dependency lists
	if len(depMap["systems/food"]) != 0 {
		t.Errorf("expected no dependencies for systems/food, got %d", len(depMap["systems/food"]))
	}
	if len(depMap["systems/health"]) != 0 {
		t.Errorf("expected no dependencies for systems/health, got %d", len(depMap["systems/health"]))
	}
}

// Test building dependency map with simple dependencies
func TestGraphBuilder_DependencyMapSimple(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{
			ID:      "systems/food",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Food",
		},
		{
			ID:      "mechanics/hunger",
			Kind:    "mechanic",
			Version: 1,
			Status:  "draft",
			Title:   "Hunger",
			Refs: domain.Ref{
				Uses: []domain.RefLink{
					{Target: "systems/food"},
				},
			},
		},
	}

	g, _ := builder.Build(nodes)
	depMap := builder.BuildDependencyMap(g)

	// Hunger should depend on Food
	deps := depMap["mechanics/hunger"]
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency for mechanics/hunger, got %d", len(deps))
	}
	if deps[0] != "systems/food" {
		t.Errorf("expected dependency on 'systems/food', got %q", deps[0])
	}

	// Food should have no dependencies
	if len(depMap["systems/food"]) != 0 {
		t.Errorf("expected no dependencies for systems/food, got %d", len(depMap["systems/food"]))
	}
}

// Test dependency map with multiple dependencies
func TestGraphBuilder_DependencyMapMultiple(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "systems/food", Kind: "system", Version: 1, Status: "draft", Title: "Food"},
		{ID: "systems/health", Kind: "system", Version: 1, Status: "draft", Title: "Health"},
		{
			ID:      "mechanics/survival",
			Kind:    "mechanic",
			Version: 1,
			Status:  "draft",
			Title:   "Survival",
			Refs: domain.Ref{
				Uses: []domain.RefLink{
					{Target: "systems/food"},
					{Target: "systems/health"},
				},
			},
		},
	}

	g, _ := builder.Build(nodes)
	depMap := builder.BuildDependencyMap(g)

	deps := depMap["mechanics/survival"]
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(deps))
	}

	// Check both dependencies are present
	foundFood := false
	foundHealth := false
	for _, dep := range deps {
		if dep == "systems/food" {
			foundFood = true
		}
		if dep == "systems/health" {
			foundHealth = true
		}
	}

	if !foundFood {
		t.Error("expected dependency on systems/food")
	}
	if !foundHealth {
		t.Error("expected dependency on systems/health")
	}
}

// Test cycle detection with no cycles
func TestGraphBuilder_NoCycle(t *testing.T) {
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
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "b"}}},
		},
	}

	g, _ := builder.Build(nodes)

	hasCycle, cycle := builder.DetectCycle(g)

	if hasCycle {
		t.Errorf("expected no cycle, but found one: %v", cycle)
	}
}

// Test cycle detection with simple cycle (A -> B -> A)
func TestGraphBuilder_SimpleCycle(t *testing.T) {
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

	hasCycle, cycle := builder.DetectCycle(g)

	if !hasCycle {
		t.Fatal("expected cycle to be detected")
	}

	if len(cycle) < 2 {
		t.Errorf("expected cycle path with at least 2 nodes, got %d", len(cycle))
	}
}

// Test cycle detection with self-reference (A -> A)
func TestGraphBuilder_SelfCycle(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{
			ID:      "a",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "A",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "a"}}},
		},
	}

	g, _ := builder.Build(nodes)

	hasCycle, cycle := builder.DetectCycle(g)

	if !hasCycle {
		t.Fatal("expected self-reference cycle to be detected")
	}

	if len(cycle) == 0 {
		t.Error("expected non-empty cycle path")
	}
}

// Test cycle detection with longer cycle (A -> B -> C -> A)
func TestGraphBuilder_LongCycle(t *testing.T) {
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
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "c"}}},
		},
		{
			ID:      "c",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "C",
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "a"}}},
		},
	}

	g, _ := builder.Build(nodes)

	hasCycle, cycle := builder.DetectCycle(g)

	if !hasCycle {
		t.Fatal("expected cycle to be detected")
	}

	if len(cycle) < 3 {
		t.Errorf("expected cycle path with at least 3 nodes, got %d: %v", len(cycle), cycle)
	}
}

// Test topological sort with no dependencies
func TestGraphBuilder_TopologicalSortNoDeps(t *testing.T) {
	builder := graph.NewBuilder()

	nodes := []domain.Node{
		{ID: "a", Kind: "system", Version: 1, Status: "draft", Title: "A"},
		{ID: "b", Kind: "system", Version: 1, Status: "draft", Title: "B"},
		{ID: "c", Kind: "system", Version: 1, Status: "draft", Title: "C"},
	}

	g, _ := builder.Build(nodes)

	sorted, err := builder.TopologicalSort(g)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sorted) != 3 {
		t.Errorf("expected 3 nodes in sorted list, got %d", len(sorted))
	}

	// All nodes should be present
	ids := make(map[string]bool)
	for _, node := range sorted {
		ids[node.ID] = true
	}

	if !ids["a"] || !ids["b"] || !ids["c"] {
		t.Error("expected all nodes to be in sorted list")
	}
}

// Test topological sort with simple dependency chain
func TestGraphBuilder_TopologicalSortChain(t *testing.T) {
	builder := graph.NewBuilder()

	// Chain: a <- b <- c (c depends on b, b depends on a)
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
			Refs:    domain.Ref{Uses: []domain.RefLink{{Target: "b"}}},
		},
	}

	g, _ := builder.Build(nodes)

	sorted, err := builder.TopologicalSort(g)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sorted) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(sorted))
	}

	// Build position map
	pos := make(map[string]int)
	for i, node := range sorted {
		pos[node.ID] = i
	}

	// a should come before b
	if pos["a"] >= pos["b"] {
		t.Errorf("expected 'a' before 'b', got positions a=%d, b=%d", pos["a"], pos["b"])
	}

	// b should come before c
	if pos["b"] >= pos["c"] {
		t.Errorf("expected 'b' before 'c', got positions b=%d, c=%d", pos["b"], pos["c"])
	}
}

// Test topological sort with diamond dependency
func TestGraphBuilder_TopologicalSortDiamond(t *testing.T) {
	builder := graph.NewBuilder()

	// Diamond: a <- b, a <- c, b <- d, c <- d
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
			Refs: domain.Ref{Uses: []domain.RefLink{
				{Target: "b"},
				{Target: "c"},
			}},
		},
	}

	g, _ := builder.Build(nodes)

	sorted, err := builder.TopologicalSort(g)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sorted) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(sorted))
	}

	// Build position map
	pos := make(map[string]int)
	for i, node := range sorted {
		pos[node.ID] = i
	}

	// a should come before b and c
	if pos["a"] >= pos["b"] {
		t.Errorf("expected 'a' before 'b'")
	}
	if pos["a"] >= pos["c"] {
		t.Errorf("expected 'a' before 'c'")
	}

	// b and c should come before d
	if pos["b"] >= pos["d"] {
		t.Errorf("expected 'b' before 'd'")
	}
	if pos["c"] >= pos["d"] {
		t.Errorf("expected 'c' before 'd'")
	}
}

// Test topological sort fails on cycle
func TestGraphBuilder_TopologicalSortWithCycle(t *testing.T) {
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

	_, err := builder.TopologicalSort(g)

	if err == nil {
		t.Fatal("expected error when cycle is present, got nil")
	}
}

// Test that Related refs don't create dependencies
func TestGraphBuilder_RelatedRefsNotDependencies(t *testing.T) {
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
				Related: []domain.RefLink{{Target: "a"}},
			},
		},
	}

	g, _ := builder.Build(nodes)
	depMap := builder.BuildDependencyMap(g)

	// 'b' should have no dependencies (Related is not a dependency)
	if len(depMap["b"]) != 0 {
		t.Errorf("expected no dependencies for 'b', got %d", len(depMap["b"]))
	}
}
