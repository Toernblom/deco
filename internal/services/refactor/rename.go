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
