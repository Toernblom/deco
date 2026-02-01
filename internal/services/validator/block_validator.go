package validator

import (
	"fmt"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
)

// knownBlockTypes defines the valid block types that can appear in content sections.
var knownBlockTypes = map[string]bool{
	"rule":     true,
	"table":    true,
	"param":    true,
	"mechanic": true,
	"list":     true,
}

// BlockValidator validates that blocks within content sections have the required
// fields for their type.
type BlockValidator struct {
	suggester *errors.Suggester
}

// NewBlockValidator creates a new block validator.
func NewBlockValidator() *BlockValidator {
	return &BlockValidator{
		suggester: errors.NewSuggester(),
	}
}

// Validate checks all blocks in a node's content sections.
func (bv *BlockValidator) Validate(node *domain.Node, collector *errors.Collector) {
	if node == nil || node.Content == nil {
		return
	}

	for _, section := range node.Content.Sections {
		for blockIdx, block := range section.Blocks {
			bv.validateBlock(&block, node.ID, section.Name, blockIdx, collector)
		}
	}
}

// validateBlock dispatches to type-specific validation.
func (bv *BlockValidator) validateBlock(block *domain.Block, nodeID, sectionName string, blockIdx int, collector *errors.Collector) {
	// Check for empty or unknown block type
	if block.Type == "" {
		collector.Add(domain.DecoError{
			Code:    "E048",
			Summary: "Block has no type",
			Detail:  bv.formatLocation(nodeID, sectionName, blockIdx),
		})
		return
	}

	if !knownBlockTypes[block.Type] {
		// Generate suggestion for similar type names
		var knownTypes []string
		for t := range knownBlockTypes {
			knownTypes = append(knownTypes, t)
		}

		err := domain.DecoError{
			Code:    "E048",
			Summary: fmt.Sprintf("Unknown block type: %s", block.Type),
			Detail:  bv.formatLocation(nodeID, sectionName, blockIdx),
		}

		suggs := bv.suggester.Suggest(block.Type, knownTypes)
		if len(suggs) > 0 {
			err.Suggestion = fmt.Sprintf("Did you mean %q?", suggs[0])
		}

		collector.Add(err)
		return
	}

	// Dispatch to type-specific validators
	switch block.Type {
	case "rule":
		bv.validateRule(block, nodeID, sectionName, blockIdx, collector)
	case "table":
		bv.validateTable(block, nodeID, sectionName, blockIdx, collector)
	case "param":
		bv.validateParam(block, nodeID, sectionName, blockIdx, collector)
	case "mechanic":
		bv.validateMechanic(block, nodeID, sectionName, blockIdx, collector)
	case "list":
		bv.validateList(block, nodeID, sectionName, blockIdx, collector)
	}
}

// validateRule checks that rule blocks have required fields.
// Required: text
func (bv *BlockValidator) validateRule(block *domain.Block, nodeID, sectionName string, blockIdx int, collector *errors.Collector) {
	bv.requireField(block, "text", nodeID, sectionName, blockIdx, collector)
}

// validateTable checks that table blocks have required fields.
// Required: columns, rows
func (bv *BlockValidator) validateTable(block *domain.Block, nodeID, sectionName string, blockIdx int, collector *errors.Collector) {
	bv.requireField(block, "columns", nodeID, sectionName, blockIdx, collector)
	bv.requireField(block, "rows", nodeID, sectionName, blockIdx, collector)

	// Validate column structure if columns exist
	if columns, ok := block.Data["columns"]; ok {
		bv.validateTableColumns(columns, nodeID, sectionName, blockIdx, collector)
	}
}

// validateTableColumns checks that each column has required fields.
// Required for each column: key
func (bv *BlockValidator) validateTableColumns(columns interface{}, nodeID, sectionName string, blockIdx int, collector *errors.Collector) {
	columnList, ok := columns.([]interface{})
	if !ok {
		return
	}

	for colIdx, col := range columnList {
		colMap, ok := col.(map[string]interface{})
		if !ok {
			continue
		}

		if _, hasKey := colMap["key"]; !hasKey {
			collector.Add(domain.DecoError{
				Code:    "E050",
				Summary: fmt.Sprintf("Table column %d missing required field: key", colIdx),
				Detail:  bv.formatLocation(nodeID, sectionName, blockIdx),
			})
		}
	}
}

// validateParam checks that param blocks have required fields.
// Required: name, datatype
func (bv *BlockValidator) validateParam(block *domain.Block, nodeID, sectionName string, blockIdx int, collector *errors.Collector) {
	bv.requireField(block, "name", nodeID, sectionName, blockIdx, collector)
	bv.requireField(block, "datatype", nodeID, sectionName, blockIdx, collector)
}

// validateMechanic checks that mechanic blocks have required fields.
// Required: name, description
func (bv *BlockValidator) validateMechanic(block *domain.Block, nodeID, sectionName string, blockIdx int, collector *errors.Collector) {
	bv.requireField(block, "name", nodeID, sectionName, blockIdx, collector)
	bv.requireField(block, "description", nodeID, sectionName, blockIdx, collector)
}

// validateList checks that list blocks have required fields.
// Required: items
func (bv *BlockValidator) validateList(block *domain.Block, nodeID, sectionName string, blockIdx int, collector *errors.Collector) {
	bv.requireField(block, "items", nodeID, sectionName, blockIdx, collector)
}

// requireField checks that a field exists in block.Data and adds an error if missing.
func (bv *BlockValidator) requireField(block *domain.Block, field, nodeID, sectionName string, blockIdx int, collector *errors.Collector) {
	if _, ok := block.Data[field]; !ok {
		collector.Add(domain.DecoError{
			Code:    "E047",
			Summary: fmt.Sprintf("Block type %q missing required field: %s", block.Type, field),
			Detail:  bv.formatLocation(nodeID, sectionName, blockIdx),
		})
	}
}

// formatLocation creates a consistent location string for error details.
func (bv *BlockValidator) formatLocation(nodeID, sectionName string, blockIdx int) string {
	return fmt.Sprintf("in node %q, section %q, block %d", nodeID, sectionName, blockIdx)
}
