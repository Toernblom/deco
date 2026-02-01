package validator

import (
	"fmt"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
)

// ContractValidator validates contract scenarios within nodes.
// It checks for proper given/when/then structure, valid field types,
// and unique scenario names within each node.
type ContractValidator struct{}

// NewContractValidator creates a new contract validator.
func NewContractValidator() *ContractValidator {
	return &ContractValidator{}
}

// Validate checks all contracts in a node for syntax errors.
// Validations include:
// - Scenario names are non-empty (E100)
// - Scenario names are unique within the node (E103)
// - Steps are non-empty (E101)
// - At least one step exists in given, when, or then (E104)
func (cv *ContractValidator) Validate(node *domain.Node, collector *errors.Collector) {
	if node == nil || len(node.Contracts) == 0 {
		return
	}

	// Track scenario names for uniqueness check
	seenNames := make(map[string]int) // name -> first occurrence index

	for i, contract := range node.Contracts {
		scenario := domain.ParseContract(contract)

		// Validate scenario name exists
		if err := scenario.Validate(); err != nil {
			if decoErr, ok := err.(*domain.DecoError); ok {
				decoErr.Detail = fmt.Sprintf("in node %s, contract at index %d: %s", node.ID, i, decoErr.Detail)
				collector.Add(*decoErr)
			}
			continue
		}

		// Check for duplicate scenario names within this node
		name := strings.TrimSpace(scenario.Name)
		if firstIdx, exists := seenNames[name]; exists {
			collector.Add(domain.DecoError{
				Code:    "E103",
				Summary: fmt.Sprintf("Duplicate contract name: %s", name),
				Detail:  fmt.Sprintf("in node %s: contract name %q appears at index %d and %d; names must be unique within a node", node.ID, name, firstIdx, i),
			})
		} else {
			seenNames[name] = i
		}

		// Validate all steps
		cv.validateSteps(node.ID, &scenario, collector)

		// Validate structure: at least one step should exist
		cv.validateStructure(node.ID, &scenario, collector)
	}
}

// validateSteps checks that all steps are non-empty.
func (cv *ContractValidator) validateSteps(nodeID string, scenario *domain.Scenario, collector *errors.Collector) {
	for _, step := range scenario.AllSteps() {
		if err := step.ValidateStep(); err != nil {
			if decoErr, ok := err.(*domain.DecoError); ok {
				decoErr.Detail = fmt.Sprintf("in node %s, contract %q, %s step %d: %s",
					nodeID, scenario.Name, step.StepType, step.StepNumber, decoErr.Detail)
				collector.Add(*decoErr)
			}
		}
	}
}

// validateStructure checks that a contract has proper given/when/then structure.
// A valid contract must have at least one step in any section.
func (cv *ContractValidator) validateStructure(nodeID string, scenario *domain.Scenario, collector *errors.Collector) {
	totalSteps := len(scenario.Given) + len(scenario.When) + len(scenario.Then)

	if totalSteps == 0 {
		collector.Add(domain.DecoError{
			Code:    "E104",
			Summary: "Contract has no steps",
			Detail:  fmt.Sprintf("in node %s, contract %q: must have at least one given, when, or then step", nodeID, scenario.Name),
		})
	}
}

// ValidateAll checks all contracts across multiple nodes.
func (cv *ContractValidator) ValidateAll(nodes []domain.Node, collector *errors.Collector) {
	for _, node := range nodes {
		cv.Validate(&node, collector)
	}
}
