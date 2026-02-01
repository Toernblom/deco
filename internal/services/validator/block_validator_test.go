package validator

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
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
