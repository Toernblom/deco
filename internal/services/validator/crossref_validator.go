package validator

import (
	"fmt"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
	"github.com/Toernblom/deco/internal/storage/config"
)

// CrossRefValidator validates that block fields with ref constraints reference
// valid values from other block types across all nodes.
type CrossRefValidator struct {
	customBlockTypes map[string]config.BlockTypeConfig
	suggester        *errors.Suggester
}

// NewCrossRefValidator creates a new cross-reference validator.
func NewCrossRefValidator(customBlockTypes map[string]config.BlockTypeConfig) *CrossRefValidator {
	return &CrossRefValidator{
		customBlockTypes: customBlockTypes,
		suggester:        errors.NewSuggester(),
	}
}

// Validate runs cross-reference validation across all nodes.
// Two passes: first collects all reference sets, then validates against them.
func (cv *CrossRefValidator) Validate(nodes []domain.Node, collector *errors.Collector) {
	if cv.customBlockTypes == nil {
		return
	}

	// Pass 1: Collect all reference sets.
	// Key: "blockType.fieldName" -> set of values
	refSets := cv.buildRefSets(nodes)

	// Pass 2: Validate all fields with ref constraints.
	for _, node := range nodes {
		if node.Content == nil {
			continue
		}

		var location *domain.Location
		if node.SourceFile != "" {
			location = &domain.Location{File: node.SourceFile}
		}

		for _, section := range node.Content.Sections {
			for blockIdx, block := range section.Blocks {
				cv.validateBlockRefs(block, node.ID, section.Name, blockIdx, location, refSets, collector)
			}
		}
	}
}

// buildRefSets collects all values for each block type + field combination.
func (cv *CrossRefValidator) buildRefSets(nodes []domain.Node) map[string]map[string]bool {
	sets := make(map[string]map[string]bool)

	for _, node := range nodes {
		if node.Content == nil {
			continue
		}
		for _, section := range node.Content.Sections {
			for _, block := range section.Blocks {
				cv.collectBlockValues(block, sets)
			}
		}
	}

	return sets
}

// collectBlockValues adds field values from a block to the reference sets.
func (cv *CrossRefValidator) collectBlockValues(block domain.Block, sets map[string]map[string]bool) {
	for fieldName, val := range block.Data {
		key := block.Type + "." + fieldName
		if sets[key] == nil {
			sets[key] = make(map[string]bool)
		}

		switch v := val.(type) {
		case string:
			sets[key][v] = true
		}
	}
}

// validateBlockRefs checks all ref-constrained fields in a block.
func (cv *CrossRefValidator) validateBlockRefs(block domain.Block, nodeID, sectionName string, blockIdx int, location *domain.Location, refSets map[string]map[string]bool, collector *errors.Collector) {
	blockCfg, ok := cv.customBlockTypes[block.Type]
	if !ok || blockCfg.Fields == nil {
		return
	}

	for fieldName, fieldDef := range blockCfg.Fields {
		if fieldDef.Ref == nil {
			continue
		}

		val, ok := block.Data[fieldName]
		if !ok {
			continue // missing field is handled by required field validation
		}

		refKey := fieldDef.Ref.BlockType + "." + fieldDef.Ref.Field
		validValues := refSets[refKey]

		// Collect valid values as a list for suggestions
		var validList []string
		for v := range validValues {
			validList = append(validList, v)
		}

		switch v := val.(type) {
		case string:
			cv.validateSingleRef(v, block.Type, fieldName, nodeID, sectionName, blockIdx, location, validValues, validList, collector)
		case []interface{}:
			for _, item := range v {
				if strItem, ok := item.(string); ok {
					cv.validateSingleRef(strItem, block.Type, fieldName, nodeID, sectionName, blockIdx, location, validValues, validList, collector)
				}
			}
		}
	}
}

// validateSingleRef checks a single value against the reference set.
func (cv *CrossRefValidator) validateSingleRef(value, blockType, fieldName, nodeID, sectionName string, blockIdx int, location *domain.Location, validValues map[string]bool, validList []string, collector *errors.Collector) {
	if validValues != nil && validValues[value] {
		return
	}

	err := domain.DecoError{
		Code:     "E054",
		Summary:  fmt.Sprintf("Cross-reference not found: %s block field %q contains %q which is not a known value", blockType, fieldName, value),
		Detail:   fmt.Sprintf("in node %q, section %q, block %d", nodeID, sectionName, blockIdx),
		Location: location,
	}

	suggs := cv.suggester.Suggest(value, validList)
	if len(suggs) > 0 {
		err.Suggestion = fmt.Sprintf("Did you mean %q?", suggs[0])
	}

	collector.Add(err)
}
