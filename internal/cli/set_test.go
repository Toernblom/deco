package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetCommand_Structure(t *testing.T) {
	t.Run("creates set command", func(t *testing.T) {
		cmd := NewSetCommand()
		if cmd == nil {
			t.Fatal("Expected set command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "set") {
			t.Errorf("Expected Use to start with 'set', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewSetCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("requires three arguments", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when no arguments provided")
		}
	})

	t.Run("requires path and value", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when only node ID provided")
		}
	})
}

func TestSetCommand_SetSimpleField(t *testing.T) {
	t.Run("sets string field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", "title", "Golden Sword", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify the change persisted
		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "Golden Sword") {
			t.Errorf("Expected title to be changed, got: %s", nodeYAML)
		}
	})

	t.Run("sets summary field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", "summary", "A legendary golden sword", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "A legendary golden sword") {
			t.Errorf("Expected summary to be changed, got: %s", nodeYAML)
		}
	})

	t.Run("sets status field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", "status", "published", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "status: published") {
			t.Errorf("Expected status to be published, got: %s", nodeYAML)
		}
	})
}

func TestSetCommand_SetNestedField(t *testing.T) {
	t.Run("sets array element by index", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", "tags[0]", "legendary", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "legendary") {
			t.Errorf("Expected tag to be changed, got: %s", nodeYAML)
		}
	})
}

func TestSetCommand_InvalidPath(t *testing.T) {
	t.Run("errors on non-existent field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", "nonexistent", "value", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for non-existent field, got nil")
		}
	})

	t.Run("errors on out-of-bounds array index", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", "tags[99]", "value", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for out-of-bounds index, got nil")
		}
	})
}

func TestSetCommand_InvalidNode(t *testing.T) {
	t.Run("errors on non-existent node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"nonexistent-999", "title", "New Title", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for non-existent node, got nil")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "not found") &&
			!strings.Contains(errMsg, "does not exist") {
			t.Errorf("Expected error about missing node, got %q", errMsg)
		}
	})
}

func TestSetCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", "title", "New Title", tmpDir})
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

func TestSetCommand_IncrementVersion(t *testing.T) {
	t.Run("increments version on set", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"sword-001", "title", "Modified Sword", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "version: 2") {
			t.Errorf("Expected version to be incremented to 2, got: %s", nodeYAML)
		}
	})
}

func TestSetCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		root := NewRootCommand()
		set := NewSetCommand()
		root.AddCommand(set)

		root.SetArgs([]string{"set", "sword-001", "title", "Root Test Sword", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "Root Test Sword") {
			t.Errorf("Expected title to be changed via root command, got: %s", nodeYAML)
		}
	})
}

func TestSetCommand_QuietFlag(t *testing.T) {
	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewSetCommand()
		flag := cmd.Flags().Lookup("quiet")
		if flag == nil {
			t.Fatal("Expected --quiet flag to be defined")
		}
		if flag.Shorthand != "q" {
			t.Errorf("Expected shorthand 'q', got %q", flag.Shorthand)
		}
	})

	t.Run("quiet flag suppresses output", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSet(t, tmpDir)

		cmd := NewSetCommand()
		cmd.SetArgs([]string{"--quiet", "sword-001", "title", "Quiet Sword", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with --quiet, got %v", err)
		}
	})
}

// Test helpers

func setupProjectForSet(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: set-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create a node to modify
	nodeYAML := `id: sword-001
kind: item
version: 1
status: draft
title: Iron Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
`
	nodePath := filepath.Join(nodesDir, "sword-001.yaml")
	if err := os.WriteFile(nodePath, []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
}

func readNodeFile(t *testing.T, dir, nodeID string) string {
	t.Helper()
	nodePath := filepath.Join(dir, ".deco", "nodes", nodeID+".yaml")
	content, err := os.ReadFile(nodePath)
	if err != nil {
		t.Fatalf("Failed to read node file: %v", err)
	}
	return string(content)
}
