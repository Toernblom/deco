package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestConstraint_Creation(t *testing.T) {
	constraint := domain.Constraint{
		Expr:    "self.needs.food.threshold == refs['systems/food'].starvation_time",
		Message: "Food threshold must match starvation time in food system",
		Scope:   "mechanic",
	}

	if constraint.Expr != "self.needs.food.threshold == refs['systems/food'].starvation_time" {
		t.Errorf("expected Expr to match, got %q", constraint.Expr)
	}
	if constraint.Message != "Food threshold must match starvation time in food system" {
		t.Errorf("expected Message to match, got %q", constraint.Message)
	}
	if constraint.Scope != "mechanic" {
		t.Errorf("expected Scope 'mechanic', got %q", constraint.Scope)
	}
}

func TestConstraint_ExpressionStorage(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{
			name: "simple comparison",
			expr: "self.value > 0",
		},
		{
			name: "reference comparison",
			expr: "self.max_health == refs['systems/health'].default_max",
		},
		{
			name: "complex expression",
			expr: "len(self.items) <= refs['systems/inventory'].max_slots && all(self.items, item.weight <= 100)",
		},
		{
			name: "nested field access",
			expr: "self.content.sections[0].type == 'rule'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := domain.Constraint{
				Expr:    tt.expr,
				Message: "Test constraint",
				Scope:   "all",
			}

			if constraint.Expr != tt.expr {
				t.Errorf("expected Expr %q, got %q", tt.expr, constraint.Expr)
			}
		})
	}
}

func TestConstraint_MessageDefinition(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "simple message",
			message: "Value must be positive",
		},
		{
			name:    "detailed message",
			message: "The maximum health value must match the default_max defined in the health system configuration",
		},
		{
			name:    "message with field names",
			message: "Field 'content.damage' must be between 1 and 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := domain.Constraint{
				Expr:    "self.value > 0",
				Message: tt.message,
				Scope:   "all",
			}

			if constraint.Message != tt.message {
				t.Errorf("expected Message %q, got %q", tt.message, constraint.Message)
			}
		})
	}
}

func TestConstraint_ScopeDefinition(t *testing.T) {
	tests := []struct {
		name  string
		scope string
	}{
		{
			name:  "all nodes",
			scope: "all",
		},
		{
			name:  "mechanic nodes only",
			scope: "mechanic",
		},
		{
			name:  "system nodes only",
			scope: "system",
		},
		{
			name:  "specific node path",
			scope: "systems/combat/*",
		},
		{
			name:  "feature nodes",
			scope: "feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := domain.Constraint{
				Expr:    "self.value > 0",
				Message: "Test constraint",
				Scope:   tt.scope,
			}

			if constraint.Scope != tt.scope {
				t.Errorf("expected Scope %q, got %q", tt.scope, constraint.Scope)
			}
		})
	}
}

func TestConstraint_Serialization(t *testing.T) {
	original := domain.Constraint{
		Expr:    "self.value >= 0 && self.value <= 100",
		Message: "Value must be between 0 and 100",
		Scope:   "mechanic",
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Constraint: %v", err)
	}

	// Unmarshal back
	var restored domain.Constraint
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal Constraint: %v", err)
	}

	// Compare
	if restored.Expr != original.Expr {
		t.Errorf("Expr mismatch: got %q, want %q", restored.Expr, original.Expr)
	}
	if restored.Message != original.Message {
		t.Errorf("Message mismatch: got %q, want %q", restored.Message, original.Message)
	}
	if restored.Scope != original.Scope {
		t.Errorf("Scope mismatch: got %q, want %q", restored.Scope, original.Scope)
	}
}

func TestConstraint_Validation(t *testing.T) {
	tests := []struct {
		name       string
		constraint domain.Constraint
		wantErr    bool
	}{
		{
			name: "valid constraint",
			constraint: domain.Constraint{
				Expr:    "self.value > 0",
				Message: "Value must be positive",
				Scope:   "all",
			},
			wantErr: false,
		},
		{
			name: "missing expression",
			constraint: domain.Constraint{
				Message: "Value must be positive",
				Scope:   "all",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			constraint: domain.Constraint{
				Expr:  "self.value > 0",
				Scope: "all",
			},
			wantErr: true,
		},
		{
			name: "missing scope",
			constraint: domain.Constraint{
				Expr:    "self.value > 0",
				Message: "Value must be positive",
			},
			wantErr: true,
		},
		{
			name: "empty expression",
			constraint: domain.Constraint{
				Expr:    "",
				Message: "Value must be positive",
				Scope:   "all",
			},
			wantErr: true,
		},
		{
			name: "empty message",
			constraint: domain.Constraint{
				Expr:    "self.value > 0",
				Message: "",
				Scope:   "all",
			},
			wantErr: true,
		},
		{
			name: "empty scope",
			constraint: domain.Constraint{
				Expr:    "self.value > 0",
				Message: "Value must be positive",
				Scope:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constraint.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Constraint.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
