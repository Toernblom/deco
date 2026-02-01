package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGraphCommand_Structure(t *testing.T) {
	t.Run("creates graph command", func(t *testing.T) {
		cmd := NewGraphCommand()
		if cmd == nil {
			t.Fatal("Expected graph command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "graph") {
			t.Errorf("Expected Use to start with 'graph', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewGraphCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("has format flag", func(t *testing.T) {
		cmd := NewGraphCommand()
		flag := cmd.Flags().Lookup("format")
		if flag == nil {
			t.Fatal("Expected --format flag to be defined")
		}
		if flag.Shorthand != "f" {
			t.Errorf("Expected shorthand 'f', got %q", flag.Shorthand)
		}
		if flag.DefValue != "dot" {
			t.Errorf("Expected default 'dot', got %q", flag.DefValue)
		}
	})
}

func TestGraphCommand_DOTOutput(t *testing.T) {
	t.Run("outputs DOT format by default", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGraphProjectWithRefs(t, tmpDir)

		var buf bytes.Buffer
		cmd := NewGraphCommand()
		cmd.SetOut(&buf)
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("outputs DOT format with explicit flag", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGraphProjectWithRefs(t, tmpDir)

		cmd := NewGraphCommand()
		cmd.SetArgs([]string{"--format", "dot", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestGraphCommand_MermaidOutput(t *testing.T) {
	t.Run("outputs Mermaid format", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGraphProjectWithRefs(t, tmpDir)

		cmd := NewGraphCommand()
		cmd.SetArgs([]string{"--format", "mermaid", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Mermaid format short flag works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGraphProjectWithRefs(t, tmpDir)

		cmd := NewGraphCommand()
		cmd.SetArgs([]string{"-f", "mermaid", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -f, got %v", err)
		}
	})
}

func TestGraphCommand_InvalidFormat(t *testing.T) {
	t.Run("errors on unknown format", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGraphProjectWithRefs(t, tmpDir)

		cmd := NewGraphCommand()
		cmd.SetArgs([]string{"--format", "invalid", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Fatal("Expected error for invalid format, got nil")
		}
		if !strings.Contains(err.Error(), "unknown format") {
			t.Errorf("Expected 'unknown format' error, got %q", err.Error())
		}
	})
}

func TestGraphCommand_EmptyProject(t *testing.T) {
	t.Run("handles empty project", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupEmptyProject(t, tmpDir)

		cmd := NewGraphCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for empty project, got %v", err)
		}
	})
}

func TestGraphCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewGraphCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for missing .deco directory, got nil")
		}
	})
}

func TestGraphCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGraphProjectWithRefs(t, tmpDir)

		root := NewRootCommand()
		graph := NewGraphCommand()
		root.AddCommand(graph)

		root.SetArgs([]string{"graph", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

// Test helper for projects with refs (graph-specific)
func setupGraphProjectWithRefs(t *testing.T, dir string) {
	t.Helper()

	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	configYAML := `version: 1
project_name: ref-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	nodes := []struct {
		id   string
		yaml string
	}{
		{
			"core",
			`id: core
kind: system
version: 1
status: approved
title: Core System
refs:
  uses:
    - target: player
      context: "Uses player entity"
    - target: enemy
`,
		},
		{
			"player",
			`id: player
kind: entity
version: 1
status: approved
title: Player Entity
refs:
  related:
    - target: enemy
      context: "Can fight enemies"
`,
		},
		{
			"enemy",
			`id: enemy
kind: entity
version: 1
status: draft
title: Enemy Entity
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
