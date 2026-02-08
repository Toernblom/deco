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

func TestObsidianExport_VaultStructure(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	vaultDir := filepath.Join(tmpDir, ".deco", "vault")

	// Check vault directory exists
	if _, err := os.Stat(vaultDir); os.IsNotExist(err) {
		t.Fatal("Expected vault directory to exist")
	}

	// Check files mirror nodes structure
	for _, path := range []string{
		"systems/core.md",
		"systems/scoring.md",
		"items/weapon.md",
	} {
		full := filepath.Join(vaultDir, path)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("Expected vault file %s to exist", path)
		}
	}
}

func TestObsidianExport_Frontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".deco", "vault", "systems", "core.md"))
	if err != nil {
		t.Fatalf("Failed to read vault file: %v", err)
	}
	content := string(data)

	// Should start with frontmatter
	if !strings.HasPrefix(content, "---\n") {
		t.Error("Expected file to start with YAML frontmatter")
	}

	// Check frontmatter fields
	if !strings.Contains(content, "id: systems/core") {
		t.Error("Expected frontmatter to contain id")
	}
	if !strings.Contains(content, "kind: system") {
		t.Error("Expected frontmatter to contain kind")
	}
	if !strings.Contains(content, "status: approved") {
		t.Error("Expected frontmatter to contain status")
	}
	if !strings.Contains(content, "version: 1") {
		t.Error("Expected frontmatter to contain version")
	}
}

func TestObsidianExport_Wikilinks(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".deco", "vault", "systems", "core.md"))
	if err != nil {
		t.Fatalf("Failed to read vault file: %v", err)
	}
	content := string(data)

	// Uses ref should be a wikilink with title alias
	if !strings.Contains(content, "[[systems/scoring|Scoring System]]") {
		t.Errorf("Expected wikilink [[systems/scoring|Scoring System]], got:\n%s", content)
	}
}

func TestObsidianExport_Tags(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".deco", "vault", "systems", "core.md"))
	if err != nil {
		t.Fatalf("Failed to read vault file: %v", err)
	}
	content := string(data)

	// Check inline Obsidian tags
	if !strings.Contains(content, "#core") {
		t.Error("Expected inline tag #core")
	}
	if !strings.Contains(content, "#shooting") {
		t.Error("Expected inline tag #shooting")
	}

	// Check tags in frontmatter
	if !strings.Contains(content, "- core") && !strings.Contains(content, "tags:") {
		t.Error("Expected tags in frontmatter")
	}
}

func TestObsidianExport_ContentBlocks(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithContent(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".deco", "vault", "systems", "core.md"))
	if err != nil {
		t.Fatalf("Failed to read vault file: %v", err)
	}
	content := string(data)

	// Table rendering
	if !strings.Contains(content, "| Input") {
		t.Error("Expected markdown table header")
	}

	// Rule rendering
	if !strings.Contains(content, "Game ends") {
		t.Error("Expected rule text")
	}

	// Section headers
	if !strings.Contains(content, "## Controls") {
		t.Error("Expected section header")
	}
}

func TestObsidianExport_Issues(t *testing.T) {
	tmpDir := t.TempDir()
	setupObsidianProjectWithIssues(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".deco", "vault", "systems", "core.md"))
	if err != nil {
		t.Fatalf("Failed to read vault file: %v", err)
	}
	content := string(data)

	// Check callout rendering
	if !strings.Contains(content, "[!warning]") {
		t.Error("Expected high severity issue as [!warning] callout")
	}
	if !strings.Contains(content, "[!info]") {
		t.Error("Expected low severity issue as [!info] callout")
	}
	if !strings.Contains(content, "Needs balancing") {
		t.Error("Expected issue description")
	}
}

func TestObsidianExport_FilterByKind(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", "--kind", "system", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	vaultDir := filepath.Join(tmpDir, ".deco", "vault")

	// System nodes should exist
	if _, err := os.Stat(filepath.Join(vaultDir, "systems", "core.md")); os.IsNotExist(err) {
		t.Error("Expected systems/core.md in filtered vault")
	}

	// Item node should NOT exist
	if _, err := os.Stat(filepath.Join(vaultDir, "items", "weapon.md")); !os.IsNotExist(err) {
		t.Error("Expected items/weapon.md to be excluded by kind filter")
	}
}

func TestObsidianExport_CleanRegeneration(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithRefs(t, tmpDir)

	vaultDir := filepath.Join(tmpDir, ".deco", "vault")

	// Create a stale file in vault
	staleDir := filepath.Join(vaultDir, "stale")
	os.MkdirAll(staleDir, 0755)
	os.WriteFile(filepath.Join(staleDir, "old.md"), []byte("stale"), 0644)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Stale file should be gone
	if _, err := os.Stat(filepath.Join(staleDir, "old.md")); !os.IsNotExist(err) {
		t.Error("Expected stale files to be cleaned on re-export")
	}
}

func TestObsidianExport_CustomBlocks(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithCustomBlocks(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".deco", "vault", "items", "powerups.md"))
	if err != nil {
		t.Fatalf("Failed to read vault file: %v", err)
	}
	content := string(data)

	// Custom block should render as structured definition list
	if !strings.Contains(content, "**Speed Boost**") {
		t.Error("Expected custom block name as bold heading")
	}
	if !strings.Contains(content, "(powerup)") {
		t.Error("Expected block type annotation")
	}
	if !strings.Contains(content, "**effect:**") {
		t.Error("Expected field name rendered as bold")
	}
}

func TestObsidianExport_Glossary(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithGlossary(t, tmpDir)

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".deco", "vault", "systems", "core.md"))
	if err != nil {
		t.Fatalf("Failed to read vault file: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "## Glossary") {
		t.Error("Expected Glossary section")
	}
	if !strings.Contains(content, "**HP**") {
		t.Error("Expected glossary term HP")
	}
	if !strings.Contains(content, ": Hit Points") {
		t.Error("Expected glossary definition")
	}
}

func TestObsidianExport_BrokenRefFallback(t *testing.T) {
	tmpDir := t.TempDir()
	setupProjectWithContent(t, tmpDir) // has ref to entities/player which doesn't exist

	cmd := NewExportCommand()
	cmd.SetArgs([]string{"--obsidian", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".deco", "vault", "systems", "core.md"))
	if err != nil {
		t.Fatalf("Failed to read vault file: %v", err)
	}
	content := string(data)

	// Broken ref should fall back to [[id]] without title alias
	if !strings.Contains(content, "[[entities/player]]") {
		t.Errorf("Expected broken ref fallback [[entities/player]], got:\n%s", content)
	}
}

// --- Test helpers ---

func setupObsidianProjectWithIssues(t *testing.T, dir string) {
	t.Helper()

	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes", "systems")
	os.MkdirAll(nodesDir, 0755)

	configYAML := `version: 1
project_name: issue-test
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configYAML), 0644)

	nodeYAML := `id: systems/core
kind: system
version: 1
status: draft
title: "Core System"
issues:
  - id: balance-001
    description: "Needs balancing"
    severity: high
    location: "content.sections[0]"
    resolved: false
  - id: minor-001
    description: "Consider adding tutorial"
    severity: low
    resolved: false
  - id: done-001
    description: "Fixed rendering bug"
    severity: medium
    resolved: true
`
	os.WriteFile(filepath.Join(nodesDir, "core.yaml"), []byte(nodeYAML), 0644)
}

func setupProjectWithCustomBlocks(t *testing.T, dir string) {
	t.Helper()

	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes", "items")
	os.MkdirAll(nodesDir, 0755)

	configYAML := `version: 1
project_name: custom-block-test
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
custom_block_types:
  powerup:
    required_fields: [name, effect, duration]
`
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configYAML), 0644)

	nodeYAML := `id: items/powerups
kind: item
version: 1
status: draft
title: "Power-Ups"
content:
  sections:
    - name: Available Power-Ups
      blocks:
        - type: powerup
          name: "Speed Boost"
          effect: "2x movement speed"
          duration: "5s"
`
	os.WriteFile(filepath.Join(nodesDir, "powerups.yaml"), []byte(nodeYAML), 0644)
}

func setupProjectWithGlossary(t *testing.T, dir string) {
	t.Helper()

	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes", "systems")
	os.MkdirAll(nodesDir, 0755)

	configYAML := `version: 1
project_name: glossary-test
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	os.WriteFile(filepath.Join(decoDir, "config.yaml"), []byte(configYAML), 0644)

	nodeYAML := `id: systems/core
kind: system
version: 1
status: draft
title: "Core System"
glossary:
  HP: "Hit Points"
  MP: "Magic Points"
`
	os.WriteFile(filepath.Join(nodesDir, "core.yaml"), []byte(nodeYAML), 0644)
}
