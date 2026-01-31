package graph

import (
	"fmt"

	"github.com/Toernblom/deco/internal/domain"
)

// Builder constructs dependency graphs from nodes.
type Builder struct{}

// NewBuilder creates a new graph builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// Build creates a Graph from a slice of nodes.
// Returns an error if duplicate node IDs are found.
func (b *Builder) Build(nodes []domain.Node) (*domain.Graph, error) {
	g := domain.NewGraph()

	for _, node := range nodes {
		if err := g.Add(node); err != nil {
			return nil, fmt.Errorf("failed to add node %q: %w", node.ID, err)
		}
	}

	return g, nil
}

// BuildDependencyMap creates a map of node dependencies.
// For each node, returns a list of node IDs it depends on (via Refs.Uses).
// Related refs are NOT considered dependencies.
func (b *Builder) BuildDependencyMap(g *domain.Graph) map[string][]string {
	depMap := make(map[string][]string)

	for _, node := range g.All() {
		deps := make([]string, 0)

		// Only Uses refs create dependencies
		for _, refLink := range node.Refs.Uses {
			deps = append(deps, refLink.Target)
		}

		depMap[node.ID] = deps
	}

	return depMap
}

// DetectCycle checks for circular dependencies in the graph.
// Returns true and the cycle path if a cycle is found, false and nil otherwise.
func (b *Builder) DetectCycle(g *domain.Graph) (bool, []string) {
	depMap := b.BuildDependencyMap(g)

	// Track visiting state: 0 = unvisited, 1 = visiting, 2 = visited
	state := make(map[string]int)
	path := make([]string, 0)

	var visit func(nodeID string) bool

	visit = func(nodeID string) bool {
		if state[nodeID] == 2 {
			// Already fully visited, no cycle here
			return false
		}

		if state[nodeID] == 1 {
			// Currently visiting - cycle detected!
			path = append(path, nodeID)
			return true
		}

		// Mark as visiting
		state[nodeID] = 1
		path = append(path, nodeID)

		// Visit all dependencies
		for _, depID := range depMap[nodeID] {
			if visit(depID) {
				return true
			}
		}

		// Mark as fully visited
		state[nodeID] = 2
		path = path[:len(path)-1] // Remove from path
		return false
	}

	// Check each node
	for _, node := range g.All() {
		if state[node.ID] == 0 {
			if visit(node.ID) {
				return true, path
			}
		}
	}

	return false, nil
}

// TopologicalSort returns nodes sorted in dependency order.
// Nodes with no dependencies come first; nodes appear only after all their dependencies.
// Returns an error if the graph contains a cycle.
func (b *Builder) TopologicalSort(g *domain.Graph) ([]domain.Node, error) {
	// First check for cycles
	if hasCycle, cycle := b.DetectCycle(g); hasCycle {
		return nil, fmt.Errorf("cannot sort graph with cycle: %v", cycle)
	}

	depMap := b.BuildDependencyMap(g)

	// Calculate in-degree for each node
	// In-degree = number of dependencies this node has
	inDegree := make(map[string]int)
	for nodeID, deps := range depMap {
		inDegree[nodeID] = len(deps)
	}

	// Find all nodes with no dependencies (in-degree 0)
	queue := make([]string, 0)
	for nodeID, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, nodeID)
		}
	}

	// Build reverse dependency map (who depends on me?)
	// This tells us which nodes to check when we process a node
	reverseDeps := make(map[string][]string)
	for nodeID, deps := range depMap {
		for _, depID := range deps {
			reverseDeps[depID] = append(reverseDeps[depID], nodeID)
		}
	}

	// Process queue
	result := make([]domain.Node, 0, g.Count())

	for len(queue) > 0 {
		// Pop from queue
		nodeID := queue[0]
		queue = queue[1:]

		// Add to result
		node, _ := g.Get(nodeID)
		result = append(result, node)

		// Reduce in-degree for nodes that depend on this one
		for _, dependentID := range reverseDeps[nodeID] {
			inDegree[dependentID]--
			if inDegree[dependentID] == 0 {
				queue = append(queue, dependentID)
			}
		}
	}

	return result, nil
}

// BuildReverseIndex creates a reverse reference index.
// For each node, returns a list of node IDs that reference it (via Uses or Related).
// This is useful for finding all nodes that depend on or are related to a given node.
func (b *Builder) BuildReverseIndex(g *domain.Graph) map[string][]string {
	reverseIndex := make(map[string][]string)

	// Initialize empty lists for all nodes
	for _, node := range g.All() {
		reverseIndex[node.ID] = make([]string, 0)
	}

	// Build reverse references
	for _, node := range g.All() {
		// Track which targets we've already added for this node (for deduplication)
		seen := make(map[string]bool)

		// Process Uses references
		for _, refLink := range node.Refs.Uses {
			target := refLink.Target
			// Only add if target exists in graph and we haven't already added this reference
			if _, exists := g.Get(target); exists && !seen[target] {
				reverseIndex[target] = append(reverseIndex[target], node.ID)
				seen[target] = true
			}
		}

		// Process Related references
		for _, refLink := range node.Refs.Related {
			target := refLink.Target
			// Only add if target exists in graph and we haven't already added this reference
			if _, exists := g.Get(target); exists && !seen[target] {
				reverseIndex[target] = append(reverseIndex[target], node.ID)
				seen[target] = true
			}
		}
	}

	return reverseIndex
}
