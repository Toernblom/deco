package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/node"
)

func TestComputeContentHash(t *testing.T) {
	n := domain.Node{
		ID:      "test-001",
		Kind:    "item",
		Version: 1,
		Status:  "draft",
		Title:   "Test Title",
		Summary: "Test summary",
		Tags:    []string{"tag1", "tag2"},
	}

	t.Run("returns consistent hash for same content", func(t *testing.T) {
		hash1 := ComputeContentHash(n)
		hash2 := ComputeContentHash(n)
		if hash1 != hash2 {
			t.Errorf("Expected consistent hash, got %s and %s", hash1, hash2)
		}
		if len(hash1) != 16 {
			t.Errorf("Expected 16 char hash, got %d chars", len(hash1))
		}
	})

	t.Run("returns deterministic hash for content with blocks", func(t *testing.T) {
		// Blocks have map[string]interface{} Data field that could be non-deterministic
		nodeWithBlocks := domain.Node{
			ID:      "test-blocks",
			Kind:    "mechanic",
			Version: 1,
			Status:  "draft",
			Title:   "Block Test",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Rules",
						Blocks: []domain.Block{
							{
								Type: "table",
								Data: map[string]interface{}{
									"zebra":   "last alphabetically",
									"alpha":   "first alphabetically",
									"middle":  "between",
									"columns": []string{"a", "b", "c"},
								},
							},
						},
					},
				},
			},
		}

		// Hash multiple times - must be identical
		hashes := make([]string, 10)
		for i := 0; i < 10; i++ {
			hashes[i] = ComputeContentHash(nodeWithBlocks)
		}

		for i := 1; i < len(hashes); i++ {
			if hashes[i] != hashes[0] {
				t.Errorf("Non-deterministic hash detected: run 0=%s, run %d=%s", hashes[0], i, hashes[i])
			}
		}
	})

	t.Run("returns deterministic hash for Glossary and Custom maps", func(t *testing.T) {
		// Glossary (map[string]string) and Custom (map[string]interface{}) must be
		// hashed deterministically regardless of Go's random map iteration order.
		// This test verifies the fix for non-deterministic content hash.
		nodeWithMaps := domain.Node{
			ID:      "test-maps",
			Kind:    "mechanic",
			Version: 1,
			Status:  "draft",
			Title:   "Map Test",
			Glossary: map[string]string{
				"zebra":    "last animal",
				"alpha":    "first letter",
				"banana":   "yellow fruit",
				"delta":    "fourth letter",
				"charlie":  "third letter",
				"echo":     "fifth letter",
			},
			Custom: map[string]interface{}{
				"zulu":     "military alphabet",
				"apple":    42,
				"bravo":    true,
				"foxtrot":  []string{"f", "o", "x"},
				"golf":     map[string]interface{}{"nested": "value"},
				"hotel":    3.14,
			},
		}

		// Hash 100 times to ensure determinism despite random map iteration
		hashes := make([]string, 100)
		for i := 0; i < 100; i++ {
			hashes[i] = ComputeContentHash(nodeWithMaps)
		}

		for i := 1; i < len(hashes); i++ {
			if hashes[i] != hashes[0] {
				t.Errorf("Non-deterministic hash for Glossary/Custom maps: run 0=%s, run %d=%s", hashes[0], i, hashes[i])
			}
		}
	})

	t.Run("returns different hash for different content", func(t *testing.T) {
		modified := n
		modified.Title = "Different Title"
		hash1 := ComputeContentHash(n)
		hash2 := ComputeContentHash(modified)
		if hash1 == hash2 {
			t.Error("Expected different hash for different content")
		}
	})

	t.Run("ignores metadata fields", func(t *testing.T) {
		modified := n
		modified.Version = 99
		modified.Status = "approved"
		modified.Reviewers = []domain.Reviewer{{Name: "alice", Version: 1}}
		hash1 := ComputeContentHash(n)
		hash2 := ComputeContentHash(modified)
		if hash1 != hash2 {
			t.Error("Expected same hash when only metadata differs")
		}
	})

	t.Run("includes Kind in hash", func(t *testing.T) {
		modified := n
		modified.Kind = "different_kind"
		hash1 := ComputeContentHash(n)
		hash2 := ComputeContentHash(modified)
		if hash1 == hash2 {
			t.Error("Expected different hash when Kind changes")
		}
	})

	t.Run("includes Glossary in hash", func(t *testing.T) {
		modified := n
		modified.Glossary = map[string]string{"term": "definition"}
		hash1 := ComputeContentHash(n)
		hash2 := ComputeContentHash(modified)
		if hash1 == hash2 {
			t.Error("Expected different hash when Glossary changes")
		}
	})

	t.Run("includes Contracts in hash", func(t *testing.T) {
		modified := n
		modified.Contracts = []domain.Contract{{Name: "contract1", Scenario: "test scenario"}}
		hash1 := ComputeContentHash(n)
		hash2 := ComputeContentHash(modified)
		if hash1 == hash2 {
			t.Error("Expected different hash when Contracts change")
		}
	})

	t.Run("includes LLMContext in hash", func(t *testing.T) {
		modified := n
		modified.LLMContext = "AI context information"
		hash1 := ComputeContentHash(n)
		hash2 := ComputeContentHash(modified)
		if hash1 == hash2 {
			t.Error("Expected different hash when LLMContext changes")
		}
	})

	t.Run("includes Constraints in hash", func(t *testing.T) {
		modified := n
		modified.Constraints = []domain.Constraint{{Expr: "size < 100", Message: "Too large"}}
		hash1 := ComputeContentHash(n)
		hash2 := ComputeContentHash(modified)
		if hash1 == hash2 {
			t.Error("Expected different hash when Constraints change")
		}
	})

	t.Run("includes Custom fields in hash", func(t *testing.T) {
		modified := n
		modified.Custom = map[string]interface{}{"custom_field": "custom_value"}
		hash1 := ComputeContentHash(n)
		hash2 := ComputeContentHash(modified)
		if hash1 == hash2 {
			t.Error("Expected different hash when Custom fields change")
		}
	})
}

func TestGetLastContentHash(t *testing.T) {
	t.Run("returns empty string when no history", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		hash := getLastContentHash(filepath.Join(tmpDir, ".deco", "history.jsonl"), "sword-001")
		if hash != "" {
			t.Errorf("Expected empty hash for no history, got %q", hash)
		}
	})

	t.Run("returns hash from most recent entry", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Write history with content hash
		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := `{"timestamp":"2026-01-01T00:00:00Z","node_id":"sword-001","operation":"create","user":"test","content_hash":"abc123def456"}
{"timestamp":"2026-01-02T00:00:00Z","node_id":"sword-001","operation":"set","user":"test","content_hash":"xyz789uvw012"}
`
		if err := os.WriteFile(historyPath, []byte(historyContent), 0644); err != nil {
			t.Fatalf("Failed to write history: %v", err)
		}

		hash := getLastContentHash(filepath.Join(tmpDir, ".deco", "history.jsonl"), "sword-001")
		if hash != "xyz789uvw012" {
			t.Errorf("Expected 'xyz789uvw012', got %q", hash)
		}
	})
}

func TestSyncCommand_Structure(t *testing.T) {
	t.Run("creates sync command", func(t *testing.T) {
		cmd := NewSyncCommand()
		if cmd == nil {
			t.Fatal("Expected sync command, got nil")
		}
		if !strings.HasPrefix(cmd.Use, "sync") {
			t.Errorf("Expected Use to start with 'sync', got %q", cmd.Use)
		}
	})

	t.Run("has description", func(t *testing.T) {
		cmd := NewSyncCommand()
		if cmd.Short == "" {
			t.Error("Expected non-empty Short description")
		}
	})

	t.Run("has dry-run flag", func(t *testing.T) {
		cmd := NewSyncCommand()
		flag := cmd.Flags().Lookup("dry-run")
		if flag == nil {
			t.Fatal("Expected --dry-run flag to be defined")
		}
	})

	t.Run("has quiet flag", func(t *testing.T) {
		cmd := NewSyncCommand()
		flag := cmd.Flags().Lookup("quiet")
		if flag == nil {
			t.Fatal("Expected --quiet flag to be defined")
		}
		if flag.Shorthand != "q" {
			t.Errorf("Expected shorthand 'q', got %q", flag.Shorthand)
		}
	})
}

func TestRunSync_NoProject(t *testing.T) {
	t.Run("errors on missing .deco directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		flags := &syncFlags{targetDir: tmpDir}
		exitCode, err := runSync(flags)

		if err == nil {
			t.Error("Expected error for missing .deco directory")
		}
		if exitCode != syncExitError {
			t.Errorf("Expected exit code %d, got %d", syncExitError, exitCode)
		}
	})
}

func TestRunSync_Baseline(t *testing.T) {
	t.Run("baselines nodes without history", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)
		// No history.jsonl exists

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitClean {
			t.Errorf("Expected exit code %d for baseline, got %d", syncExitClean, exitCode)
		}

		// Verify baseline was recorded
		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		content, err := os.ReadFile(historyPath)
		if err != nil {
			t.Fatalf("Expected history file to be created: %v", err)
		}
		if !strings.Contains(string(content), `"operation":"baseline"`) {
			t.Errorf("Expected baseline operation in history, got: %s", string(content))
		}
		if !strings.Contains(string(content), `"content_hash"`) {
			t.Errorf("Expected content_hash in history, got: %s", string(content))
		}
	})
}

func TestRunSync_NoChanges(t *testing.T) {
	t.Run("returns clean exit when hashes match", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Create baseline history entry with current hash
		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := ComputeContentHash(n)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"sword-001","operation":"create","user":"test","content_hash":"%s"}`, hash)
		os.WriteFile(historyPath, []byte(historyContent+"\n"), 0644)

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitClean {
			t.Errorf("Expected exit code %d, got %d", syncExitClean, exitCode)
		}
	})
}

func TestRunSync_MetadataOnlyChange(t *testing.T) {
	t.Run("ignores metadata-only changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Baseline with current hash
		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := ComputeContentHash(n)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"sword-001","operation":"create","user":"test","content_hash":"%s"}`, hash)
		os.WriteFile(historyPath, []byte(historyContent+"\n"), 0644)

		// Change only metadata (version, status) - hash should still match
		nodeYAML := `id: sword-001
kind: item
version: 5
status: approved
title: Iron Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
`
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		os.WriteFile(nodePath, []byte(nodeYAML), 0644)

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitClean {
			t.Errorf("Expected exit code %d (metadata-only change ignored), got %d", syncExitClean, exitCode)
		}
	})
}

func TestRunSync_ContentChange(t *testing.T) {
	t.Run("syncs nodes with content changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Baseline with original hash
		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := ComputeContentHash(n)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"sword-001","operation":"create","user":"test","content_hash":"%s"}`, hash)
		os.WriteFile(historyPath, []byte(historyContent+"\n"), 0644)

		// Change content (title) - should trigger sync
		nodeYAML := `id: sword-001
kind: item
version: 1
status: approved
title: Golden Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
reviewers:
  - name: alice@example.com
    timestamp: 2026-01-01T00:00:00Z
    version: 1
`
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		os.WriteFile(nodePath, []byte(nodeYAML), 0644)

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitModified {
			t.Errorf("Expected exit code %d, got %d", syncExitModified, exitCode)
		}

		// Verify version was bumped
		content, _ := os.ReadFile(nodePath)
		if !strings.Contains(string(content), "version: 2") {
			t.Errorf("Expected version to be bumped to 2, got: %s", string(content))
		}

		// Verify status was reset to draft
		if !strings.Contains(string(content), "status: draft") {
			t.Errorf("Expected status to be reset to draft, got: %s", string(content))
		}
	})
}

func TestRunSync_DryRun(t *testing.T) {
	t.Run("dry-run does not modify files but returns modified exit code", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Baseline with original hash
		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := ComputeContentHash(n)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"sword-001","operation":"create","user":"test","content_hash":"%s"}`, hash)
		os.WriteFile(historyPath, []byte(historyContent+"\n"), 0644)

		// Change content
		nodeYAML := `id: sword-001
kind: item
version: 1
status: approved
title: Golden Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
`
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		os.WriteFile(nodePath, []byte(nodeYAML), 0644)

		flags := &syncFlags{targetDir: tmpDir, dryRun: true, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Dry-run returns modified (1) when changes WOULD be made
		if exitCode != syncExitModified {
			t.Errorf("Expected exit code %d for dry-run with changes, got %d", syncExitModified, exitCode)
		}

		// Verify file was NOT modified
		content, _ := os.ReadFile(nodePath)
		if strings.Contains(string(content), "version: 2") {
			t.Error("Dry-run should not modify files")
		}
	})

	t.Run("dry-run returns clean when no changes needed", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Baseline with current hash (no changes)
		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := ComputeContentHash(n)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"sword-001","operation":"create","user":"test","content_hash":"%s"}`, hash)
		os.WriteFile(historyPath, []byte(historyContent+"\n"), 0644)

		flags := &syncFlags{targetDir: tmpDir, dryRun: true, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Dry-run returns clean (0) when no changes would be made
		if exitCode != syncExitClean {
			t.Errorf("Expected exit code %d for dry-run with no changes, got %d", syncExitClean, exitCode)
		}
	})
}

func TestRunSync_HistoryLogging(t *testing.T) {
	t.Run("logs sync operation to history with content hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Baseline with original hash
		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := ComputeContentHash(n)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"sword-001","operation":"create","user":"test","content_hash":"%s"}`, hash)
		os.WriteFile(historyPath, []byte(historyContent+"\n"), 0644)

		// Change content
		nodeYAML := `id: sword-001
kind: item
version: 1
status: approved
title: Golden Sword
summary: A basic iron sword
tags:
  - weapon
  - combat
`
		nodePath := filepath.Join(tmpDir, ".deco", "nodes", "sword-001.yaml")
		os.WriteFile(nodePath, []byte(nodeYAML), 0644)

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		_, err := runSync(flags)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify history entry was created
		content, err := os.ReadFile(historyPath)
		if err != nil {
			t.Fatalf("Failed to read history: %v", err)
		}

		historyStr := string(content)
		if !strings.Contains(historyStr, `"operation":"sync"`) {
			t.Errorf("Expected sync operation in history, got: %s", historyStr)
		}
		if !strings.Contains(historyStr, `"node_id":"sword-001"`) {
			t.Errorf("Expected node_id in history, got: %s", historyStr)
		}
		if !strings.Contains(historyStr, `"content_hash"`) {
			t.Errorf("Expected content_hash in history, got: %s", historyStr)
		}
	})
}

func TestRunSync_ErrorAccumulation(t *testing.T) {
	t.Run("returns error exit when node fails to load", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Create an invalid node (malformed YAML)
		invalidNodePath := filepath.Join(tmpDir, ".deco", "nodes", "broken-001.yaml")
		invalidYAML := `id: broken-001
kind: item
version: invalid_should_be_int
status: draft
title: Broken Node
`
		os.WriteFile(invalidNodePath, []byte(invalidYAML), 0644)

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err == nil {
			t.Error("Expected error when node fails to load")
		}
		if exitCode != syncExitError {
			t.Errorf("Expected exit code %d, got %d", syncExitError, exitCode)
		}
	})

	t.Run("accumulates multiple errors and still exits with error", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Create multiple invalid nodes
		for i := 1; i <= 3; i++ {
			invalidNodePath := filepath.Join(tmpDir, ".deco", "nodes", fmt.Sprintf("broken-%03d.yaml", i))
			invalidYAML := fmt.Sprintf(`id: broken-%03d
kind: item
version: not_a_number
status: draft
title: Broken Node %d
`, i, i)
			os.WriteFile(invalidNodePath, []byte(invalidYAML), 0644)
		}

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err == nil {
			t.Error("Expected error when multiple nodes fail to load")
		}
		if exitCode != syncExitError {
			t.Errorf("Expected exit code %d, got %d", syncExitError, exitCode)
		}
		// Error message should indicate multiple errors
		if err != nil && !strings.Contains(err.Error(), "3") {
			t.Errorf("Expected error to mention 3 errors, got: %v", err)
		}
	})

	t.Run("valid nodes still process when some nodes fail", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Create an invalid node alongside the valid one
		invalidNodePath := filepath.Join(tmpDir, ".deco", "nodes", "broken-001.yaml")
		invalidYAML := `id: broken-001
kind: item
version: not_a_number
status: draft
title: Broken Node
`
		os.WriteFile(invalidNodePath, []byte(invalidYAML), 0644)

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, _ := runSync(flags)

		// Should return error because one node failed
		if exitCode != syncExitError {
			t.Errorf("Expected exit code %d (some nodes failed), got %d", syncExitError, exitCode)
		}

		// But the valid node should still be baselined
		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		content, err := os.ReadFile(historyPath)
		if err != nil {
			t.Fatalf("Expected history file to exist: %v", err)
		}
		if !strings.Contains(string(content), "sword-001") {
			t.Error("Expected valid node to still be baselined despite other failures")
		}
	})
}

func TestRunSync_RenameDetection(t *testing.T) {
	t.Run("detects manual rename and updates references", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSyncWithRefs(t, tmpDir)

		// Create history for old ID (gameplay-001)
		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)

		// Load combat-001 (which was renamed from gameplay-001) to get its hash
		combat, _ := nodeRepo.Load("combat-001")
		combatHash := ComputeContentHash(combat)

		// Write history as if gameplay-001 existed with the same content hash
		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"gameplay-001","operation":"create","user":"test","content_hash":"%s"}
{"timestamp":"2026-01-01T00:00:00Z","node_id":"player-001","operation":"create","user":"test","content_hash":"abc123def456"}
`, combatHash)
		os.WriteFile(historyPath, []byte(historyContent), 0644)

		flags := &syncFlags{targetDir: tmpDir, quiet: false}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitModified {
			t.Errorf("Expected exit code %d (rename detected), got %d", syncExitModified, exitCode)
		}

		// Verify player-001 had its reference updated from gameplay-001 to combat-001
		player, err := nodeRepo.Load("player-001")
		if err != nil {
			t.Fatalf("Failed to load player-001: %v", err)
		}

		if len(player.Refs.Uses) == 0 {
			t.Fatal("Expected player-001 to have Uses references")
		}
		if player.Refs.Uses[0].Target != "combat-001" {
			t.Errorf("Expected reference updated to combat-001, got %s", player.Refs.Uses[0].Target)
		}

		// Verify version was bumped
		if player.Version != 2 {
			t.Errorf("Expected player version to be bumped to 2, got %d", player.Version)
		}

		// Verify move operation was logged
		historyContent2, _ := os.ReadFile(historyPath)
		historyStr := string(historyContent2)
		if !strings.Contains(historyStr, `"operation":"move"`) {
			t.Error("Expected move operation in history")
		}
		if !strings.Contains(historyStr, `"node_id":"combat-001"`) {
			t.Error("Expected combat-001 in move history entry")
		}
	})

	t.Run("dry-run shows rename without applying", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSyncWithRefs(t, tmpDir)

		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)
		combat, _ := nodeRepo.Load("combat-001")
		combatHash := ComputeContentHash(combat)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"gameplay-001","operation":"create","user":"test","content_hash":"%s"}
{"timestamp":"2026-01-01T00:00:00Z","node_id":"player-001","operation":"create","user":"test","content_hash":"abc123def456"}
`, combatHash)
		os.WriteFile(historyPath, []byte(historyContent), 0644)

		flags := &syncFlags{targetDir: tmpDir, dryRun: true, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitModified {
			t.Errorf("Expected exit code %d for dry-run with rename, got %d", syncExitModified, exitCode)
		}

		// Verify player-001 was NOT modified
		player, _ := nodeRepo.Load("player-001")
		if player.Refs.Uses[0].Target != "gameplay-001" {
			t.Error("Dry-run should not update references")
		}
		if player.Version != 1 {
			t.Error("Dry-run should not bump version")
		}
	})

	t.Run("no-refactor flag skips reference updates", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSyncWithRefs(t, tmpDir)

		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)
		combat, _ := nodeRepo.Load("combat-001")
		combatHash := ComputeContentHash(combat)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"gameplay-001","operation":"create","user":"test","content_hash":"%s"}
{"timestamp":"2026-01-01T00:00:00Z","node_id":"player-001","operation":"create","user":"test","content_hash":"abc123def456"}
`, combatHash)
		os.WriteFile(historyPath, []byte(historyContent), 0644)

		flags := &syncFlags{targetDir: tmpDir, noRefactor: true, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Should still return modified because rename was detected
		if exitCode != syncExitModified {
			t.Errorf("Expected exit code %d, got %d", syncExitModified, exitCode)
		}

		// Verify player-001 was NOT modified (no-refactor)
		player, _ := nodeRepo.Load("player-001")
		if player.Refs.Uses[0].Target != "gameplay-001" {
			t.Error("--no-refactor should not update references")
		}
	})

	t.Run("no false positive on genuine new nodes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Create history for existing node only
		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := ComputeContentHash(n)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"sword-001","operation":"create","user":"test","content_hash":"%s"}
`, hash)
		os.WriteFile(historyPath, []byte(historyContent), 0644)

		// Add a genuinely new node with different content
		newNodeYAML := `id: shield-001
kind: item
version: 1
status: draft
title: Iron Shield
summary: A basic iron shield
tags:
  - armor
  - defense
`
		newNodePath := filepath.Join(nodesDir, "shield-001.yaml")
		os.WriteFile(newNodePath, []byte(newNodeYAML), 0644)

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Should return clean (baseline doesn't count as modified)
		if exitCode != syncExitClean {
			t.Errorf("Expected exit code %d, got %d", syncExitClean, exitCode)
		}

		// Verify shield-001 was baselined (not treated as rename)
		historyContent2, _ := os.ReadFile(historyPath)
		historyStr := string(historyContent2)
		if !strings.Contains(historyStr, `"node_id":"shield-001"`) {
			t.Error("Expected shield-001 to be baselined")
		}
		if !strings.Contains(historyStr, `"operation":"baseline"`) {
			t.Error("Expected baseline operation for new node")
		}
		// Should NOT have move operation
		if strings.Contains(historyStr, `"operation":"move"`) {
			t.Error("Should not detect move for genuinely new node")
		}
	})

	t.Run("rename with no references to update", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		nodesDir := filepath.Join(tmpDir, ".deco", "nodes")
		nodeRepo := node.NewYAMLRepository(nodesDir)

		// Rename sword-001 to blade-001 manually
		swordContent, _ := os.ReadFile(filepath.Join(nodesDir, "sword-001.yaml"))
		bladeContent := strings.Replace(string(swordContent), "id: sword-001", "id: blade-001", 1)
		os.WriteFile(filepath.Join(nodesDir, "blade-001.yaml"), []byte(bladeContent), 0644)
		os.Remove(filepath.Join(nodesDir, "sword-001.yaml"))

		// Create history for old ID
		blade, _ := nodeRepo.Load("blade-001")
		bladeHash := ComputeContentHash(blade)

		historyPath := filepath.Join(tmpDir, ".deco", "history.jsonl")
		historyContent := fmt.Sprintf(`{"timestamp":"2026-01-01T00:00:00Z","node_id":"sword-001","operation":"create","user":"test","content_hash":"%s"}
`, bladeHash)
		os.WriteFile(historyPath, []byte(historyContent), 0644)

		flags := &syncFlags{targetDir: tmpDir, quiet: true}
		exitCode, err := runSync(flags)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitCode != syncExitModified {
			t.Errorf("Expected exit code %d, got %d", syncExitModified, exitCode)
		}

		// Verify move operation was logged
		historyContent2, _ := os.ReadFile(historyPath)
		historyStr := string(historyContent2)
		if !strings.Contains(historyStr, `"operation":"move"`) {
			t.Error("Expected move operation in history")
		}
	})
}

func TestSyncCommand_NoRefactorFlag(t *testing.T) {
	t.Run("has no-refactor flag", func(t *testing.T) {
		cmd := NewSyncCommand()
		flag := cmd.Flags().Lookup("no-refactor")
		if flag == nil {
			t.Fatal("Expected --no-refactor flag to be defined")
		}
	})
}

// Test helpers

func setupProjectForSyncWithRefs(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: sync-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create combat-001 (simulating a rename from gameplay-001)
	combatYAML := `id: combat-001
kind: mechanic
version: 1
status: draft
title: Combat System
summary: Core combat mechanics
tags:
  - mechanic
  - combat
`
	combatPath := filepath.Join(nodesDir, "combat-001.yaml")
	if err := os.WriteFile(combatPath, []byte(combatYAML), 0644); err != nil {
		t.Fatalf("Failed to create combat node: %v", err)
	}

	// Create player-001 with reference to old ID (gameplay-001)
	playerYAML := `id: player-001
kind: entity
version: 1
status: draft
title: Player Character
summary: The player's character
tags:
  - entity
  - player
refs:
  uses:
    - target: gameplay-001
      context: uses gameplay mechanics
`
	playerPath := filepath.Join(nodesDir, "player-001.yaml")
	if err := os.WriteFile(playerPath, []byte(playerYAML), 0644); err != nil {
		t.Fatalf("Failed to create player node: %v", err)
	}
}

func setupProjectForSync(t *testing.T, dir string) {
	t.Helper()

	// Create .deco structure
	decoDir := filepath.Join(dir, ".deco")
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		t.Fatalf("Failed to create nodes directory: %v", err)
	}

	// Create config.yaml
	configYAML := `version: 1
project_name: sync-test-project
nodes_path: .deco/nodes
history_path: .deco/history.jsonl
`
	configPath := filepath.Join(decoDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}

	// Create a node
	nodeYAML := `id: sword-001
kind: item
version: 1
status: approved
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

func TestComputeContentHashWithDir_IncludesDocFiles(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "narrative.md")
	os.WriteFile(mdPath, []byte("The story begins."), 0644)

	node := domain.Node{
		ID:      "test/node",
		Kind:    "system",
		Version: 1,
		Status:  "draft",
		Title:   "Test",
		Docs: []domain.DocRef{
			{Path: "narrative.md"},
		},
	}

	t.Run("hash changes when doc file content changes", func(t *testing.T) {
		hash1 := ComputeContentHashWithDir(node, dir)

		// Change the .md file content
		os.WriteFile(mdPath, []byte("The story continues with a twist."), 0644)

		hash2 := ComputeContentHashWithDir(node, dir)

		if hash1 == hash2 {
			t.Error("expected hash to change when doc file content changes")
		}
	})

	t.Run("hash is consistent for same content", func(t *testing.T) {
		hash1 := ComputeContentHashWithDir(node, dir)
		hash2 := ComputeContentHashWithDir(node, dir)
		if hash1 != hash2 {
			t.Errorf("expected consistent hash, got %s and %s", hash1, hash2)
		}
	})

	t.Run("missing doc file is silently skipped", func(t *testing.T) {
		nodeWithMissing := domain.Node{
			ID:      "test/node2",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Test",
			Docs: []domain.DocRef{
				{Path: "nonexistent.md"},
			},
		}

		hash := ComputeContentHashWithDir(nodeWithMissing, dir)
		if hash == "" {
			t.Error("expected non-empty hash even with missing doc file")
		}
	})

	t.Run("includes doc block file paths in hash", func(t *testing.T) {
		blockMdPath := filepath.Join(dir, "scene.md")
		os.WriteFile(blockMdPath, []byte("A dark and stormy night."), 0644)

		nodeWithDocBlock := domain.Node{
			ID:      "test/blocks",
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   "Test",
			Content: &domain.Content{
				Sections: []domain.Section{
					{
						Name: "Narrative",
						Blocks: []domain.Block{
							{
								Type: "doc",
								Data: map[string]interface{}{
									"path": "scene.md",
								},
							},
						},
					},
				},
			},
		}

		hash1 := ComputeContentHashWithDir(nodeWithDocBlock, dir)
		os.WriteFile(blockMdPath, []byte("A bright and sunny morning."), 0644)
		hash2 := ComputeContentHashWithDir(nodeWithDocBlock, dir)

		if hash1 == hash2 {
			t.Error("expected hash to change when doc block file changes")
		}
	})

	t.Run("without projectRoot falls back to yaml-only hash", func(t *testing.T) {
		hash1 := ComputeContentHash(node)
		hash2 := ComputeContentHashWithDir(node, "")
		if hash1 != hash2 {
			t.Errorf("expected same hash without project root, got %s and %s", hash1, hash2)
		}
	})
}
