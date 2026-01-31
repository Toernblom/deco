package validator

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
)

// SchemaValidator validates that nodes have all required fields and correct types.
type SchemaValidator struct{}

// NewSchemaValidator creates a new schema validator.
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{}
}

// Validate checks that a node has all required fields.
func (sv *SchemaValidator) Validate(node *domain.Node, collector *errors.Collector) {
	if node == nil {
		collector.Add(domain.DecoError{
			Code:    "E008",
			Summary: "Node is nil",
			Detail:  "Cannot validate a nil node",
		})
		return
	}

	// Check required fields
	if node.ID == "" {
		collector.Add(domain.DecoError{
			Code:    "E008",
			Summary: "Missing required field: ID",
			Detail:  "Node ID is required",
		})
	}

	if node.Kind == "" {
		collector.Add(domain.DecoError{
			Code:    "E008",
			Summary: "Missing required field: Kind",
			Detail:  "Node Kind is required",
		})
	}

	if node.Version == 0 {
		collector.Add(domain.DecoError{
			Code:    "E008",
			Summary: "Missing required field: Version",
			Detail:  "Node Version is required and must be > 0",
		})
	}

	if node.Status == "" {
		collector.Add(domain.DecoError{
			Code:    "E008",
			Summary: "Missing required field: Status",
			Detail:  "Node Status is required",
		})
	}

	if node.Title == "" {
		collector.Add(domain.DecoError{
			Code:    "E008",
			Summary: "Missing required field: Title",
			Detail:  "Node Title is required",
		})
	}
}

// ReferenceValidator validates that all node references resolve to existing nodes.
type ReferenceValidator struct {
	suggester *errors.Suggester
}

// NewReferenceValidator creates a new reference validator.
func NewReferenceValidator() *ReferenceValidator {
	return &ReferenceValidator{
		suggester: errors.NewSuggester(),
	}
}

// Validate checks that all references in nodes resolve correctly.
// Generates suggestions for broken references that look like typos.
func (rv *ReferenceValidator) Validate(nodes []domain.Node, collector *errors.Collector) {
	// Build set of existing node IDs
	nodeIDs := make(map[string]bool)
	for _, node := range nodes {
		nodeIDs[node.ID] = true
	}

	// Collect all IDs for suggestion generation
	var allIDs []string
	for _, node := range nodes {
		allIDs = append(allIDs, node.ID)
	}

	// Check each node's references
	for _, node := range nodes {
		// Check Uses references
		for _, refLink := range node.Refs.Uses {
			if !nodeIDs[refLink.Target] {
				err := domain.DecoError{
					Code:    "E020",
					Summary: "Reference not found: " + refLink.Target,
					Detail:  "Referenced node '" + refLink.Target + "' does not exist",
				}

				// Generate suggestion for similar IDs
				suggs := rv.suggester.Suggest(refLink.Target, allIDs)
				if len(suggs) > 0 {
					err.Suggestion = "Did you mean '" + suggs[0] + "'?"
				}

				collector.Add(err)
			}
		}

		// Check Related references
		for _, refLink := range node.Refs.Related {
			if !nodeIDs[refLink.Target] {
				err := domain.DecoError{
					Code:    "E020",
					Summary: "Reference not found: " + refLink.Target,
					Detail:  "Referenced node '" + refLink.Target + "' does not exist",
				}

				// Generate suggestion for similar IDs
				suggs := rv.suggester.Suggest(refLink.Target, allIDs)
				if len(suggs) > 0 {
					err.Suggestion = "Did you mean '" + suggs[0] + "'?"
				}

				collector.Add(err)
			}
		}
	}
}

// ConstraintValidator validates CEL expression constraints on nodes.
type ConstraintValidator struct{}

// NewConstraintValidator creates a new constraint validator.
func NewConstraintValidator() *ConstraintValidator {
	return &ConstraintValidator{}
}

// Validate evaluates all constraints on a node.
// The allNodes parameter is provided for cross-node constraints.
func (cv *ConstraintValidator) Validate(node *domain.Node, allNodes []domain.Node, collector *errors.Collector) {
	if node == nil {
		return
	}

	// Evaluate each constraint
	for _, constraint := range node.Constraints {
		if err := cv.evaluateConstraint(node, constraint, collector); err != nil {
			// If there's an error parsing or evaluating the CEL expression,
			// add it as an E042 error (CEL expression error)
			collector.Add(domain.DecoError{
				Code:    "E042",
				Summary: "CEL expression error: " + constraint.Expr,
				Detail:  err.Error(),
			})
		}
	}
}

// evaluateConstraint evaluates a single constraint using CEL
func (cv *ConstraintValidator) evaluateConstraint(node *domain.Node, constraint domain.Constraint, collector *errors.Collector) error {
	// Create CEL environment
	env, err := cel.NewEnv(
		cel.Variable("id", cel.StringType),
		cel.Variable("kind", cel.StringType),
		cel.Variable("version", cel.IntType),
		cel.Variable("status", cel.StringType),
		cel.Variable("title", cel.StringType),
		cel.Variable("tags", cel.ListType(cel.StringType)),
	)
	if err != nil {
		return fmt.Errorf("failed to create CEL environment: %w", err)
	}

	// Parse the expression
	ast, issues := env.Compile(constraint.Expr)
	if issues != nil && issues.Err() != nil {
		return fmt.Errorf("failed to compile CEL expression: %w", issues.Err())
	}

	// Create program
	prg, err := env.Program(ast)
	if err != nil {
		return fmt.Errorf("failed to create CEL program: %w", err)
	}

	// Prepare input data
	inputData := map[string]interface{}{
		"id":      node.ID,
		"kind":    node.Kind,
		"version": int64(node.Version),
		"status":  node.Status,
		"title":   node.Title,
		"tags":    node.Tags,
	}

	// Evaluate the expression
	result, _, err := prg.Eval(inputData)
	if err != nil {
		return fmt.Errorf("failed to evaluate CEL expression: %w", err)
	}

	// Check if the result is a boolean and false
	if boolResult, ok := result.Value().(bool); ok {
		if !boolResult {
			// Constraint violated
			collector.Add(domain.DecoError{
				Code:    "E041",
				Summary: "Constraint violation: " + constraint.Expr,
				Detail:  constraint.Message,
			})
		}
	} else {
		return fmt.Errorf("CEL expression did not evaluate to a boolean")
	}

	return nil
}

// Orchestrator coordinates all validators and aggregates errors.
type Orchestrator struct {
	schemaValidator     *SchemaValidator
	referenceValidator  *ReferenceValidator
	constraintValidator *ConstraintValidator
}

// NewOrchestrator creates a new validator orchestrator.
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		schemaValidator:     NewSchemaValidator(),
		referenceValidator:  NewReferenceValidator(),
		constraintValidator: NewConstraintValidator(),
	}
}

// ValidateAll runs all validators on the provided nodes and returns aggregated errors.
func (o *Orchestrator) ValidateAll(nodes []domain.Node) *errors.Collector {
	collector := errors.NewCollectorWithLimit(1000)

	// Run schema validation on each node
	for _, node := range nodes {
		o.schemaValidator.Validate(&node, collector)
	}

	// Run reference validation on all nodes
	o.referenceValidator.Validate(nodes, collector)

	// Run constraint validation on each node
	for _, node := range nodes {
		o.constraintValidator.Validate(&node, nodes, collector)
	}

	return collector
}
