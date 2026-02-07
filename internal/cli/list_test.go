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

func TestListCommand_Structure(t *testing.T) {
	t.Run("creates list command", func(t *testing.T) {
		cmd := NewListCommand()
		if cmd == nil {
			t.Fatal("Expected list command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "list") {
			t.Errorf("Expected Use to start with 'list', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewListCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})
}

func TestListCommand_ListAllNodes(t *testing.T) {
	t.Run("lists all nodes in project", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("handles empty project", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupEmptyProject(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for empty project, got %v", err)
		}
	})
}

func TestListCommand_KindFilter(t *testing.T) {
	t.Run("has kind flag", func(t *testing.T) {
		cmd := NewListCommand()
		flag := cmd.Flags().Lookup("kind")
		if flag == nil {
			t.Fatal("Expected --kind flag to be defined")
		}
		if flag.Shorthand != "k" {
			t.Errorf("Expected shorthand 'k', got %q", flag.Shorthand)
		}
	})

	t.Run("filters nodes by kind", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{"--kind", "character", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("kind filter short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{"-k", "item", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -k, got %v", err)
		}
	})
}

func TestListCommand_StatusFilter(t *testing.T) {
	t.Run("has status flag", func(t *testing.T) {
		cmd := NewListCommand()
		flag := cmd.Flags().Lookup("status")
		if flag == nil {
			t.Fatal("Expected --status flag to be defined")
		}
		if flag.Shorthand != "s" {
			t.Errorf("Expected shorthand 's', got %q", flag.Shorthand)
		}
	})

	t.Run("filters nodes by status", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{"--status", "draft", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("status filter short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{"-s", "approved", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -s, got %v", err)
		}
	})
}

func TestListCommand_TagFilter(t *testing.T) {
	t.Run("has tag flag", func(t *testing.T) {
		cmd := NewListCommand()
		flag := cmd.Flags().Lookup("tag")
		if flag == nil {
			t.Fatal("Expected --tag flag to be defined")
		}
		if flag.Shorthand != "t" {
			t.Errorf("Expected shorthand 't', got %q", flag.Shorthand)
		}
	})

	t.Run("filters nodes by tag", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{"--tag", "combat", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("tag filter short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{"-t", "test", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -t, got %v", err)
		}
	})
}

func TestListCommand_CombinedFilters(t *testing.T) {
	t.Run("combines kind and status filters", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{"--kind", "item", "--status", "draft", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("combines all filters", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		cmd := NewListCommand()
		cmd.SetArgs([]string{"--kind", "item", "--status", "draft", "--tag", "combat", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestListCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewListCommand()
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

func TestListCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		root := NewRootCommand()
		list := NewListCommand()
		root.AddCommand(list)

		root.SetArgs([]string{"list", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("integrates with root command with filters", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithMultipleNodes(t, tmpDir)

		root := NewRootCommand()
		list := NewListCommand()
		root.AddCommand(list)

		root.SetArgs([]string{"list", "--kind", "item", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

// Test helper

func setupProjectWithMultipleNodes(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: multi-node-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create multiple nodes with different kinds, statuses, and tags
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
tags:
  - protagonist
  - combat
`,
		},
		{
			"potion-001",
			`id: potion-001
kind: item
version: 1
status: approved
title: Health Potion
tags:
  - consumable
  - healing
`,
		},
		{
			"quest-001",
			`id: quest-001
kind: quest
version: 1
status: draft
title: Main Quest
tags:
  - story
  - test
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
