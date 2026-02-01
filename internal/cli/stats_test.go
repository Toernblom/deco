package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStatsCommand_Structure(t *testing.T) {
	t.Run("creates stats command", func(t *testing.T) {
		cmd := NewStatsCommand()
		if cmd == nil {
			t.Fatal("Expected stats command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "stats") {
			t.Errorf("Expected Use to start with 'stats', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewStatsCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})
}

func TestStatsCommand_BasicExecution(t *testing.T) {
	t.Run("shows stats for project with nodes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewStatsCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("handles empty project", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupEmptyProject(t, tmpDir)

		cmd := NewStatsCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for empty project, got %v", err)
		}
	})
}

func TestStatsCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewStatsCommand()
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

func TestStatsCommand_WithIssues(t *testing.T) {
	t.Run("counts open issues by severity", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithIssues(t, tmpDir)

		cmd := NewStatsCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestStatsCommand_WithDanglingRefs(t *testing.T) {
	t.Run("counts dangling references", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDanglingRefs(t, tmpDir)

		cmd := NewStatsCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestStatsCommand_WithConstraints(t *testing.T) {
	t.Run("counts constraint violations", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithConstraintViolation(t, tmpDir)

		cmd := NewStatsCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestStatsCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		root := NewRootCommand()
		stats := NewStatsCommand()
		root.AddCommand(stats)

		root.SetArgs([]string{"stats", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

// Test helpers

func setupProjectWithIssues(t *testing.T, dir string) {
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

	// Create node with issues
	nodeYAML := `id: feature-001
kind: feature
version: 1
status: draft
title: Feature with Issues
issues:
  - id: issue-1
    description: Critical bug needs fixing
    severity: critical
    location: content.sections[0]
    resolved: false
  - id: issue-2
    description: Minor enhancement
    severity: low
    location: content.sections[1]
    resolved: false
  - id: issue-3
    description: Already fixed
    severity: high
    location: content.sections[2]
    resolved: true
`
	nodePath := filepath.Join(nodesDir, "feature-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
}

func setupProjectWithDanglingRefs(t *testing.T, dir string) {
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

	// Create node with dangling references
	nodeYAML := `id: feature-001
kind: feature
version: 1
status: draft
title: Feature with Dangling Refs
refs:
  uses:
    - target: nonexistent-node
      context: This node does not exist
  related:
    - target: another-missing
`
	nodePath := filepath.Join(nodesDir, "feature-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
}

func setupProjectWithConstraintViolation(t *testing.T, dir string) {
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

	// Create node with failing constraint
	nodeYAML := `id: item-001
kind: item
version: 1
status: draft
title: Item with Constraint
constraints:
  - expr: "version > 5"
    message: "Version must be greater than 5"
`
	nodePath := filepath.Join(nodesDir, "item-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
}
