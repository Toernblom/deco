package query

import (
	"fmt"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
)

// FilterCriteria defines the criteria for filtering nodes and blocks.
type FilterCriteria struct {
	Kind         *string           // Filter by node kind (exact match)
	Status       *string           // Filter by node status (exact match)
	Tags         []string          // Filter by tags (node must have all specified tags)
	BlockType    *string           // Filter by block type within content sections
	FieldFilters map[string]string // Filter by block field key=value (AND logic)
}

// BlockMatch represents a matched block with its parent context.
type BlockMatch struct {
	NodeID      string
	NodeTitle   string
	SectionName string
	BlockIndex  int
	Block       domain.Block
}

// QueryEngine provides query and search capabilities for nodes.
type QueryEngine struct {
	// No state needed for now
}

// New creates a new QueryEngine instance.
func New() *QueryEngine {
	return &QueryEngine{}
}

// Filter returns nodes that match the given criteria.
// Multiple criteria are combined with AND logic.
// Empty criteria returns all nodes.
func (qe *QueryEngine) Filter(nodes []domain.Node, criteria FilterCriteria) []domain.Node {
	var results []domain.Node

	for _, node := range nodes {
		if qe.matchesCriteria(node, criteria) {
			results = append(results, node)
		}
	}

	return results
}

// FilterBlocks returns blocks that match the given criteria, including block-level filters.
// Node-level filters (Kind, Status, Tags) narrow which nodes are searched.
// BlockType and FieldFilters filter blocks within matching nodes.
func (qe *QueryEngine) FilterBlocks(nodes []domain.Node, criteria FilterCriteria) []BlockMatch {
	var results []BlockMatch

	for _, node := range nodes {
		// Apply node-level filters first
		if !qe.matchesNodeCriteria(node, criteria) {
			continue
		}

		if node.Content == nil {
			continue
		}

		for _, section := range node.Content.Sections {
			for blockIdx, block := range section.Blocks {
				if qe.matchesBlockCriteria(block, criteria) {
					results = append(results, BlockMatch{
						NodeID:      node.ID,
						NodeTitle:   node.Title,
						SectionName: section.Name,
						BlockIndex:  blockIdx,
						Block:       block,
					})
				}
			}
		}
	}

	return results
}

// matchesBlockCriteria checks if a block matches block-level filter criteria.
func (qe *QueryEngine) matchesBlockCriteria(block domain.Block, criteria FilterCriteria) bool {
	if criteria.BlockType != nil && block.Type != *criteria.BlockType {
		return false
	}

	for key, value := range criteria.FieldFilters {
		fieldVal, ok := block.Data[key]
		if !ok {
			return false
		}
		if fmt.Sprintf("%v", fieldVal) != value {
			return false
		}
	}

	return true
}

// matchesNodeCriteria checks if a node matches node-level filter criteria (Kind, Status, Tags).
func (qe *QueryEngine) matchesNodeCriteria(node domain.Node, criteria FilterCriteria) bool {
	if criteria.Kind != nil && node.Kind != *criteria.Kind {
		return false
	}
	if criteria.Status != nil && node.Status != *criteria.Status {
		return false
	}
	if len(criteria.Tags) > 0 {
		if !qe.hasAllTags(node.Tags, criteria.Tags) {
			return false
		}
	}
	return true
}

// matchesCriteria checks if a node matches all filter criteria
func (qe *QueryEngine) matchesCriteria(node domain.Node, criteria FilterCriteria) bool {
	// Check kind filter
	if criteria.Kind != nil && node.Kind != *criteria.Kind {
		return false
	}

	// Check status filter
	if criteria.Status != nil && node.Status != *criteria.Status {
		return false
	}

	// Check tags filter (node must have ALL specified tags)
	if len(criteria.Tags) > 0 {
		if !qe.hasAllTags(node.Tags, criteria.Tags) {
			return false
		}
	}

	return true
}

// hasAllTags checks if nodeTags contains all required tags
func (qe *QueryEngine) hasAllTags(nodeTags []string, requiredTags []string) bool {
	for _, required := range requiredTags {
		found := false
		for _, nodeTag := range nodeTags {
			if nodeTag == required {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Search returns nodes that contain the search term in their title or summary.
// Search is case-insensitive and supports partial matches.
// Empty search term returns all nodes.
func (qe *QueryEngine) Search(nodes []domain.Node, term string) []domain.Node {
	var results []domain.Node

	// Empty term matches all nodes
	if term == "" {
		return nodes
	}

	// Convert term to lowercase for case-insensitive search
	lowerTerm := strings.ToLower(term)

	for _, node := range nodes {
		if qe.matchesSearchTerm(node, lowerTerm) {
			results = append(results, node)
		}
	}

	return results
}

// matchesSearchTerm checks if a node contains the search term in title or summary
func (qe *QueryEngine) matchesSearchTerm(node domain.Node, lowerTerm string) bool {
	// Check title
	if strings.Contains(strings.ToLower(node.Title), lowerTerm) {
		return true
	}

	// Check summary
	if strings.Contains(strings.ToLower(node.Summary), lowerTerm) {
		return true
	}

	return false
}
