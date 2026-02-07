package validator

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
	"github.com/Toernblom/deco/internal/storage/config"
)

func TestCrossRef_SingleField_Valid(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name":     {Type: "string", Required: true},
				"material": {Type: "string", Ref: &config.RefConstraint{BlockType: "resource", Field: "name"}},
			},
		},
		"resource": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
	}

	validator := NewCrossRefValidator(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	nodes := []domain.Node{
		{
			ID: "buildings", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smithy", "material": "Stone"}},
						},
					},
				},
			},
		},
		{
			ID: "resources", Kind: "system", Version: 1, Status: "draft", Title: "Resources",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Materials",
						Blocks: []domain.Block{
							{Type: "resource", Data: map[string]interface{}{"name": "Stone"}},
							{Type: "resource", Data: map[string]interface{}{"name": "Iron"}},
						},
					},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid cross-ref, got: %v", collector.Errors())
	}
}

func TestCrossRef_SingleField_Invalid(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name":     {Type: "string", Required: true},
				"material": {Type: "string", Ref: &config.RefConstraint{BlockType: "resource", Field: "name"}},
			},
		},
		"resource": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
	}

	validator := NewCrossRefValidator(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	nodes := []domain.Node{
		{
			ID: "buildings", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smithy", "material": "Stoone"}}, // typo
						},
					},
				},
			},
		},
		{
			ID: "resources", Kind: "system", Version: 1, Status: "draft", Title: "Resources",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Materials",
						Blocks: []domain.Block{
							{Type: "resource", Data: map[string]interface{}{"name": "Stone"}},
							{Type: "resource", Data: map[string]interface{}{"name": "Iron"}},
						},
					},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for invalid cross-ref")
	}
	errs := collector.Errors()
	foundE054 := false
	for _, err := range errs {
		if err.Code == "E054" {
			foundE054 = true
			if err.Suggestion == "" {
				t.Error("expected did-you-mean suggestion")
			}
		}
	}
	if !foundE054 {
		t.Errorf("expected E054 for cross-ref not found, got codes: %v", errorCodes(errs))
	}
}

func TestCrossRef_ListField_Valid(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name":      {Type: "string", Required: true},
				"materials": {Type: "list", Ref: &config.RefConstraint{BlockType: "resource", Field: "name"}},
			},
		},
		"resource": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
	}

	validator := NewCrossRefValidator(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	nodes := []domain.Node{
		{
			ID: "buildings", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{
								"name":      "Smithy",
								"materials": []interface{}{"Stone", "Bronze"},
							}},
						},
					},
				},
			},
		},
		{
			ID: "resources", Kind: "system", Version: 1, Status: "draft", Title: "Resources",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Materials",
						Blocks: []domain.Block{
							{Type: "resource", Data: map[string]interface{}{"name": "Stone"}},
							{Type: "resource", Data: map[string]interface{}{"name": "Bronze"}},
							{Type: "resource", Data: map[string]interface{}{"name": "Iron"}},
						},
					},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid list cross-ref, got: %v", collector.Errors())
	}
}

func TestCrossRef_ListField_PartialInvalid(t *testing.T) {
	customTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name":      {Type: "string", Required: true},
				"materials": {Type: "list", Ref: &config.RefConstraint{BlockType: "resource", Field: "name"}},
			},
		},
		"resource": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
	}

	validator := NewCrossRefValidator(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	nodes := []domain.Node{
		{
			ID: "buildings", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{
								"name":      "Smithy",
								"materials": []interface{}{"Stone", "Brnze"}, // one valid, one typo
							}},
						},
					},
				},
			},
		},
		{
			ID: "resources", Kind: "system", Version: 1, Status: "draft", Title: "Resources",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Materials",
						Blocks: []domain.Block{
							{Type: "resource", Data: map[string]interface{}{"name": "Stone"}},
							{Type: "resource", Data: map[string]interface{}{"name": "Bronze"}},
						},
					},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for partial invalid list cross-ref")
	}
	errs := collector.Errors()
	if len(errs) != 1 {
		t.Errorf("expected exactly 1 error (for 'Brnze'), got %d", len(errs))
	}
	if errs[0].Code != "E054" {
		t.Errorf("expected E054, got %s", errs[0].Code)
	}
}

func TestCrossRef_NoInstances_NoError(t *testing.T) {
	// When the referenced block type has no instances, field values should still be checked
	// but missing reference set means all values fail
	customTypes := map[string]config.BlockTypeConfig{
		"recipe": {
			Fields: map[string]config.FieldDef{
				"building": {Type: "string", Ref: &config.RefConstraint{BlockType: "building", Field: "name"}},
			},
		},
		"building": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
	}

	validator := NewCrossRefValidator(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	// No building blocks exist anywhere
	nodes := []domain.Node{
		{
			ID: "recipes", Kind: "system", Version: 1, Status: "draft", Title: "Recipes",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Crafting",
						Blocks: []domain.Block{
							{Type: "recipe", Data: map[string]interface{}{"building": "Smithy"}},
						},
					},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	// Should produce an error since "Smithy" can't be resolved
	if !collector.HasErrors() {
		t.Fatal("expected error when referenced block type has no instances")
	}
}

func TestCrossRef_SelfReferencing(t *testing.T) {
	// A block type can reference itself (e.g., resource depends on other resources)
	customTypes := map[string]config.BlockTypeConfig{
		"resource": {
			Fields: map[string]config.FieldDef{
				"name":    {Type: "string", Required: true},
				"refined": {Type: "string", Ref: &config.RefConstraint{BlockType: "resource", Field: "name"}},
			},
		},
	}

	validator := NewCrossRefValidator(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	nodes := []domain.Node{
		{
			ID: "resources", Kind: "system", Version: 1, Status: "draft", Title: "Resources",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Materials",
						Blocks: []domain.Block{
							{Type: "resource", Data: map[string]interface{}{"name": "Iron Ore"}},
							{Type: "resource", Data: map[string]interface{}{"name": "Iron Bar", "refined": "Iron Ore"}},
						},
					},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid self-reference, got: %v", collector.Errors())
	}
}

func TestCrossRef_AcrossMultipleNodes(t *testing.T) {
	// Cross-refs should work across different nodes
	customTypes := map[string]config.BlockTypeConfig{
		"recipe": {
			Fields: map[string]config.FieldDef{
				"output":   {Type: "string", Required: true, Ref: &config.RefConstraint{BlockType: "resource", Field: "name"}},
				"building": {Type: "string", Required: true, Ref: &config.RefConstraint{BlockType: "building", Field: "name"}},
				"inputs":   {Type: "list", Ref: &config.RefConstraint{BlockType: "resource", Field: "name"}},
			},
		},
		"resource": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
		"building": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
	}

	validator := NewCrossRefValidator(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	nodes := []domain.Node{
		{
			ID: "resources", Kind: "system", Version: 1, Status: "draft", Title: "Resources",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Materials",
						Blocks: []domain.Block{
							{Type: "resource", Data: map[string]interface{}{"name": "Iron Ore"}},
							{Type: "resource", Data: map[string]interface{}{"name": "Iron Bar"}},
						},
					},
				},
			},
		},
		{
			ID: "buildings", Kind: "system", Version: 1, Status: "draft", Title: "Buildings",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Structures",
						Blocks: []domain.Block{
							{Type: "building", Data: map[string]interface{}{"name": "Smelter"}},
						},
					},
				},
			},
		},
		{
			ID: "recipes", Kind: "system", Version: 1, Status: "draft", Title: "Recipes",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Crafting",
						Blocks: []domain.Block{
							{Type: "recipe", Data: map[string]interface{}{
								"output":   "Iron Bar",
								"building": "Smelter",
								"inputs":   []interface{}{"Iron Ore"},
							}},
						},
					},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid cross-node refs, got: %v", collector.Errors())
	}
}

func TestCrossRef_NilCustomTypes(t *testing.T) {
	validator := NewCrossRefValidator(nil)
	collector := errors.NewCollectorWithLimit(100)

	nodes := []domain.Node{
		{ID: "node1", Kind: "system", Version: 1, Status: "draft", Title: "Node 1"},
	}

	validator.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Error("expected no errors for nil custom types")
	}
}

func TestCrossRef_NoFieldsConfig(t *testing.T) {
	// Old-style config without Fields should not trigger cross-ref validation
	customTypes := map[string]config.BlockTypeConfig{
		"quest": {
			RequiredFields: []string{"name"},
		},
	}

	validator := NewCrossRefValidator(customTypes)
	collector := errors.NewCollectorWithLimit(100)

	nodes := []domain.Node{
		{
			ID: "quests", Kind: "system", Version: 1, Status: "draft", Title: "Quests",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Main",
						Blocks: []domain.Block{
							{Type: "quest", Data: map[string]interface{}{"name": "Main Quest"}},
						},
					},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Error("expected no errors for config without Fields")
	}
}
