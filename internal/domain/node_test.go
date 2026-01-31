package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
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
