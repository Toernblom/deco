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

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateCommand_Structure(t *testing.T) {
	t.Run("creates validate command", func(t *testing.T) {
		cmd := NewValidateCommand()
		if cmd == nil {
			t.Fatal("Expected validate command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "validate") {
			t.Errorf("Expected Use to start with 'validate', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewValidateCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})
}

func TestValidateCommand_ValidProject(t *testing.T) {
	t.Run("exits with code 0 for valid nodes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupValidProject(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for valid project, got %v", err)
		}
	})

	t.Run("handles empty project", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupEmptyProject(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for empty project, got %v", err)
		}
	})
}

func TestValidateCommand_InvalidNodes(t *testing.T) {
	t.Run("exits with code 1 for schema errors", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithSchemaErrors(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for invalid nodes, got nil")
		}
	})

	t.Run("exits with code 1 for reference errors", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferenceErrors(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for broken references, got nil")
		}
	})

	t.Run("exits with code 1 for constraint violations", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithConstraintViolations(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for constraint violations, got nil")
		}
	})

	t.Run("exits with code 1 for contract syntax errors", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithContractSyntaxErrors(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for contract syntax errors, got nil")
		}
	})

	t.Run("exits with code 1 for contract reference errors", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithContractReferenceErrors(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for contract reference errors, got nil")
		}
	})
}

func TestValidateCommand_ErrorOutput(t *testing.T) {
	t.Run("prints error details", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithSchemaErrors(t, tmpDir)

		// Capture output by executing command
		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})

		// Execute will return error, which is expected
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for invalid nodes")
		}

		// Check error message contains expected information
		errMsg := err.Error()
		if !strings.Contains(errMsg, "validation") &&
			!strings.Contains(errMsg, "error") &&
			!strings.Contains(errMsg, "failed") {
			t.Errorf("Expected error message to mention validation or errors, got %q", errMsg)
		}
	})
}

func TestValidateCommand_QuietFlag(t *testing.T) {
	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewValidateCommand()
		flag := cmd.Flags().Lookup("quiet")
		if flag == nil {
			t.Fatal("Expected --quiet flag to be defined")
		}
		if flag.Shorthand != "q" {
			t.Errorf("Expected shorthand 'q', got %q", flag.Shorthand)
		}
		if flag.DefValue != "false" {
			t.Errorf("Expected default 'false', got %q", flag.DefValue)
		}
	})

	t.Run("quiet flag suppresses output on valid project", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupValidProject(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{"--quiet", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for valid project, got %v", err)
		}
	})

	t.Run("quiet flag short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupValidProject(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{"-q", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error with -q, got %v", err)
		}
	})

	t.Run("quiet flag maintains exit code 1 on errors", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithSchemaErrors(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{"--quiet", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for invalid nodes even with --quiet")
		}
	})
}

func TestValidateCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for missing .deco directory, got nil")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, ".deco") &&
			!strings.Contains(errMsg, "not initialized") &&
			!strings.Contains(errMsg, "not found") {
			t.Errorf("Expected error about missing .deco directory, got %q", errMsg)
		}
	})
}

func TestValidateCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupValidProject(t, tmpDir)

		root := NewRootCommand()
		validate := NewValidateCommand()
		root.AddCommand(validate)

		root.SetArgs([]string{"validate", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("integrates with root command on invalid project", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithSchemaErrors(t, tmpDir)

		root := NewRootCommand()
		validate := NewValidateCommand()
		root.AddCommand(validate)

		root.SetArgs([]string{"validate", tmpDir})
		err := root.Execute()
		if err == nil {
			t.Error("Expected error for invalid project via root command")
		}
	})
}

// Test helpers

func setupValidProject(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create valid node
	nodeYAML := `id: test-item-001
kind: item
version: 1
status: draft
title: Test Item
tags:
  - test
`
	nodePath := filepath.Join(nodesDir, "test-item-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create valid node: %v", err)
	}
}

func setupEmptyProject(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure with no nodes
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: empty-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}
}

func setupProjectWithSchemaErrors(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: invalid-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create invalid node (missing required fields)
	invalidNode := `id: test-item-001
kind: item
# Missing version, status, title - should fail schema validation
tags:
  - test
`
	nodePath := filepath.Join(nodesDir, "test-item-001.yaml")
	if err := os.WriteFile(nodePath, []byte(invalidNode), 0644); err != nil {
		t.Fatalf("Failed to write invalid node: %v", err)
	}
}

func setupProjectWithReferenceErrors(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: broken-refs-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create node with broken reference
	nodeYAML := `id: test-item-001
kind: item
version: 1
status: draft
title: Test Item
tags:
  - test
refs:
  uses:
    - target: nonexistent-item-999
`
	nodePath := filepath.Join(nodesDir, "test-item-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node with broken reference: %v", err)
	}
}

func setupProjectWithConstraintViolations(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: constraint-violation-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create node with constraint that will fail
	nodeYAML := `id: test-item-001
kind: item
version: 1
status: draft
title: Test Item
tags:
  - test
constraints:
  - expr: "status == 'published'"
    message: "Item must be published"
`
	nodePath := filepath.Join(nodesDir, "test-item-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node with constraint: %v", err)
	}
}

func setupProjectWithContractSyntaxErrors(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: contract-syntax-error-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create node with contract syntax errors (missing name, empty step)
	nodeYAML := `id: test-feature-001
kind: feature
version: 1
status: draft
title: Test Feature
contracts:
  - name: ""
    scenario: "Contract with missing name"
    given:
      - "some precondition"
    when:
      - ""
    then:
      - "expected result"
`
	nodePath := filepath.Join(nodesDir, "test-feature-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node with contract syntax errors: %v", err)
	}
}

func setupProjectWithContractReferenceErrors(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: contract-ref-error-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create node with contract referencing non-existent node
	nodeYAML := `id: test-feature-001
kind: feature
version: 1
status: draft
title: Test Feature
contracts:
  - name: "Test Flow"
    scenario: "Contract with invalid node reference"
    given:
      - "@systems/nonexistent is active"
    when:
      - "player does something"
    then:
      - "expected result"
`
	nodePath := filepath.Join(nodesDir, "test-feature-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node with contract reference errors: %v", err)
	}
}

func TestValidateCommand_ApprovalValidation(t *testing.T) {
	t.Run("validates approval requirements from config", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithApprovalConfig(t, tmpDir, 2)
		createApprovedNodeWithOneApproval(t, tmpDir)

		cmd := NewValidateCommand()
		cmd.SetArgs([]string{tmpDir})

		err := cmd.Execute()
		// Should have validation errors (E050 for insufficient approvals)
		if err == nil {
			t.Error("Expected validation error E050 for insufficient approvals")
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
content:
  sections:
    - name: Overview
      blocks:
        - type: rule
          text: A test rule
reviewers:
  - name: alice@example.com
    timestamp: 2026-01-01T00:00:00Z
    version: 1
`
	nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
	os.WriteFile(filepath.Join(nodesDir, "test-node.yaml"), []byte(nodeContent), 0644)
}
