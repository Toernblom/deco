package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMvCommand_Structure(t *testing.T) {
	t.Run("creates mv command", func(t *testing.T) {
		cmd := NewMvCommand()
		if cmd == nil {
			t.Fatal("Expected mv command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "mv") {
			t.Errorf("Expected Use to start with 'mv', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewMvCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("requires two arguments", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when no arguments provided")
		}
	})

	t.Run("requires new-id argument", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when only old ID provided")
		}
	})
}

func TestMvCommand_RenamesNode(t *testing.T) {
	t.Run("renames node ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "blade-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify new file exists
		newPath := filepath.Join(tmpDir, ".deco", "nodes", "blade-001.yaml")
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			t.Error("Expected new node file to exist")
		}

		// Verify old file is gone
		oldPath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
			t.Error("Expected old node file to be deleted")
		}

		// Verify ID field was updated
		nodeYAML := readMvNodeFile(t, tmpDir, "blade-001")
		if !strings.Contains(nodeYAML, "id: blade-001") {
			t.Errorf("Expected node ID to be updated to blade-001, got: %s", nodeYAML)
		}
	})

	t.Run("preserves other node fields", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "blade-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		nodeYAML := readMvNodeFile(t, tmpDir, "blade-001")
		if !strings.Contains(nodeYAML, "kind: item") {
			t.Errorf("Expected kind to be preserved, got: %s", nodeYAML)
		}
		if !strings.Contains(nodeYAML, "title: Iron Sword") {
			t.Errorf("Expected title to be preserved, got: %s", nodeYAML)
		}
	})
}

func TestMvCommand_UpdatesReferences(t *testing.T) {
	t.Run("updates Uses references in other nodes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithRefs(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "blade-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Check that hero-001's Uses reference was updated
		heroYAML := readMvNodeFile(t, tmpDir, "hero-001")
		if !strings.Contains(heroYAML, "blade-001") {
			t.Errorf("Expected hero's Uses reference to be updated to blade-001, got: %s", heroYAML)
		}
		if strings.Contains(heroYAML, "sword-001") {
			t.Errorf("Expected sword-001 reference to be replaced, got: %s", heroYAML)
		}
	})

	t.Run("updates Related references in other nodes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithRelated(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "blade-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Check that armor-001's Related reference was updated
		armorYAML := readMvNodeFile(t, tmpDir, "armor-001")
		if !strings.Contains(armorYAML, "blade-001") {
			t.Errorf("Expected armor's Related reference to be updated to blade-001, got: %s", armorYAML)
		}
	})

	t.Run("increments version on nodes with updated refs", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithRefs(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "blade-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Hero had version 1, should now be version 2
		heroYAML := readMvNodeFile(t, tmpDir, "hero-001")
		if !strings.Contains(heroYAML, "version: 2") {
			t.Errorf("Expected hero's version to be incremented to 2, got: %s", heroYAML)
		}
	})

	t.Run("does not modify nodes without refs to renamed node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithRefs(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "blade-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Potion has no refs to sword, should stay at version 1
		potionYAML := readMvNodeFile(t, tmpDir, "potion-001")
		if !strings.Contains(potionYAML, "version: 1") {
			t.Errorf("Expected potion's version to remain 1, got: %s", potionYAML)
		}
	})
}

func TestMvCommand_CreatesHistoryEntry(t *testing.T) {
	t.Run("records rename in history", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "blade-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Check history file
		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		content, err := os.ReadFile(historyPath)
		if err != nil {
			t.Fatalf("Failed to read history file: %v", err)
		}

		historyStr := string(content)
		// History should contain a move/rename operation
		if !strings.Contains(historyStr, "move") && !strings.Contains(historyStr, "rename") {
			t.Errorf("Expected history entry for move/rename, got: %s", historyStr)
		}
		// Should reference the new ID
		if !strings.Contains(historyStr, "blade-001") {
			t.Errorf("Expected history entry to contain new node ID, got: %s", historyStr)
		}
	})
}

func TestMvCommand_Errors(t *testing.T) {
	t.Run("errors on non-existent node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"nonexistent-999", "new-id", tmpDir})
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

	t.Run("errors when target ID already exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithRefs(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "hero-001", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when target ID exists, got nil")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "exists") &&
			!strings.Contains(errMsg, "already") {
			t.Errorf("Expected error about existing node, got %q", errMsg)
		}
	})

	t.Run("errors on empty old ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"", "new-id", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for empty old ID, got nil")
		}
	})

	t.Run("errors on empty new ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for empty new ID, got nil")
		}
	})

	t.Run("errors when old and new ID are the same", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "sword-001", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when IDs are the same, got nil")
		}
	})
}

func TestMvCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"sword-001", "blade-001", tmpDir})
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

func TestMvCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForMv(t, tmpDir)

		root := NewRootCommand()
		mv := NewMvCommand()
		root.AddCommand(mv)

		root.SetArgs([]string{"mv", "sword-001", "blade-001", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify rename happened
		newPath := filepath.Join(tmpDir, ".deco", "nodes", "blade-001.yaml")
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			t.Error("Expected new node file to exist")
		}
	})
}

func TestMvCommand_QuietFlag(t *testing.T) {
	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewMvCommand()
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
		setupProjectForMv(t, tmpDir)

		cmd := NewMvCommand()
		cmd.SetArgs([]string{"--quiet", "sword-001", "blade-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with --quiet, got %v", err)
		}
	})
}

// Test helpers

func setupProjectForMv(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: mv-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create empty history file
	historyPath := filepath.Join(decoDir, "history.jsonl")
	if err := os.WriteFile(historyPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create history.jsonl: %v", err)
	}

	// Create a node to rename
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

func setupProjectWithRefs(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: mv-refs-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create empty history file
	historyPath := filepath.Join(decoDir, "history.jsonl")
	if err := os.WriteFile(historyPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create history.jsonl: %v", err)
	}

	// Create sword node (to be renamed)
	swordYAML := `id: sword-001
kind: item
version: 1
status: draft
title: Iron Sword
summary: A basic iron sword
`
	swordPath := filepath.Join(nodesDir, "sword-001.yaml")
	if err := os.WriteFile(swordPath, []byte(swordYAML), 0644); err != nil {
		t.Fatalf("Failed to create sword node: %v", err)
	}

	// Create hero node that uses sword
	heroYAML := `id: hero-001
kind: character
version: 1
status: draft
title: Brave Hero
summary: A brave hero
refs:
  uses:
    - target: sword-001
      context: equipped weapon
`
	heroPath := filepath.Join(nodesDir, "hero-001.yaml")
	if err := os.WriteFile(heroPath, []byte(heroYAML), 0644); err != nil {
		t.Fatalf("Failed to create hero node: %v", err)
	}

	// Create potion node (no refs to sword)
	potionYAML := `id: potion-001
kind: item
version: 1
status: draft
title: Health Potion
summary: Restores health
`
	potionPath := filepath.Join(nodesDir, "potion-001.yaml")
	if err := os.WriteFile(potionPath, []byte(potionYAML), 0644); err != nil {
		t.Fatalf("Failed to create potion node: %v", err)
	}
}

func setupProjectWithRelated(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: mv-related-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create empty history file
	historyPath := filepath.Join(decoDir, "history.jsonl")
	if err := os.WriteFile(historyPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create history.jsonl: %v", err)
	}

	// Create sword node (to be renamed)
	swordYAML := `id: sword-001
kind: item
version: 1
status: draft
title: Iron Sword
summary: A basic iron sword
`
	swordPath := filepath.Join(nodesDir, "sword-001.yaml")
	if err := os.WriteFile(swordPath, []byte(swordYAML), 0644); err != nil {
		t.Fatalf("Failed to create sword node: %v", err)
	}

	// Create armor node with related ref to sword
	armorYAML := `id: armor-001
kind: item
version: 1
status: draft
title: Knight Armor
summary: Heavy armor for knights
refs:
  related:
    - target: sword-001
      context: part of knight equipment set
`
	armorPath := filepath.Join(nodesDir, "armor-001.yaml")
	if err := os.WriteFile(armorPath, []byte(armorYAML), 0644); err != nil {
		t.Fatalf("Failed to create armor node: %v", err)
	}
}

func readMvNodeFile(t *testing.T, dir, nodeID string) string {
	t.Helper()
	nodePath := filepath.Join(dir, ".deco", "nodes", nodeID+".yaml")
	content, err := os.ReadFile(nodePath)
	if err != nil {
		t.Fatalf("Failed to read node file: %v", err)
	}
	return string(content)
}
