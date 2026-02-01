package history

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/Toernblom/deco/internal/domain"
)

// JSONLRepository implements Repository using JSONL (JSON Lines) format.
// This is an append-only log stored at .deco/history.jsonl
type JSONLRepository struct {
	rootDir string
	mu      sync.Mutex // Protects concurrent writes
}

// NewYAMLRepository creates a new JSONL-based history repository.
// The name is kept as NewYAMLRepository for consistency with test expectations,
// even though it uses JSONL format internally.
func NewYAMLRepository(rootDir string) *JSONLRepository {
	return &JSONLRepository{
		rootDir: rootDir,
	}
}

// historyFile returns the path to the history log file
func (r *JSONLRepository) historyFile() string {
	return filepath.Join(r.rootDir, ".deco", "history.jsonl")
}

// Append adds a new entry to the audit log.
// Entries are immutable once appended.
func (r *JSONLRepository) Append(entry domain.AuditEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Ensure .deco directory exists
	decoDir := filepath.Join(r.rootDir, ".deco")
	err := os.MkdirAll(decoDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create .deco directory: %w", err)
	}

	// Open file in append mode (create if doesn't exist)
	file, err := os.OpenFile(r.historyFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	// Marshal entry to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	// Append JSON line (with newline)
	_, err = file.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	return nil
}

// Query retrieves audit entries matching the filter criteria.
// Returns entries in chronological order (oldest first).
func (r *JSONLRepository) Query(filter Filter) ([]domain.AuditEntry, error) {
	historyPath := r.historyFile()

	// Check if file exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		// No history file = empty results
		return []domain.AuditEntry{}, nil
	}

	// Open file for reading
	file, err := os.Open(historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	var entries []domain.AuditEntry
	scanner := bufio.NewScanner(file)

	// Read each line
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue // Skip empty lines
		}

		var entry domain.AuditEntry
		err := json.Unmarshal(line, &entry)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal entry: %w", err)
		}

		// Apply filters
		if !matchesFilter(entry, filter) {
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading history file: %w", err)
	}

	// Sort by timestamp (oldest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	// Apply limit if specified
	if filter.Limit > 0 && len(entries) > filter.Limit {
		entries = entries[:filter.Limit]
	}

	return entries, nil
}

// QueryLatestHashes returns a map of nodeID -> latest content hash for all nodes.
// This reads the history file once, providing O(history) complexity instead of
// O(nodes × history) when querying each node individually.
func (r *JSONLRepository) QueryLatestHashes() (map[string]string, error) {
	historyPath := r.historyFile()

	// Check if file exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return map[string]string{}, nil
	}

	file, err := os.Open(historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	// Track latest hash and timestamp per node
	type nodeState struct {
		hash      string
		timestamp int64
	}
	latestByNode := make(map[string]*nodeState)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry domain.AuditEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal entry: %w", err)
		}

		// Skip entries without content hash
		if entry.ContentHash == "" {
			continue
		}

		ts := entry.Timestamp.Unix()
		if existing, ok := latestByNode[entry.NodeID]; ok {
			if ts > existing.timestamp {
				existing.hash = entry.ContentHash
				existing.timestamp = ts
			}
		} else {
			latestByNode[entry.NodeID] = &nodeState{
				hash:      entry.ContentHash,
				timestamp: ts,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading history file: %w", err)
	}

	// Convert to simple map
	result := make(map[string]string, len(latestByNode))
	for nodeID, state := range latestByNode {
		result[nodeID] = state.hash
	}

	return result, nil
}

// matchesFilter checks if an entry matches the filter criteria
func matchesFilter(entry domain.AuditEntry, filter Filter) bool {
	// Filter by NodeID
	if filter.NodeID != "" && entry.NodeID != filter.NodeID {
		return false
	}

	// Filter by Operation
	if filter.Operation != "" && entry.Operation != filter.Operation {
		return false
	}

	// Filter by User
	if filter.User != "" && entry.User != filter.User {
		return false
	}

	// Filter by Since (after this timestamp)
	if filter.Since > 0 && entry.Timestamp.Unix() < filter.Since {
		return false
	}

	// Filter by Until (before this timestamp)
	if filter.Until > 0 && entry.Timestamp.Unix() > filter.Until {
		return false
	}

	return true
}
