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
