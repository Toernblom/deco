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

package domain

import "fmt"

// Graph represents a collection of nodes with efficient lookup.
// It provides methods for adding, retrieving, and iterating over nodes.
type Graph struct {
	nodes map[string]Node
}

// NewGraph creates a new empty graph.
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]Node),
	}
}

// Count returns the number of nodes in the graph.
func (g *Graph) Count() int {
	return len(g.nodes)
}

// All returns a slice of all nodes in the graph.
func (g *Graph) All() []Node {
	result := make([]Node, 0, len(g.nodes))
	for _, node := range g.nodes {
		result = append(result, node)
	}
	return result
}

// Get retrieves a node by ID. Returns the node and true if found,
// or a zero-value node and false if not found.
func (g *Graph) Get(id string) (Node, bool) {
	node, exists := g.nodes[id]
	return node, exists
}

// Add adds a node to the graph. Returns an error if a node with
// the same ID already exists.
func (g *Graph) Add(node Node) error {
	if _, exists := g.nodes[node.ID]; exists {
		return fmt.Errorf("node with ID %q already exists", node.ID)
	}
	g.nodes[node.ID] = node
	return nil
}

// Remove removes a node from the graph by ID.
// Returns true if the node was removed, false if it didn't exist.
func (g *Graph) Remove(id string) bool {
	if _, exists := g.nodes[id]; !exists {
		return false
	}
	delete(g.nodes, id)
	return true
}

// Update updates an existing node in the graph.
// Returns an error if the node doesn't exist.
func (g *Graph) Update(node Node) error {
	if _, exists := g.nodes[node.ID]; !exists {
		return fmt.Errorf("node with ID %q does not exist", node.ID)
	}
	g.nodes[node.ID] = node
	return nil
}
