package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"gopkg.in/yaml.v3"
)

func TestNode_Creation(t *testing.T) {
	node := domain.Node{
		ID:      "systems/settlement/housing",
		Kind:    "mechanic",
		Version: 1,
		Status:  "draft",
		Title:   "Housing System",
	}

	if node.ID != "systems/settlement/housing" {
		t.Errorf("expected ID 'systems/settlement/housing', got %q", node.ID)
	}
	if node.Kind != "mechanic" {
		t.Errorf("expected Kind 'mechanic', got %q", node.Kind)
	}
	if node.Version != 1 {
		t.Errorf("expected Version 1, got %d", node.Version)
	}
	if node.Status != "draft" {
		t.Errorf("expected Status 'draft', got %q", node.Status)
	}
	if node.Title != "Housing System" {
		t.Errorf("expected Title 'Housing System', got %q", node.Title)
	}
}

func TestNode_RequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		node    domain.Node
		wantErr bool
	}{
		{
			name: "valid node with all required fields",
			node: domain.Node{
				ID:      "test/node",
				Kind:    "mechanic",
				Version: 1,
				Status:  "draft",
				Title:   "Test Node",
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			node: domain.Node{
				Kind:    "mechanic",
				Version: 1,
				Status:  "draft",
				Title:   "Test Node",
			},
			wantErr: true,
		},
		{
			name: "missing Kind",
			node: domain.Node{
				ID:      "test/node",
				Version: 1,
				Status:  "draft",
				Title:   "Test Node",
			},
			wantErr: true,
		},
		{
			name: "missing Version",
			node: domain.Node{
				ID:     "test/node",
				Kind:   "mechanic",
				Status: "draft",
				Title:  "Test Node",
			},
			wantErr: true,
		},
		{
			name: "missing Status",
			node: domain.Node{
				ID:      "test/node",
				Kind:    "mechanic",
				Version: 1,
				Title:   "Test Node",
			},
			wantErr: true,
		},
		{
			name: "missing Title",
			node: domain.Node{
				ID:      "test/node",
				Kind:    "mechanic",
				Version: 1,
				Status:  "draft",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNode_Serialization(t *testing.T) {
	original := domain.Node{
		ID:      "test/node",
		Kind:    "mechanic",
		Version: 1,
		Status:  "draft",
		Title:   "Test Node",
		Tags:    []string{"core", "settlement"},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Node: %v", err)
	}

	// Unmarshal back
	var restored domain.Node
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal Node: %v", err)
	}

	// Compare
	if restored.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", restored.ID, original.ID)
	}
	if restored.Kind != original.Kind {
		t.Errorf("Kind mismatch: got %q, want %q", restored.Kind, original.Kind)
	}
	if restored.Version != original.Version {
		t.Errorf("Version mismatch: got %d, want %d", restored.Version, original.Version)
	}
	if restored.Status != original.Status {
		t.Errorf("Status mismatch: got %q, want %q", restored.Status, original.Status)
	}
	if restored.Title != original.Title {
		t.Errorf("Title mismatch: got %q, want %q", restored.Title, original.Title)
	}
	if len(restored.Tags) != len(original.Tags) {
		t.Errorf("Tags length mismatch: got %d, want %d", len(restored.Tags), len(original.Tags))
	}
}

func TestNode_FieldAccess(t *testing.T) {
	node := domain.Node{
		ID:      "test/node",
		Kind:    "mechanic",
		Version: 2,
		Status:  "approved",
		Title:   "Test Node",
		Tags:    []string{"tag1", "tag2"},
	}

	// Test field access
	if node.ID != "test/node" {
		t.Errorf("ID field access failed")
	}
	if node.Kind != "mechanic" {
		t.Errorf("Kind field access failed")
	}
	if node.Version != 2 {
		t.Errorf("Version field access failed")
	}
	if node.Status != "approved" {
		t.Errorf("Status field access failed")
	}
	if node.Title != "Test Node" {
		t.Errorf("Title field access failed")
	}
	if len(node.Tags) != 2 {
		t.Errorf("Tags field access failed")
	}
	if node.Tags[0] != "tag1" || node.Tags[1] != "tag2" {
		t.Errorf("Tags values incorrect")
	}
}

// ===== BLOCK YAML TESTS =====

// Test Block YAML unmarshaling with inline fields
func TestBlock_UnmarshalYAML_InlineFields(t *testing.T) {
	yamlData := `
type: table
id: test_table
columns:
  - name: column1
  - name: column2
rows:
  - [a, b]
  - [c, d]
`
	var block domain.Block
	err := yaml.Unmarshal([]byte(yamlData), &block)
	if err != nil {
		t.Fatalf("failed to unmarshal block: %v", err)
	}

	if block.Type != "table" {
		t.Errorf("expected Type 'table', got %q", block.Type)
	}

	if block.Data == nil {
		t.Fatal("expected Data to be non-nil")
	}

	if block.Data["id"] != "test_table" {
		t.Errorf("expected Data[id] 'test_table', got %v", block.Data["id"])
	}

	columns, ok := block.Data["columns"].([]interface{})
	if !ok {
		t.Fatalf("expected Data[columns] to be []interface{}, got %T", block.Data["columns"])
	}
	if len(columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(columns))
	}

	rows, ok := block.Data["rows"].([]interface{})
	if !ok {
		t.Fatalf("expected Data[rows] to be []interface{}, got %T", block.Data["rows"])
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}
}

// Test Block YAML unmarshaling with explicit data field
func TestBlock_UnmarshalYAML_ExplicitData(t *testing.T) {
	yamlData := `
type: rule
data:
  id: test_rule
  text: "This is a rule"
`
	var block domain.Block
	err := yaml.Unmarshal([]byte(yamlData), &block)
	if err != nil {
		t.Fatalf("failed to unmarshal block: %v", err)
	}

	if block.Type != "rule" {
		t.Errorf("expected Type 'rule', got %q", block.Type)
	}

	if block.Data["id"] != "test_rule" {
		t.Errorf("expected Data[id] 'test_rule', got %v", block.Data["id"])
	}

	if block.Data["text"] != "This is a rule" {
		t.Errorf("expected Data[text] 'This is a rule', got %v", block.Data["text"])
	}
}

// Test Block YAML marshaling preserves inline fields
func TestBlock_MarshalYAML_InlineFields(t *testing.T) {
	block := domain.Block{
		Type: "param",
		Data: map[string]interface{}{
			"id":    "test_param",
			"name":  "Test Parameter",
			"value": 42,
		},
	}

	data, err := yaml.Marshal(&block)
	if err != nil {
		t.Fatalf("failed to marshal block: %v", err)
	}

	// Unmarshal back to verify round-trip
	var restored domain.Block
	err = yaml.Unmarshal(data, &restored)
	if err != nil {
		t.Fatalf("failed to unmarshal marshaled block: %v", err)
	}

	if restored.Type != block.Type {
		t.Errorf("Type mismatch: got %q, want %q", restored.Type, block.Type)
	}

	if restored.Data["id"] != block.Data["id"] {
		t.Errorf("Data[id] mismatch: got %v, want %v", restored.Data["id"], block.Data["id"])
	}

	if restored.Data["name"] != block.Data["name"] {
		t.Errorf("Data[name] mismatch: got %v, want %v", restored.Data["name"], block.Data["name"])
	}

	// Note: YAML may convert int to int64 or float64, so compare as numbers
	restoredValue, ok := restored.Data["value"].(int)
	if !ok {
		// Try float64 (YAML sometimes decodes as float)
		if fv, ok := restored.Data["value"].(float64); ok {
			restoredValue = int(fv)
		} else {
			t.Fatalf("Data[value] has unexpected type %T", restored.Data["value"])
		}
	}
	if restoredValue != 42 {
		t.Errorf("Data[value] mismatch: got %v, want 42", restored.Data["value"])
	}
}

// Test Block YAML round-trip preserves all data
func TestBlock_YAMLRoundTrip(t *testing.T) {
	yamlData := `
type: table
id: combat_results
columns:
  - key: attacker
    display: Attacker
  - key: defender
    display: Defender
  - key: result
    display: Result
rows:
  - attacker: Knight
    defender: Goblin
    result: Victory
  - attacker: Peasant
    defender: Dragon
    result: Defeat
`
	var original domain.Block
	err := yaml.Unmarshal([]byte(yamlData), &original)
	if err != nil {
		t.Fatalf("failed to unmarshal original: %v", err)
	}

	// Marshal back to YAML
	marshaled, err := yaml.Marshal(&original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal the marshaled data
	var restored domain.Block
	err = yaml.Unmarshal(marshaled, &restored)
	if err != nil {
		t.Fatalf("failed to unmarshal marshaled: %v", err)
	}

	// Compare
	if restored.Type != original.Type {
		t.Errorf("Type mismatch after round-trip")
	}

	if len(restored.Data) != len(original.Data) {
		t.Errorf("Data length mismatch: got %d, want %d", len(restored.Data), len(original.Data))
	}

	// Check that key fields are preserved
	if restored.Data["id"] != original.Data["id"] {
		t.Errorf("Data[id] mismatch after round-trip")
	}

	// Check columns preserved
	restoredCols, ok1 := restored.Data["columns"].([]interface{})
	originalCols, ok2 := original.Data["columns"].([]interface{})
	if !ok1 || !ok2 {
		t.Fatalf("columns type mismatch")
	}
	if len(restoredCols) != len(originalCols) {
		t.Errorf("columns length mismatch: got %d, want %d", len(restoredCols), len(originalCols))
	}

	// Check rows preserved
	restoredRows, ok1 := restored.Data["rows"].([]interface{})
	originalRows, ok2 := original.Data["rows"].([]interface{})
	if !ok1 || !ok2 {
		t.Fatalf("rows type mismatch")
	}
	if len(restoredRows) != len(originalRows) {
		t.Errorf("rows length mismatch: got %d, want %d", len(restoredRows), len(originalRows))
	}
}
