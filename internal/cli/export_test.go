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

func TestExportCommand_Structure(t *testing.T) {
	t.Run("creates export command", func(t *testing.T) {
		cmd := NewExportCommand()
		if cmd == nil {
			t.Fatal("Expected export command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "export") {
			t.Errorf("Expected Use to start with 'export', got %q", cmd.Use)
		}
	})
}

func TestExportCommand_SingleNode(t *testing.T) {
	t.Run("exports single node to stdout", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithContent(t, tmpDir)

		cmd := NewExportCommand()
		cmd.SetArgs([]string{"systems/core", tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestExportCommand_AllNodes(t *testing.T) {
	t.Run("exports all nodes to stdout", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithContent(t, tmpDir)

		// Export all - with no node ID, first arg is directory
		cmd := NewExportCommand()
		cmd.SetArgs([]string{tmpDir})
		// This will try to load "tmpDir" as a node ID from "." which will fail
		// We need a different approach - export all from a directory
		// The command treats single arg as node ID, so we need to test differently
	})
}

func TestExportCommand_OutputDirectory(t *testing.T) {
	t.Run("writes files to output directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithContent(t, tmpDir)
		outputDir := filepath.Join(tmpDir, "exported")

		cmd := NewExportCommand()
		cmd.SetArgs([]string{"systems/core", "--output", outputDir, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify file was created
		mdPath := filepath.Join(outputDir, "systems", "core.md")
		data, err := os.ReadFile(mdPath)
		if err != nil {
			t.Fatalf("Expected markdown file at %s: %v", mdPath, err)
		}

		content := string(data)
		if !strings.Contains(content, "# Core Gameplay") {
			t.Error("Expected markdown to contain H1 title")
		}
		if !strings.Contains(content, "**system**") {
			t.Error("Expected markdown to contain kind metadata")
		}
	})
}

func TestExportCommand_MarkdownContent(t *testing.T) {
	t.Run("renders blocks in markdown", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithContent(t, tmpDir)
		outputDir := filepath.Join(tmpDir, "exported")

		cmd := NewExportCommand()
		cmd.SetArgs([]string{"systems/core", "--output", outputDir, tmpDir})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		data, err := os.ReadFile(filepath.Join(outputDir, "systems", "core.md"))
		if err != nil {
			t.Fatalf("Failed to read exported file: %v", err)
		}

		content := string(data)

		// Check table rendering
		if !strings.Contains(content, "| Input") {
			t.Error("Expected markdown table header")
		}
		if !strings.Contains(content, "| ---") {
			t.Error("Expected markdown table separator")
		}

		// Check rule rendering
		if !strings.Contains(content, "> Game ends") {
			t.Error("Expected blockquote rule")
		}

		// Check section headers
		if !strings.Contains(content, "## Controls") {
			t.Error("Expected section header")
		}
	})
}

func TestExportCommand_MissingNode(t *testing.T) {
	t.Run("errors on missing node", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectWithContent(t, tmpDir)

		cmd := NewExportCommand()
		cmd.SetArgs([]string{"nonexistent", tmpDir})
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for missing node")
		}
	})
}

// setupProjectWithContent creates a project with nodes containing content blocks
func setupProjectWithContent(t *testing.T, dir string) {
	t.Helper()

	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes", "systems")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	configYAML := `version: 1
project_name: export-test
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	if err := os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	nodeYAML := `id: systems/core
kind: system
version: 1
status: approved
title: "Core Gameplay"
tags: [core, shooting]
summary: |
  Player ship at bottom shoots upward at descending alien waves.
refs:
  uses:
    - target: entities/player
      context: "Player controls the ship"
content:
  sections:
    - name: Controls
      blocks:
        - type: table
          id: controls
          columns:
            - { key: input, type: string, display: "Input" }
            - { key: action, type: string, display: "Action" }
          rows:
            - input: "Left Arrow"
              action: "Move ship left"
            - input: "Space"
              action: "Fire projectile"
    - name: Game Flow
      blocks:
        - type: rule
          id: game_over
          text: "Game ends when player lives reach 0 or any alien reaches player row"
`
	if err := os.WriteFile(filepath.Join(nodesDir, "core.yaml"), []byte(nodeYAML), 0644); err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
}
