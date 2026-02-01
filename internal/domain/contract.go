package domain

import (
	"regexp"
	"strings"
)

// Step represents a single Given/When/Then step in a contract scenario.
// Steps can reference nodes using the @node.id syntax.
type Step struct {
	Text       string   // The raw step text
	NodeRefs   []string // Node IDs referenced in this step (extracted from @node.id patterns)
	StepType   StepType // Given, When, or Then
	StepNumber int      // Position within its type (1-indexed)
}

// StepType indicates which phase of a scenario a step belongs to.
type StepType string

const (
	StepTypeGiven StepType = "given"
	StepTypeWhen  StepType = "when"
	StepTypeThen  StepType = "then"
)

// Scenario represents a parsed contract scenario with structured steps.
// This is the validation-ready form of the YAML Contract type.
type Scenario struct {
	Name        string // Scenario name
	Description string // Scenario description text
	Given       []Step // Preconditions
	When        []Step // Actions
	Then        []Step // Expected outcomes
}

// nodeRefPattern matches @node.id references in step text.
// Examples: @systems/core, @mechanics/combat, @events/player_death
var nodeRefPattern = regexp.MustCompile(`@([a-zA-Z0-9_/.-]+)`)

// ParseContract converts a YAML Contract into a validation-ready Scenario.
func ParseContract(c Contract) Scenario {
	return Scenario{
		Name:        c.Name,
		Description: c.Scenario,
		Given:       parseSteps(c.Given, StepTypeGiven),
		When:        parseSteps(c.When, StepTypeWhen),
		Then:        parseSteps(c.Then, StepTypeThen),
	}
}

// ParseContracts converts multiple YAML Contracts into Scenarios.
func ParseContracts(contracts []Contract) []Scenario {
	scenarios := make([]Scenario, len(contracts))
	for i, c := range contracts {
		scenarios[i] = ParseContract(c)
	}
	return scenarios
}

// parseSteps converts raw step strings into structured Step objects.
func parseSteps(texts []string, stepType StepType) []Step {
	steps := make([]Step, len(texts))
	for i, text := range texts {
		steps[i] = Step{
			Text:       text,
			NodeRefs:   extractNodeRefs(text),
			StepType:   stepType,
			StepNumber: i + 1,
		}
	}
	return steps
}

// extractNodeRefs finds all @node.id references in step text.
func extractNodeRefs(text string) []string {
	matches := nodeRefPattern.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil
	}

	refs := make([]string, len(matches))
	for i, match := range matches {
		refs[i] = match[1] // Capture group contains the node ID without @
	}
	return refs
}

// AllSteps returns all steps in the scenario in order (Given, When, Then).
func (s *Scenario) AllSteps() []Step {
	all := make([]Step, 0, len(s.Given)+len(s.When)+len(s.Then))
	all = append(all, s.Given...)
	all = append(all, s.When...)
	all = append(all, s.Then...)
	return all
}

// AllNodeRefs returns all unique node references across all steps.
func (s *Scenario) AllNodeRefs() []string {
	seen := make(map[string]bool)
	var refs []string

	for _, step := range s.AllSteps() {
		for _, ref := range step.NodeRefs {
			if !seen[ref] {
				seen[ref] = true
				refs = append(refs, ref)
			}
		}
	}
	return refs
}

// Validate checks basic structural requirements of a scenario.
// Returns an error if the scenario is invalid.
func (s *Scenario) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return &DecoError{
			Code:    "E100",
			Summary: "contract has no name",
			Detail:  "every contract must have a non-empty name",
		}
	}
	return nil
}

// ValidateStep checks that a step is well-formed.
func (step *Step) ValidateStep() error {
	if strings.TrimSpace(step.Text) == "" {
		return &DecoError{
			Code:    "E101",
			Summary: "empty step",
			Detail:  "steps cannot be empty strings",
		}
	}
	return nil
}
