package history

import "github.com/Toernblom/deco/internal/domain"

// Filter defines criteria for querying the audit log.
type Filter struct {
	// NodeID filters entries for a specific node (empty means all nodes).
	NodeID string

	// Operation filters by operation type (e.g., "create", "update").
	Operation string

	// User filters by user who made the change.
	User string

	// Since filters entries after this timestamp (Unix seconds).
	Since int64

	// Until filters entries before this timestamp (Unix seconds).
	Until int64

	// Limit restricts the number of results (0 means no limit).
	Limit int
}

// Repository defines the interface for audit log persistence.
// The audit log is append-only for data integrity.
type Repository interface {
	// Append adds a new entry to the audit log.
	// Entries are immutable once appended.
	Append(entry domain.AuditEntry) error

	// Query retrieves audit entries matching the filter criteria.
	// Returns entries in chronological order (oldest first).
	Query(filter Filter) ([]domain.AuditEntry, error)
}
