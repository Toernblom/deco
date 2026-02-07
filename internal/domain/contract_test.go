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

package domain

import (
	"testing"
)

func TestParseContract(t *testing.T) {
	contract := Contract{
		Name:     "Eating food grows snake",
		Scenario: "When snake head collides with food, snake grows by one segment",
		Given:    []string{"snake has length 3", "food exists at position (6,5)"},
		When:     []string{"snake moves right"},
		Then:     []string{"snake length becomes 4", "score increases"},
	}

	scenario := ParseContract(contract)

	if scenario.Name != contract.Name {
		t.Errorf("Name = %q, want %q", scenario.Name, contract.Name)
	}
	if scenario.Description != contract.Scenario {
		t.Errorf("Description = %q, want %q", scenario.Description, contract.Scenario)
	}
	if len(scenario.Given) != 2 {
		t.Errorf("len(Given) = %d, want 2", len(scenario.Given))
	}
	if len(scenario.When) != 1 {
		t.Errorf("len(When) = %d, want 1", len(scenario.When))
	}
	if len(scenario.Then) != 2 {
		t.Errorf("len(Then) = %d, want 2", len(scenario.Then))
	}
}

func TestParseContractStepTypes(t *testing.T) {
	contract := Contract{
		Name:  "Test",
		Given: []string{"precondition"},
		When:  []string{"action"},
		Then:  []string{"outcome"},
	}

	scenario := ParseContract(contract)

	if scenario.Given[0].StepType != StepTypeGiven {
		t.Errorf("Given step type = %q, want %q", scenario.Given[0].StepType, StepTypeGiven)
	}
	if scenario.When[0].StepType != StepTypeWhen {
		t.Errorf("When step type = %q, want %q", scenario.When[0].StepType, StepTypeWhen)
	}
	if scenario.Then[0].StepType != StepTypeThen {
		t.Errorf("Then step type = %q, want %q", scenario.Then[0].StepType, StepTypeThen)
	}
}

func TestParseContractStepNumbers(t *testing.T) {
	contract := Contract{
		Name:  "Test",
		Given: []string{"first", "second", "third"},
	}

	scenario := ParseContract(contract)

	for i, step := range scenario.Given {
		expected := i + 1
		if step.StepNumber != expected {
			t.Errorf("Given[%d].StepNumber = %d, want %d", i, step.StepNumber, expected)
		}
	}
}

func TestExtractNodeRefs(t *testing.T) {
	tests := []struct {
		text     string
		expected []string
	}{
		{"no refs here", nil},
		{"@systems/core is referenced", []string{"systems/core"}},
		{"@a and @b", []string{"a", "b"}},
		{"@systems/settlement/colonists", []string{"systems/settlement/colonists"}},
		{"@events/player_death emits", []string{"events/player_death"}},
		{"text @ref1 middle @ref2 end", []string{"ref1", "ref2"}},
	}

	for _, tt := range tests {
		refs := extractNodeRefs(tt.text)
		if len(refs) != len(tt.expected) {
			t.Errorf("extractNodeRefs(%q) = %v, want %v", tt.text, refs, tt.expected)
			continue
		}
		for i, ref := range refs {
			if ref != tt.expected[i] {
				t.Errorf("extractNodeRefs(%q)[%d] = %q, want %q", tt.text, i, ref, tt.expected[i])
			}
		}
	}
}

func TestScenarioAllSteps(t *testing.T) {
	scenario := Scenario{
		Given: []Step{{Text: "g1"}, {Text: "g2"}},
		When:  []Step{{Text: "w1"}},
		Then:  []Step{{Text: "t1"}, {Text: "t2"}},
	}

	all := scenario.AllSteps()

	if len(all) != 5 {
		t.Errorf("len(AllSteps()) = %d, want 5", len(all))
	}

	expected := []string{"g1", "g2", "w1", "t1", "t2"}
	for i, step := range all {
		if step.Text != expected[i] {
			t.Errorf("AllSteps()[%d].Text = %q, want %q", i, step.Text, expected[i])
		}
	}
}

func TestScenarioAllNodeRefs(t *testing.T) {
	scenario := Scenario{
		Given: []Step{{NodeRefs: []string{"a", "b"}}},
		When:  []Step{{NodeRefs: []string{"b", "c"}}}, // b is duplicate
		Then:  []Step{{NodeRefs: []string{"d"}}},
	}

	refs := scenario.AllNodeRefs()

	// Should deduplicate
	if len(refs) != 4 {
		t.Errorf("len(AllNodeRefs()) = %d, want 4", len(refs))
	}

	seen := make(map[string]bool)
	for _, ref := range refs {
		if seen[ref] {
			t.Errorf("AllNodeRefs() contains duplicate: %q", ref)
		}
		seen[ref] = true
	}
}

func TestScenarioValidate(t *testing.T) {
	tests := []struct {
		name     string
		scenario Scenario
		wantErr  bool
	}{
		{"valid", Scenario{Name: "Test"}, false},
		{"empty name", Scenario{Name: ""}, true},
		{"whitespace name", Scenario{Name: "   "}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.scenario.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStepValidate(t *testing.T) {
	tests := []struct {
		name    string
		step    Step
		wantErr bool
	}{
		{"valid", Step{Text: "some step"}, false},
		{"empty", Step{Text: ""}, true},
		{"whitespace", Step{Text: "   "}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.step.ValidateStep()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStep() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseContracts(t *testing.T) {
	contracts := []Contract{
		{Name: "First", Given: []string{"a"}},
		{Name: "Second", When: []string{"b"}},
	}

	scenarios := ParseContracts(contracts)

	if len(scenarios) != 2 {
		t.Errorf("len(ParseContracts) = %d, want 2", len(scenarios))
	}
	if scenarios[0].Name != "First" {
		t.Errorf("scenarios[0].Name = %q, want %q", scenarios[0].Name, "First")
	}
	if scenarios[1].Name != "Second" {
		t.Errorf("scenarios[1].Name = %q, want %q", scenarios[1].Name, "Second")
	}
}
