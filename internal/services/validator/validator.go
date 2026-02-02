package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

// validStatuses defines the allowed node status values.
var validStatuses = map[string]bool{
	"draft":      true,
	"review":     true,
	"approved":   true,
	"deprecated": true,
	"archived":   true,
}

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

	// Helper to create location from node source file
	nodeLocation := func() *domain.Location {
		if node.SourceFile == "" {
			return nil
		}
		return &domain.Location{File: node.SourceFile}
	}

	// Check required fields
	if node.ID == "" {
		collector.Add(domain.DecoError{
			Code:     "E008",
			Summary:  "Missing required field: ID",
			Detail:   "Node ID is required",
			Location: nodeLocation(),
		})
	}

	if node.Kind == "" {
		collector.Add(domain.DecoError{
			Code:     "E008",
			Summary:  "Missing required field: Kind",
			Detail:   "Node Kind is required",
			Location: nodeLocation(),
		})
	}

	if node.Version == 0 {
		collector.Add(domain.DecoError{
			Code:     "E008",
			Summary:  "Missing required field: Version",
			Detail:   "Node Version is required and must be > 0",
			Location: nodeLocation(),
		})
	}

	if node.Status == "" {
		collector.Add(domain.DecoError{
			Code:     "E008",
			Summary:  "Missing required field: Status",
			Detail:   "Node Status is required",
			Location: nodeLocation(),
		})
	} else if !validStatuses[node.Status] {
		collector.Add(domain.DecoError{
			Code:       "E011",
			Summary:    fmt.Sprintf("Invalid status: %q", node.Status),
			Detail:     fmt.Sprintf("Status must be one of: draft, review, approved, deprecated, archived. Got: %q", node.Status),
			Suggestion: "Change status to a valid value: draft, review, approved, deprecated, or archived",
			Location:   nodeLocation(),
		})
	}

	if node.Title == "" {
		collector.Add(domain.DecoError{
			Code:     "E008",
			Summary:  "Missing required field: Title",
			Detail:   "Node Title is required",
			Location: nodeLocation(),
		})
	}
}

// SchemaRulesValidator validates nodes against per-kind schema rules defined in config.
type SchemaRulesValidator struct {
	rules map[string]config.SchemaRuleConfig
}

// NewSchemaRulesValidator creates a validator with the given per-kind rules.
func NewSchemaRulesValidator(rules map[string]config.SchemaRuleConfig) *SchemaRulesValidator {
	return &SchemaRulesValidator{rules: rules}
}

// Validate checks that a node has all required fields for its kind.
func (srv *SchemaRulesValidator) Validate(node *domain.Node, collector *errors.Collector) {
	if node == nil || srv.rules == nil {
		return
	}

	rule, exists := srv.rules[node.Kind]
	if !exists {
		return
	}

	// Helper to create location from node source file
	var location *domain.Location
	if node.SourceFile != "" {
		location = &domain.Location{File: node.SourceFile}
	}

	// Check required fields in node.Custom
	for _, field := range rule.RequiredFields {
		if node.Custom == nil {
			collector.Add(domain.DecoError{
				Code:       "E051",
				Summary:    fmt.Sprintf("Node %q (kind=%s) missing required field: %s", node.ID, node.Kind, field),
				Detail:     fmt.Sprintf("Schema rules require %q nodes to have field %q in custom data", node.Kind, field),
				Suggestion: fmt.Sprintf("Add 'custom: { %s: ... }' to the node", field),
				Location:   location,
			})
			continue
		}

		if _, ok := node.Custom[field]; !ok {
			collector.Add(domain.DecoError{
				Code:       "E051",
				Summary:    fmt.Sprintf("Node %q (kind=%s) missing required field: %s", node.ID, node.Kind, field),
				Detail:     fmt.Sprintf("Schema rules require %q nodes to have field %q in custom data", node.Kind, field),
				Suggestion: fmt.Sprintf("Add '%s: ...' to the node's custom section", field),
				Location:   location,
			})
		}
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

	// Helper to create location from node source file
	var location *domain.Location
	if node.SourceFile != "" {
		location = &domain.Location{File: node.SourceFile}
	}

	// Check if content exists and has at least one section
	if node.Content == nil || len(node.Content.Sections) == 0 {
		collector.Add(domain.DecoError{
			Code:       "E046",
			Summary:    fmt.Sprintf("Node %q with status %q requires content", node.ID, node.Status),
			Detail:     "Approved and published nodes must have content with at least one section. Draft nodes may omit content.",
			Suggestion: "Add a content section with at least one block, or change status to 'draft' while content is being developed.",
			Location:   location,
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
		// Helper to create location from node source file
		var location *domain.Location
		if node.SourceFile != "" {
			location = &domain.Location{File: node.SourceFile}
		}

		// Check Uses references
		for _, refLink := range node.Refs.Uses {
			if !nodeIDs[refLink.Target] {
				err := domain.DecoError{
					Code:     "E020",
					Summary:  "Reference not found: " + refLink.Target,
					Detail:   "Referenced node '" + refLink.Target + "' does not exist",
					Location: location,
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
					Code:     "E020",
					Summary:  "Reference not found: " + refLink.Target,
					Detail:   "Referenced node '" + refLink.Target + "' does not exist",
					Location: location,
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
type ConstraintValidator struct {
	env      *cel.Env
	programs map[string]cel.Program // cache compiled programs by expression
}

// NewConstraintValidator creates a new constraint validator.
func NewConstraintValidator() *ConstraintValidator {
	// Create a shared CEL environment - this is expensive so we do it once
	env, err := cel.NewEnv(
		cel.Variable("id", cel.StringType),
		cel.Variable("kind", cel.StringType),
		cel.Variable("version", cel.IntType),
		cel.Variable("status", cel.StringType),
		cel.Variable("title", cel.StringType),
		cel.Variable("tags", cel.ListType(cel.StringType)),
	)
	if err != nil {
		// Should never happen with our static variable definitions
		panic(fmt.Sprintf("failed to create CEL environment: %v", err))
	}
	return &ConstraintValidator{
		env:      env,
		programs: make(map[string]cel.Program),
	}
}

// Validate evaluates all constraints on a node.
// The allNodes parameter is provided for cross-node constraints.
func (cv *ConstraintValidator) Validate(node *domain.Node, allNodes []domain.Node, collector *errors.Collector) {
	if node == nil {
		return
	}

	// Helper to create location from node source file
	var location *domain.Location
	if node.SourceFile != "" {
		location = &domain.Location{File: node.SourceFile}
	}

	// Evaluate each constraint
	for _, constraint := range node.Constraints {
		// Skip constraints that don't match the node's scope
		if !cv.matchesScope(constraint.Scope, node) {
			continue
		}

		if err := cv.evaluateConstraint(node, constraint, location, collector); err != nil {
			// If there's an error parsing or evaluating the CEL expression,
			// add it as an E042 error (CEL expression error)
			collector.Add(domain.DecoError{
				Code:     "E042",
				Summary:  "CEL expression error: " + constraint.Expr,
				Detail:   err.Error(),
				Location: location,
			})
		}
	}
}

// matchesScope checks if a constraint's scope applies to the given node.
// Scope patterns:
//   - "all" matches any node
//   - exact kind match (e.g., "mechanic") matches nodes with that Kind
//   - path pattern with glob (e.g., "systems/*") matches node IDs using filepath.Match
func (cv *ConstraintValidator) matchesScope(scope string, node *domain.Node) bool {
	if scope == "" || scope == "all" {
		return true
	}

	// Try exact kind match first
	if scope == node.Kind {
		return true
	}

	// Try glob pattern match against node ID
	if strings.Contains(scope, "*") || strings.Contains(scope, "?") {
		// Use filepath.Match for glob pattern matching
		matched, err := filepath.Match(scope, node.ID)
		if err == nil && matched {
			return true
		}
	}

	return false
}

// evaluateConstraint evaluates a single constraint using CEL
func (cv *ConstraintValidator) evaluateConstraint(node *domain.Node, constraint domain.Constraint, location *domain.Location, collector *errors.Collector) error {
	// Get or compile the program (cached)
	prg, err := cv.getOrCompileProgram(constraint.Expr)
	if err != nil {
		return err
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
				Code:     "E041",
				Summary:  "Constraint violation: " + constraint.Expr,
				Detail:   constraint.Message,
				Location: location,
			})
		}
	} else {
		return fmt.Errorf("CEL expression did not evaluate to a boolean")
	}

	return nil
}

// getOrCompileProgram returns a cached CEL program or compiles a new one
func (cv *ConstraintValidator) getOrCompileProgram(expr string) (cel.Program, error) {
	// Check cache first
	if prg, ok := cv.programs[expr]; ok {
		return prg, nil
	}

	// Compile the expression
	ast, issues := cv.env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("failed to compile CEL expression: %w", issues.Err())
	}

	// Create program
	prg, err := cv.env.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("failed to create CEL program: %w", err)
	}

	// Cache for reuse
	cv.programs[expr] = prg
	return prg, nil
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
	"reviewers":   true, // Review workflow approvals
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
	type nodeInfo struct {
		index      int
		sourceFile string
	}
	seen := make(map[string]nodeInfo) // ID -> first occurrence info

	for i, node := range nodes {
		if node.ID == "" {
			continue // Empty IDs are handled by schema validator
		}

		if first, exists := seen[node.ID]; exists {
			// Create location for the duplicate
			var location *domain.Location
			if node.SourceFile != "" {
				location = &domain.Location{File: node.SourceFile}
			}

			detail := fmt.Sprintf("Node ID %q appears multiple times. Node IDs must be unique.", node.ID)
			if first.sourceFile != "" && node.SourceFile != "" {
				detail = fmt.Sprintf("Node ID %q appears in multiple files: %s and %s", node.ID, first.sourceFile, node.SourceFile)
			}

			collector.Add(domain.DecoError{
				Code:     "E009",
				Summary:  fmt.Sprintf("Duplicate node ID: %s", node.ID),
				Detail:   detail,
				Location: location,
			})
		} else {
			seen[node.ID] = nodeInfo{index: i, sourceFile: node.SourceFile}
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
// The nodeID and filePath are used for error reporting.
func (uf *UnknownFieldValidator) ValidateYAML(nodeID string, filePath string, rawKeys []string, collector *errors.Collector) {
	// Collect known keys for suggestions
	var knownKeys []string
	for k := range knownTopLevelKeys {
		knownKeys = append(knownKeys, k)
	}

	// Create location from file path
	var location *domain.Location
	if filePath != "" {
		location = &domain.Location{File: filePath}
	}

	for _, key := range rawKeys {
		if !knownTopLevelKeys[key] {
			err := domain.DecoError{
				Code:     "E010",
				Summary:  fmt.Sprintf("Unknown field %q in node %s", key, nodeID),
				Detail:   fmt.Sprintf("Field %q is not a recognized top-level field. Use 'custom:' for extension data.", key),
				Location: location,
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

		// Validate keys (pass the file path for location reporting)
		uf.ValidateYAML(nodeID, path, keys, collector)

		return nil
	})
}

// ApprovalValidator validates that approved nodes have sufficient approvals.
type ApprovalValidator struct {
	requiredApprovals int
}

// NewApprovalValidator creates a new approval validator.
func NewApprovalValidator(requiredApprovals int) *ApprovalValidator {
	return &ApprovalValidator{requiredApprovals: requiredApprovals}
}

// Validate checks that approved nodes have enough current-version approvals.
func (av *ApprovalValidator) Validate(node *domain.Node, collector *errors.Collector) {
	if node == nil {
		return
	}

	// Only check approved nodes
	if node.Status != "approved" {
		return
	}

	// Helper to create location from node source file
	var location *domain.Location
	if node.SourceFile != "" {
		location = &domain.Location{File: node.SourceFile}
	}

	// Count approvals for current version
	validApprovals := 0
	for _, r := range node.Reviewers {
		if r.Version == node.Version {
			validApprovals++
		}
	}

	if validApprovals < av.requiredApprovals {
		collector.Add(domain.DecoError{
			Code:       "E050",
			Summary:    fmt.Sprintf("Node %q requires %d approval(s), has %d", node.ID, av.requiredApprovals, validApprovals),
			Detail:     fmt.Sprintf("Approved nodes must have at least %d approval(s) for version %d. Current approvals: %d.", av.requiredApprovals, node.Version, validApprovals),
			Suggestion: "Use 'deco review approve' to add approvals, or change status back to 'draft' or 'review'.",
			Location:   location,
		})
	}
}

// Orchestrator coordinates all validators and aggregates errors.
type Orchestrator struct {
	schemaValidator       *SchemaValidator
	schemaRulesValidator  *SchemaRulesValidator
	contentValidator      *ContentValidator
	referenceValidator    *ReferenceValidator
	constraintValidator   *ConstraintValidator
	duplicateIDValidator  *DuplicateIDValidator
	unknownFieldValidator *UnknownFieldValidator
	contractValidator     *ContractValidator
	blockValidator        *BlockValidator
	approvalValidator     *ApprovalValidator
}

// NewOrchestratorWithConfig creates a validator orchestrator with config-based settings.
func NewOrchestratorWithConfig(requiredApprovals int) *Orchestrator {
	return &Orchestrator{
		schemaValidator:       NewSchemaValidator(),
		contentValidator:      NewContentValidator(),
		referenceValidator:    NewReferenceValidator(),
		constraintValidator:   NewConstraintValidator(),
		duplicateIDValidator:  NewDuplicateIDValidator(),
		unknownFieldValidator: NewUnknownFieldValidator(),
		contractValidator:     NewContractValidator(),
		blockValidator:        NewBlockValidator(),
		approvalValidator:     NewApprovalValidator(requiredApprovals),
	}
}

// NewOrchestratorWithFullConfig creates a validator orchestrator with full config support.
// This includes custom block types, schema rules, and other config-driven validation rules.
func NewOrchestratorWithFullConfig(requiredApprovals int, customBlockTypes map[string]config.BlockTypeConfig, schemaRules map[string]config.SchemaRuleConfig) *Orchestrator {
	return &Orchestrator{
		schemaValidator:       NewSchemaValidator(),
		schemaRulesValidator:  NewSchemaRulesValidator(schemaRules),
		contentValidator:      NewContentValidator(),
		referenceValidator:    NewReferenceValidator(),
		constraintValidator:   NewConstraintValidator(),
		duplicateIDValidator:  NewDuplicateIDValidator(),
		unknownFieldValidator: NewUnknownFieldValidator(),
		contractValidator:     NewContractValidator(),
		blockValidator:        NewBlockValidatorWithConfig(customBlockTypes),
		approvalValidator:     NewApprovalValidator(requiredApprovals),
	}
}

// NewOrchestrator creates a new validator orchestrator.
func NewOrchestrator() *Orchestrator {
	return NewOrchestratorWithConfig(1)
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

	// Run schema rules validation on each node (per-kind required fields)
	if o.schemaRulesValidator != nil {
		for _, node := range nodes {
			o.schemaRulesValidator.Validate(&node, collector)
		}
	}

	// Run content validation on each node (approved/published require content)
	for _, node := range nodes {
		o.contentValidator.Validate(&node, collector)
	}

	// Run block validation on each node
	for _, node := range nodes {
		o.blockValidator.Validate(&node, collector)
	}

	// Run reference validation on all nodes
	o.referenceValidator.Validate(nodes, collector)

	// Run constraint validation on each node
	for _, node := range nodes {
		o.constraintValidator.Validate(&node, nodes, collector)
	}

	// Run contract validation on all nodes
	o.contractValidator.ValidateAll(nodes, collector)

	// Run approval validator on each node
	if o.approvalValidator != nil {
		for i := range nodes {
			o.approvalValidator.Validate(&nodes[i], collector)
		}
	}

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

// ValidateNode validates a single node without cross-node checks.
// This is useful for pre-save validation after patches/rewrites.
// It runs schema, content, block, and constraint validation but skips
// reference and duplicate ID checks which require all nodes.
func (o *Orchestrator) ValidateNode(node *domain.Node) *errors.Collector {
	collector := errors.NewCollectorWithLimit(100)

	// Run schema validation
	o.schemaValidator.Validate(node, collector)

	// Run schema rules validation (per-kind required fields)
	if o.schemaRulesValidator != nil {
		o.schemaRulesValidator.Validate(node, collector)
	}

	// Run content validation (approved/published require content)
	o.contentValidator.Validate(node, collector)

	// Run block validation
	o.blockValidator.Validate(node, collector)

	// Run constraint validation (single node, no cross-node checks)
	o.constraintValidator.Validate(node, []domain.Node{*node}, collector)

	// Run approval validator
	if o.approvalValidator != nil {
		o.approvalValidator.Validate(node, collector)
	}

	return collector
}
