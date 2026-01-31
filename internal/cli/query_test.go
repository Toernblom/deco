package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestQueryCommand_Structure(t *testing.T) {
	t.Run("creates query command", func(t *testing.T) {
		cmd := NewQueryCommand()
		if cmd == nil {
			t.Fatal("Expected query command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "query") {
			t.Errorf("Expected Use to start with 'query', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewQueryCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})
}

func TestQueryCommand_TextSearch(t *testing.T) {
	t.Run("searches nodes by title", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"sword", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("searches nodes by summary", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"protagonist", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("search is case-insensitive", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"SWORD", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("handles no search results", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"nonexistent-term-xyz", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error for no results, got %v", err)
		}
	})
}

func TestQueryCommand_KindFilter(t *testing.T) {
	t.Run("has kind flag", func(t *testing.T) {
		cmd := NewQueryCommand()
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
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"--kind", "item", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("kind filter short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"-k", "character", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -k, got %v", err)
		}
	})
}

func TestQueryCommand_StatusFilter(t *testing.T) {
	t.Run("has status flag", func(t *testing.T) {
		cmd := NewQueryCommand()
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
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"--status", "draft", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("status filter short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"-s", "published", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -s, got %v", err)
		}
	})
}

func TestQueryCommand_TagFilter(t *testing.T) {
	t.Run("has tag flag", func(t *testing.T) {
		cmd := NewQueryCommand()
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
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"--tag", "combat", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("tag filter short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"-t", "weapon", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -t, got %v", err)
		}
	})
}

func TestQueryCommand_CombinedSearchAndFilters(t *testing.T) {
	t.Run("combines search with kind filter", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"sword", "--kind", "item", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("combines search with status filter", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"hero", "--status", "published", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("combines search with tag filter", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"dragon", "--tag", "main-story", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("combines search with all filters", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"sword", "--kind", "item", "--status", "draft", "--tag", "weapon", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("no search term lists all with filters", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"--kind", "item", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestQueryCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"search-term", tmpDir})
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

func TestQueryCommand_EmptyProject(t *testing.T) {
	t.Run("handles empty project", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupEmptyProject(t, tmpDir)

		cmd := NewQueryCommand()
		cmd.SetArgs([]string{"search-term", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for empty project, got %v", err)
		}
	})
}

func TestQueryCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		root := NewRootCommand()
		query := NewQueryCommand()
		root.AddCommand(query)

		root.SetArgs([]string{"query", "sword", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("integrates with root command with filters", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		root := NewRootCommand()
		query := NewQueryCommand()
		root.AddCommand(query)

		root.SetArgs([]string{"query", "--kind", "item", "sword", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("integrates with root command no search term", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForQuery(t, tmpDir)

		root := NewRootCommand()
		query := NewQueryCommand()
		root.AddCommand(query)

		root.SetArgs([]string{"query", "--status", "draft", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

// Test helper

func setupProjectForQuery(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: query-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create multiple nodes with varying titles, summaries, kinds, statuses, and tags
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
summary: A basic iron sword for beginners
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
status: published
title: Hero Character
summary: The main protagonist of the story
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
status: published
title: Health Potion
summary: Restores health when consumed
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
title: Defeat the Dragon
summary: Find the dragon in its lair and defeat it
tags:
  - main-story
  - combat
`,
		},
		{
			"npc-001",
			`id: npc-001
kind: character
version: 1
status: draft
title: Village Elder
summary: The wise elder who gives quests
tags:
  - npc
  - quest-giver
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
