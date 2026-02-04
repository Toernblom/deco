package validator

import (
	"fmt"
	"sort"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
	"github.com/Toernblom/deco/internal/storage/config"
)

// knownBlockTypes defines the valid block types that can appear in content sections.
var knownBlockTypes = map[string]bool{
	"rule":     true,
	"table":    true,
	"param":    true,
	"mechanic": true,
	"list":     true,
}

var builtInBlockFields = map[string][]string{
	"rule":     {"id", "text"},
	"table":    {"id", "columns", "rows"},
	"param":    {"id", "name", "datatype", "min", "max", "default", "unit", "description"},
	"mechanic": {"id", "name", "description", "conditions", "outputs", "inputs"},
	"list":     {"id", "items"},
}

var allowedTableColumnFields = map[string]bool{
	"key":     true,
	"type":    true,
	"enum":    true,
	"display": true,
}

// BlockValidator validates that blocks within content sections have the required
// fields for their type.
type BlockValidator struct {
	suggester        *errors.Suggester
	customBlockTypes map[string]config.BlockTypeConfig
}

// NewBlockValidator creates a new block validator.
func NewBlockValidator() *BlockValidator {
	return &BlockValidator{
		suggester:        errors.NewSuggester(),
		customBlockTypes: nil,
	}
}

// NewBlockValidatorWithConfig creates a block validator with custom type support.
func NewBlockValidatorWithConfig(customBlockTypes map[string]config.BlockTypeConfig) *BlockValidator {
	return &BlockValidator{
		suggester:        errors.NewSuggester(),
		customBlockTypes: customBlockTypes,
	}
}

// Validate checks all blocks in a node's content sections.
func (bv *BlockValidator) Validate(node *domain.Node, collector *errors.Collector) {
	if node == nil || node.Content == nil {
		return
	}

	// Create location from node source file
	var location *domain.Location
	if node.SourceFile != "" {
		location = &domain.Location{File: node.SourceFile}
	}

	for _, section := range node.Content.Sections {
		for blockIdx, block := range section.Blocks {
			bv.validateBlock(&block, node.ID, section.Name, blockIdx, location, collector)
		}
	}
}

// validateBlock dispatches to type-specific validation.
func (bv *BlockValidator) validateBlock(block *domain.Block, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	// Check for empty or unknown block type
	if block.Type == "" {
		collector.Add(domain.DecoError{
			Code:     "E048",
			Summary:  "Block has no type",
			Detail:   bv.formatLocation(nodeID, sectionName, blockIdx),
			Location: location,
		})
		return
	}

	// Check if it's a known built-in type
	isBuiltIn := knownBlockTypes[block.Type]

	// Check if it's a custom type
	var customTypeConfig *config.BlockTypeConfig
	if bv.customBlockTypes != nil {
		if cfg, ok := bv.customBlockTypes[block.Type]; ok {
			customTypeConfig = &cfg
		}
	}

	if !isBuiltIn && customTypeConfig == nil {
		// Unknown type - generate helpful error with suggestions
		var allTypes []string
		for t := range knownBlockTypes {
			allTypes = append(allTypes, t)
		}
		for t := range bv.customBlockTypes {
			allTypes = append(allTypes, t)
		}

		err := domain.DecoError{
			Code:     "E048",
			Summary:  fmt.Sprintf("Unknown block type: %s", block.Type),
			Detail:   bv.formatLocation(nodeID, sectionName, blockIdx),
			Location: location,
		}

		suggs := bv.suggester.Suggest(block.Type, allTypes)
		if len(suggs) > 0 {
			err.Suggestion = fmt.Sprintf("Did you mean %q?", suggs[0])
		}

		collector.Add(err)
		return
	}

	// Run built-in validation if applicable
	if isBuiltIn {
		switch block.Type {
		case "rule":
			bv.validateRule(block, nodeID, sectionName, blockIdx, location, collector)
		case "table":
			bv.validateTable(block, nodeID, sectionName, blockIdx, location, collector)
		case "param":
			bv.validateParam(block, nodeID, sectionName, blockIdx, location, collector)
		case "mechanic":
			bv.validateMechanic(block, nodeID, sectionName, blockIdx, location, collector)
		case "list":
			bv.validateList(block, nodeID, sectionName, blockIdx, location, collector)
		}
	}

	// Run custom type validation if configured (extends built-in validation)
	if customTypeConfig != nil {
		bv.validateCustomType(block, customTypeConfig, nodeID, sectionName, blockIdx, location, collector)
	}

	allowedFields := bv.allowedFieldsForBlock(block.Type, isBuiltIn, customTypeConfig)
	bv.validateUnknownBlockFields(block, allowedFields, nodeID, sectionName, blockIdx, location, collector)
}

// validateRule checks that rule blocks have required fields.
// Required: text
func (bv *BlockValidator) validateRule(block *domain.Block, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	bv.requireField(block, "text", nodeID, sectionName, blockIdx, location, collector)
}

// validateTable checks that table blocks have required fields.
// Required: columns, rows
func (bv *BlockValidator) validateTable(block *domain.Block, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	bv.requireField(block, "columns", nodeID, sectionName, blockIdx, location, collector)
	bv.requireField(block, "rows", nodeID, sectionName, blockIdx, location, collector)

	// Validate column structure if columns exist
	if columns, ok := block.Data["columns"]; ok {
		bv.validateTableColumns(columns, nodeID, sectionName, blockIdx, location, collector)
	}
}

// validateTableColumns checks that each column has required fields.
// Required for each column: key
// This function consolidates related errors: if an unknown field's suggestion
// matches a missing required field (typo scenario), only one error is reported.
func (bv *BlockValidator) validateTableColumns(columns interface{}, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	columnList, ok := columns.([]interface{})
	if !ok {
		return
	}

	allowed := allowedFieldList(allowedTableColumnFields)

	for colIdx, col := range columnList {
		colMap, ok := col.(map[string]interface{})
		if !ok {
			continue
		}

		// Collect unknown fields and their suggestions
		type unknownField struct {
			name       string
			suggestion string
		}
		var unknownFields []unknownField

		columnKeys := make([]string, 0, len(colMap))
		for key := range colMap {
			columnKeys = append(columnKeys, key)
		}
		sort.Strings(columnKeys)

		for _, key := range columnKeys {
			if !allowedTableColumnFields[key] {
				uf := unknownField{name: key}
				suggs := bv.suggester.Suggest(key, allowed)
				if len(suggs) > 0 {
					uf.suggestion = suggs[0]
				}
				unknownFields = append(unknownFields, uf)
			}
		}

		// Check for missing "key" field
		_, hasKey := colMap["key"]

		// Check if any unknown field is a typo for "key"
		typoForKey := ""
		for _, uf := range unknownFields {
			if uf.suggestion == "key" {
				typoForKey = uf.name
				break
			}
		}

		// Build column contents string for context
		colContents := bv.formatColumnContents(colMap)

		// If missing "key" AND there's a typo for "key", consolidate into single error
		if !hasKey && typoForKey != "" {
			err := domain.DecoError{
				Code:     "E049",
				Summary:  fmt.Sprintf("Unknown field %q in table column %d (did you mean \"key\"?)", typoForKey, colIdx),
				Detail:   bv.formatLocation(nodeID, sectionName, blockIdx),
				Location: location,
				Context: []string{
					fmt.Sprintf("Column %d contains: %s", colIdx, colContents),
					"This also causes: missing required field \"key\"",
				},
			}
			collector.Add(err)

			// Report other unknown fields (not the typo)
			for _, uf := range unknownFields {
				if uf.name == typoForKey {
					continue
				}
				err := domain.DecoError{
					Code:     "E049",
					Summary:  fmt.Sprintf("Unknown table column field %q in column %d", uf.name, colIdx),
					Detail:   bv.formatLocation(nodeID, sectionName, blockIdx),
					Location: location,
					Context:  []string{fmt.Sprintf("Column %d contains: %s", colIdx, colContents)},
				}
				if uf.suggestion != "" {
					err.Suggestion = fmt.Sprintf("Did you mean %q?", uf.suggestion)
				}
				collector.Add(err)
			}
		} else {
			// No consolidation - report errors separately
			if !hasKey {
				collector.Add(domain.DecoError{
					Code:     "E050",
					Summary:  fmt.Sprintf("Table column %d missing required field: key", colIdx),
					Detail:   bv.formatLocation(nodeID, sectionName, blockIdx),
					Location: location,
					Context:  []string{fmt.Sprintf("Column %d contains: %s", colIdx, colContents)},
				})
			}

			for _, uf := range unknownFields {
				err := domain.DecoError{
					Code:     "E049",
					Summary:  fmt.Sprintf("Unknown table column field %q in column %d", uf.name, colIdx),
					Detail:   bv.formatLocation(nodeID, sectionName, blockIdx),
					Location: location,
					Context:  []string{fmt.Sprintf("Column %d contains: %s", colIdx, colContents)},
				}
				if uf.suggestion != "" {
					err.Suggestion = fmt.Sprintf("Did you mean %q?", uf.suggestion)
				}
				collector.Add(err)
			}
		}
	}
}

// formatColumnContents creates a brief representation of a column's contents.
func (bv *BlockValidator) formatColumnContents(colMap map[string]interface{}) string {
	if len(colMap) == 0 {
		return "{}"
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(colMap))
	for k := range colMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		v := colMap[k]
		// Format value - truncate long strings
		var vStr string
		switch val := v.(type) {
		case string:
			if len(val) > 20 {
				vStr = fmt.Sprintf("%q...", val[:17])
			} else {
				vStr = fmt.Sprintf("%q", val)
			}
		default:
			vStr = fmt.Sprintf("%v", val)
		}
		parts = append(parts, fmt.Sprintf("%s: %s", k, vStr))
	}

	return "{" + fmt.Sprintf("%s", joinStrings(parts, ", ")) + "}"
}

func joinStrings(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += sep + parts[i]
	}
	return result
}

// validateParam checks that param blocks have required fields.
// Required: name, datatype
func (bv *BlockValidator) validateParam(block *domain.Block, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	bv.requireField(block, "name", nodeID, sectionName, blockIdx, location, collector)
	bv.requireField(block, "datatype", nodeID, sectionName, blockIdx, location, collector)
}

// validateMechanic checks that mechanic blocks have required fields.
// Required: name, description
func (bv *BlockValidator) validateMechanic(block *domain.Block, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	bv.requireField(block, "name", nodeID, sectionName, blockIdx, location, collector)
	bv.requireField(block, "description", nodeID, sectionName, blockIdx, location, collector)
}

// validateList checks that list blocks have required fields.
// Required: items
func (bv *BlockValidator) validateList(block *domain.Block, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	bv.requireField(block, "items", nodeID, sectionName, blockIdx, location, collector)
}

// requireField checks that a field exists in block.Data and adds an error if missing.
func (bv *BlockValidator) requireField(block *domain.Block, field, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	if _, ok := block.Data[field]; !ok {
		collector.Add(domain.DecoError{
			Code:     "E047",
			Summary:  fmt.Sprintf("Block type %q missing required field: %s", block.Type, field),
			Detail:   bv.formatLocation(nodeID, sectionName, blockIdx),
			Location: location,
		})
	}
}

// formatLocation creates a consistent location string for error details.
func (bv *BlockValidator) formatLocation(nodeID, sectionName string, blockIdx int) string {
	return fmt.Sprintf("in node %q, section %q, block %d", nodeID, sectionName, blockIdx)
}

// validateCustomType validates a block against a custom type's required fields.
func (bv *BlockValidator) validateCustomType(block *domain.Block, cfg *config.BlockTypeConfig, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	for _, field := range cfg.RequiredFields {
		bv.requireField(block, field, nodeID, sectionName, blockIdx, location, collector)
	}
}

func (bv *BlockValidator) allowedFieldsForBlock(blockType string, isBuiltIn bool, cfg *config.BlockTypeConfig) map[string]bool {
	allowed := map[string]bool{}
	if isBuiltIn {
		for _, field := range builtInBlockFields[blockType] {
			allowed[field] = true
		}
	}
	if cfg != nil {
		for _, field := range cfg.RequiredFields {
			allowed[field] = true
		}
		for _, field := range cfg.OptionalFields {
			allowed[field] = true
		}
	}
	if isBuiltIn || cfg != nil {
		allowed["id"] = true
	}
	return allowed
}

func (bv *BlockValidator) validateUnknownBlockFields(block *domain.Block, allowed map[string]bool, nodeID, sectionName string, blockIdx int, location *domain.Location, collector *errors.Collector) {
	if block == nil || len(block.Data) == 0 || len(allowed) == 0 {
		return
	}

	blockKeys := make([]string, 0, len(block.Data))
	for key := range block.Data {
		blockKeys = append(blockKeys, key)
	}
	sort.Strings(blockKeys)

	allowedList := allowedFieldList(allowed)
	for _, key := range blockKeys {
		if !allowed[key] {
			err := domain.DecoError{
				Code:     "E049",
				Summary:  fmt.Sprintf("Unknown field %q in %s block", key, block.Type),
				Detail:   bv.formatLocation(nodeID, sectionName, blockIdx),
				Location: location,
			}

			suggs := bv.suggester.Suggest(key, allowedList)
			if len(suggs) > 0 {
				err.Suggestion = fmt.Sprintf("Did you mean %q?", suggs[0])
			}

			collector.Add(err)
		}
	}
}

func allowedFieldList(allowed map[string]bool) []string {
	allowedList := make([]string, 0, len(allowed))
	for key := range allowed {
		allowedList = append(allowedList, key)
	}
	sort.Strings(allowedList)
	return allowedList
}
