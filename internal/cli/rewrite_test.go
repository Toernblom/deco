package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRewriteCommand_Structure(t *testing.T) {
	t.Run("creates rewrite command", func(t *testing.T) {
		cmd := NewRewriteCommand()
		if cmd == nil {
			t.Fatal("Expected rewrite command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "rewrite") {
			t.Errorf("Expected Use to start with 'rewrite', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewRewriteCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("has dry-run flag", func(t *testing.T) {
		cmd := NewRewriteCommand()
		flag := cmd.Flags().Lookup("dry-run")
		if flag == nil {
			t.Fatal("Expected --dry-run flag to be defined")
		}
	})
}

func TestRewriteCommand_Rewrite(t *testing.T) {
	t.Run("rewrites node from YAML file", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForRewrite(t, tmpDir)

		// Create new node content file
		newContent := `id: sword-001
kind: item
version: 2
status: draft
title: Enchanted Blade
summary: A blade imbued with fire magic
tags:
  - weapon
  - combat
  - magic
`
		inputFile := filepath.Join(tmpDir, "new-sword.yaml")
		if err := os.WriteFile(inputFile, []byte(newContent), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		cmd := NewRewriteCommand()
		cmd.SetArgs([]string{"sword-001", inputFile, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileRewrite(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "Enchanted Blade") {
			t.Errorf("Expected title to be changed, got: %s", nodeYAML)
		}
		if !strings.Contains(nodeYAML, "magic") {
			t.Errorf("Expected magic tag, got: %s", nodeYAML)
		}
	})

	t.Run("dry-run shows diff without modifying", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForRewrite(t, tmpDir)

		newContent := `id: sword-001
kind: item
version: 2
status: draft
title: Should Not Apply
`
		inputFile := filepath.Join(tmpDir, "dry-run.yaml")
		if err := os.WriteFile(inputFile, []byte(newContent), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		cmd := NewRewriteCommand()
		cmd.SetArgs([]string{"--dry-run", "sword-001", inputFile, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileRewrite(t, tmpDir, "sword-001")
		if strings.Contains(nodeYAML, "Should Not Apply") {
			t.Errorf("Dry run should not modify node")
		}
		if !strings.Contains(nodeYAML, "Iron Sword") {
			t.Errorf("Expected original title, got: %s", nodeYAML)
		}
	})
}

func TestRewriteCommand_Validation(t *testing.T) {
	t.Run("rejects rewrite with invalid node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForRewrite(t, tmpDir)

		// Missing required title field
		newContent := `id: sword-001
kind: item
version: 2
status: draft
`
		inputFile := filepath.Join(tmpDir, "invalid.yaml")
		if err := os.WriteFile(inputFile, []byte(newContent), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		cmd := NewRewriteCommand()
		cmd.SetArgs([]string{"sword-001", inputFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected validation error for missing title")
		}
		if !strings.Contains(err.Error(), "validation failed") {
			t.Errorf("Expected validation failure, got: %v", err)
		}

		// Original node should be unchanged
		nodeYAML := readNodeFileRewrite(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "Iron Sword") {
			t.Errorf("Original node should be unchanged")
		}
	})

	t.Run("rejects mismatched node ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForRewrite(t, tmpDir)

		newContent := `id: wrong-id
kind: item
version: 1
status: draft
title: Wrong ID Node
`
		inputFile := filepath.Join(tmpDir, "wrong-id.yaml")
		if err := os.WriteFile(inputFile, []byte(newContent), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		cmd := NewRewriteCommand()
		cmd.SetArgs([]string{"sword-001", inputFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected error for mismatched node ID")
		}
		if !strings.Contains(err.Error(), "does not match") {
			t.Errorf("Expected ID mismatch error, got: %v", err)
		}
	})

	t.Run("rejects published status without content", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForRewrite(t, tmpDir)

		newContent := `id: sword-001
kind: item
version: 2
status: published
title: Published Without Content
`
		inputFile := filepath.Join(tmpDir, "no-content.yaml")
		if err := os.WriteFile(inputFile, []byte(newContent), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		cmd := NewRewriteCommand()
		cmd.SetArgs([]string{"sword-001", inputFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected validation error for published without content")
		}
	})
}

func TestRewriteCommand_Errors(t *testing.T) {
	t.Run("errors on missing input file", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForRewrite(t, tmpDir)

		cmd := NewRewriteCommand()
		cmd.SetArgs([]string{"sword-001", "/nonexistent/file.yaml", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for missing input file")
		}
	})

	t.Run("errors on invalid YAML", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForRewrite(t, tmpDir)

		inputFile := filepath.Join(tmpDir, "invalid.yaml")
		if err := os.WriteFile(inputFile, []byte("not: valid: yaml:"), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		cmd := NewRewriteCommand()
		cmd.SetArgs([]string{"sword-001", inputFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for invalid YAML")
		}
	})

	t.Run("errors on nonexistent node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForRewrite(t, tmpDir)

		newContent := `id: nonexistent-001
kind: item
version: 1
status: draft
title: New Node
`
		inputFile := filepath.Join(tmpDir, "new.yaml")
		if err := os.WriteFile(inputFile, []byte(newContent), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		cmd := NewRewriteCommand()
		cmd.SetArgs([]string{"nonexistent-001", inputFile, tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for nonexistent node")
		}
	})
}

func TestRewriteCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForRewrite(t, tmpDir)

		newContent := `id: sword-001
kind: item
version: 2
status: draft
title: Root Command Sword
`
		inputFile := filepath.Join(tmpDir, "root-test.yaml")
		if err := os.WriteFile(inputFile, []byte(newContent), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		root := NewRootCommand()
		rewrite := NewRewriteCommand()
		root.AddCommand(rewrite)

		root.SetArgs([]string{"rewrite", "sword-001", inputFile, tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readNodeFileRewrite(t, tmpDir, "sword-001")
		if !strings.Contains(nodeYAML, "Root Command Sword") {
			t.Errorf("Expected title change via root command, got: %s", nodeYAML)
		}
	})
}

// Test helpers

func setupProjectForRewrite(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: rewrite-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create a node to rewrite
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

func readNodeFileRewrite(t *testing.T, dir, nodeID string) string {
	t.Helper()
	nodePath := filepath.Join(dir, ".deco", "nodes", nodeID+".yaml")
	data, err := os.ReadFile(nodePath)
	if err != nil {
		t.Fatalf("Failed to read node file: %v", err)
	}
	return string(data)
}
