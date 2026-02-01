# Review Workflow Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a review workflow with approvals, status transitions, and audit trail for GDD nodes.

**Architecture:** Extend Node struct with `Reviewers` field. Add new audit operations (submit, approve, reject). Create ApprovalValidator to enforce approval requirements before status transitions. Add `deco review` CLI command with subcommands.

**Tech Stack:** Go, Cobra CLI, YAML storage, JSONL history

---

## Task 1: Add Reviewer Struct and Reviewers Field to Node

**Files:**
- Modify: `internal/domain/node.go`
- Test: `internal/domain/node_test.go`

**Step 1: Write the failing test**

Add to `internal/domain/node_test.go`:

```go
func TestNode_ReviewerStruct(t *testing.T) {
	t.Run("reviewer has required fields", func(t *testing.T) {
		reviewer := domain.Reviewer{
			Name:      "alice@example.com",
			Timestamp: time.Now(),
			Version:   1,
			Note:      "LGTM",
		}
		if reviewer.Name == "" {
			t.Error("Expected Name to be set")
		}
		if reviewer.Version != 1 {
			t.Error("Expected Version to be 1")
		}
	})

	t.Run("node can have reviewers", func(t *testing.T) {
		node := domain.Node{
			ID:      "test/node",
			Kind:    "mechanic",
			Version: 1,
			Status:  "review",
			Title:   "Test Node",
			Reviewers: []domain.Reviewer{
				{Name: "alice@example.com", Timestamp: time.Now(), Version: 1},
			},
		}
		if len(node.Reviewers) != 1 {
			t.Errorf("Expected 1 reviewer, got %d", len(node.Reviewers))
		}
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/domain/... -run TestNode_ReviewerStruct -v`
Expected: FAIL with "undefined: domain.Reviewer"

**Step 3: Write minimal implementation**

Add to `internal/domain/node.go` after the Contract struct:

```go
// Reviewer represents an approval record for a node version.
type Reviewer struct {
	Name      string    `json:"name" yaml:"name"`           // reviewer email/username
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"` // when approved
	Version   int       `json:"version" yaml:"version"`     // version that was approved
	Note      string    `json:"note,omitempty" yaml:"note,omitempty"` // optional comment
}
```

Add to Node struct (after Constraints field):

```go
	Reviewers   []Reviewer             `json:"reviewers,omitempty" yaml:"reviewers,omitempty"`
```

Add import for "time" if not present.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/domain/... -run TestNode_ReviewerStruct -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/domain/node.go internal/domain/node_test.go
git commit -m "feat(domain): add Reviewer struct and Reviewers field to Node"
```

---

## Task 2: Add required_approvals to Config

**Files:**
- Modify: `internal/storage/config/repository.go`
- Test: `internal/storage/config/yaml_repository_test.go`

**Step 1: Write the failing test**

Add to `internal/storage/config/yaml_repository_test.go`:

```go
func TestConfig_RequiredApprovals(t *testing.T) {
	t.Run("loads required_approvals from config", func(t *testing.T) {
		tmpDir := t.TempDir()
		decoDir := filepath.Join(tmpDir, ".deco")
		os.MkdirAll(decoDir, 0755)

		configContent := `project_name: TestProject
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
version: 1
required_approvals: 2
`
		os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644)

		repo := NewYAMLRepository(tmpDir)
		cfg, err := repo.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if cfg.RequiredApprovals != 2 {
			t.Errorf("Expected RequiredApprovals=2, got %d", cfg.RequiredApprovals)
		}
	})

	t.Run("defaults to 1 if not specified", func(t *testing.T) {
		tmpDir := t.TempDir()
		decoDir := filepath.Join(tmpDir, ".deco")
		os.MkdirAll(decoDir, 0755)

		configContent := `project_name: TestProject
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
version: 1
`
		os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644)

		repo := NewYAMLRepository(tmpDir)
		cfg, err := repo.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if cfg.RequiredApprovals != 1 {
			t.Errorf("Expected RequiredApprovals=1 (default), got %d", cfg.RequiredApprovals)
		}
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/storage/config/... -run TestConfig_RequiredApprovals -v`
Expected: FAIL with "cfg.RequiredApprovals undefined"

**Step 3: Write minimal implementation**

Add to `internal/storage/config/repository.go` Config struct:

```go
	// RequiredApprovals is the number of approvals needed for a node to be approved.
	// Defaults to 1.
	RequiredApprovals int `yaml:"required_approvals" json:"required_approvals"`
```

In `internal/storage/config/yaml_repository.go`, in the Load() function, after loading the config add default:

```go
	// Default required_approvals to 1 if not set
	if cfg.RequiredApprovals == 0 {
		cfg.RequiredApprovals = 1
	}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/storage/config/... -run TestConfig_RequiredApprovals -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/storage/config/repository.go internal/storage/config/yaml_repository.go internal/storage/config/yaml_repository_test.go
git commit -m "feat(config): add required_approvals setting with default of 1"
```

---

## Task 3: Add Review Operations to Audit

**Files:**
- Modify: `internal/domain/audit.go`
- Test: `internal/domain/audit_test.go`

**Step 1: Write the failing test**

Add to `internal/domain/audit_test.go`:

```go
func TestAuditEntry_ReviewOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		wantErr   bool
	}{
		{"submit is valid", "submit", false},
		{"approve is valid", "approve", false},
		{"reject is valid", "reject", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := AuditEntry{
				Timestamp: time.Now(),
				NodeID:    "test/node",
				Operation: tt.operation,
				User:      "alice",
			}
			err := entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/domain/... -run TestAuditEntry_ReviewOperations -v`
Expected: FAIL with "Operation must be one of..."

**Step 3: Write minimal implementation**

In `internal/domain/audit.go`, update the validOperations map in Validate():

```go
	validOperations := map[string]bool{
		"create":  true,
		"update":  true,
		"delete":  true,
		"set":     true,
		"append":  true,
		"unset":   true,
		"move":    true,
		"submit":  true,  // draft -> review
		"approve": true,  // add approval
		"reject":  true,  // review -> draft
	}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/domain/... -run TestAuditEntry_ReviewOperations -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/domain/audit.go internal/domain/audit_test.go
git commit -m "feat(audit): add submit, approve, reject operations for review workflow"
```

---

## Task 4: Create ApprovalValidator

**Files:**
- Modify: `internal/services/validator/validator.go`
- Test: `internal/services/validator/validator_test.go`

**Step 1: Write the failing test**

Add to `internal/services/validator/validator_test.go`:

```go
func TestApprovalValidator(t *testing.T) {
	t.Run("approved node without enough approvals fails", func(t *testing.T) {
		validator := NewApprovalValidator(2) // require 2 approvals
		collector := errors.NewCollector()

		node := &domain.Node{
			ID:         "test/node",
			Kind:       "mechanic",
			Version:    1,
			Status:     "approved",
			Title:      "Test",
			SourceFile: "test.yaml",
			Reviewers: []domain.Reviewer{
				{Name: "alice", Timestamp: time.Now(), Version: 1},
			},
		}

		validator.Validate(node, collector)

		if !collector.HasErrors() {
			t.Error("Expected validation error for insufficient approvals")
		}
	})

	t.Run("approved node with enough approvals passes", func(t *testing.T) {
		validator := NewApprovalValidator(2)
		collector := errors.NewCollector()

		node := &domain.Node{
			ID:         "test/node",
			Kind:       "mechanic",
			Version:    1,
			Status:     "approved",
			Title:      "Test",
			SourceFile: "test.yaml",
			Reviewers: []domain.Reviewer{
				{Name: "alice", Timestamp: time.Now(), Version: 1},
				{Name: "bob", Timestamp: time.Now(), Version: 1},
			},
		}

		validator.Validate(node, collector)

		if collector.HasErrors() {
			t.Errorf("Expected no errors, got: %v", collector.Errors())
		}
	})

	t.Run("draft node skips approval check", func(t *testing.T) {
		validator := NewApprovalValidator(2)
		collector := errors.NewCollector()

		node := &domain.Node{
			ID:         "test/node",
			Kind:       "mechanic",
			Version:    1,
			Status:     "draft",
			Title:      "Test",
			SourceFile: "test.yaml",
		}

		validator.Validate(node, collector)

		if collector.HasErrors() {
			t.Errorf("Expected no errors for draft node, got: %v", collector.Errors())
		}
	})

	t.Run("approvals must match current version", func(t *testing.T) {
		validator := NewApprovalValidator(1)
		collector := errors.NewCollector()

		node := &domain.Node{
			ID:         "test/node",
			Kind:       "mechanic",
			Version:    2, // current version is 2
			Status:     "approved",
			Title:      "Test",
			SourceFile: "test.yaml",
			Reviewers: []domain.Reviewer{
				{Name: "alice", Timestamp: time.Now(), Version: 1}, // approved version 1
			},
		}

		validator.Validate(node, collector)

		if !collector.HasErrors() {
			t.Error("Expected error for stale approval")
		}
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/services/validator/... -run TestApprovalValidator -v`
Expected: FAIL with "undefined: NewApprovalValidator"

**Step 3: Write minimal implementation**

Add to `internal/services/validator/validator.go` (before Orchestrator):

```go
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/services/validator/... -run TestApprovalValidator -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/services/validator/validator.go internal/services/validator/validator_test.go
git commit -m "feat(validator): add ApprovalValidator for review workflow"
```

---

## Task 5: Integrate ApprovalValidator into Orchestrator

**Files:**
- Modify: `internal/services/validator/validator.go`
- Test: `internal/services/validator/validator_test.go`

**Step 1: Write the failing test**

Add to `internal/services/validator/validator_test.go`:

```go
func TestOrchestrator_ApprovalValidation(t *testing.T) {
	t.Run("orchestrator validates approvals", func(t *testing.T) {
		orch := NewOrchestratorWithConfig(2) // 2 required approvals

		nodes := []domain.Node{
			{
				ID:         "test/approved",
				Kind:       "mechanic",
				Version:    1,
				Status:     "approved",
				Title:      "Test",
				SourceFile: "test.yaml",
				Reviewers: []domain.Reviewer{
					{Name: "alice", Timestamp: time.Now(), Version: 1},
				},
			},
		}

		collector := orch.ValidateAll(nodes)

		// Should have error for insufficient approvals
		hasApprovalError := false
		for _, err := range collector.Errors() {
			if err.Code == "E050" {
				hasApprovalError = true
				break
			}
		}
		if !hasApprovalError {
			t.Error("Expected approval validation error E050")
		}
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/services/validator/... -run TestOrchestrator_ApprovalValidation -v`
Expected: FAIL with "undefined: NewOrchestratorWithConfig"

**Step 3: Write minimal implementation**

Update `internal/services/validator/validator.go`:

Add field to Orchestrator struct:
```go
type Orchestrator struct {
	schemaValidator       *SchemaValidator
	contentValidator      *ContentValidator
	referenceValidator    *ReferenceValidator
	constraintValidator   *ConstraintValidator
	duplicateIDValidator  *DuplicateIDValidator
	unknownFieldValidator *UnknownFieldValidator
	contractValidator     *ContractValidator
	blockValidator        *BlockValidator
	approvalValidator     *ApprovalValidator
}
```

Add new constructor:
```go
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
```

Update NewOrchestrator to use default of 1:
```go
func NewOrchestrator() *Orchestrator {
	return NewOrchestratorWithConfig(1)
}
```

In ValidateAll(), add after single-node validators loop:
```go
	// Run approval validator on each node
	if o.approvalValidator != nil {
		for i := range nodes {
			o.approvalValidator.Validate(&nodes[i], collector)
		}
	}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/services/validator/... -run TestOrchestrator_ApprovalValidation -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/services/validator/validator.go internal/services/validator/validator_test.go
git commit -m "feat(validator): integrate ApprovalValidator into Orchestrator"
```

---

## Task 6: Update CLI validate command to use config

**Files:**
- Modify: `internal/cli/validate.go`
- Test: `internal/cli/validate_test.go`

**Step 1: Write the failing test**

Add to `internal/cli/validate_test.go`:

```go
func TestValidateCommand_ApprovalValidation(t *testing.T) {
	t.Run("validates approval requirements from config", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithApprovalConfig(t, tmpDir, 2)
		createApprovedNodeWithOneApproval(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})

		// Capture output
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)

		err := cmd.Execute()
		// Should have validation errors (exit code or error return)
		output := buf.String()
		if !strings.Contains(output, "E050") && err == nil {
			t.Error("Expected approval validation error E050")
		}
	})
}

func setupProjectWithApprovalConfig(t *testing.T, tmpDir string, requiredApprovals int) {
	t.Helper()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	os.MkdirAll(nodesDir, 0755)

	configContent := fmt.Sprintf(`project_name: TestProject
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
version: 1
required_approvals: %d
`, requiredApprovals)
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644)
}

func createApprovedNodeWithOneApproval(t *testing.T, tmpDir string) {
	t.Helper()
	nodeContent := `id: test/node
kind: mechanic
version: 1
status: approved
title: Test Node
reviewers:
  - name: alice@example.com
    timestamp: 2026-01-01T00:00:00Z
    version: 1
`
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	os.WriteFile(filepath.Join(nodesDir, "test-node.yaml"), []byte(nodeContent), 0644)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli/... -run TestValidateCommand_ApprovalValidation -v`
Expected: FAIL (validate command doesn't use config's required_approvals yet)

**Step 3: Write minimal implementation**

In `internal/cli/validate.go`, update runValidate() to load config and pass required_approvals:

```go
func runValidate(flags *validateFlags) error {
	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load all nodes
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	// Validate with config-based settings
	orchestrator := validator.NewOrchestratorWithConfig(cfg.RequiredApprovals)
	collector := orchestrator.ValidateAll(nodes)
	// ... rest of function
}
```

Add import for config package if not present.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli/... -run TestValidateCommand_ApprovalValidation -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/cli/validate.go internal/cli/validate_test.go
git commit -m "feat(cli): use config required_approvals in validate command"
```

---

## Task 7: Create review.go CLI with submit subcommand

**Files:**
- Create: `internal/cli/review.go`
- Test: `internal/cli/review_test.go`

**Step 1: Write the failing test**

Create `internal/cli/review_test.go`:

```go
package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReviewCommand_Structure(t *testing.T) {
	t.Run("creates review command with subcommands", func(t *testing.T) {
		cmd := NewReviewCommand()
		if cmd == nil {
			t.Fatal("Expected review command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "review") {
			t.Errorf("Expected Use to start with 'review', got %q", cmd.Use)
		}

		// Should have submit subcommand
		submitCmd, _, _ := cmd.Find([]string{"submit"})
		if submitCmd == nil || submitCmd.Use == "review" {
			t.Error("Expected submit subcommand")
		}
	})
}

func TestReviewCommand_Submit(t *testing.T) {
	t.Run("submit changes status from draft to review", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createDraftNode(t, tmpDir, "test/node")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"submit", "test/node", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify status changed
		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: review") {
			t.Errorf("Expected status to be 'review', got: %s", nodeYAML)
		}
	})

	t.Run("submit fails if not in draft status", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "approved")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"submit", "test/node", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when submitting non-draft node")
		}
	})
}

func setupReviewProject(t *testing.T, tmpDir string) {
	t.Helper()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	os.MkdirAll(nodesDir, 0755)

	configContent := `project_name: TestProject
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
version: 1
required_approvals: 1
`
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644)
}

func createDraftNode(t *testing.T, tmpDir, nodeID string) {
	t.Helper()
	createNodeWithStatus(t, tmpDir, nodeID, "draft")
}

func createNodeWithStatus(t *testing.T, tmpDir, nodeID, status string) {
	t.Helper()
	nodeContent := `id: ` + nodeID + `
kind: mechanic
version: 1
status: ` + status + `
title: Test Node
`
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	filename := strings.ReplaceAll(nodeID, "/", "-") + ".yaml"
	os.WriteFile(filepath.Join(nodesDir, filename), []byte(nodeContent), 0644)
}

func readNodeFileByID(t *testing.T, tmpDir, nodeID string) string {
	t.Helper()
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	filename := strings.ReplaceAll(nodeID, "/", "-") + ".yaml"
	content, err := os.ReadFile(filepath.Join(nodesDir, filename))
	if err != nil {
		t.Fatalf("Failed to read node file: %v", err)
	}
	return string(content)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli/... -run TestReviewCommand -v`
Expected: FAIL with "undefined: NewReviewCommand"

**Step 3: Write minimal implementation**

Create `internal/cli/review.go`:

```go
package cli

import (
	"fmt"
	"os/user"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type reviewFlags struct {
	targetDir string
	nodeID    string
	quiet     bool
}

// NewReviewCommand creates the review command with subcommands
func NewReviewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Manage review workflow for nodes",
		Long: `Manage review workflow for nodes.

Subcommands:
  submit   - Submit a draft node for review
  approve  - Approve a node under review
  reject   - Reject a node back to draft
  status   - Show review status of a node`,
	}

	cmd.AddCommand(newSubmitCommand())

	return cmd
}

func newSubmitCommand() *cobra.Command {
	flags := &reviewFlags{}

	cmd := &cobra.Command{
		Use:   "submit <node-id> [directory]",
		Short: "Submit a draft node for review",
		Long:  `Submit a draft node for review. Changes status from 'draft' to 'review'.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runSubmit(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runSubmit(flags *reviewFlags) error {
	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Validate current status
	if n.Status != "draft" {
		return fmt.Errorf("cannot submit node %q: status is %q, must be 'draft'", flags.nodeID, n.Status)
	}

	// Update status
	oldStatus := n.Status
	n.Status = "review"

	// Save the node
	if err := nodeRepo.Save(n); err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log submit operation
	if err := logReviewOperation(flags.targetDir, n.ID, "submit", oldStatus, n.Status, ""); err != nil {
		fmt.Printf("Warning: failed to log submit operation: %v\n", err)
	}

	if !flags.quiet {
		fmt.Printf("Submitted %s for review (status: draft -> review)\n", flags.nodeID)
	}

	return nil
}

func logReviewOperation(targetDir, nodeID, operation, oldStatus, newStatus, note string) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Operation: operation,
		User:      username,
		Before:    map[string]interface{}{"status": oldStatus},
		After:     map[string]interface{}{"status": newStatus},
	}

	if note != "" {
		entry.After["note"] = note
	}

	return historyRepo.Append(entry)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli/... -run TestReviewCommand -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/cli/review.go internal/cli/review_test.go
git commit -m "feat(cli): add review command with submit subcommand"
```

---

## Task 8: Add approve subcommand

**Files:**
- Modify: `internal/cli/review.go`
- Test: `internal/cli/review_test.go`

**Step 1: Write the failing test**

Add to `internal/cli/review_test.go`:

```go
func TestReviewCommand_Approve(t *testing.T) {
	t.Run("approve adds reviewer to node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"approve", "test/node", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "reviewers:") {
			t.Error("Expected reviewers field to be added")
		}
	})

	t.Run("approve with note includes note in reviewer", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"approve", "test/node", "--note", "LGTM", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "LGTM") {
			t.Errorf("Expected note 'LGTM' in node, got: %s", nodeYAML)
		}
	})

	t.Run("approve transitions to approved when requirements met", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir) // requires 1 approval
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"approve", "test/node", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: approved") {
			t.Errorf("Expected status 'approved', got: %s", nodeYAML)
		}
	})

	t.Run("approve fails if not in review status", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "draft")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"approve", "test/node", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when approving non-review node")
		}
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli/... -run TestReviewCommand_Approve -v`
Expected: FAIL (approve subcommand doesn't exist)

**Step 3: Write minimal implementation**

Add to `internal/cli/review.go` in NewReviewCommand():

```go
	cmd.AddCommand(newApproveCommand())
```

Add function:

```go
type approveFlags struct {
	targetDir string
	nodeID    string
	note      string
	quiet     bool
}

func newApproveCommand() *cobra.Command {
	flags := &approveFlags{}

	cmd := &cobra.Command{
		Use:   "approve <node-id> [directory]",
		Short: "Approve a node under review",
		Long:  `Approve a node under review. Adds your approval and transitions to 'approved' if requirements are met.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runApprove(flags)
		},
	}

	cmd.Flags().StringVar(&flags.note, "note", "", "Optional approval note")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runApprove(flags *approveFlags) error {
	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Validate current status
	if n.Status != "review" {
		return fmt.Errorf("cannot approve node %q: status is %q, must be 'review'", flags.nodeID, n.Status)
	}

	// Get current user
	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	// Add reviewer
	reviewer := domain.Reviewer{
		Name:      username,
		Timestamp: time.Now(),
		Version:   n.Version,
		Note:      flags.note,
	}
	n.Reviewers = append(n.Reviewers, reviewer)

	// Count approvals for current version
	validApprovals := 0
	for _, r := range n.Reviewers {
		if r.Version == n.Version {
			validApprovals++
		}
	}

	// Transition to approved if requirements met
	oldStatus := n.Status
	if validApprovals >= cfg.RequiredApprovals {
		n.Status = "approved"
	}

	// Save the node
	if err := nodeRepo.Save(n); err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log approve operation
	if err := logReviewOperation(flags.targetDir, n.ID, "approve", oldStatus, n.Status, flags.note); err != nil {
		fmt.Printf("Warning: failed to log approve operation: %v\n", err)
	}

	if !flags.quiet {
		if n.Status == "approved" {
			fmt.Printf("Approved %s (status: review -> approved, %d/%d approvals)\n", flags.nodeID, validApprovals, cfg.RequiredApprovals)
		} else {
			fmt.Printf("Added approval to %s (%d/%d approvals needed)\n", flags.nodeID, validApprovals, cfg.RequiredApprovals)
		}
	}

	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli/... -run TestReviewCommand_Approve -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/cli/review.go internal/cli/review_test.go
git commit -m "feat(cli): add approve subcommand to review"
```

---

## Task 9: Add reject subcommand

**Files:**
- Modify: `internal/cli/review.go`
- Test: `internal/cli/review_test.go`

**Step 1: Write the failing test**

Add to `internal/cli/review_test.go`:

```go
func TestReviewCommand_Reject(t *testing.T) {
	t.Run("reject changes status from review to draft", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"reject", "test/node", "--note", "Needs more detail", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: draft") {
			t.Errorf("Expected status 'draft', got: %s", nodeYAML)
		}
	})

	t.Run("reject requires note", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "review")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"reject", "test/node", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when rejecting without note")
		}
	})

	t.Run("reject fails if not in review status", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithStatus(t, tmpDir, "test/node", "draft")

		cmd := NewReviewCommand()
		cmd.SetArgs([]string{"reject", "test/node", "--note", "reason", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when rejecting non-review node")
		}
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli/... -run TestReviewCommand_Reject -v`
Expected: FAIL (reject subcommand doesn't exist)

**Step 3: Write minimal implementation**

Add to `internal/cli/review.go` in NewReviewCommand():

```go
	cmd.AddCommand(newRejectCommand())
```

Add function:

```go
type rejectFlags struct {
	targetDir string
	nodeID    string
	note      string
	quiet     bool
}

func newRejectCommand() *cobra.Command {
	flags := &rejectFlags{}

	cmd := &cobra.Command{
		Use:   "reject <node-id> [directory]",
		Short: "Reject a node back to draft",
		Long:  `Reject a node back to draft. Requires a note explaining the rejection reason.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runReject(flags)
		},
	}

	cmd.Flags().StringVar(&flags.note, "note", "", "Rejection reason (required)")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runReject(flags *rejectFlags) error {
	// Validate note is provided
	if flags.note == "" {
		return fmt.Errorf("rejection note is required (use --note)")
	}

	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Validate current status
	if n.Status != "review" {
		return fmt.Errorf("cannot reject node %q: status is %q, must be 'review'", flags.nodeID, n.Status)
	}

	// Update status
	oldStatus := n.Status
	n.Status = "draft"

	// Save the node
	if err := nodeRepo.Save(n); err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log reject operation
	if err := logReviewOperation(flags.targetDir, n.ID, "reject", oldStatus, n.Status, flags.note); err != nil {
		fmt.Printf("Warning: failed to log reject operation: %v\n", err)
	}

	if !flags.quiet {
		fmt.Printf("Rejected %s (status: review -> draft)\nReason: %s\n", flags.nodeID, flags.note)
	}

	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli/... -run TestReviewCommand_Reject -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/cli/review.go internal/cli/review_test.go
git commit -m "feat(cli): add reject subcommand to review"
```

---

## Task 10: Add status subcommand

**Files:**
- Modify: `internal/cli/review.go`
- Test: `internal/cli/review_test.go`

**Step 1: Write the failing test**

Add to `internal/cli/review_test.go`:

```go
func TestReviewCommand_Status(t *testing.T) {
	t.Run("status shows review state", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createNodeWithReviewers(t, tmpDir, "test/node", "review", 1)

		cmd := NewReviewCommand()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetArgs([]string{"status", "test/node", tmpDir})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "review") {
			t.Errorf("Expected output to contain status, got: %s", output)
		}
	})
}

func createNodeWithReviewers(t *testing.T, tmpDir, nodeID, status string, numReviewers int) {
	t.Helper()
	nodeContent := `id: ` + nodeID + `
kind: mechanic
version: 1
status: ` + status + `
title: Test Node
`
	if numReviewers > 0 {
		nodeContent += "reviewers:\n"
		for i := 0; i < numReviewers; i++ {
			nodeContent += fmt.Sprintf(`  - name: reviewer%d@example.com
    timestamp: 2026-01-01T00:00:00Z
    version: 1
`, i+1)
		}
	}
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	filename := strings.ReplaceAll(nodeID, "/", "-") + ".yaml"
	os.WriteFile(filepath.Join(nodesDir, filename), []byte(nodeContent), 0644)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli/... -run TestReviewCommand_Status -v`
Expected: FAIL (status subcommand doesn't exist)

**Step 3: Write minimal implementation**

Add to `internal/cli/review.go` in NewReviewCommand():

```go
	cmd.AddCommand(newStatusCommand())
```

Add function:

```go
type statusFlags struct {
	targetDir string
	nodeID    string
}

func newStatusCommand() *cobra.Command {
	flags := &statusFlags{}

	cmd := &cobra.Command{
		Use:   "status <node-id> [directory]",
		Short: "Show review status of a node",
		Long:  `Show review status of a node including current status, approvals, and requirements.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			if len(args) > 1 {
				flags.targetDir = args[1]
			} else {
				flags.targetDir = "."
			}
			return runStatus(cmd, flags)
		},
	}

	return cmd
}

func runStatus(cmd *cobra.Command, flags *statusFlags) error {
	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Count approvals for current version
	validApprovals := 0
	for _, r := range n.Reviewers {
		if r.Version == n.Version {
			validApprovals++
		}
	}

	// Print status
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Node: %s\n", n.ID)
	fmt.Fprintf(out, "Version: %d\n", n.Version)
	fmt.Fprintf(out, "Status: %s\n", n.Status)
	fmt.Fprintf(out, "Approvals: %d/%d\n", validApprovals, cfg.RequiredApprovals)

	if len(n.Reviewers) > 0 {
		fmt.Fprintf(out, "\nReviewers:\n")
		for _, r := range n.Reviewers {
			versionNote := ""
			if r.Version != n.Version {
				versionNote = " (stale - v" + fmt.Sprint(r.Version) + ")"
			}
			fmt.Fprintf(out, "  - %s at %s%s\n", r.Name, r.Timestamp.Format("2006-01-02 15:04"), versionNote)
			if r.Note != "" {
				fmt.Fprintf(out, "    Note: %s\n", r.Note)
			}
		}
	}

	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli/... -run TestReviewCommand_Status -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/cli/review.go internal/cli/review_test.go
git commit -m "feat(cli): add status subcommand to review"
```

---

## Task 11: Register review command in main.go

**Files:**
- Modify: `cmd/deco/main.go`

**Step 1: Add registration**

In `cmd/deco/main.go`, add after other command registrations:

```go
	root.AddCommand(cli.NewReviewCommand())
```

**Step 2: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 3: Manual test**

Run: `go build -o deco ./cmd/deco && ./deco review --help`
Expected: Shows review command help with subcommands

**Step 4: Commit**

```bash
git add cmd/deco/main.go
git commit -m "feat(cli): register review command in main"
```

---

## Task 12: Auto-reset status on edit

**Files:**
- Modify: `internal/cli/set.go`
- Test: `internal/cli/set_test.go`

**Step 1: Write the failing test**

Add to `internal/cli/set_test.go`:

```go
func TestSetCommand_ReviewReset(t *testing.T) {
	t.Run("editing approved node resets status to draft", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupApprovedNode(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"test-node", "title", "New Title", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "test-node")
		if !strings.Contains(nodeYAML, "status: draft") {
			t.Errorf("Expected status to reset to 'draft', got: %s", nodeYAML)
		}
	})

	t.Run("editing approved node clears reviewers", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupApprovedNode(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"test-node", "title", "New Title", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "test-node")
		if strings.Contains(nodeYAML, "reviewers:") {
			t.Errorf("Expected reviewers to be cleared, got: %s", nodeYAML)
		}
	})
}

func setupApprovedNode(t *testing.T, tmpDir string) {
	t.Helper()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	os.MkdirAll(nodesDir, 0755)

	configContent := `project_name: TestProject
nodes_path: .deco/nodes
version: 1
`
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configContent), 0644)

	nodeContent := `id: test-node
kind: mechanic
version: 1
status: approved
title: Test Node
reviewers:
  - name: alice@example.com
    timestamp: 2026-01-01T00:00:00Z
    version: 1
`
	os.WriteFile(filepath.Join(nodesDir, "test-node.yaml"), []byte(nodeContent), 0644)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli/... -run TestSetCommand_ReviewReset -v`
Expected: FAIL (status not reset, reviewers not cleared)

**Step 3: Write minimal implementation**

In `internal/cli/set.go`, in runSet() after loading the node and before applying the patch:

```go
	// Auto-reset review status on edit
	statusResetNeeded := n.Status == "approved" || n.Status == "review"
	if statusResetNeeded {
		n.Status = "draft"
		n.Reviewers = nil // Clear stale approvals
	}
```

Add import for domain if needed (already imported).

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli/... -run TestSetCommand_ReviewReset -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/cli/set.go internal/cli/set_test.go
git commit -m "feat(cli): auto-reset status to draft on node edit"
```

---

## Task 13: Final Integration Test

**Files:**
- Test: `internal/cli/review_test.go`

**Step 1: Write integration test**

Add to `internal/cli/review_test.go`:

```go
func TestReviewWorkflow_Integration(t *testing.T) {
	t.Run("full workflow: draft -> review -> approved -> edit -> draft", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupReviewProject(t, tmpDir)
		createDraftNode(t, tmpDir, "test/node")

		// 1. Submit for review
		submitCmd := NewReviewCommand()
		submitCmd.SetArgs([]string{"submit", "test/node", tmpDir})
		if err := submitCmd.Execute(); err != nil {
			t.Fatalf("Submit failed: %v", err)
		}

		nodeYAML := readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: review") {
			t.Fatalf("Expected status 'review' after submit")
		}

		// 2. Approve
		approveCmd := NewReviewCommand()
		approveCmd.SetArgs([]string{"approve", "test/node", "--note", "LGTM", tmpDir})
		if err := approveCmd.Execute(); err != nil {
			t.Fatalf("Approve failed: %v", err)
		}

		nodeYAML = readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: approved") {
			t.Fatalf("Expected status 'approved' after approve")
		}

		// 3. Edit (should reset to draft)
		setCmd := NewSetCommand()
		setCmd.SetArgs([]string{"test/node", "title", "Updated Title", tmpDir})
		if err := setCmd.Execute(); err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		nodeYAML = readNodeFileByID(t, tmpDir, "test/node")
		if !strings.Contains(nodeYAML, "status: draft") {
			t.Fatalf("Expected status 'draft' after edit")
		}
		if strings.Contains(nodeYAML, "reviewers:") {
			t.Fatalf("Expected reviewers to be cleared after edit")
		}
	})
}
```

**Step 2: Run integration test**

Run: `go test ./internal/cli/... -run TestReviewWorkflow_Integration -v`
Expected: PASS

**Step 3: Run all tests**

Run: `go test ./... -count=1`
Expected: All tests pass

**Step 4: Commit**

```bash
git add internal/cli/review_test.go
git commit -m "test(cli): add review workflow integration test"
```

---

## Summary

This plan implements the review workflow in 13 tasks:

1. Add Reviewer struct and Reviewers field to Node
2. Add required_approvals to Config
3. Add review operations to Audit
4. Create ApprovalValidator
5. Integrate ApprovalValidator into Orchestrator
6. Update validate command to use config
7. Create review.go with submit subcommand
8. Add approve subcommand
9. Add reject subcommand
10. Add status subcommand
11. Register review command in main.go
12. Auto-reset status on edit
13. Final integration test

Each task follows TDD with failing test first, minimal implementation, verification, and commit.
