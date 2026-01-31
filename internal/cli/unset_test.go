package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnsetCommand_Structure(t *testing.T) {
	t.Run("creates unset command", func(t *testing.T) {
		cmd := NewUnsetCommand()
		if cmd == nil {
			t.Fatal("Expected unset command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "unset") {
			t.Errorf("Expected Use to start with 'unset', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewUnsetCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("requires two arguments", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when no arguments provided")
		}
	})
}

func TestUnsetCommand_UnsetOptionalField(t *testing.T) {
	t.Run("unsets summary field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "summary", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		// Summary should be empty or gone
		if strings.Contains(nodeYAML, "summary: A basic iron sword") {
			t.Errorf("Expected summary to be unset, got: %s", nodeYAML)
		}
	})

	t.Run("unsets tags field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "tags", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		// Tags should be empty
		if strings.Contains(nodeYAML, "- weapon") {
			t.Errorf("Expected tags to be unset, got: %s", nodeYAML)
		}
	})
}

func TestUnsetCommand_UnsetArrayElement(t *testing.T) {
	t.Run("unsets single array element", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "tags[0]", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		// First tag (weapon) should be removed
		if strings.Contains(nodeYAML, "weapon") {
			t.Errorf("Expected first tag to be unset, got: %s", nodeYAML)
		}
		// Second tag (combat) should remain
		if !strings.Contains(nodeYAML, "combat") {
			t.Errorf("Expected remaining tags to persist, got: %s", nodeYAML)
		}
	})
}

func TestUnsetCommand_ErrorOnRequiredField(t *testing.T) {
	t.Run("errors on unset id", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "id", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when unsetting required field 'id', got nil")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "required") {
			t.Errorf("Expected error about required field, got %q", errMsg)
		}
	})

	t.Run("errors on unset kind", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "kind", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when unsetting required field 'kind', got nil")
		}
	})

	t.Run("errors on unset version", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "version", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when unsetting required field 'version', got nil")
		}
	})

	t.Run("errors on unset status", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "status", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when unsetting required field 'status', got nil")
		}
	})

	t.Run("errors on unset title", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "title", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when unsetting required field 'title', got nil")
		}
	})
}

func TestUnsetCommand_InvalidPath(t *testing.T) {
	t.Run("errors on non-existent field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "nonexistent", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for non-existent field, got nil")
		}
	})
}

func TestUnsetCommand_InvalidNode(t *testing.T) {
	t.Run("errors on non-existent node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"nonexistent-999", "summary", tmpDir})
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

func TestUnsetCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "summary", tmpDir})
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

func TestUnsetCommand_IncrementVersion(t *testing.T) {
	t.Run("increments version on unset", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"sword-001", "summary", tmpDir})
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

func TestUnsetCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForUnset(t, tmpDir)

		root := NewRootCommand()
		unset := NewUnsetCommand()
		root.AddCommand(unset)

		root.SetArgs([]string{"unset", "sword-001", "summary", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if strings.Contains(nodeYAML, "summary: A basic iron sword") {
			t.Errorf("Expected summary to be unset via root command, got: %s", nodeYAML)
		}
	})
}

func TestUnsetCommand_QuietFlag(t *testing.T) {
	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewUnsetCommand()
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
		setupProjectForUnset(t, tmpDir)

		cmd := NewUnsetCommand()
		cmd.SetArgs([]string{"--quiet", "sword-001", "summary", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with --quiet, got %v", err)
		}
	})
}

// Test helper

func setupProjectForUnset(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: unset-test-project
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
