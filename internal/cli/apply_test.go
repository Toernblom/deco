package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyCommand_Structure(t *testing.T) {
	t.Run("creates apply command", func(t *testing.T) {
		cmd := NewApplyCommand()
		if cmd == nil {
			t.Fatal("Expected apply command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "apply") {
			t.Errorf("Expected Use to start with 'apply', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewApplyCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("requires node ID and patch file", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when no arguments provided")
		}
	})
}

func TestApplyCommand_ApplyPatchFile(t *testing.T) {
	t.Run("applies single set operation", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "single-set.json", `[
			{"op": "set", "path": "title", "value": "Modified Sword"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "Modified Sword") {
			t.Errorf("Expected title to be modified, got: %s", nodeYAML)
		}
	})

	t.Run("applies multiple operations", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "multi-op.json", `[
			{"op": "set", "path": "title", "value": "Golden Sword"},
			{"op": "set", "path": "status", "value": "review"},
			{"op": "append", "path": "tags", "value": "legendary"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "Golden Sword") {
			t.Errorf("Expected title to be modified, got: %s", nodeYAML)
		}
		if !strings.Contains(nodeYAML, "status: review") {
			t.Errorf("Expected status to be modified, got: %s", nodeYAML)
		}
		if !strings.Contains(nodeYAML, "legendary") {
			t.Errorf("Expected tag to be appended, got: %s", nodeYAML)
		}
	})

	t.Run("applies unset operation", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "unset-op.json", `[
			{"op": "unset", "path": "summary"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if strings.Contains(nodeYAML, "summary: A basic iron sword") {
			t.Errorf("Expected summary to be unset, got: %s", nodeYAML)
		}
	})

	t.Run("writes content hash to history", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "hash-test.json", `[
			{"op": "set", "path": "title", "value": "Hash Test Sword"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify history contains content_hash
		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		history, err := os.ReadFile(historyPath)
		if err != nil {
			t.Fatalf("Failed to read history: %v", err)
		}
		if !strings.Contains(string(history), "content_hash") {
			t.Errorf("Expected history to contain content_hash, got: %s", string(history))
		}
	})
}

func TestApplyCommand_DryRun(t *testing.T) {
	t.Run("has dry-run flag", func(t *testing.T) {
		cmd := NewApplyCommand()
		flag := cmd.Flags().Lookup("dry-run")
		if flag == nil {
			t.Fatal("Expected --dry-run flag to be defined")
		}
	})

	t.Run("dry-run does not modify node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "dry-run.json", `[
			{"op": "set", "path": "title", "value": "Should Not Change"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"--dry-run", "sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if strings.Contains(nodeYAML, "Should Not Change") {
			t.Errorf("Expected dry-run to not modify node, got: %s", nodeYAML)
		}
		if !strings.Contains(nodeYAML, "Iron Sword") {
			t.Errorf("Expected original title to remain, got: %s", nodeYAML)
		}
	})
}

func TestApplyCommand_Rollback(t *testing.T) {
	t.Run("rolls back on error", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		// Second operation will fail (invalid field)
		patchFile := createPatchFile(t, tmpDir, "rollback.json", `[
			{"op": "set", "path": "title", "value": "Should Rollback"},
			{"op": "set", "path": "nonexistent", "value": "fail"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected error on invalid operation")
		}

		// Title should NOT be changed due to rollback
		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if strings.Contains(nodeYAML, "Should Rollback") {
			t.Errorf("Expected rollback to restore original title, got: %s", nodeYAML)
		}
		if !strings.Contains(nodeYAML, "Iron Sword") {
			t.Errorf("Expected original title to be preserved, got: %s", nodeYAML)
		}
	})
}

func TestApplyCommand_InvalidPatchFile(t *testing.T) {
	t.Run("errors on invalid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "invalid.json", `not valid json`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})

	t.Run("errors on missing patch file", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", "/nonexistent/patch.json", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for missing patch file, got nil")
		}
	})

	t.Run("errors on unknown operation", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "unknown-op.json", `[
			{"op": "invalid", "path": "title", "value": "test"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for unknown operation, got nil")
		}
	})
}

func TestApplyCommand_InvalidNode(t *testing.T) {
	t.Run("errors on non-existent node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "valid.json", `[
			{"op": "set", "path": "title", "value": "test"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"nonexistent-999", patchFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for non-existent node, got nil")
		}
	})
}

func TestApplyCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		patchFile := createPatchFile(t, tmpDir, "valid.json", `[
			{"op": "set", "path": "title", "value": "test"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for missing .deco directory, got nil")
		}
	})
}

func TestApplyCommand_IncrementVersion(t *testing.T) {
	t.Run("increments version on apply", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "version.json", `[
			{"op": "set", "path": "summary", "value": "Updated summary"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
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

func TestApplyCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "root.json", `[
			{"op": "set", "path": "title", "value": "Root Test Sword"}
		]`)

		root := NewRootCommand()
		apply := NewApplyCommand()
		root.AddCommand(apply)

		root.SetArgs([]string{"apply", "sword-001", patchFile, tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "Root Test Sword") {
			t.Errorf("Expected title to be modified via root command, got: %s", nodeYAML)
		}
	})
}

func TestApplyCommand_QuietFlag(t *testing.T) {
	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewApplyCommand()
		flag := cmd.Flags().Lookup("quiet")
		if flag == nil {
			t.Fatal("Expected --quiet flag to be defined")
		}
	})

	t.Run("quiet flag suppresses output", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "quiet.json", `[
			{"op": "set", "path": "title", "value": "Quiet Sword"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"--quiet", "sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with --quiet, got %v", err)
		}
	})
}

func TestApplyCommand_ValidationGate(t *testing.T) {
	t.Run("rejects patch that creates invalid node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		// Setting status to "published" without content should fail validation
		patchFile := createPatchFile(t, tmpDir, "invalid.json", `[
			{"op": "set", "path": "status", "value": "published"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected validation error for published status without content")
		}
		if !strings.Contains(err.Error(), "validation failed") {
			t.Errorf("Expected validation failure error, got: %v", err)
		}

		// Node should not have been modified
		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if strings.Contains(nodeYAML, "status: published") {
			t.Errorf("Node should not have been modified on validation failure")
		}
	})

	t.Run("rejects patch that removes required field", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		// Unsetting title should fail validation (required field)
		patchFile := createPatchFile(t, tmpDir, "unset-required.json", `[
			{"op": "unset", "path": "title"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected error for unsetting required field")
		}
	})

	t.Run("dry-run shows validation failure without modifying", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForApply(t, tmpDir)
		patchFile := createPatchFile(t, tmpDir, "invalid-dry.json", `[
			{"op": "set", "path": "status", "value": "published"}
		]`)

		cmd := NewApplyCommand()
		cmd.SetArgs([]string{"--dry-run", "sword-001", patchFile, tmpDir})
		err := cmd.Execute()
		// dry-run should succeed even if validation fails
		if err != nil {
			t.Fatalf("Expected dry-run to succeed, got %v", err)
		}

		// Node should not have been modified
		nodeYAML := readNodeFile(t, tmpDir, "sword-001")
		if strings.Contains(nodeYAML, "status: published") {
			t.Errorf("Dry run should not modify node")
		}
	})
}

// Test helpers

func setupProjectForApply(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: apply-test-project
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

func createPatchFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create patch file: %v", err)
	}
	return path
}
