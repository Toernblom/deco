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

package yaml_test

import (
	"testing"

	yaml_errors "github.com/Toernblom/deco/internal/errors/yaml"
)

// Test basic line number tracking for simple YAML
func TestLocationTracker_SimpleFields(t *testing.T) {
	yamlContent := `id: systems/food
kind: system
version: 1
status: draft
title: Food System`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	tests := []struct {
		path         string
		expectedLine int
		expectedCol  int
	}{
		{"id", 1, 1},
		{"kind", 2, 1},
		{"version", 3, 1},
		{"status", 4, 1},
		{"title", 5, 1},
	}

	for _, tt := range tests {
		loc := tracker.GetLocation(tt.path)
		if loc.Line != tt.expectedLine {
			t.Errorf("Path %q: expected line %d, got %d", tt.path, tt.expectedLine, loc.Line)
		}
		if loc.Column != tt.expectedCol {
			t.Errorf("Path %q: expected column %d, got %d", tt.path, tt.expectedCol, loc.Column)
		}
	}
}

// Test column positions for keys and values
func TestLocationTracker_ColumnPositions(t *testing.T) {
	yamlContent := `id: systems/food
metadata:
  kind: system
  details:
    version: 1`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	tests := []struct {
		path         string
		expectedLine int
		expectedCol  int
	}{
		{"id", 1, 1},
		{"metadata", 2, 1},
		{"metadata.kind", 3, 3},            // Indented by 2 spaces
		{"metadata.details", 4, 3},         // Indented by 2 spaces
		{"metadata.details.version", 5, 5}, // Indented by 4 spaces
	}

	for _, tt := range tests {
		loc := tracker.GetLocation(tt.path)
		if loc.Line != tt.expectedLine {
			t.Errorf("Path %q: expected line %d, got %d", tt.path, tt.expectedLine, loc.Line)
		}
		if loc.Column != tt.expectedCol {
			t.Errorf("Path %q: expected column %d, got %d", tt.path, tt.expectedCol, loc.Column)
		}
	}
}

// Test nested structure tracking
func TestLocationTracker_NestedStructures(t *testing.T) {
	yamlContent := `id: systems/food
metadata:
  author: alice
  created: 2024-01-01
  tags:
    - survival
    - resource
constraints:
  requires:
    - systems/water`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	tests := []struct {
		path         string
		expectedLine int
	}{
		{"id", 1},
		{"metadata", 2},
		{"metadata.author", 3},
		{"metadata.created", 4},
		{"metadata.tags", 5},
		{"metadata.tags[0]", 6},
		{"metadata.tags[1]", 7},
		{"constraints", 8},
		{"constraints.requires", 9},
		{"constraints.requires[0]", 10},
	}

	for _, tt := range tests {
		loc := tracker.GetLocation(tt.path)
		if loc.Line != tt.expectedLine {
			t.Errorf("Path %q: expected line %d, got %d", tt.path, tt.expectedLine, loc.Line)
		}
	}
}

// Test multiline value tracking
func TestLocationTracker_MultilineValues(t *testing.T) {
	yamlContent := `id: systems/food
description: |
  This is a multiline
  description that spans
  multiple lines
title: Food System`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	// The description key should be on line 2
	loc := tracker.GetLocation("description")
	if loc.Line != 2 {
		t.Errorf("Expected description on line 2, got %d", loc.Line)
	}

	// The title should be on line 6 (after the multiline value)
	loc = tracker.GetLocation("title")
	if loc.Line != 6 {
		t.Errorf("Expected title on line 6, got %d", loc.Line)
	}
}

// Test deeply nested paths
func TestLocationTracker_DeeplyNested(t *testing.T) {
	yamlContent := `systems:
  food:
    production:
      farming:
        crops:
          - wheat
          - corn`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	tests := []struct {
		path         string
		expectedLine int
	}{
		{"systems", 1},
		{"systems.food", 2},
		{"systems.food.production", 3},
		{"systems.food.production.farming", 4},
		{"systems.food.production.farming.crops", 5},
		{"systems.food.production.farming.crops[0]", 6},
		{"systems.food.production.farming.crops[1]", 7},
	}

	for _, tt := range tests {
		loc := tracker.GetLocation(tt.path)
		if loc.Line != tt.expectedLine {
			t.Errorf("Path %q: expected line %d, got %d", tt.path, tt.expectedLine, loc.Line)
		}
	}
}

// Test array of objects
func TestLocationTracker_ArrayOfObjects(t *testing.T) {
	yamlContent := `items:
  - id: item1
    name: Sword
  - id: item2
    name: Shield`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	tests := []struct {
		path         string
		expectedLine int
	}{
		{"items", 1},
		{"items[0]", 2},
		{"items[0].id", 2},
		{"items[0].name", 3},
		{"items[1]", 4},
		{"items[1].id", 4},
		{"items[1].name", 5},
	}

	for _, tt := range tests {
		loc := tracker.GetLocation(tt.path)
		if loc.Line != tt.expectedLine {
			t.Errorf("Path %q: expected line %d, got %d", tt.path, tt.expectedLine, loc.Line)
		}
	}
}

// Test invalid path returns zero location
func TestLocationTracker_InvalidPath(t *testing.T) {
	yamlContent := `id: systems/food
kind: system`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	loc := tracker.GetLocation("nonexistent.path")
	if loc.Line != 0 || loc.Column != 0 {
		t.Errorf("Expected zero location for invalid path, got line=%d col=%d", loc.Line, loc.Column)
	}
}

// Test empty YAML
func TestLocationTracker_EmptyYAML(t *testing.T) {
	yamlContent := ``

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	loc := tracker.GetLocation("any.path")
	if loc.Line != 0 {
		t.Errorf("Expected zero location for empty YAML, got line=%d", loc.Line)
	}
}

// Test malformed YAML returns error
func TestLocationTracker_MalformedYAML(t *testing.T) {
	yamlContent := `id: systems/food
kind: [invalid yaml structure`

	_, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err == nil {
		t.Error("Expected error for malformed YAML")
	}
}

// Test mixed content (maps and arrays)
func TestLocationTracker_MixedContent(t *testing.T) {
	yamlContent := `id: systems/food
tags:
  - survival
  - resource
metadata:
  author: alice
  versions:
    - 1
    - 2
    - 3`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	tests := []struct {
		path         string
		expectedLine int
	}{
		{"id", 1},
		{"tags", 2},
		{"tags[0]", 3},
		{"tags[1]", 4},
		{"metadata", 5},
		{"metadata.author", 6},
		{"metadata.versions", 7},
		{"metadata.versions[0]", 8},
		{"metadata.versions[1]", 9},
		{"metadata.versions[2]", 10},
	}

	for _, tt := range tests {
		loc := tracker.GetLocation(tt.path)
		if loc.Line != tt.expectedLine {
			t.Errorf("Path %q: expected line %d, got %d", tt.path, tt.expectedLine, loc.Line)
		}
	}
}

// Test root-level array
func TestLocationTracker_RootArray(t *testing.T) {
	yamlContent := `- id: item1
  name: First
- id: item2
  name: Second`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	tests := []struct {
		path         string
		expectedLine int
	}{
		{"[0]", 1},
		{"[0].id", 1},
		{"[0].name", 2},
		{"[1]", 3},
		{"[1].id", 3},
		{"[1].name", 4},
	}

	for _, tt := range tests {
		loc := tracker.GetLocation(tt.path)
		if loc.Line != tt.expectedLine {
			t.Errorf("Path %q: expected line %d, got %d", tt.path, tt.expectedLine, loc.Line)
		}
	}
}

// Test value location vs key location
func TestLocationTracker_ValueLocation(t *testing.T) {
	yamlContent := `id: systems/food
title: Food System`

	tracker, err := yaml_errors.NewLocationTracker([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	// GetLocation should return the key location
	loc := tracker.GetLocation("id")
	if loc.Line != 1 || loc.Column != 1 {
		t.Errorf("Expected key location at 1:1, got %d:%d", loc.Line, loc.Column)
	}

	// GetValueLocation should return the value location
	valLoc := tracker.GetValueLocation("id")
	if valLoc.Line != 1 {
		t.Errorf("Expected value on line 1, got %d", valLoc.Line)
	}
	// Value should be after "id: " (column 5)
	if valLoc.Column < 4 {
		t.Errorf("Expected value column >= 4, got %d", valLoc.Column)
	}
}

// Test with file path
func TestLocationTracker_WithFilePath(t *testing.T) {
	yamlContent := `id: systems/food`

	tracker, err := yaml_errors.NewLocationTrackerWithFile([]byte(yamlContent), "test.yaml")
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	loc := tracker.GetLocation("id")
	if loc.File != "test.yaml" {
		t.Errorf("Expected file 'test.yaml', got %q", loc.File)
	}
	if loc.Line != 1 {
		t.Errorf("Expected line 1, got %d", loc.Line)
	}
}
