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
	"bytes"
	"io"
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

// captureStdout captures stdout output from fn and returns it as a string.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// setupProjectWithRefs creates a project with 3 nodes forming a ref chain:
// systems/core uses systems/scoring, systems/scoring uses items/weapon
func setupProjectWithRefs(t *testing.T, dir string) {
	t.Helper()

	decoDir := filepath.Join(dir, ".deco")
	systemsDir := filepath.Join(decoDir, "nodes", "systems")
	itemsDir := filepath.Join(decoDir, "nodes", "items")
	if err := os.MkdirAll(systemsDir, 0755); err != nil {
		t.Fatalf("Failed to create systems directory: %v", err)
	}
	if err := os.MkdirAll(itemsDir, 0755); err != nil {
		t.Fatalf("Failed to create items directory: %v", err)
	}

	configYAML := `version: 1
project_name: ref-test
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	if err := os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	coreYAML := `id: systems/core
kind: system
version: 1
status: approved
title: "Core Gameplay"
tags: [core, shooting]
summary: "Main gameplay loop"
refs:
  uses:
    - target: systems/scoring
      context: "Tracks player score"
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
          text: "Game ends when player lives reach 0"
`
	scoringYAML := `id: systems/scoring
kind: system
version: 2
status: draft
title: "Scoring System"
tags: [core, scoring]
summary: "Handles point tracking"
refs:
  uses:
    - target: items/weapon
      context: "Weapon damage affects score"
content:
  sections:
    - name: Parameters
      blocks:
        - type: param
          id: base_score
          name: base_score
          value: "100"
          min: 0
          max: 10000
          unit: points
`
	weaponYAML := `id: items/weapon
kind: item
version: 1
status: approved
title: "Weapon"
tags: [combat, items]
summary: "Player weapon definition"
content:
  sections:
    - name: Stats
      blocks:
        - type: param
          id: damage
          name: damage
          value: "25"
          min: 1
          max: 100
`
	if err := os.WriteFile(filepath.Join(systemsDir, "core.yaml"), []byte(coreYAML), 0644); err != nil {
		t.Fatalf("Failed to create core node: %v", err)
	}
	if err := os.WriteFile(filepath.Join(systemsDir, "scoring.yaml"), []byte(scoringYAML), 0644); err != nil {
		t.Fatalf("Failed to create scoring node: %v", err)
	}
	if err := os.WriteFile(filepath.Join(itemsDir, "weapon.yaml"), []byte(weaponYAML), 0644); err != nil {
		t.Fatalf("Failed to create weapon node: %v", err)
	}
}

func TestCompactExport_SingleNode(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithContent(t, tmpDir)

	output := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "systems/core", tmpDir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Header line: id, version, status, tags
	if !strings.Contains(output, "# systems/core (v1, approved) [core, shooting]") {
		t.Errorf("Expected compact header line, got:\n%s", output)
	}

	// Title with summary
	if !strings.Contains(output, "Core Gameplay") {
		t.Error("Expected title in output")
	}

	// Refs line
	if !strings.Contains(output, "uses: entities/player") {
		t.Errorf("Expected refs line with 'uses: entities/player', got:\n%s", output)
	}

	// Section header
	if !strings.Contains(output, "## Controls") {
		t.Error("Expected section header '## Controls'")
	}

	// Block types
	if !strings.Contains(output, "[table]") {
		t.Error("Expected [table] block in output")
	}
	if !strings.Contains(output, "[rule]") {
		t.Error("Expected [rule] block in output")
	}

	// Separator
	if !strings.Contains(output, "---") {
		t.Error("Expected '---' separator at end of node")
	}
}

func TestCompactExport_FilterByKind(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir) // has system and item nodes

	// Filter-only (no node ID) requires the working dir to be the project root
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	output := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "--kind", "system"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Should contain the system nodes
	if !strings.Contains(output, "# systems/core") {
		t.Error("Expected systems/core in kind-filtered output")
	}
	if !strings.Contains(output, "# systems/scoring") {
		t.Error("Expected systems/scoring in kind-filtered output")
	}

	// Should NOT contain the item node
	if strings.Contains(output, "# items/weapon") {
		t.Error("Expected items/weapon to be excluded by kind filter")
	}
}

func TestCompactExport_FilterByTag(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	output := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "--tag", "combat"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Only the weapon has the "combat" tag
	if !strings.Contains(output, "# items/weapon") {
		t.Error("Expected items/weapon with 'combat' tag")
	}
	if strings.Contains(output, "# systems/core") {
		t.Error("Expected systems/core to be excluded by tag filter")
	}
}

func TestCompactExport_FollowUses(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	output := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "systems/core", "--follow", tmpDir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Root node should appear
	if !strings.Contains(output, "# systems/core") {
		t.Error("Expected root node systems/core")
	}

	// Followed node should appear (depth=1 default, so only direct uses)
	if !strings.Contains(output, "# systems/scoring") {
		t.Error("Expected followed node systems/scoring")
	}

	// Followed annotation
	if !strings.Contains(output, "uses from systems/core") {
		t.Errorf("Expected follow annotation 'uses from systems/core', got:\n%s", output)
	}
}

func TestCompactExport_FollowDepth(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	// Depth 1: core + scoring (but not weapon)
	output1 := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "systems/core", "--follow", "--depth", "1", tmpDir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Depth 1: expected no error, got %v", err)
		}
	})

	if !strings.Contains(output1, "# systems/core") {
		t.Error("Depth 1: expected systems/core")
	}
	if !strings.Contains(output1, "# systems/scoring") {
		t.Error("Depth 1: expected systems/scoring")
	}
	if strings.Contains(output1, "# items/weapon") {
		t.Error("Depth 1: expected items/weapon to NOT be included")
	}

	// Depth 2: core + scoring + weapon
	output2 := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "systems/core", "--follow", "--depth", "2", tmpDir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Depth 2: expected no error, got %v", err)
		}
	})

	if !strings.Contains(output2, "# systems/core") {
		t.Error("Depth 2: expected systems/core")
	}
	if !strings.Contains(output2, "# systems/scoring") {
		t.Error("Depth 2: expected systems/scoring")
	}
	if !strings.Contains(output2, "# items/weapon") {
		t.Error("Depth 2: expected items/weapon at depth 2")
	}
}

func TestCompactExport_FollowDeduplication(t *testing.T) {
	// A uses B and C, B uses C => C should appear only once
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes", "systems")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatal(err)
	}
	configYAML := `version: 1
project_name: dedup-test
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configYAML), 0644)

	nodeA := `id: systems/a
kind: system
version: 1
status: draft
title: "Node A"
refs:
  uses:
    - target: systems/b
    - target: systems/c
`
	nodeB := `id: systems/b
kind: system
version: 1
status: draft
title: "Node B"
refs:
  uses:
    - target: systems/c
`
	nodeC := `id: systems/c
kind: system
version: 1
status: draft
title: "Node C"
`
	os.WriteFile(filepath.Join(nodesDir, "a.yaml"), []byte(nodeA), 0644)
	os.WriteFile(filepath.Join(nodesDir, "b.yaml"), []byte(nodeB), 0644)
	os.WriteFile(filepath.Join(nodesDir, "c.yaml"), []byte(nodeC), 0644)

	output := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "systems/a", "--follow", "--depth", "2", tmpDir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Count occurrences of "# systems/c" — should be exactly 1
	count := strings.Count(output, "# systems/c")
	if count != 1 {
		t.Errorf("Expected systems/c to appear exactly once, got %d times:\n%s", count, output)
	}

	// All three nodes should appear
	if !strings.Contains(output, "# systems/a") {
		t.Error("Expected systems/a")
	}
	if !strings.Contains(output, "# systems/b") {
		t.Error("Expected systems/b")
	}
}

func TestCompactExport_NoFollowDefault(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	output := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "systems/core", tmpDir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Only root node, no followed nodes
	if !strings.Contains(output, "# systems/core") {
		t.Error("Expected root node systems/core")
	}
	if strings.Contains(output, "# systems/scoring") {
		t.Error("Without --follow, systems/scoring should NOT appear")
	}
}

func TestCompactExport_OutputToFile(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	outputFile := filepath.Join(tmpDir, "output.md")

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--compact", "--kind", "system", "--output", outputFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Expected output file at %s: %v", outputFile, err)
	}

	content := string(data)
	if !strings.Contains(content, "# systems/core") {
		t.Error("Expected systems/core in output file")
	}
	if !strings.Contains(content, "# systems/scoring") {
		t.Error("Expected systems/scoring in output file")
	}
}

func TestCompactExport_FollowRequiresCompact(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--follow=uses", "systems/core", tmpDir})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error when using --follow without --compact")
	}
	if !strings.Contains(err.Error(), "--follow") {
		t.Errorf("Expected error to mention --follow, got: %v", err)
	}
}

func TestCompactExport_BlockRendering(t *testing.T) {
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes", "systems")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatal(err)
	}
	configYAML := `version: 1
project_name: block-test
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configYAML), 0644)

	nodeYAML := `id: systems/blocks
kind: system
version: 1
status: draft
title: "Block Test"
content:
  sections:
    - name: Tables
      blocks:
        - type: table
          id: stats
          columns:
            - { key: name, type: string, display: "Name" }
            - { key: value, type: int, display: "Value" }
          rows:
            - name: "HP"
              value: 100
            - name: "MP"
              value: 50
    - name: Rules
      blocks:
        - type: rule
          id: death_rule
          name: death
          text: "Player dies when HP reaches 0"
    - name: Params
      blocks:
        - type: param
          id: speed
          name: speed
          value: "10"
          min: 1
          max: 100
          unit: "m/s"
`
	os.WriteFile(filepath.Join(nodesDir, "blocks.yaml"), []byte(nodeYAML), 0644)

	output := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "systems/blocks", tmpDir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Table block: should contain [table], column names, and row data
	if !strings.Contains(output, "[table] stats") {
		t.Errorf("Expected [table] stats, got:\n%s", output)
	}
	if !strings.Contains(output, "columns(Name, Value)") {
		t.Errorf("Expected columns(Name, Value), got:\n%s", output)
	}
	if !strings.Contains(output, "(HP, 100)") {
		t.Errorf("Expected row data (HP, 100), got:\n%s", output)
	}

	// Rule block
	if !strings.Contains(output, "[rule] death: Player dies when HP reaches 0") {
		t.Errorf("Expected [rule] with name and text, got:\n%s", output)
	}

	// Param block with value and constraints
	if !strings.Contains(output, "[param] speed = 10") {
		t.Errorf("Expected [param] speed = 10, got:\n%s", output)
	}
	if !strings.Contains(output, "min=1") {
		t.Errorf("Expected min=1 constraint, got:\n%s", output)
	}
	if !strings.Contains(output, "max=100") {
		t.Errorf("Expected max=100 constraint, got:\n%s", output)
	}
	if !strings.Contains(output, "unit=m/s") {
		t.Errorf("Expected unit=m/s constraint, got:\n%s", output)
	}
}

func TestCompactExport_FollowAll(t *testing.T) {
	// Test that --follow all includes both uses and related refs
	tmpDir := t.TempDir()
	decoDir := filepath.Join(tmpDir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes", "systems")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatal(err)
	}
	configYAML := `version: 1
project_name: follow-all-test
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configYAML), 0644)

	nodeA := `id: systems/a
kind: system
version: 1
status: draft
title: "Node A"
refs:
  uses:
    - target: systems/b
  related:
    - target: systems/c
`
	nodeB := `id: systems/b
kind: system
version: 1
status: draft
title: "Node B"
`
	nodeC := `id: systems/c
kind: system
version: 1
status: draft
title: "Node C"
`
	os.WriteFile(filepath.Join(nodesDir, "a.yaml"), []byte(nodeA), 0644)
	os.WriteFile(filepath.Join(nodesDir, "b.yaml"), []byte(nodeB), 0644)
	os.WriteFile(filepath.Join(nodesDir, "c.yaml"), []byte(nodeC), 0644)

	output := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "systems/a", "--follow=all", tmpDir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// All three should appear
	if !strings.Contains(output, "# systems/a") {
		t.Error("Expected systems/a")
	}
	if !strings.Contains(output, "# systems/b") {
		t.Error("Expected systems/b (via uses)")
	}
	if !strings.Contains(output, "# systems/c") {
		t.Error("Expected systems/c (via related)")
	}

	// Check annotations
	if !strings.Contains(output, "uses from systems/a") {
		t.Error("Expected 'uses from systems/a' annotation for B")
	}
	if !strings.Contains(output, "related from systems/a") {
		t.Error("Expected 'related from systems/a' annotation for C")
	}
}

func TestCompactExport_FilterByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir) // core=approved, scoring=draft, weapon=approved

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	output := captureStdout(t, func() {
		cmd := NewExportCommand()
		cmd.SetArgs([]string{"--compact", "--status", "draft"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Only scoring is draft
	if !strings.Contains(output, "# systems/scoring") {
		t.Error("Expected systems/scoring (status=draft)")
	}
	if strings.Contains(output, "# systems/core") {
		t.Error("Expected systems/core (status=approved) to be excluded")
	}
	if strings.Contains(output, "# items/weapon") {
		t.Error("Expected items/weapon (status=approved) to be excluded")
	}
}
