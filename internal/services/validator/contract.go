package validator

import (
	"fmt"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
)

// ContractValidator validates contract scenarios within nodes.
// It checks for proper given/when/then structure, valid field types,
// unique scenario names within each node, and valid node references.
type ContractValidator struct {
	suggester *errors.Suggester
}

// NewContractValidator creates a new contract validator.
func NewContractValidator() *ContractValidator {
	return &ContractValidator{
		suggester: errors.NewSuggester(),
	}
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

	// Helper to create location from node source file
	var location *domain.Location
	if node.SourceFile != "" {
		location = &domain.Location{File: node.SourceFile}
	}

	// Track scenario names for uniqueness check
	seenNames := make(map[string]int) // name -> first occurrence index

	for i, contract := range node.Contracts {
		scenario := domain.ParseContract(contract)

		// Validate scenario name exists
		if err := scenario.Validate(); err != nil {
			if decoErr, ok := err.(*domain.DecoError); ok {
				decoErr.Detail = fmt.Sprintf("in node %s, contract at index %d: %s", node.ID, i, decoErr.Detail)
				decoErr.Location = location
				collector.Add(*decoErr)
			}
			continue
		}

		// Check for duplicate scenario names within this node
		name := strings.TrimSpace(scenario.Name)
		if firstIdx, exists := seenNames[name]; exists {
			collector.Add(domain.DecoError{
				Code:     "E103",
				Summary:  fmt.Sprintf("Duplicate contract name: %s", name),
				Detail:   fmt.Sprintf("in node %s: contract name %q appears at index %d and %d; names must be unique within a node", node.ID, name, firstIdx, i),
				Location: location,
			})
		} else {
			seenNames[name] = i
		}

		// Validate all steps
		cv.validateSteps(node.ID, &scenario, location, collector)

		// Validate structure: at least one step should exist
		cv.validateStructure(node.ID, &scenario, location, collector)
	}
}

// validateSteps checks that all steps are non-empty.
func (cv *ContractValidator) validateSteps(nodeID string, scenario *domain.Scenario, location *domain.Location, collector *errors.Collector) {
	for _, step := range scenario.AllSteps() {
		if err := step.ValidateStep(); err != nil {
			if decoErr, ok := err.(*domain.DecoError); ok {
				decoErr.Detail = fmt.Sprintf("in node %s, contract %q, %s step %d: %s",
					nodeID, scenario.Name, step.StepType, step.StepNumber, decoErr.Detail)
				decoErr.Location = location
				collector.Add(*decoErr)
			}
		}
	}
}

// validateStructure checks that a contract has proper given/when/then structure.
// A valid contract must have at least one step in any section.
func (cv *ContractValidator) validateStructure(nodeID string, scenario *domain.Scenario, location *domain.Location, collector *errors.Collector) {
	totalSteps := len(scenario.Given) + len(scenario.When) + len(scenario.Then)

	if totalSteps == 0 {
		collector.Add(domain.DecoError{
			Code:     "E104",
			Summary:  "Contract has no steps",
			Detail:   fmt.Sprintf("in node %s, contract %q: must have at least one given, when, or then step", nodeID, scenario.Name),
			Location: location,
		})
	}
}

// ValidateAll checks all contracts across multiple nodes.
// This includes syntax validation and node reference validation.
func (cv *ContractValidator) ValidateAll(nodes []domain.Node, collector *errors.Collector) {
	// Build set of existing node IDs for reference validation
	nodeIDs := make(map[string]bool)
	var allIDs []string
	for _, node := range nodes {
		nodeIDs[node.ID] = true
		allIDs = append(allIDs, node.ID)
	}

	// Validate each node's contracts
	for _, node := range nodes {
		cv.Validate(&node, collector)
		cv.validateNodeRefs(&node, nodeIDs, allIDs, collector)
	}
}

// validateNodeRefs checks that all @node.id references in contract steps
// point to existing nodes in the design graph.
func (cv *ContractValidator) validateNodeRefs(node *domain.Node, nodeIDs map[string]bool, allIDs []string, collector *errors.Collector) {
	if node == nil || len(node.Contracts) == 0 {
		return
	}

	// Helper to create location from node source file
	var location *domain.Location
	if node.SourceFile != "" {
		location = &domain.Location{File: node.SourceFile}
	}

	for _, contract := range node.Contracts {
		scenario := domain.ParseContract(contract)

		// Check each step's node references
		for _, step := range scenario.AllSteps() {
			for _, ref := range step.NodeRefs {
				if !nodeIDs[ref] {
					err := domain.DecoError{
						Code:     "E102",
						Summary:  fmt.Sprintf("Invalid node reference: @%s", ref),
						Detail:   fmt.Sprintf("in node %s, contract %q, %s step %d: referenced node %q does not exist", node.ID, scenario.Name, step.StepType, step.StepNumber, ref),
						Location: location,
					}

					// Generate suggestion for similar IDs
					suggs := cv.suggester.Suggest(ref, allIDs)
					if len(suggs) > 0 {
						err.Suggestion = fmt.Sprintf("Did you mean '@%s'?", suggs[0])
					}

					collector.Add(err)
				}
			}
		}
	}
}
