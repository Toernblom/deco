package query

import (
	"strings"

	"github.com/Toernblom/deco/internal/domain"
)

// FilterCriteria defines the criteria for filtering nodes.
type FilterCriteria struct {
	Kind   *string  // Filter by node kind (exact match)
	Status *string  // Filter by node status (exact match)
	Tags   []string // Filter by tags (node must have all specified tags)
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
