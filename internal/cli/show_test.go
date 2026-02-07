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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShowCommand_Structure(t *testing.T) {
	t.Run("creates show command", func(t *testing.T) {
		cmd := NewShowCommand()
		if cmd == nil {
			t.Fatal("Expected show command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "show") {
			t.Errorf("Expected Use to start with 'show', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewShowCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("requires node ID argument", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		cmd := NewShowCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when no node ID provided")
		}
	})
}

func TestShowCommand_DisplayNode(t *testing.T) {
	t.Run("displays node details", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		cmd := NewShowCommand()
		cmd.SetArgs([]string{"sword-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("displays node with all fields", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		cmd := NewShowCommand()
		cmd.SetArgs([]string{"quest-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestShowCommand_ReverseReferences(t *testing.T) {
	t.Run("shows reverse references when node is referenced", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		cmd := NewShowCommand()
		cmd.SetArgs([]string{"sword-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("shows no reverse references for unreferenced node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		cmd := NewShowCommand()
		cmd.SetArgs([]string{"hero-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestShowCommand_MissingNode(t *testing.T) {
	t.Run("errors on missing node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		cmd := NewShowCommand()
		cmd.SetArgs([]string{"nonexistent-999", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for missing node, got nil")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "not found") &&
			!strings.Contains(errMsg, "does not exist") {
			t.Errorf("Expected error about missing node, got %q", errMsg)
		}
	})
}

func TestShowCommand_JSONOutput(t *testing.T) {
	t.Run("has json flag", func(t *testing.T) {
		cmd := NewShowCommand()
		flag := cmd.Flags().Lookup("json")
		if flag == nil {
			t.Fatal("Expected --json flag to be defined")
		}
		if flag.Shorthand != "j" {
			t.Errorf("Expected shorthand 'j', got %q", flag.Shorthand)
		}
		if flag.DefValue != "false" {
			t.Errorf("Expected default 'false', got %q", flag.DefValue)
		}
	})

	t.Run("outputs JSON with --json flag", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		cmd := NewShowCommand()
		cmd.SetArgs([]string{"--json", "sword-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with --json, got %v", err)
		}
	})

	t.Run("json flag short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		cmd := NewShowCommand()
		cmd.SetArgs([]string{"-j", "hero-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -j, got %v", err)
		}
	})
}

func TestShowCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewShowCommand()
		cmd.SetArgs([]string{"any-node", tmpDir})
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

func TestShowCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		root := NewRootCommand()
		show := NewShowCommand()
		root.AddCommand(show)

		root.SetArgs([]string{"show", "sword-001", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("integrates with root command with --json", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithReferences(t, tmpDir)

		root := NewRootCommand()
		show := NewShowCommand()
		root.AddCommand(show)

		root.SetArgs([]string{"show", "--json", "hero-001", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

// Test helper

func setupProjectWithReferences(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: ref-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create nodes with references
	nodes := []struct {
		id   string
		yaml string
	}{
		{
			"sword-001",
			`id: sword-001
kind: item
version: 1
status: draft
title: Iron Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
`,
		},
		{
			"hero-001",
			`id: hero-001
kind: character
version: 1
status: approved
title: Hero Character
summary: The main protagonist
tags:
  - protagonist
refs:
  uses:
    - target: sword-001
      context: Equipped weapon
`,
		},
		{
			"quest-001",
			`id: quest-001
kind: quest
version: 1
status: draft
title: Defeat the Dragon
summary: Find the dragon and defeat it
tags:
  - main-story
  - combat
refs:
  uses:
    - target: sword-001
      context: Required item
  related:
    - target: hero-001
      context: Quest giver
`,
		},
	}

	for _, node := range nodes {
		nodePath := filepath.Join(nodesDir, node.id+".yaml")
		if err := os.WriteFile(nodePath, []byte(node.yaml), 0644); err != nil {
			t.Fatalf("Failed to create node %s: %v", node.id, err)
		}
	}
}
