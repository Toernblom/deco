package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
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

// ContentValidator validates that approved/published nodes have content.
// Draft nodes are allowed to be content-free.
type ContentValidator struct{}

// NewContentValidator creates a new content validator.
func NewContentValidator() *ContentValidator {
	return &ContentValidator{}
}

// Validate checks that approved/published nodes have content with at least one section.
func (cv *ContentValidator) Validate(node *domain.Node, collector *errors.Collector) {
	if node == nil {
		return
	}

	// Only require content for approved or published nodes
	if node.Status != "approved" && node.Status != "published" {
		return
	}

	// Check if content exists and has at least one section
	if node.Content == nil || len(node.Content.Sections) == 0 {
		collector.Add(domain.DecoError{
			Code:       "E046",
			Summary:    fmt.Sprintf("Node %q with status %q requires content", node.ID, node.Status),
			Detail:     "Approved and published nodes must have content with at least one section. Draft nodes may omit content.",
			Suggestion: "Add a content section with at least one block, or change status to 'draft' while content is being developed.",
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

// knownTopLevelKeys defines the valid top-level keys in a node YAML file.
// Any key not in this set (except for explicit extension mechanisms) is considered unknown.
var knownTopLevelKeys = map[string]bool{
	"id":          true,
	"kind":        true,
	"version":     true,
	"status":      true,
	"title":       true,
	"tags":        true,
	"refs":        true,
	"content":     true,
	"issues":      true,
	"summary":     true,
	"glossary":    true,
	"contracts":   true,
	"llm_context": true,
	"constraints": true,
	"custom":      true, // Explicit extension namespace
}

// DuplicateIDValidator detects nodes with duplicate IDs.
type DuplicateIDValidator struct{}

// NewDuplicateIDValidator creates a new duplicate ID validator.
func NewDuplicateIDValidator() *DuplicateIDValidator {
	return &DuplicateIDValidator{}
}

// Validate checks for duplicate node IDs across all nodes.
func (dv *DuplicateIDValidator) Validate(nodes []domain.Node, collector *errors.Collector) {
	// Track which IDs we've seen and where
	seen := make(map[string]int) // ID -> index of first occurrence

	for i, node := range nodes {
		if node.ID == "" {
			continue // Empty IDs are handled by schema validator
		}

		if firstIdx, exists := seen[node.ID]; exists {
			collector.Add(domain.DecoError{
				Code:    "E009",
				Summary: fmt.Sprintf("Duplicate node ID: %s", node.ID),
				Detail:  fmt.Sprintf("Node ID %q appears multiple times (first at index %d, duplicate at index %d). Node IDs must be unique.", node.ID, firstIdx, i),
			})
		} else {
			seen[node.ID] = i
		}
	}
}

// UnknownFieldValidator detects unknown top-level keys in node YAML files.
// This helps catch typos and prevents silent data loss from misspelled fields.
type UnknownFieldValidator struct {
	suggester *errors.Suggester
}

// NewUnknownFieldValidator creates a new unknown field validator.
func NewUnknownFieldValidator() *UnknownFieldValidator {
	return &UnknownFieldValidator{
		suggester: errors.NewSuggester(),
	}
}

// ValidateYAML checks a raw YAML map for unknown top-level keys.
// The nodeID is used for error reporting.
func (uf *UnknownFieldValidator) ValidateYAML(nodeID string, rawKeys []string, collector *errors.Collector) {
	// Collect known keys for suggestions
	var knownKeys []string
	for k := range knownTopLevelKeys {
		knownKeys = append(knownKeys, k)
	}

	for _, key := range rawKeys {
		if !knownTopLevelKeys[key] {
			err := domain.DecoError{
				Code:    "E010",
				Summary: fmt.Sprintf("Unknown field %q in node %s", key, nodeID),
				Detail:  fmt.Sprintf("Field %q is not a recognized top-level field. Use 'custom:' for extension data.", key),
			}

			// Generate suggestion for similar field names
			suggs := uf.suggester.Suggest(key, knownKeys)
			if len(suggs) > 0 {
				err.Suggestion = fmt.Sprintf("Did you mean %q?", suggs[0])
			}

			collector.Add(err)
		}
	}
}

// ValidateDirectory reads all YAML files in a nodes directory and checks for unknown fields.
func (uf *UnknownFieldValidator) ValidateDirectory(rootDir string, collector *errors.Collector) {
	nodesDir := filepath.Join(rootDir, ".deco", "nodes")

	// Check if nodes directory exists
	if _, err := os.Stat(nodesDir); os.IsNotExist(err) {
		return
	}

	// Walk the nodes directory recursively
	_ = filepath.Walk(nodesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files with access errors
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .yaml files
		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		// Read and parse the file
		data, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip unreadable files
		}

		// Parse into a raw map to get all top-level keys
		var rawMap map[string]interface{}
		if err := yaml.Unmarshal(data, &rawMap); err != nil {
			return nil // Skip unparseable files (will be caught by other validators)
		}

		// Get node ID from the file (for error reporting)
		nodeID := ""
		if id, ok := rawMap["id"].(string); ok {
			nodeID = id
		} else {
			// Use relative path as fallback
			relPath, _ := filepath.Rel(nodesDir, path)
			nodeID = strings.TrimSuffix(relPath, ".yaml")
		}

		// Extract top-level keys
		var keys []string
		for k := range rawMap {
			keys = append(keys, k)
		}

		// Validate keys
		uf.ValidateYAML(nodeID, keys, collector)

		return nil
	})
}

// Orchestrator coordinates all validators and aggregates errors.
type Orchestrator struct {
	schemaValidator       *SchemaValidator
	contentValidator      *ContentValidator
	referenceValidator    *ReferenceValidator
	constraintValidator   *ConstraintValidator
	duplicateIDValidator  *DuplicateIDValidator
	unknownFieldValidator *UnknownFieldValidator
	contractValidator     *ContractValidator
}

// NewOrchestrator creates a new validator orchestrator.
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		schemaValidator:       NewSchemaValidator(),
		contentValidator:      NewContentValidator(),
		referenceValidator:    NewReferenceValidator(),
		constraintValidator:   NewConstraintValidator(),
		duplicateIDValidator:  NewDuplicateIDValidator(),
		unknownFieldValidator: NewUnknownFieldValidator(),
		contractValidator:     NewContractValidator(),
	}
}

// ValidateAll runs all validators on the provided nodes and returns aggregated errors.
func (o *Orchestrator) ValidateAll(nodes []domain.Node) *errors.Collector {
	collector := errors.NewCollectorWithLimit(1000)

	// Check for duplicate IDs first (critical error)
	o.duplicateIDValidator.Validate(nodes, collector)

	// Run schema validation on each node
	for _, node := range nodes {
		o.schemaValidator.Validate(&node, collector)
	}

	// Run content validation on each node (approved/published require content)
	for _, node := range nodes {
		o.contentValidator.Validate(&node, collector)
	}

	// Run reference validation on all nodes
	o.referenceValidator.Validate(nodes, collector)

	// Run constraint validation on each node
	for _, node := range nodes {
		o.constraintValidator.Validate(&node, nodes, collector)
	}

	// Run contract validation on all nodes
	o.contractValidator.ValidateAll(nodes, collector)

	return collector
}

// ValidateAllWithDir runs all validators including unknown field detection.
// The rootDir is needed to read raw YAML files for unknown field checking.
func (o *Orchestrator) ValidateAllWithDir(nodes []domain.Node, rootDir string) *errors.Collector {
	collector := o.ValidateAll(nodes)

	// Check for unknown top-level fields in YAML files
	o.unknownFieldValidator.ValidateDirectory(rootDir, collector)

	return collector
}
