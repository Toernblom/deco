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

func TestDiffCommand_Structure(t *testing.T) {
	t.Run("creates diff command", func(t *testing.T) {
		cmd := NewDiffCommand()
		if cmd == nil {
			t.Fatal("Expected diff command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "diff") {
			t.Errorf("Expected Use to start with 'diff', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewDiffCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})
}

func TestDiffCommand_RequiresNodeID(t *testing.T) {
	t.Run("errors without node ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{tmpDir})
		err := cmd.Execute()
		// cobra should error because node ID is required
		if err == nil {
			t.Error("Expected error when node ID not provided")
		}
	})
}

func TestDiffCommand_ShowsNodeHistory(t *testing.T) {
	t.Run("shows changes for node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{"sword-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("handles node with no history", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{"nonexistent-001", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for node with no history, got %v", err)
		}
	})
}

func TestDiffCommand_LastFlag(t *testing.T) {
	t.Run("has last flag", func(t *testing.T) {
		cmd := NewDiffCommand()
		flag := cmd.Flags().Lookup("last")
		if flag == nil {
			t.Fatal("Expected --last flag to be defined")
		}
	})

	t.Run("limits to last N changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{"sword-001", "--last", "1", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestDiffCommand_SinceFlag(t *testing.T) {
	t.Run("has since flag", func(t *testing.T) {
		cmd := NewDiffCommand()
		flag := cmd.Flags().Lookup("since")
		if flag == nil {
			t.Fatal("Expected --since flag to be defined")
		}
	})

	t.Run("accepts date format", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{"sword-001", "--since", "2024-01-01", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with date format, got %v", err)
		}
	})

	t.Run("accepts RFC3339 format", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{"sword-001", "--since", "2024-01-01T00:00:00Z", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with RFC3339 format, got %v", err)
		}
	})

	t.Run("accepts relative duration", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{"sword-001", "--since", "2h", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error with relative duration, got %v", err)
		}
	})

	t.Run("errors on invalid since format", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{"sword-001", "--since", "invalid", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for invalid --since format")
		}
	})
}

func TestDiffCommand_CombinedFilters(t *testing.T) {
	t.Run("combines since and last filters", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{"sword-001", "--since", "1d", "--last", "2", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestDiffCommand_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := NewDiffCommand()
		cmd.SetArgs([]string{"sword-001", tmpDir})
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

func TestDiffCommand_WithRootCommand(t *testing.T) {
	t.Run("integrates with root command", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithDetailedHistory(t, tmpDir)

		root := NewRootCommand()
		diff := NewDiffCommand()
		root.AddCommand(diff)

		root.SetArgs([]string{"diff", "sword-001", tmpDir})
		err := root.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestParseSince(t *testing.T) {
	t.Run("parses RFC3339", func(t *testing.T) {
		result, err := parseSince("2024-06-15T10:30:00Z")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result.Year() != 2024 || result.Month() != 6 || result.Day() != 15 {
			t.Errorf("Unexpected parsed date: %v", result)
		}
	})

	t.Run("parses date only", func(t *testing.T) {
		result, err := parseSince("2024-06-15")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result.Year() != 2024 || result.Month() != 6 || result.Day() != 15 {
			t.Errorf("Unexpected parsed date: %v", result)
		}
	})

	t.Run("parses hours", func(t *testing.T) {
		before := time.Now()
		result, err := parseSince("2h")
		after := time.Now()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		// Should be roughly 2 hours ago
		expected := before.Add(-2 * time.Hour)
		if result.Before(expected.Add(-time.Minute)) || result.After(after.Add(-2*time.Hour+time.Minute)) {
			t.Errorf("Result %v not within expected range", result)
		}
	})

	t.Run("parses days", func(t *testing.T) {
		before := time.Now()
		result, err := parseSince("1d")
		after := time.Now()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expected := before.Add(-24 * time.Hour)
		if result.Before(expected.Add(-time.Minute)) || result.After(after.Add(-24*time.Hour+time.Minute)) {
			t.Errorf("Result %v not within expected range", result)
		}
	})

	t.Run("parses weeks", func(t *testing.T) {
		before := time.Now()
		result, err := parseSince("1w")
		after := time.Now()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expected := before.Add(-7 * 24 * time.Hour)
		if result.Before(expected.Add(-time.Minute)) || result.After(after.Add(-7*24*time.Hour+time.Minute)) {
			t.Errorf("Result %v not within expected range", result)
		}
	})

	t.Run("parses minutes", func(t *testing.T) {
		before := time.Now()
		result, err := parseSince("30m")
		after := time.Now()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expected := before.Add(-30 * time.Minute)
		if result.Before(expected.Add(-time.Minute)) || result.After(after.Add(-30*time.Minute+time.Minute)) {
			t.Errorf("Result %v not within expected range", result)
		}
	})

	t.Run("errors on invalid format", func(t *testing.T) {
		_, err := parseSince("invalid")
		if err == nil {
			t.Error("Expected error for invalid format")
		}
	})

	t.Run("errors on unknown unit", func(t *testing.T) {
		_, err := parseSince("5x")
		if err == nil {
			t.Error("Expected error for unknown unit")
		}
	})
}

// Test helper

func setupProjectWithDetailedHistory(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: diff-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create history.jsonl with detailed before/after data
	now := time.Now()
	historyContent := ""

	// sword-001: create
	historyContent += `{"timestamp":"` + now.Add(-3*time.Hour).Format(time.RFC3339) + `","node_id":"sword-001","operation":"create","user":"alice","after":{"id":"sword-001","kind":"item","title":"Iron Sword","content":{"damage":10}}}` + "\n"

	// hero-001: create
	historyContent += `{"timestamp":"` + now.Add(-2*time.Hour).Format(time.RFC3339) + `","node_id":"hero-001","operation":"create","user":"bob","after":{"id":"hero-001","kind":"character","title":"Hero"}}` + "\n"

	// sword-001: update (damage changed)
	historyContent += `{"timestamp":"` + now.Add(-1*time.Hour).Format(time.RFC3339) + `","node_id":"sword-001","operation":"update","user":"alice","before":{"content":{"damage":10}},"after":{"content":{"damage":15}}}` + "\n"

	// sword-001: set (add enchantment)
	historyContent += `{"timestamp":"` + now.Add(-30*time.Minute).Format(time.RFC3339) + `","node_id":"sword-001","operation":"set","user":"alice","before":{},"after":{"content":{"enchantment":"fire"}}}` + "\n"

	// hero-001: update
	historyContent += `{"timestamp":"` + now.Format(time.RFC3339) + `","node_id":"hero-001","operation":"update","user":"bob","before":{"title":"Hero"},"after":{"title":"Brave Hero"}}` + "\n"

	historyPath := filepath.Join(decoDir, "history.jsonl")
	if err := os.WriteFile(historyPath, []byte(historyContent), 0644); err != nil {
		t.Fatalf("Failed to create history.jsonl: %v", err)
	}
}
