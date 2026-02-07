package query

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
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
		if !matchesFieldValue(fieldVal, value) {
			return false
		}
	}

	return true
}

// matchesFieldValue checks if a field value matches the filter value.
// For lists, it checks if any element matches (contains semantics).
// For scalars, it does exact string comparison.
func matchesFieldValue(fieldVal interface{}, filterVal string) bool {
	switch v := fieldVal.(type) {
	case string:
		return v == filterVal
	case []interface{}:
		for _, item := range v {
			if fmt.Sprintf("%v", item) == filterVal {
				return true
			}
		}
		return false
	default:
		return fmt.Sprintf("%v", fieldVal) == filterVal
	}
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

// FollowTarget specifies a block type and field to follow into.
type FollowTarget struct {
	BlockType string
	Field     string
}

// FollowResult groups followed blocks by their source value.
type FollowResult struct {
	Value    string       // The followed value (e.g., "Planks")
	RefCount int          // How many source blocks reference this value
	Matches  []BlockMatch // Blocks in target types that provide this value
}

// FollowBlocks takes matched source blocks, extracts values from the follow field,
// and finds blocks in target types that provide those values.
// If targets is nil, it looks up the ref config for the source block type.
// Returns results grouped by followed value, sorted alphabetically.
func (qe *QueryEngine) FollowBlocks(sourceMatches []BlockMatch, followField string, targets []FollowTarget, allNodes []domain.Node, blockTypes map[string]config.BlockTypeConfig) ([]FollowResult, error) {
	if len(sourceMatches) == 0 {
		return nil, nil
	}

	// Resolve targets from ref config if not explicitly provided
	if len(targets) == 0 {
		sourceBlockType := sourceMatches[0].Block.Type
		resolved, err := qe.resolveFollowTargets(sourceBlockType, followField, blockTypes)
		if err != nil {
			return nil, err
		}
		targets = resolved
	}

	// Extract all values from the follow field across source blocks
	// Track which values come from which source blocks for ref counting
	valueRefCount := make(map[string]int)
	fieldFound := false
	for _, match := range sourceMatches {
		val, ok := match.Block.Data[followField]
		if !ok {
			continue
		}
		fieldFound = true
		for _, v := range extractStringValues(val) {
			valueRefCount[v]++
		}
	}

	if !fieldFound {
		return nil, fmt.Errorf("field %q not found in matched blocks", followField)
	}

	// For each target, find blocks that match any of the extracted values
	// Group results by value
	resultMap := make(map[string]*FollowResult)
	for value, count := range valueRefCount {
		resultMap[value] = &FollowResult{
			Value:    value,
			RefCount: count,
		}
	}

	for _, target := range targets {
		targetBlockType := target.BlockType
		targetField := target.Field

		for _, node := range allNodes {
			if node.Content == nil {
				continue
			}
			for _, section := range node.Content.Sections {
				for blockIdx, block := range section.Blocks {
					if block.Type != targetBlockType {
						continue
					}
					blockValues := extractStringValues(block.Data[targetField])
					for _, bv := range blockValues {
						if result, ok := resultMap[bv]; ok {
							result.Matches = append(result.Matches, BlockMatch{
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
		}
	}

	// Convert map to sorted slice
	var results []FollowResult
	for _, r := range resultMap {
		results = append(results, *r)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Value < results[j].Value
	})

	return results, nil
}

// resolveFollowTargets looks up the ref config for a block type's field
// and returns the follow targets.
func (qe *QueryEngine) resolveFollowTargets(blockType, fieldName string, blockTypes map[string]config.BlockTypeConfig) ([]FollowTarget, error) {
	if blockTypes == nil {
		return nil, fmt.Errorf("field %q has no ref constraint; use --follow %s:<block_type>.<field>", fieldName, fieldName)
	}

	btConfig, ok := blockTypes[blockType]
	if !ok || btConfig.Fields == nil {
		return nil, fmt.Errorf("field %q has no ref constraint; use --follow %s:<block_type>.<field>", fieldName, fieldName)
	}

	fieldDef, ok := btConfig.Fields[fieldName]
	if !ok || len(fieldDef.Refs) == 0 {
		return nil, fmt.Errorf("field %q has no ref constraint; use --follow %s:<block_type>.<field>", fieldName, fieldName)
	}

	var targets []FollowTarget
	for _, ref := range fieldDef.Refs {
		targets = append(targets, FollowTarget{
			BlockType: ref.BlockType,
			Field:     ref.Field,
		})
	}
	return targets, nil
}

// extractStringValues extracts string values from a field value.
// For strings, returns a single-element slice. For lists, returns all string elements.
func extractStringValues(val interface{}) []string {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case string:
		return []string{v}
	case []interface{}:
		var result []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	default:
		return []string{fmt.Sprintf("%v", v)}
	}
}
