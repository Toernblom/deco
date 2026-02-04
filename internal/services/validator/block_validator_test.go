package validator

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
	"github.com/Toernblom/deco/internal/storage/config"
)

func TestBlockValidator_ValidRuleBlock(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Rules",
					Blocks: []domain.Block{
						{
							Type: "rule",
							Data: map[string]interface{}{
								"id":   "test_rule",
								"text": "This is a valid rule",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid rule block, got: %v", collector.Errors())
	}
}

func TestBlockValidator_RuleMissingText(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Rules",
					Blocks: []domain.Block{
						{
							Type: "rule",
							Data: map[string]interface{}{
								"id": "test_rule",
								// missing "text" field
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for rule block missing text")
	}

	errs := collector.Errors()
	if errs[0].Code != "E047" {
		t.Errorf("expected E047, got %s", errs[0].Code)
	}
}

func TestBlockValidator_ValidTableBlock(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Data",
					Blocks: []domain.Block{
						{
							Type: "table",
							Data: map[string]interface{}{
								"id": "test_table",
								"columns": []interface{}{
									map[string]interface{}{"key": "name", "type": "string", "display": "Name"},
								},
								"rows": []interface{}{
									map[string]interface{}{"name": "Test"},
								},
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid table block, got: %v", collector.Errors())
	}
}

func TestBlockValidator_TableMissingColumns(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Data",
					Blocks: []domain.Block{
						{
							Type: "table",
							Data: map[string]interface{}{
								"id": "test_table",
								// missing "columns"
								"rows": []interface{}{
									map[string]interface{}{"name": "Test"},
								},
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for table block missing columns")
	}

	errs := collector.Errors()
	if errs[0].Code != "E047" {
		t.Errorf("expected E047, got %s", errs[0].Code)
	}
}

func TestBlockValidator_TableMissingRows(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Data",
					Blocks: []domain.Block{
						{
							Type: "table",
							Data: map[string]interface{}{
								"id": "test_table",
								"columns": []interface{}{
									map[string]interface{}{"key": "name", "type": "string", "display": "Name"},
								},
								// missing "rows"
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for table block missing rows")
	}

	errs := collector.Errors()
	if errs[0].Code != "E047" {
		t.Errorf("expected E047, got %s", errs[0].Code)
	}
}

func TestBlockValidator_TableColumnMissingKey(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Data",
					Blocks: []domain.Block{
						{
							Type: "table",
							Data: map[string]interface{}{
								"id": "test_table",
								"columns": []interface{}{
									map[string]interface{}{"type": "string", "display": "Name"}, // missing key
								},
								"rows": []interface{}{
									map[string]interface{}{"name": "Test"},
								},
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for table column missing key")
	}

	errs := collector.Errors()
	if errs[0].Code != "E050" {
		t.Errorf("expected E050, got %s", errs[0].Code)
	}
}

func TestBlockValidator_TableColumnUnknownField(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Data",
					Blocks: []domain.Block{
						{
							Type: "table",
							Data: map[string]interface{}{
								"id": "test_table",
								"columns": []interface{}{
									map[string]interface{}{
										"key":     "name",
										"type":    "string",
										"display": "Name",
										"dispay":  "typo",
									},
								},
								"rows": []interface{}{
									map[string]interface{}{"name": "Test"},
								},
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for unknown table column field")
	}

	found := false
	for _, err := range collector.Errors() {
		if err.Code == "E049" {
			found = true
		}
	}
	if !found {
		t.Error("expected E049 for unknown table column field")
	}
}

func TestBlockValidator_ValidParamBlock(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Parameters",
					Blocks: []domain.Block{
						{
							Type: "param",
							Data: map[string]interface{}{
								"id":       "test_param",
								"name":     "Test Parameter",
								"datatype": "int",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid param block, got: %v", collector.Errors())
	}
}

func TestBlockValidator_ParamMissingName(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Parameters",
					Blocks: []domain.Block{
						{
							Type: "param",
							Data: map[string]interface{}{
								"id":       "test_param",
								"datatype": "int",
								// missing "name"
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for param block missing name")
	}

	errs := collector.Errors()
	if errs[0].Code != "E047" {
		t.Errorf("expected E047, got %s", errs[0].Code)
	}
}

func TestBlockValidator_ParamMissingDatatype(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Parameters",
					Blocks: []domain.Block{
						{
							Type: "param",
							Data: map[string]interface{}{
								"id":   "test_param",
								"name": "Test Parameter",
								// missing "datatype"
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for param block missing datatype")
	}

	errs := collector.Errors()
	if errs[0].Code != "E047" {
		t.Errorf("expected E047, got %s", errs[0].Code)
	}
}

func TestBlockValidator_UnknownBlockField(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Parameters",
					Blocks: []domain.Block{
						{
							Type: "param",
							Data: map[string]interface{}{
								"id":       "test_param",
								"name":     "Test Parameter",
								"datatype": "int",
								"datatyp":  "typo",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for unknown block field")
	}

	var unknownErr domain.DecoError
	found := false
	for _, err := range collector.Errors() {
		if err.Code == "E049" {
			unknownErr = err
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected E049 for unknown block field")
	}
	if unknownErr.Suggestion == "" {
		t.Error("expected suggestion for unknown block field")
	}
}

func TestBlockValidator_ValidMechanicBlock(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Mechanics",
					Blocks: []domain.Block{
						{
							Type: "mechanic",
							Data: map[string]interface{}{
								"id":          "test_mechanic",
								"name":        "Test Mechanic",
								"description": "A test mechanic description",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid mechanic block, got: %v", collector.Errors())
	}
}

func TestBlockValidator_MechanicMissingName(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Mechanics",
					Blocks: []domain.Block{
						{
							Type: "mechanic",
							Data: map[string]interface{}{
								"id":          "test_mechanic",
								"description": "A test mechanic description",
								// missing "name"
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for mechanic block missing name")
	}

	errs := collector.Errors()
	if errs[0].Code != "E047" {
		t.Errorf("expected E047, got %s", errs[0].Code)
	}
}

func TestBlockValidator_MechanicMissingDescription(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Mechanics",
					Blocks: []domain.Block{
						{
							Type: "mechanic",
							Data: map[string]interface{}{
								"id":   "test_mechanic",
								"name": "Test Mechanic",
								// missing "description"
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for mechanic block missing description")
	}

	errs := collector.Errors()
	if errs[0].Code != "E047" {
		t.Errorf("expected E047, got %s", errs[0].Code)
	}
}

func TestBlockValidator_ValidListBlock(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Lists",
					Blocks: []domain.Block{
						{
							Type: "list",
							Data: map[string]interface{}{
								"id": "test_list",
								"items": []interface{}{
									"Item 1",
									"Item 2",
								},
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid list block, got: %v", collector.Errors())
	}
}

func TestBlockValidator_ListMissingItems(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Lists",
					Blocks: []domain.Block{
						{
							Type: "list",
							Data: map[string]interface{}{
								"id": "test_list",
								// missing "items"
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for list block missing items")
	}

	errs := collector.Errors()
	if errs[0].Code != "E047" {
		t.Errorf("expected E047, got %s", errs[0].Code)
	}
}

func TestBlockValidator_UnknownBlockType(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Unknown",
					Blocks: []domain.Block{
						{
							Type: "unknown_type",
							Data: map[string]interface{}{
								"id": "test",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for unknown block type")
	}

	errs := collector.Errors()
	if errs[0].Code != "E048" {
		t.Errorf("expected E048, got %s", errs[0].Code)
	}
}

func TestBlockValidator_NilNode(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	validator.Validate(nil, collector)

	if collector.HasErrors() {
		t.Error("expected no errors for nil node")
	}
}

func TestBlockValidator_NilContent(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID:      "test-node",
		Content: nil,
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Error("expected no errors for nil content")
	}
}

func TestBlockValidator_EmptySections(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Error("expected no errors for empty sections")
	}
}

func TestBlockValidator_MultipleBlocks(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Mixed",
					Blocks: []domain.Block{
						{
							Type: "rule",
							Data: map[string]interface{}{
								"id":   "rule1",
								"text": "A valid rule",
							},
						},
						{
							Type: "param",
							Data: map[string]interface{}{
								"id": "param1",
								// missing name and datatype
							},
						},
						{
							Type: "table",
							Data: map[string]interface{}{
								"id": "table1",
								// missing columns and rows
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected errors for invalid blocks")
	}

	// Should have multiple errors: param missing name & datatype, table missing columns & rows
	errs := collector.Errors()
	if len(errs) < 4 {
		t.Errorf("expected at least 4 errors, got %d", len(errs))
	}
}

func TestBlockValidator_ErrorIncludesLocation(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "systems/core",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Movement",
					Blocks: []domain.Block{
						{
							Type: "rule",
							Data: map[string]interface{}{
								"id": "test_rule",
								// missing text
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error")
	}

	err := collector.Errors()[0]
	// Error detail should include node ID, section name, and block index
	if err.Detail == "" {
		t.Error("expected error detail to include location info")
	}
}

func TestBlockValidator_EmptyBlockType(t *testing.T) {
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Test",
					Blocks: []domain.Block{
						{
							Type: "", // empty type
							Data: map[string]interface{}{
								"id": "test",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for empty block type")
	}

	errs := collector.Errors()
	if errs[0].Code != "E048" {
		t.Errorf("expected E048, got %s", errs[0].Code)
	}
}

// Custom block type tests

func TestBlockValidator_CustomType_Valid(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"quest": {
			RequiredFields: []string{"name", "reward"},
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Quests",
					Blocks: []domain.Block{
						{
							Type: "quest",
							Data: map[string]interface{}{
								"name":   "Defeat the Dragon",
								"reward": "100 gold",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid custom block, got: %v", collector.Errors())
	}
}

func TestBlockValidator_CustomType_OptionalFields(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"quest": {
			RequiredFields: []string{"name"},
			OptionalFields: []string{"reward"},
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Quests",
					Blocks: []domain.Block{
						{
							Type: "quest",
							Data: map[string]interface{}{
								"id":     "main_quest",
								"name":   "Defeat the Dragon",
								"reward": "100 gold",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for custom block with optional fields, got: %v", collector.Errors())
	}
}

func TestBlockValidator_CustomType_MissingRequiredField(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"quest": {
			RequiredFields: []string{"name", "reward"},
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Quests",
					Blocks: []domain.Block{
						{
							Type: "quest",
							Data: map[string]interface{}{
								"name": "Defeat the Dragon",
								// missing "reward"
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for custom block missing required field")
	}

	errs := collector.Errors()
	if errs[0].Code != "E047" {
		t.Errorf("expected E047, got %s", errs[0].Code)
	}
}

func TestBlockValidator_CustomType_AllFieldsMissing(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"quest": {
			RequiredFields: []string{"name", "reward", "description"},
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Quests",
					Blocks: []domain.Block{
						{
							Type: "quest",
							Data: map[string]interface{}{
								// all required fields missing
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected errors for missing required fields")
	}

	errs := collector.Errors()
	if len(errs) != 3 {
		t.Errorf("expected 3 errors (one per missing field), got %d", len(errs))
	}
}

func TestBlockValidator_CustomType_NoRequiredFields(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"note": {
			RequiredFields: []string{}, // no required fields
			OptionalFields: []string{"text"},
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Notes",
					Blocks: []domain.Block{
						{
							Type: "note",
							Data: map[string]interface{}{
								"text": "Just a note",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for custom type with no required fields, got: %v", collector.Errors())
	}
}

func TestBlockValidator_CustomType_MultipleCustomTypes(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"quest": {
			RequiredFields: []string{"name"},
		},
		"achievement": {
			RequiredFields: []string{"title", "points"},
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Content",
					Blocks: []domain.Block{
						{
							Type: "quest",
							Data: map[string]interface{}{
								"name": "Main Quest",
							},
						},
						{
							Type: "achievement",
							Data: map[string]interface{}{
								"title":  "First Blood",
								"points": 10,
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid multiple custom types, got: %v", collector.Errors())
	}
}

func TestBlockValidator_CustomType_SuggestsBuiltInAndCustom(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"quest": {
			RequiredFields: []string{"name"},
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Content",
					Blocks: []domain.Block{
						{
							Type: "quets", // typo - should suggest "quest"
							Data: map[string]interface{}{
								"name": "Main Quest",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for unknown block type")
	}

	errs := collector.Errors()
	if errs[0].Code != "E048" {
		t.Errorf("expected E048, got %s", errs[0].Code)
	}
	// Should suggest "quest" since it's close to "quets"
	if errs[0].Suggestion == "" {
		t.Error("expected suggestion for typo")
	}
}

func TestBlockValidator_CustomType_BuiltInStillWorks(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"quest": {
			RequiredFields: []string{"name"},
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Content",
					Blocks: []domain.Block{
						{
							Type: "rule", // built-in type
							Data: map[string]interface{}{
								"text": "This is a rule",
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid built-in type with custom types defined, got: %v", collector.Errors())
	}
}

func TestBlockValidator_CustomType_ExtendsBuiltIn(t *testing.T) {
	// Custom config can extend built-in types by adding more required fields
	customTypes := map[string]config.BlockTypeConfig{
		"rule": {
			RequiredFields: []string{"priority"}, // adds "priority" requirement on top of built-in "text"
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Rules",
					Blocks: []domain.Block{
						{
							Type: "rule",
							Data: map[string]interface{}{
								"text": "This is a rule",
								// missing "priority" which custom config requires
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	// Custom extends built-in, so missing "priority" should cause error
	if !collector.HasErrors() {
		t.Fatal("expected error when custom extends built-in and field is missing")
	}

	errs := collector.Errors()
	if len(errs) != 1 {
		t.Errorf("expected 1 error for missing priority, got %d", len(errs))
	}
}

func TestBlockValidator_CustomType_ExtendsBuiltIn_BothRequired(t *testing.T) {
	// Verify that both built-in and custom requirements are enforced
	customTypes := map[string]config.BlockTypeConfig{
		"rule": {
			RequiredFields: []string{"priority"},
		},
	}
	validator := NewBlockValidatorWithConfig(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Rules",
					Blocks: []domain.Block{
						{
							Type: "rule",
							Data: map[string]interface{}{
								// missing both "text" (built-in) and "priority" (custom)
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected errors for missing fields")
	}

	errs := collector.Errors()
	if len(errs) != 2 {
		t.Errorf("expected 2 errors (text from built-in, priority from custom), got %d", len(errs))
	}
}

// Test error consolidation for table column typos

func TestBlockValidator_TableColumnTypo_Consolidates(t *testing.T) {
	// When a typo causes both "unknown field" and "missing required field" errors,
	// they should be consolidated into a single error
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Data",
					Blocks: []domain.Block{
						{
							Type: "table",
							Data: map[string]interface{}{
								"id": "test_table",
								"columns": []interface{}{
									map[string]interface{}{
										"ke":      "name", // typo for "key"
										"type":    "string",
										"display": "Name",
									},
								},
								"rows": []interface{}{
									map[string]interface{}{"name": "Test"},
								},
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for table column typo")
	}

	errs := collector.Errors()
	// Should be exactly 1 error (consolidated), not 2
	if len(errs) != 1 {
		t.Errorf("expected 1 consolidated error, got %d", len(errs))
		for i, err := range errs {
			t.Logf("  error %d: [%s] %s", i, err.Code, err.Summary)
		}
	}

	// The error should be E049 with consolidated message
	if errs[0].Code != "E049" {
		t.Errorf("expected E049, got %s", errs[0].Code)
	}

	// Should have context explaining the consolidation
	if len(errs[0].Context) < 2 {
		t.Errorf("expected at least 2 context items, got %d", len(errs[0].Context))
	}

	// First context should show column contents
	foundContents := false
	foundAlsoCauses := false
	for _, ctx := range errs[0].Context {
		if len(ctx) > 0 && ctx[0:6] == "Column" {
			foundContents = true
		}
		if len(ctx) > 10 && ctx[0:9] == "This also" {
			foundAlsoCauses = true
		}
	}
	if !foundContents {
		t.Error("expected context to include column contents")
	}
	if !foundAlsoCauses {
		t.Error("expected context to include 'also causes' note")
	}
}

func TestBlockValidator_TableColumnTypo_NoConsolidationWhenUnrelated(t *testing.T) {
	// When unknown field doesn't suggest the missing field, no consolidation
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Data",
					Blocks: []domain.Block{
						{
							Type: "table",
							Data: map[string]interface{}{
								"id": "test_table",
								"columns": []interface{}{
									map[string]interface{}{
										"xyz":     "name", // not a typo for "key"
										"type":    "string",
										"display": "Name",
									},
								},
								"rows": []interface{}{
									map[string]interface{}{"name": "Test"},
								},
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected errors for table column issues")
	}

	errs := collector.Errors()
	// Should be 2 errors (not consolidated since "xyz" doesn't suggest "key")
	if len(errs) != 2 {
		t.Errorf("expected 2 separate errors, got %d", len(errs))
	}

	// Should have one E050 (missing key) and one E049 (unknown field)
	hasE050 := false
	hasE049 := false
	for _, err := range errs {
		if err.Code == "E050" {
			hasE050 = true
		}
		if err.Code == "E049" {
			hasE049 = true
		}
	}
	if !hasE050 {
		t.Error("expected E050 for missing key")
	}
	if !hasE049 {
		t.Error("expected E049 for unknown field")
	}
}

func TestBlockValidator_TableColumnShowsContents(t *testing.T) {
	// Verify that column contents are shown in error context
	validator := NewBlockValidator()
	collector := errors.NewCollectorWithLimit(100)

	node := domain.Node{
		ID: "test-node",
		Content: &domain.Content{
			Sections: []domain.Section{
				{
					Name: "Data",
					Blocks: []domain.Block{
						{
							Type: "table",
							Data: map[string]interface{}{
								"id": "test_table",
								"columns": []interface{}{
									map[string]interface{}{
										"type":    "string",
										"display": "Name",
										// missing "key"
									},
								},
								"rows": []interface{}{
									map[string]interface{}{"name": "Test"},
								},
							},
						},
					},
				},
			},
		},
	}

	validator.Validate(&node, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for missing key")
	}

	errs := collector.Errors()
	if len(errs[0].Context) == 0 {
		t.Fatal("expected context with column contents")
	}

	// Context should mention the column contents
	foundContents := false
	for _, ctx := range errs[0].Context {
		if len(ctx) > 6 && ctx[0:6] == "Column" {
			foundContents = true
			// Should contain the actual field values
			if !contains(ctx, "display") || !contains(ctx, "type") {
				t.Errorf("expected context to show field values, got: %s", ctx)
			}
		}
	}
	if !foundContents {
		t.Error("expected context to include column contents")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}
