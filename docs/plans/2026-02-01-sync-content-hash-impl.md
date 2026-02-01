# Content Hash-Based Sync Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace git-based sync detection with content hashing stored in history.jsonl

**Architecture:** Compute SHA-256 hash of content fields (title, summary, tags, refs, issues, content). Store hash in history entries. Sync compares current hash against last recorded hash to detect drift.

**Tech Stack:** Go, crypto/sha256, gopkg.in/yaml.v3

---

### Task 1: Add ContentHash to AuditEntry

**Files:**
- Modify: `internal/domain/audit.go:10-17`
- Modify: `internal/domain/audit.go:32-44`
- Test: `internal/domain/audit_test.go`

**Step 1: Write the failing test**

Add to `internal/domain/audit_test.go`:

```go
func TestAuditEntry_BaselineOperation(t *testing.T) {
	entry := AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      "test-001",
		Operation:   "baseline",
		User:        "testuser",
		ContentHash: "a1b2c3d4e5f67890",
	}

	err := entry.Validate()
	if err != nil {
		t.Errorf("Expected baseline operation to be valid, got: %v", err)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/domain/... -run TestAuditEntry_BaselineOperation -v`
Expected: FAIL with "unknown field ContentHash" or validation error

**Step 3: Add ContentHash field and baseline operation**

In `internal/domain/audit.go`, update AuditEntry struct:

```go
type AuditEntry struct {
	Timestamp   time.Time              `json:"timestamp" yaml:"timestamp"`
	NodeID      string                 `json:"node_id" yaml:"node_id"`
	Operation   string                 `json:"operation" yaml:"operation"`
	User        string                 `json:"user" yaml:"user"`
	ContentHash string                 `json:"content_hash,omitempty" yaml:"content_hash,omitempty"`
	Before      map[string]interface{} `json:"before,omitempty" yaml:"before,omitempty"`
	After       map[string]interface{} `json:"after,omitempty" yaml:"after,omitempty"`
}
```

Add "baseline" to validOperations map:

```go
validOperations := map[string]bool{
	"create":   true,
	"update":   true,
	"delete":   true,
	"set":      true,
	"append":   true,
	"unset":    true,
	"move":     true,
	"submit":   true,
	"approve":  true,
	"reject":   true,
	"sync":     true,
	"baseline": true, // record current state without modification
}
```

Update error message to include baseline.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/domain/... -run TestAuditEntry_BaselineOperation -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/domain/audit.go internal/domain/audit_test.go
git commit -m "feat(domain): add ContentHash field and baseline operation to AuditEntry"
```

---

### Task 2: Add computeContentHash function

**Files:**
- Modify: `internal/cli/sync.go`
- Test: `internal/cli/sync_test.go`

**Step 1: Write the failing test**

Add to `internal/cli/sync_test.go`:

```go
func TestComputeContentHash(t *testing.T) {
	node := domain.Node{
		ID:      "test-001",
		Kind:    "item",
		Version: 1,
		Status:  "draft",
		Title:   "Test Title",
		Summary: "Test summary",
		Tags:    []string{"tag1", "tag2"},
	}

	t.Run("returns consistent hash for same content", func(t *testing.T) {
		hash1 := computeContentHash(node)
		hash2 := computeContentHash(node)
		if hash1 != hash2 {
			t.Errorf("Expected consistent hash, got %s and %s", hash1, hash2)
		}
		if len(hash1) != 16 {
			t.Errorf("Expected 16 char hash, got %d chars", len(hash1))
		}
	})

	t.Run("returns different hash for different content", func(t *testing.T) {
		modified := node
		modified.Title = "Different Title"
		hash1 := computeContentHash(node)
		hash2 := computeContentHash(modified)
		if hash1 == hash2 {
			t.Error("Expected different hash for different content")
		}
	})

	t.Run("ignores metadata fields", func(t *testing.T) {
		modified := node
		modified.Version = 99
		modified.Status = "approved"
		modified.Reviewers = []domain.Reviewer{{Name: "alice", Version: 1}}
		hash1 := computeContentHash(node)
		hash2 := computeContentHash(modified)
		if hash1 != hash2 {
			t.Error("Expected same hash when only metadata differs")
		}
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli/... -run TestComputeContentHash -v`
Expected: FAIL with "undefined: computeContentHash"

**Step 3: Implement computeContentHash**

Add to `internal/cli/sync.go` (add "crypto/sha256" and "encoding/hex" to imports):

```go
// contentFields holds only the fields that affect content hash
type contentFields struct {
	Title   string           `yaml:"title"`
	Summary string           `yaml:"summary"`
	Tags    []string         `yaml:"tags,omitempty"`
	Refs    domain.Ref       `yaml:"refs,omitempty"`
	Issues  []domain.Issue   `yaml:"issues,omitempty"`
	Content *domain.Content  `yaml:"content,omitempty"`
}

// computeContentHash computes a SHA-256 hash of the content fields
// Returns 16 hex characters (first 64 bits of the hash)
func computeContentHash(n domain.Node) string {
	fields := contentFields{
		Title:   n.Title,
		Summary: n.Summary,
		Tags:    n.Tags,
		Refs:    n.Refs,
		Issues:  n.Issues,
		Content: n.Content,
	}

	data, err := yaml.Marshal(fields)
	if err != nil {
		return ""
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:8])
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli/... -run TestComputeContentHash -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/sync.go internal/cli/sync_test.go
git commit -m "feat(cli): add computeContentHash function for content fingerprinting"
```

---

### Task 3: Add getLastContentHash helper

**Files:**
- Modify: `internal/cli/sync.go`
- Test: `internal/cli/sync_test.go`

**Step 1: Write the failing test**

Add to `internal/cli/sync_test.go`:

```go
func TestGetLastContentHash(t *testing.T) {
	t.Run("returns empty string when no history", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		hash := getLastContentHash(tmpDir, "sword-001")
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

		hash := getLastContentHash(tmpDir, "sword-001")
		if hash != "xyz789uvw012" {
			t.Errorf("Expected 'xyz789uvw012', got %q", hash)
		}
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli/... -run TestGetLastContentHash -v`
Expected: FAIL with "undefined: getLastContentHash"

**Step 3: Implement getLastContentHash**

Add to `internal/cli/sync.go`:

```go
// getLastContentHash retrieves the most recent content hash for a node from history
// Returns empty string if no hash found
func getLastContentHash(targetDir, nodeID string) string {
	historyRepo := history.NewYAMLRepository(targetDir)

	entries, err := historyRepo.Query(history.Filter{NodeID: nodeID})
	if err != nil || len(entries) == 0 {
		return ""
	}

	// Entries are sorted oldest first, get the last one
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].ContentHash != "" {
			return entries[i].ContentHash
		}
	}

	return ""
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli/... -run TestGetLastContentHash -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/sync.go internal/cli/sync_test.go
git commit -m "feat(cli): add getLastContentHash to retrieve hash from history"
```

---

### Task 4: Rewrite runSync with hash-based detection

**Files:**
- Modify: `internal/cli/sync.go`
- Test: `internal/cli/sync_test.go`

**Step 1: Update test helpers to not require git**

In `internal/cli/sync_test.go`, update `TestRunSync_NotGitRepo` test:

```go
// DELETE this entire test - sync no longer requires git
func TestRunSync_NotGitRepo(t *testing.T) {
	// REMOVE - no longer applicable
}
```

**Step 2: Rewrite runSync**

Replace the `runSync` function in `internal/cli/sync.go`:

```go
func runSync(flags *syncFlags) (int, error) {
	// Verify we're in a deco project
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return syncExitError, fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Discover all nodes
	discovery := node.NewFileDiscovery(flags.targetDir)
	nodePaths, err := discovery.DiscoverAll()
	if err != nil {
		return syncExitError, fmt.Errorf("failed to discover nodes: %w", err)
	}

	if len(nodePaths) == 0 {
		return syncExitClean, nil
	}

	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	var syncResults []syncResult
	var baselinedNodes []string

	for _, nodePath := range nodePaths {
		nodeID := discovery.PathToID(nodePath)
		if nodeID == "" {
			continue
		}

		currentNode, err := nodeRepo.Load(nodeID)
		if err != nil {
			continue
		}

		currentHash := computeContentHash(currentNode)
		lastHash := getLastContentHash(flags.targetDir, nodeID)

		if lastHash == "" {
			// No history - baseline this node
			if !flags.dryRun {
				if err := logBaselineOperation(flags.targetDir, nodeID, currentHash); err != nil {
					if !flags.quiet {
						fmt.Fprintf(os.Stderr, "Warning: failed to baseline %s: %v\n", nodeID, err)
					}
					continue
				}
			}
			baselinedNodes = append(baselinedNodes, nodeID)
			continue
		}

		if currentHash == lastHash {
			// No change
			continue
		}

		// Content changed - need to sync
		result := syncResult{
			nodeID:     currentNode.ID,
			oldVersion: currentNode.Version,
			newVersion: currentNode.Version + 1,
			oldStatus:  currentNode.Status,
		}

		if !flags.dryRun {
			if err := applySyncWithHash(flags.targetDir, &currentNode, nodeRepo, currentHash); err != nil {
				if !flags.quiet {
					fmt.Fprintf(os.Stderr, "Warning: failed to sync %s: %v\n", nodeID, err)
				}
				continue
			}
		}

		syncResults = append(syncResults, result)
	}

	// Output results
	if !flags.quiet {
		if len(baselinedNodes) > 0 {
			if flags.dryRun {
				fmt.Printf("Would baseline: %s (%d nodes)\n", strings.Join(baselinedNodes, ", "), len(baselinedNodes))
			} else {
				fmt.Printf("Baselined: %s (%d nodes)\n", strings.Join(baselinedNodes, ", "), len(baselinedNodes))
			}
		}

		if len(syncResults) > 0 {
			if flags.dryRun {
				fmt.Print("Would sync: ")
			} else {
				fmt.Print("Synced: ")
			}
			parts := make([]string, len(syncResults))
			for i, r := range syncResults {
				parts[i] = fmt.Sprintf("%s (v%d→v%d)", r.nodeID, r.oldVersion, r.newVersion)
			}
			fmt.Println(strings.Join(parts, ", "))
		}
	}

	if flags.dryRun {
		return syncExitClean, nil
	}

	if len(syncResults) > 0 {
		return syncExitModified, nil
	}

	return syncExitClean, nil
}
```

**Step 3: Add logBaselineOperation and applySyncWithHash helpers**

Add to `internal/cli/sync.go`:

```go
// logBaselineOperation records initial state for a node without modification
func logBaselineOperation(targetDir, nodeID, contentHash string) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      nodeID,
		Operation:   "baseline",
		User:        username,
		ContentHash: contentHash,
	}

	return historyRepo.Append(entry)
}

// applySyncWithHash applies sync changes and logs with content hash
func applySyncWithHash(targetDir string, n *domain.Node, nodeRepo *node.YAMLRepository, contentHash string) error {
	oldVersion := n.Version
	oldStatus := n.Status

	// Bump version
	n.Version++

	// Reset status if was approved or review
	if n.Status == "approved" || n.Status == "review" {
		n.Status = "draft"
	}

	// Clear reviewers
	n.Reviewers = nil

	// Save the node
	if err := nodeRepo.Save(*n); err != nil {
		return err
	}

	// Log to history with content hash
	return logSyncOperationWithHash(targetDir, n.ID, oldVersion, n.Version, oldStatus, n.Status, contentHash)
}

// logSyncOperationWithHash adds a sync entry with content hash
func logSyncOperationWithHash(targetDir, nodeID string, oldVersion, newVersion int, oldStatus, newStatus, contentHash string) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      nodeID,
		Operation:   "sync",
		User:        username,
		ContentHash: contentHash,
		Before: map[string]interface{}{
			"version": oldVersion,
			"status":  oldStatus,
		},
		After: map[string]interface{}{
			"version": newVersion,
			"status":  newStatus,
		},
	}

	return historyRepo.Append(entry)
}
```

**Step 4: Remove old git-based functions**

Remove these functions from `internal/cli/sync.go`:
- `isGitRepo`
- `getModifiedNodeFiles`
- `extractNodeID`
- `getNodeFromHEAD`
- `applySync` (replaced by applySyncWithHash)
- `logSyncOperation` (replaced by logSyncOperationWithHash)

Remove `"os/exec"` and `"path/filepath"` from imports (if no longer needed).

**Step 5: Update imports**

Ensure imports include:
```go
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)
```

**Step 6: Run all sync tests**

Run: `go test ./internal/cli/... -run Sync -v`
Expected: Some tests may fail - we'll fix them in Task 5

**Step 7: Commit**

```bash
git add internal/cli/sync.go
git commit -m "feat(cli): rewrite sync to use content hash detection instead of git"
```

---

### Task 5: Update sync tests for hash-based detection

**Files:**
- Modify: `internal/cli/sync_test.go`

**Step 1: Remove git-dependent tests**

Delete these tests from `internal/cli/sync_test.go`:
- `TestExtractNodeID`
- `TestIsGitRepo`
- `TestRunSync_NotGitRepo`

**Step 2: Update remaining tests to use hash-based detection**

Replace `TestRunSync_NoChanges`:

```go
func TestRunSync_NoChanges(t *testing.T) {
	t.Run("returns clean exit when hashes match", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Create baseline history entry with current hash
		nodeRepo := node.NewYAMLRepository(tmpDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := computeContentHash(n)

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
```

Replace `TestRunSync_MetadataOnlyChange`:

```go
func TestRunSync_MetadataOnlyChange(t *testing.T) {
	t.Run("ignores metadata-only changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Baseline with current hash
		nodeRepo := node.NewYAMLRepository(tmpDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := computeContentHash(n)

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
```

Replace `TestRunSync_ContentChange`:

```go
func TestRunSync_ContentChange(t *testing.T) {
	t.Run("syncs nodes with content changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupProjectForSync(t, tmpDir)

		// Baseline with original hash
		nodeRepo := node.NewYAMLRepository(tmpDir)
		n, _ := nodeRepo.Load("sword-001")
		hash := computeContentHash(n)

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
```

Add baseline test:

```go
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
```

Update `TestRunSync_DryRun` similarly.

**Step 3: Remove git helper from tests**

Remove `initGitRepo` and `commitAllChanges` functions from `sync_test.go` (no longer needed).

Remove `"os/exec"` from imports if no longer used.

**Step 4: Run all sync tests**

Run: `go test ./internal/cli/... -run Sync -v`
Expected: PASS

**Step 5: Run all tests**

Run: `go test ./...`
Expected: PASS

**Step 6: Commit**

```bash
git add internal/cli/sync_test.go
git commit -m "test(cli): update sync tests for hash-based detection"
```

---

### Task 6: Update command description

**Files:**
- Modify: `internal/cli/sync.go`

**Step 1: Update Long description**

Replace the Long description in NewSyncCommand:

```go
Long: `Detect manually-edited nodes and fix their metadata.

When nodes are edited directly (bypassing the CLI), their version and
review status may become stale. This command detects such changes by
comparing content hashes and:

1. Bumps the version number
2. Resets status to "draft" (if was approved/review)
3. Clears reviewers
4. Logs the sync operation to history

Nodes without history are automatically baselined (their current state
is recorded without modification).

Exit codes:
  0 - No changes needed (or baseline only)
  1 - Files modified, re-commit needed
  2 - Error (invalid project)

Examples:
  deco sync              # Sync nodes in current directory
  deco sync --dry-run    # Show what would change
  deco sync /path/to/project`,
```

**Step 2: Commit**

```bash
git add internal/cli/sync.go
git commit -m "docs(cli): update sync command description for hash-based detection"
```

---

### Task 7: Manual verification

**Step 1: Build and test manually**

```bash
cd examples/snake
go run ../../cmd/deco sync --dry-run
```

Expected: Should show "Would baseline: items/food, systems/core, systems/scoring (3 nodes)" or similar

**Step 2: Run actual sync**

```bash
go run ../../cmd/deco sync
```

Expected: Should baseline all nodes

**Step 3: Make a change and verify detection**

Edit `examples/snake/.deco/nodes/items/food.yaml`, change title.

```bash
go run ../../cmd/deco sync --dry-run
```

Expected: Should show "Would sync: items/food (v1→v2)"

**Step 4: Final test run**

```bash
go test ./...
```

Expected: All tests pass

**Step 5: Final commit**

```bash
git add -A
git commit -m "feat(cli): complete hash-based sync detection

Replaces git-based detection with content hashing stored in history.
Works retroactively on already-committed changes.

- Add ContentHash field to AuditEntry
- Add baseline operation for recording initial state
- Remove git dependency from sync command
- Auto-baseline nodes without history"
```
