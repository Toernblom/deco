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

package history_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/history"
)

func TestYAMLRepository_Append(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    "systems/food",
		Operation: "create",
		User:      "alice",
		After: map[string]interface{}{
			"title": "Food System",
		},
	}

	err := repo.Append(entry)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		t.Error("Expected history file to be created")
	}
}

func TestYAMLRepository_Append_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	// Don't create .deco directory - test should create it
	repo := history.NewYAMLRepository(historyFile)

	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    "systems/food",
		Operation: "create",
		User:      "alice",
	}

	err := repo.Append(entry)
	if err != nil {
		t.Fatalf("Append should create directory: %v", err)
	}

	// Verify directory and file were created
	decoDir := filepath.Join(tmpDir, ".deco")
	if _, err := os.Stat(decoDir); os.IsNotExist(err) {
		t.Error("Expected .deco directory to be created")
	}
}

func TestYAMLRepository_Append_MultipleEntries(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Append multiple entries
	entries := []domain.AuditEntry{
		{
			Timestamp: time.Now(),
			NodeID:    "systems/food",
			Operation: "create",
			User:      "alice",
		},
		{
			Timestamp: time.Now().Add(time.Second),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "bob",
		},
		{
			Timestamp: time.Now().Add(2 * time.Second),
			NodeID:    "mechanics/combat",
			Operation: "create",
			User:      "alice",
		},
	}

	for _, entry := range entries {
		err := repo.Append(entry)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	// Query all and verify count
	results, err := repo.Query(history.Filter{})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(results))
	}
}

func TestYAMLRepository_Query_All(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Append test entries
	entries := []domain.AuditEntry{
		{
			Timestamp: time.Now(),
			NodeID:    "systems/food",
			Operation: "create",
			User:      "alice",
		},
		{
			Timestamp: time.Now().Add(time.Second),
			NodeID:    "mechanics/combat",
			Operation: "update",
			User:      "bob",
		},
	}

	for _, entry := range entries {
		err := repo.Append(entry)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	// Query all entries (empty filter)
	results, err := repo.Query(history.Filter{})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(results))
	}

	// Verify chronological order (oldest first)
	if !results[0].Timestamp.Before(results[1].Timestamp) &&
		!results[0].Timestamp.Equal(results[1].Timestamp) {
		t.Error("Expected entries in chronological order (oldest first)")
	}
}

func TestYAMLRepository_Query_EmptyLog(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Query empty log
	results, err := repo.Query(history.Filter{})
	if err != nil {
		t.Fatalf("Query should succeed on empty log: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 entries in empty log, got %d", len(results))
	}
}

func TestYAMLRepository_Query_FilterByNodeID(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Append entries for different nodes
	entries := []domain.AuditEntry{
		{
			Timestamp: time.Now(),
			NodeID:    "systems/food",
			Operation: "create",
			User:      "alice",
		},
		{
			Timestamp: time.Now().Add(time.Second),
			NodeID:    "systems/water",
			Operation: "create",
			User:      "alice",
		},
		{
			Timestamp: time.Now().Add(2 * time.Second),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "bob",
		},
	}

	for _, entry := range entries {
		err := repo.Append(entry)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	// Query for specific node
	results, err := repo.Query(history.Filter{NodeID: "systems/food"})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries for systems/food, got %d", len(results))
	}

	// Verify all results match the filter
	for _, entry := range results {
		if entry.NodeID != "systems/food" {
			t.Errorf("Expected NodeID 'systems/food', got %q", entry.NodeID)
		}
	}
}

func TestYAMLRepository_Query_FilterByOperation(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Append entries with different operations
	entries := []domain.AuditEntry{
		{
			Timestamp: time.Now(),
			NodeID:    "systems/food",
			Operation: "create",
			User:      "alice",
		},
		{
			Timestamp: time.Now().Add(time.Second),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "alice",
		},
		{
			Timestamp: time.Now().Add(2 * time.Second),
			NodeID:    "systems/water",
			Operation: "create",
			User:      "bob",
		},
	}

	for _, entry := range entries {
		err := repo.Append(entry)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	// Query for specific operation
	results, err := repo.Query(history.Filter{Operation: "create"})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 create operations, got %d", len(results))
	}

	// Verify all results match the filter
	for _, entry := range results {
		if entry.Operation != "create" {
			t.Errorf("Expected Operation 'create', got %q", entry.Operation)
		}
	}
}

func TestYAMLRepository_Query_FilterByUser(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Append entries from different users
	entries := []domain.AuditEntry{
		{
			Timestamp: time.Now(),
			NodeID:    "systems/food",
			Operation: "create",
			User:      "alice",
		},
		{
			Timestamp: time.Now().Add(time.Second),
			NodeID:    "systems/water",
			Operation: "create",
			User:      "bob",
		},
		{
			Timestamp: time.Now().Add(2 * time.Second),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "alice",
		},
	}

	for _, entry := range entries {
		err := repo.Append(entry)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	// Query for specific user
	results, err := repo.Query(history.Filter{User: "alice"})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries for alice, got %d", len(results))
	}

	// Verify all results match the filter
	for _, entry := range results {
		if entry.User != "alice" {
			t.Errorf("Expected User 'alice', got %q", entry.User)
		}
	}
}

func TestYAMLRepository_Query_FilterByTimeRange(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Create entries with specific timestamps
	baseTime := time.Now()
	entries := []domain.AuditEntry{
		{
			Timestamp: baseTime,
			NodeID:    "systems/food",
			Operation: "create",
			User:      "alice",
		},
		{
			Timestamp: baseTime.Add(time.Hour),
			NodeID:    "systems/water",
			Operation: "create",
			User:      "bob",
		},
		{
			Timestamp: baseTime.Add(2 * time.Hour),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "alice",
		},
		{
			Timestamp: baseTime.Add(3 * time.Hour),
			NodeID:    "mechanics/combat",
			Operation: "create",
			User:      "charlie",
		},
	}

	for _, entry := range entries {
		err := repo.Append(entry)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	// Query with Since filter
	results, err := repo.Query(history.Filter{
		Since: baseTime.Add(90 * time.Minute).Unix(),
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries after Since, got %d", len(results))
	}

	// Query with Until filter
	results, err = repo.Query(history.Filter{
		Until: baseTime.Add(90 * time.Minute).Unix(),
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries before Until, got %d", len(results))
	}

	// Query with both Since and Until
	results, err = repo.Query(history.Filter{
		Since: baseTime.Add(30 * time.Minute).Unix(),
		Until: baseTime.Add(150 * time.Minute).Unix(),
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries in time range, got %d", len(results))
	}
}

func TestYAMLRepository_Query_Limit(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Append multiple entries
	for i := 0; i < 5; i++ {
		entry := domain.AuditEntry{
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "alice",
		}
		err := repo.Append(entry)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	// Query with limit
	results, err := repo.Query(history.Filter{Limit: 3})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 entries with limit, got %d", len(results))
	}
}

func TestYAMLRepository_Query_CombinedFilters(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Append diverse entries
	baseTime := time.Now()
	entries := []domain.AuditEntry{
		{
			Timestamp: baseTime,
			NodeID:    "systems/food",
			Operation: "create",
			User:      "alice",
		},
		{
			Timestamp: baseTime.Add(time.Hour),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "alice",
		},
		{
			Timestamp: baseTime.Add(2 * time.Hour),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "bob",
		},
		{
			Timestamp: baseTime.Add(3 * time.Hour),
			NodeID:    "systems/water",
			Operation: "update",
			User:      "alice",
		},
	}

	for _, entry := range entries {
		err := repo.Append(entry)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	// Query with multiple filters
	results, err := repo.Query(history.Filter{
		NodeID:    "systems/food",
		Operation: "update",
		User:      "alice",
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	// Should match only the second entry
	if len(results) != 1 {
		t.Errorf("Expected 1 entry matching all filters, got %d", len(results))
	}

	if len(results) > 0 {
		if results[0].NodeID != "systems/food" {
			t.Errorf("Expected NodeID 'systems/food', got %q", results[0].NodeID)
		}
		if results[0].Operation != "update" {
			t.Errorf("Expected Operation 'update', got %q", results[0].Operation)
		}
		if results[0].User != "alice" {
			t.Errorf("Expected User 'alice', got %q", results[0].User)
		}
	}
}

func TestYAMLRepository_Query_NoMatches(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Append some entries
	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    "systems/food",
		Operation: "create",
		User:      "alice",
	}
	err := repo.Append(entry)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Query for non-existent node
	results, err := repo.Query(history.Filter{NodeID: "nonexistent/node"})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 entries for non-existent node, got %d", len(results))
	}
}

func TestYAMLRepository_AppendOnly_ChronologicalOrder(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Append entries out of chronological order (by timestamp value)
	// But they should appear in the order they were appended
	entries := []domain.AuditEntry{
		{
			Timestamp: time.Now().Add(2 * time.Hour),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "alice",
		},
		{
			Timestamp: time.Now(),
			NodeID:    "systems/food",
			Operation: "create",
			User:      "alice",
		},
		{
			Timestamp: time.Now().Add(time.Hour),
			NodeID:    "systems/food",
			Operation: "update",
			User:      "bob",
		},
	}

	for _, entry := range entries {
		err := repo.Append(entry)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	// Query all and verify they're returned in append order
	results, err := repo.Query(history.Filter{})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(results))
	}

	// The implementation should return entries sorted by timestamp (oldest first)
	// So we should see: create (now), update (now+1h), update (now+2h)
	// NOT in append order
	expectedOrder := []string{"create", "update", "update"}
	for i, entry := range results {
		if entry.Operation != expectedOrder[i] {
			t.Errorf("Entry %d: expected operation %q, got %q", i, expectedOrder[i], entry.Operation)
		}
	}
}

func TestYAMLRepository_ConcurrentAppend(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Concurrent appends
	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			entry := domain.AuditEntry{
				Timestamp: time.Now(),
				NodeID:    "systems/test",
				Operation: "update",
				User:      "concurrent",
			}
			err := repo.Append(entry)
			if err != nil {
				t.Errorf("Concurrent append %d failed: %v", idx, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify all entries were appended
	results, err := repo.Query(history.Filter{})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != numGoroutines {
		t.Errorf("Expected %d entries from concurrent appends, got %d", numGoroutines, len(results))
	}
}

func TestYAMLRepository_QueryLatestHashes(t *testing.T) {
	t.Run("returns empty map for no history", func(t *testing.T) {
		tmpDir := t.TempDir()
		historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
		repo := history.NewYAMLRepository(historyFile)

		hashes, err := repo.QueryLatestHashes()
		if err != nil {
			t.Fatalf("QueryLatestHashes failed: %v", err)
		}
		if len(hashes) != 0 {
			t.Errorf("Expected empty map, got %d entries", len(hashes))
		}
	})

	t.Run("returns latest hash for each node", func(t *testing.T) {
		tmpDir := t.TempDir()
		historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
		repo := history.NewYAMLRepository(historyFile)

		// Append entries with content hashes
		entries := []domain.AuditEntry{
			{
				Timestamp:   time.Now(),
				NodeID:      "node-a",
				Operation:   "create",
				User:        "alice",
				ContentHash: "hash-a-1",
			},
			{
				Timestamp:   time.Now().Add(time.Second),
				NodeID:      "node-b",
				Operation:   "create",
				User:        "alice",
				ContentHash: "hash-b-1",
			},
			{
				Timestamp:   time.Now().Add(2 * time.Second),
				NodeID:      "node-a",
				Operation:   "update",
				User:        "alice",
				ContentHash: "hash-a-2", // More recent hash for node-a
			},
		}

		for _, entry := range entries {
			if err := repo.Append(entry); err != nil {
				t.Fatalf("Append failed: %v", err)
			}
		}

		hashes, err := repo.QueryLatestHashes()
		if err != nil {
			t.Fatalf("QueryLatestHashes failed: %v", err)
		}

		if len(hashes) != 2 {
			t.Errorf("Expected 2 nodes, got %d", len(hashes))
		}
		if hashes["node-a"] != "hash-a-2" {
			t.Errorf("Expected hash-a-2 for node-a, got %q", hashes["node-a"])
		}
		if hashes["node-b"] != "hash-b-1" {
			t.Errorf("Expected hash-b-1 for node-b, got %q", hashes["node-b"])
		}
	})

	t.Run("skips entries without content hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
		repo := history.NewYAMLRepository(historyFile)

		entries := []domain.AuditEntry{
			{
				Timestamp:   time.Now(),
				NodeID:      "node-a",
				Operation:   "create",
				User:        "alice",
				ContentHash: "hash-a-1",
			},
			{
				Timestamp:   time.Now().Add(time.Second),
				NodeID:      "node-a",
				Operation:   "set", // No content hash
				User:        "alice",
			},
		}

		for _, entry := range entries {
			if err := repo.Append(entry); err != nil {
				t.Fatalf("Append failed: %v", err)
			}
		}

		hashes, err := repo.QueryLatestHashes()
		if err != nil {
			t.Fatalf("QueryLatestHashes failed: %v", err)
		}

		// Should still return hash-a-1 since the later entry has no hash
		if hashes["node-a"] != "hash-a-1" {
			t.Errorf("Expected hash-a-1 for node-a, got %q", hashes["node-a"])
		}
	})
}

func TestYAMLRepository_PreserveComplexData(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".deco", "history.jsonl")
	repo := history.NewYAMLRepository(historyFile)

	// Create entry with complex Before/After data
	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    "systems/food",
		Operation: "update",
		User:      "alice",
		Before: map[string]interface{}{
			"title":  "Old Food System",
			"status": "draft",
			"tags":   []interface{}{"old", "deprecated"},
		},
		After: map[string]interface{}{
			"title":  "New Food System",
			"status": "approved",
			"tags":   []interface{}{"new", "active"},
			"metadata": map[string]interface{}{
				"version": 2,
				"author":  "alice",
			},
		},
	}

	err := repo.Append(entry)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Query and verify data integrity
	results, err := repo.Query(history.Filter{NodeID: "systems/food"})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(results))
	}

	result := results[0]
	if result.Before["title"] != "Old Food System" {
		t.Errorf("Before data not preserved correctly")
	}
	if result.After["title"] != "New Food System" {
		t.Errorf("After data not preserved correctly")
	}

	// Verify nested metadata
	if afterMeta, ok := result.After["metadata"].(map[string]interface{}); ok {
		if afterMeta["author"] != "alice" {
			t.Errorf("Nested metadata not preserved correctly")
		}
	} else {
		t.Error("Expected metadata to be a map")
	}
}
