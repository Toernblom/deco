package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendCommand_Structure(t *testing.T) {
	t.Run("creates append command", func(t *testing.T) {
		cmd := NewAppendCommand()
		if cmd == nil {
			t.Fatal("Expected append command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "append") {
			t.Errorf("Expected Use to start with 'append', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewAppendCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("requires three arguments", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when no arguments provided")
		}
	})

	t.Run("requires path and value", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"sword-001", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when only node ID provided")
		}
	})
}

func TestAppendCommand_AppendToArray(t *testing.T) {
	t.Run("appends single value to tags", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"sword-001", "tags", "legendary", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "legendary") {
			t.Errorf("Expected 'legendary' tag to be appended, got: %s", nodeYAML)
		}
	})

	t.Run("appends to existing tags", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"sword-001", "tags", "rare", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		// Original tags should still be there
		if !strings.Contains(nodeYAML, "weapon") {
			t.Errorf("Expected original 'weapon' tag to remain, got: %s", nodeYAML)
		}
		if !strings.Contains(nodeYAML, "rare") {
			t.Errorf("Expected 'rare' tag to be appended, got: %s", nodeYAML)
		}
	})

	t.Run("multiple appends accumulate", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		// First append
		cmd1 := NewAppendCommand()
		cmd1.SetArgs([]string{"sword-001", "tags", "rare", tmpDir})
		if err := cmd1.Execute(); err != nil {
			t.Fatalf("First append failed: %v", err)
		}

		// Second append
		cmd2 := NewAppendCommand()
		cmd2.SetArgs([]string{"sword-001", "tags", "magical", tmpDir})
		if err := cmd2.Execute(); err != nil {
			t.Fatalf("Second append failed: %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "rare") {
			t.Errorf("Expected 'rare' tag, got: %s", nodeYAML)
		}
		if !strings.Contains(nodeYAML, "magical") {
			t.Errorf("Expected 'magical' tag, got: %s", nodeYAML)
		}
	})
}

func TestAppendCommand_ErrorOnNonArray(t *testing.T) {
	t.Run("errors on non-array field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"sword-001", "title", "extra value", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when appending to non-array field, got nil")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "non-array") &&
			!strings.Contains(errMsg, "cannot append") &&
			!strings.Contains(errMsg, "not a slice") {
			t.Errorf("Expected error about non-array field, got %q", errMsg)
		}
	})
}

func TestAppendCommand_InvalidPath(t *testing.T) {
	t.Run("errors on non-existent field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"sword-001", "nonexistent", "value", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for non-existent field, got nil")
		}
	})
}

func TestAppendCommand_InvalidNode(t *testing.T) {
	t.Run("errors on non-existent node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"nonexistent-999", "tags", "new-tag", tmpDir})
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

func TestAppendCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"sword-001", "tags", "new-tag", tmpDir})
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

func TestAppendCommand_IncrementVersion(t *testing.T) {
	t.Run("increments version on append", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"sword-001", "tags", "new-tag", tmpDir})
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

func TestAppendCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForAppend(t, tmpDir)

		root := NewRootCommand()
		append := NewAppendCommand()
		root.AddCommand(append)

		root.SetArgs([]string{"append", "sword-001", "tags", "root-test-tag", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "root-test-tag") {
			t.Errorf("Expected tag to be appended via root command, got: %s", nodeYAML)
		}
	})
}

func TestAppendCommand_QuietFlag(t *testing.T) {
	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewAppendCommand()
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
		setupProjectForAppend(t, tmpDir)

		cmd := NewAppendCommand()
		cmd.SetArgs([]string{"--quiet", "sword-001", "tags", "quiet-tag", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with --quiet, got %v", err)
		}
	})
}

// Test helper

func setupProjectForAppend(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: append-test-project
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
