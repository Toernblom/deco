// Copyright (C) 2026 Anton Törnblom
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

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
				"material": {Type: "string", Refs: []config.RefConstraint{{BlockType: "resource", Field: "name"}}},
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
				"material": {Type: "string", Refs: []config.RefConstraint{{BlockType: "resource", Field: "name"}}},
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
				"materials": {Type: "list", Refs: []config.RefConstraint{{BlockType: "resource", Field: "name"}}},
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
				"materials": {Type: "list", Refs: []config.RefConstraint{{BlockType: "resource", Field: "name"}}},
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
				"building": {Type: "string", Refs: []config.RefConstraint{{BlockType: "building", Field: "name"}}},
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
				"refined": {Type: "string", Refs: []config.RefConstraint{{BlockType: "resource", Field: "name"}}},
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
				"output":   {Type: "string", Required: true, Refs: []config.RefConstraint{{BlockType: "resource", Field: "name"}}},
				"building": {Type: "string", Required: true, Refs: []config.RefConstraint{{BlockType: "building", Field: "name"}}},
				"inputs":   {Type: "list", Refs: []config.RefConstraint{{BlockType: "resource", Field: "name"}}},
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

// ===== UNION REF TESTS =====

func TestCrossRef_UnionRef_Valid(t *testing.T) {
	// building.materials refs both resource.name and recipe.output (OR logic)
	customTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
				"materials": {Type: "list", Refs: []config.RefConstraint{
					{BlockType: "resource", Field: "name"},
					{BlockType: "recipe", Field: "output"},
				}},
			},
		},
		"resource": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
		"recipe": {
			Fields: map[string]config.FieldDef{
				"output": {Type: "string", Required: true},
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
								"materials": []interface{}{"Stone", "Planks"}, // Stone=resource, Planks=recipe output
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
					{Name: "Raw", Blocks: []domain.Block{
						{Type: "resource", Data: map[string]interface{}{"name": "Stone"}},
					}},
				},
			},
		},
		{
			ID: "recipes", Kind: "system", Version: 1, Status: "draft", Title: "Recipes",
			Content: &domain.Content{
				Sections: []domain.Section{
					{Name: "Crafting", Blocks: []domain.Block{
						{Type: "recipe", Data: map[string]interface{}{"output": "Planks"}},
					}},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for valid union refs, got: %v", collector.Errors())
	}
}

func TestCrossRef_UnionRef_Invalid(t *testing.T) {
	// "Plonks" exists in neither resource.name nor recipe.output
	customTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
				"materials": {Type: "list", Refs: []config.RefConstraint{
					{BlockType: "resource", Field: "name"},
					{BlockType: "recipe", Field: "output"},
				}},
			},
		},
		"resource": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
		"recipe": {
			Fields: map[string]config.FieldDef{
				"output": {Type: "string", Required: true},
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
								"materials": []interface{}{"Stone", "Plonks"}, // Plonks is a typo
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
					{Name: "Raw", Blocks: []domain.Block{
						{Type: "resource", Data: map[string]interface{}{"name": "Stone"}},
					}},
				},
			},
		},
		{
			ID: "recipes", Kind: "system", Version: 1, Status: "draft", Title: "Recipes",
			Content: &domain.Content{
				Sections: []domain.Section{
					{Name: "Crafting", Blocks: []domain.Block{
						{Type: "recipe", Data: map[string]interface{}{"output": "Planks"}},
					}},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for invalid union ref value 'Plonks'")
	}
	errs := collector.Errors()
	if len(errs) != 1 {
		t.Errorf("expected exactly 1 error, got %d", len(errs))
	}
	if errs[0].Code != "E054" {
		t.Errorf("expected E054, got %s", errs[0].Code)
	}
	// Error should mention checked targets
	if errs[0].Suggestion == "" {
		t.Error("expected did-you-mean suggestion for 'Plonks' -> 'Planks'")
	}
}

func TestCrossRef_UnionRef_OverlappingValues(t *testing.T) {
	// "Stone" exists in both resource.name and recipe.output — should validate fine
	customTypes := map[string]config.BlockTypeConfig{
		"building": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
				"materials": {Type: "list", Refs: []config.RefConstraint{
					{BlockType: "resource", Field: "name"},
					{BlockType: "recipe", Field: "output"},
				}},
			},
		},
		"resource": {
			Fields: map[string]config.FieldDef{
				"name": {Type: "string", Required: true},
			},
		},
		"recipe": {
			Fields: map[string]config.FieldDef{
				"output": {Type: "string", Required: true},
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
					{Name: "Structures", Blocks: []domain.Block{
						{Type: "building", Data: map[string]interface{}{
							"name":      "Quarry",
							"materials": []interface{}{"Stone"},
						}},
					}},
				},
			},
		},
		{
			ID: "resources", Kind: "system", Version: 1, Status: "draft", Title: "Resources",
			Content: &domain.Content{
				Sections: []domain.Section{
					{Name: "Raw", Blocks: []domain.Block{
						{Type: "resource", Data: map[string]interface{}{"name": "Stone"}},
					}},
				},
			},
		},
		{
			ID: "recipes", Kind: "system", Version: 1, Status: "draft", Title: "Recipes",
			Content: &domain.Content{
				Sections: []domain.Section{
					{Name: "Crafting", Blocks: []domain.Block{
						{Type: "recipe", Data: map[string]interface{}{"output": "Stone"}}, // Also outputs Stone
					}},
				},
			},
		},
	}

	validator.Validate(nodes, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for overlapping union ref values, got: %v", collector.Errors())
	}
}
