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

package refactor

import (
	"fmt"

	"github.com/Toernblom/deco/internal/domain"
)

// Renamer provides node renaming functionality with reference updates.
type Renamer struct{}

// NewRenamer creates a new Renamer.
func NewRenamer() *Renamer {
	return &Renamer{}
}

// Rename changes a node's ID and updates all references pointing to it.
// Returns a new slice with the renamed node and updated references.
// The original nodes slice is not modified.
//
// Parameters:
//   - nodes: slice of all nodes in the graph
//   - oldID: the current ID of the node to rename
//   - newID: the new ID for the node
//
// Returns:
//   - updated slice of nodes with the rename applied
//   - error if oldID doesn't exist, newID already exists, or IDs are invalid
func (r *Renamer) Rename(nodes []domain.Node, oldID, newID string) ([]domain.Node, error) {
	// Validate inputs
	if nodes == nil {
		return nil, fmt.Errorf("nodes slice cannot be nil")
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("nodes slice cannot be empty")
	}
	if oldID == "" {
		return nil, fmt.Errorf("oldID cannot be empty")
	}
	if newID == "" {
		return nil, fmt.Errorf("newID cannot be empty")
	}
	if oldID == newID {
		return nil, fmt.Errorf("newID must be different from oldID")
	}

	// Build index to check existence
	nodeIndex := make(map[string]int)
	for i, node := range nodes {
		nodeIndex[node.ID] = i
	}

	// Check oldID exists
	if _, exists := nodeIndex[oldID]; !exists {
		return nil, fmt.Errorf("node with ID %q does not exist", oldID)
	}

	// Check newID doesn't already exist
	if _, exists := nodeIndex[newID]; exists {
		return nil, fmt.Errorf("node with ID %q already exists", newID)
	}

	// Create deep copy of nodes
	result := make([]domain.Node, len(nodes))
	for i, node := range nodes {
		result[i] = copyNode(node)
	}

	// Update the target node's ID
	for i := range result {
		if result[i].ID == oldID {
			result[i].ID = newID
			break
		}
	}

	// Update all references pointing to oldID
	for i := range result {
		updated := false

		// Update Uses references
		for j := range result[i].Refs.Uses {
			if result[i].Refs.Uses[j].Target == oldID {
				result[i].Refs.Uses[j].Target = newID
				updated = true
			}
		}

		// Update Related references
		for j := range result[i].Refs.Related {
			if result[i].Refs.Related[j].Target == oldID {
				result[i].Refs.Related[j].Target = newID
				updated = true
			}
		}

		// Update EmitsEvents references
		for j := range result[i].Refs.EmitsEvents {
			if result[i].Refs.EmitsEvents[j] == oldID {
				result[i].Refs.EmitsEvents[j] = newID
				updated = true
			}
		}

		// Update Vocabulary references
		for j := range result[i].Refs.Vocabulary {
			if result[i].Refs.Vocabulary[j] == oldID {
				result[i].Refs.Vocabulary[j] = newID
				updated = true
			}
		}

		// Increment version if this node's references were updated
		if updated {
			result[i].Version++
		}
	}

	return result, nil
}

// copyNode creates a deep copy of a node.
func copyNode(n domain.Node) domain.Node {
	copy := domain.Node{
		ID:         n.ID,
		Kind:       n.Kind,
		Version:    n.Version,
		Status:     n.Status,
		Title:      n.Title,
		Summary:    n.Summary,
		LLMContext: n.LLMContext,
		Content:    n.Content, // Content is a pointer, shallow copy is OK for our use case
	}

	// Copy Tags
	if n.Tags != nil {
		copy.Tags = make([]string, len(n.Tags))
		for i, tag := range n.Tags {
			copy.Tags[i] = tag
		}
	}

	// Copy Refs
	copy.Refs = copyRef(n.Refs)

	// Copy Glossary
	if n.Glossary != nil {
		copy.Glossary = make(map[string]string)
		for k, v := range n.Glossary {
			copy.Glossary[k] = v
		}
	}

	// Copy Contracts
	if n.Contracts != nil {
		copy.Contracts = make([]domain.Contract, len(n.Contracts))
		for i, c := range n.Contracts {
			copy.Contracts[i] = c
		}
	}

	// Copy Constraints
	if n.Constraints != nil {
		copy.Constraints = make([]domain.Constraint, len(n.Constraints))
		for i, c := range n.Constraints {
			copy.Constraints[i] = c
		}
	}

	// Copy Issues
	if n.Issues != nil {
		copy.Issues = make([]domain.Issue, len(n.Issues))
		for i, issue := range n.Issues {
			copy.Issues[i] = issue
		}
	}

	// Copy Custom
	if n.Custom != nil {
		copy.Custom = make(map[string]interface{})
		for k, v := range n.Custom {
			copy.Custom[k] = v
		}
	}

	return copy
}

// UpdateReferences updates all references from oldID to newID across all nodes.
// Unlike Rename, this does not require oldID to exist as a node - useful for
// detecting manual renames where the old node file was already deleted/renamed.
// Returns a new slice with updated references. Nodes whose references were
// updated will have their version incremented.
func (r *Renamer) UpdateReferences(nodes []domain.Node, oldID, newID string) ([]domain.Node, error) {
	if nodes == nil {
		return nil, fmt.Errorf("nodes slice cannot be nil")
	}
	if oldID == "" {
		return nil, fmt.Errorf("oldID cannot be empty")
	}
	if newID == "" {
		return nil, fmt.Errorf("newID cannot be empty")
	}
	if oldID == newID {
		return nil, fmt.Errorf("newID must be different from oldID")
	}

	// Create deep copy of nodes
	result := make([]domain.Node, len(nodes))
	for i, node := range nodes {
		result[i] = copyNode(node)
	}

	// Update all references pointing to oldID
	for i := range result {
		updated := false

		// Update Uses references
		for j := range result[i].Refs.Uses {
			if result[i].Refs.Uses[j].Target == oldID {
				result[i].Refs.Uses[j].Target = newID
				updated = true
			}
		}

		// Update Related references
		for j := range result[i].Refs.Related {
			if result[i].Refs.Related[j].Target == oldID {
				result[i].Refs.Related[j].Target = newID
				updated = true
			}
		}

		// Update EmitsEvents references
		for j := range result[i].Refs.EmitsEvents {
			if result[i].Refs.EmitsEvents[j] == oldID {
				result[i].Refs.EmitsEvents[j] = newID
				updated = true
			}
		}

		// Update Vocabulary references
		for j := range result[i].Refs.Vocabulary {
			if result[i].Refs.Vocabulary[j] == oldID {
				result[i].Refs.Vocabulary[j] = newID
				updated = true
			}
		}

		// Increment version if this node's references were updated
		if updated {
			result[i].Version++
		}
	}

	return result, nil
}

// copyRef creates a deep copy of a Ref.
func copyRef(r domain.Ref) domain.Ref {
	copy := domain.Ref{}

	// Copy Uses
	if r.Uses != nil {
		copy.Uses = make([]domain.RefLink, len(r.Uses))
		for i, ref := range r.Uses {
			copy.Uses[i] = domain.RefLink{
				Target:   ref.Target,
				Context:  ref.Context,
				Resolved: ref.Resolved,
			}
		}
	}

	// Copy Related
	if r.Related != nil {
		copy.Related = make([]domain.RefLink, len(r.Related))
		for i, ref := range r.Related {
			copy.Related[i] = domain.RefLink{
				Target:   ref.Target,
				Context:  ref.Context,
				Resolved: ref.Resolved,
			}
		}
	}

	// Copy EmitsEvents
	if r.EmitsEvents != nil {
		copy.EmitsEvents = make([]string, len(r.EmitsEvents))
		for i, e := range r.EmitsEvents {
			copy.EmitsEvents[i] = e
		}
	}

	// Copy Vocabulary
	if r.Vocabulary != nil {
		copy.Vocabulary = make([]string, len(r.Vocabulary))
		for i, v := range r.Vocabulary {
			copy.Vocabulary[i] = v
		}
	}

	return copy
}
