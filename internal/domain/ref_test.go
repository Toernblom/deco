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

package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestRef_Creation(t *testing.T) {
	ref := domain.Ref{
		Uses:        []domain.RefLink{{Target: "systems/food"}},
		Related:     []domain.RefLink{{Target: "systems/health"}},
		EmitsEvents: []string{"colonist.starving", "colonist.fed"},
		Vocabulary:  []string{"hunger", "satiety"},
	}

	if len(ref.Uses) != 1 || ref.Uses[0].Target != "systems/food" {
		t.Errorf("expected Uses to contain systems/food")
	}
	if len(ref.Related) != 1 || ref.Related[0].Target != "systems/health" {
		t.Errorf("expected Related to contain systems/health")
	}
	if len(ref.EmitsEvents) != 2 {
		t.Errorf("expected 2 emitted events, got %d", len(ref.EmitsEvents))
	}
	if len(ref.Vocabulary) != 2 {
		t.Errorf("expected 2 vocabulary terms, got %d", len(ref.Vocabulary))
	}
}

func TestRefLink_Creation(t *testing.T) {
	link := domain.RefLink{
		Target:  "systems/combat/damage",
		Context: "for calculating combat damage",
	}

	if link.Target != "systems/combat/damage" {
		t.Errorf("expected Target 'systems/combat/damage', got %q", link.Target)
	}
	if link.Context != "for calculating combat damage" {
		t.Errorf("expected Context 'for calculating combat damage', got %q", link.Context)
	}
}

func TestRefLink_ResolutionState(t *testing.T) {
	link := domain.RefLink{
		Target:   "systems/food",
		Resolved: false,
	}

	if link.Resolved {
		t.Errorf("expected new RefLink to be unresolved")
	}

	// Mark as resolved
	link.Resolved = true

	if !link.Resolved {
		t.Errorf("expected RefLink to be resolved after setting Resolved=true")
	}
}

func TestRef_UsesType(t *testing.T) {
	ref := domain.Ref{
		Uses: []domain.RefLink{
			{Target: "systems/food", Context: "food consumption"},
			{Target: "systems/water", Context: "water consumption"},
			{Target: "systems/shelter", Context: "shelter requirement"},
		},
	}

	if len(ref.Uses) != 3 {
		t.Errorf("expected 3 uses references, got %d", len(ref.Uses))
	}

	expectedTargets := map[string]bool{
		"systems/food":    true,
		"systems/water":   true,
		"systems/shelter": true,
	}

	for _, link := range ref.Uses {
		if !expectedTargets[link.Target] {
			t.Errorf("unexpected target in Uses: %q", link.Target)
		}
	}
}

func TestRef_RelatedType(t *testing.T) {
	ref := domain.Ref{
		Related: []domain.RefLink{
			{Target: "features/survival", Context: "related survival feature"},
			{Target: "mechanics/needs", Context: "needs mechanic"},
		},
	}

	if len(ref.Related) != 2 {
		t.Errorf("expected 2 related references, got %d", len(ref.Related))
	}
}

func TestRef_EmitsEvents(t *testing.T) {
	ref := domain.Ref{
		EmitsEvents: []string{
			"colonist.health.changed",
			"colonist.died",
			"colonist.revived",
		},
	}

	if len(ref.EmitsEvents) != 3 {
		t.Errorf("expected 3 emitted events, got %d", len(ref.EmitsEvents))
	}

	expectedEvents := map[string]bool{
		"colonist.health.changed": true,
		"colonist.died":           true,
		"colonist.revived":        true,
	}

	for _, event := range ref.EmitsEvents {
		if !expectedEvents[event] {
			t.Errorf("unexpected event: %q", event)
		}
	}
}

func TestRef_Vocabulary(t *testing.T) {
	ref := domain.Ref{
		Vocabulary: []string{
			"hunger",
			"satiety",
			"starvation",
			"nourishment",
		},
	}

	if len(ref.Vocabulary) != 4 {
		t.Errorf("expected 4 vocabulary terms, got %d", len(ref.Vocabulary))
	}

	expectedTerms := map[string]bool{
		"hunger":      true,
		"satiety":     true,
		"starvation":  true,
		"nourishment": true,
	}

	for _, term := range ref.Vocabulary {
		if !expectedTerms[term] {
			t.Errorf("unexpected vocabulary term: %q", term)
		}
	}
}

func TestRef_Serialization(t *testing.T) {
	original := domain.Ref{
		Uses: []domain.RefLink{
			{Target: "systems/food", Context: "food system"},
		},
		Related: []domain.RefLink{
			{Target: "systems/health", Context: "health system"},
		},
		EmitsEvents: []string{"colonist.starving"},
		Vocabulary:  []string{"hunger"},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Ref: %v", err)
	}

	// Unmarshal back
	var restored domain.Ref
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal Ref: %v", err)
	}

	// Compare
	if len(restored.Uses) != len(original.Uses) {
		t.Errorf("Uses length mismatch: got %d, want %d", len(restored.Uses), len(original.Uses))
	}
	if len(restored.Related) != len(original.Related) {
		t.Errorf("Related length mismatch: got %d, want %d", len(restored.Related), len(original.Related))
	}
	if len(restored.EmitsEvents) != len(original.EmitsEvents) {
		t.Errorf("EmitsEvents length mismatch: got %d, want %d", len(restored.EmitsEvents), len(original.EmitsEvents))
	}
	if len(restored.Vocabulary) != len(original.Vocabulary) {
		t.Errorf("Vocabulary length mismatch: got %d, want %d", len(restored.Vocabulary), len(original.Vocabulary))
	}
}

func TestRefLink_Serialization(t *testing.T) {
	original := domain.RefLink{
		Target:   "systems/combat",
		Context:  "damage calculation",
		Resolved: true,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal RefLink: %v", err)
	}

	// Unmarshal back
	var restored domain.RefLink
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal RefLink: %v", err)
	}

	// Compare
	if restored.Target != original.Target {
		t.Errorf("Target mismatch: got %q, want %q", restored.Target, original.Target)
	}
	if restored.Context != original.Context {
		t.Errorf("Context mismatch: got %q, want %q", restored.Context, original.Context)
	}
	if restored.Resolved != original.Resolved {
		t.Errorf("Resolved mismatch: got %v, want %v", restored.Resolved, original.Resolved)
	}
}

func TestRefLink_Validation(t *testing.T) {
	tests := []struct {
		name    string
		link    domain.RefLink
		wantErr bool
	}{
		{
			name: "valid reflink",
			link: domain.RefLink{
				Target: "systems/food",
			},
			wantErr: false,
		},
		{
			name: "valid reflink with context",
			link: domain.RefLink{
				Target:  "systems/food",
				Context: "food consumption",
			},
			wantErr: false,
		},
		{
			name: "missing target",
			link: domain.RefLink{
				Context: "some context",
			},
			wantErr: true,
		},
		{
			name: "empty target",
			link: domain.RefLink{
				Target: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.link.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RefLink.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
