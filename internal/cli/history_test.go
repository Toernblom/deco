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
	"time"
)

func TestHistoryCommand_Structure(t *testing.T) {
	t.Run("creates history command", func(t *testing.T) {
		cmd := NewHistoryCommand()
		if cmd == nil {
			t.Fatal("Expected history command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "history") {
			t.Errorf("Expected Use to start with 'history', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewHistoryCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})
}

func TestHistoryCommand_ShowAllHistory(t *testing.T) {
	t.Run("shows history entries", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithHistory(t, tmpDir)

		cmd := NewHistoryCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("handles empty history", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithEmptyHistory(t, tmpDir)

		cmd := NewHistoryCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for empty history, got %v", err)
		}
	})
}

func TestHistoryCommand_NodeFilter(t *testing.T) {
	t.Run("has node flag", func(t *testing.T) {
		cmd := NewHistoryCommand()
		flag := cmd.Flags().Lookup("node")
		if flag == nil {
			t.Fatal("Expected --node flag to be defined")
		}
		if flag.Shorthand != "n" {
			t.Errorf("Expected shorthand 'n', got %q", flag.Shorthand)
		}
	})

	t.Run("filters history by node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithHistory(t, tmpDir)

		cmd := NewHistoryCommand()
		cmd.SetArgs([]string{"--node", "sword-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("node filter short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithHistory(t, tmpDir)

		cmd := NewHistoryCommand()
		cmd.SetArgs([]string{"-n", "hero-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -n, got %v", err)
		}
	})
}

func TestHistoryCommand_LimitFlag(t *testing.T) {
	t.Run("has limit flag", func(t *testing.T) {
		cmd := NewHistoryCommand()
		flag := cmd.Flags().Lookup("limit")
		if flag == nil {
			t.Fatal("Expected --limit flag to be defined")
		}
		if flag.Shorthand != "l" {
			t.Errorf("Expected shorthand 'l', got %q", flag.Shorthand)
		}
	})

	t.Run("limits number of entries", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithHistory(t, tmpDir)

		cmd := NewHistoryCommand()
		cmd.SetArgs([]string{"--limit", "2", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("limit flag short version works", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithHistory(t, tmpDir)

		cmd := NewHistoryCommand()
		cmd.SetArgs([]string{"-l", "1", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with -l, got %v", err)
		}
	})
}

func TestHistoryCommand_CombinedFilters(t *testing.T) {
	t.Run("combines node and limit filters", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithHistory(t, tmpDir)

		cmd := NewHistoryCommand()
		cmd.SetArgs([]string{"--node", "sword-001", "--limit", "5", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestHistoryCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewHistoryCommand()
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

func TestHistoryCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithHistory(t, tmpDir)

		root := NewRootCommand()
		history := NewHistoryCommand()
		root.AddCommand(history)

		root.SetArgs([]string{"history", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("integrates with root command with filters", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithHistory(t, tmpDir)

		root := NewRootCommand()
		history := NewHistoryCommand()
		root.AddCommand(history)

		root.SetArgs([]string{"history", "--node", "sword-001", "--limit", "10", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

// Test helpers

func setupProjectWithHistory(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: history-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create history.jsonl with some sample entries
	now := time.Now()
	historyContent := ""
	historyContent += `{"timestamp":"` + now.Add(-3*time.Hour).Format(time.RFC3339) + `","node_id":"sword-001","operation":"create","user":"alice"}` + "\n"
	historyContent += `{"timestamp":"` + now.Add(-2*time.Hour).Format(time.RFC3339) + `","node_id":"hero-001","operation":"create","user":"bob"}` + "\n"
	historyContent += `{"timestamp":"` + now.Add(-1*time.Hour).Format(time.RFC3339) + `","node_id":"sword-001","operation":"update","user":"alice"}` + "\n"
	historyContent += `{"timestamp":"` + now.Format(time.RFC3339) + `","node_id":"hero-001","operation":"update","user":"bob"}` + "\n"

	historyPath := filepath.Join(decoDir, "history.jsonl")
	if err := os.WriteFile(historyPath, []byte(historyContent), 0644); err != nil {
		t.Fatalf("Failed to create history.jsonl: %v", err)
	}
}

func setupProjectWithEmptyHistory(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: empty-history-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create empty history.jsonl
	historyPath := filepath.Join(decoDir, "history.jsonl")
	if err := os.WriteFile(historyPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create history.jsonl: %v", err)
	}
}
