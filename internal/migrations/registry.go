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

package migrations

import (
	"fmt"
	"sync"

	"github.com/Toernblom/deco/internal/domain"
)

// TransformFunc is a function that transforms a node during migration.
// It receives the node and returns the transformed node (or error).
// The function should not modify the input node directly; return a new copy.
type TransformFunc func(node domain.Node) (domain.Node, error)

// Migration represents a single schema migration.
type Migration struct {
	// Name is a unique identifier for this migration.
	Name string
	// Description explains what this migration does.
	Description string
	// SourceHash is the schema hash this migration applies to.
	// Empty string matches any source (useful for initial migrations).
	SourceHash string
	// TargetHash is the schema hash after this migration.
	TargetHash string
	// Transform is the function that transforms each node.
	// If nil, this is an identity migration (updates schema version only).
	Transform TransformFunc
}

// Registry manages available migrations.
type Registry struct {
	mu         sync.RWMutex
	migrations []Migration
}

// NewRegistry creates a new migration registry.
func NewRegistry() *Registry {
	return &Registry{
		migrations: make([]Migration, 0),
	}
}

// DefaultRegistry is the global migration registry.
var DefaultRegistry = NewRegistry()

// Register adds a migration to the registry.
func (r *Registry) Register(m Migration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate names
	for _, existing := range r.migrations {
		if existing.Name == m.Name {
			return fmt.Errorf("migration %q already registered", m.Name)
		}
	}

	r.migrations = append(r.migrations, m)
	return nil
}

// Find returns a migration that transforms from sourceHash to targetHash.
// Returns nil if no direct migration exists.
func (r *Registry) Find(sourceHash, targetHash string) *Migration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.migrations {
		m := &r.migrations[i]
		// Match source (empty source matches any)
		if m.SourceHash != "" && m.SourceHash != sourceHash {
			continue
		}
		// Match target
		if m.TargetHash == targetHash {
			return m
		}
	}
	return nil
}

// FindPath returns a sequence of migrations to get from sourceHash to targetHash.
// Uses breadth-first search to find the shortest path.
// Returns nil if no path exists.
func (r *Registry) FindPath(sourceHash, targetHash string) []Migration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if sourceHash == targetHash {
		return nil // Already at target
	}

	// BFS to find path
	type state struct {
		hash string
		path []Migration
	}

	visited := make(map[string]bool)
	queue := []state{{hash: sourceHash, path: nil}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.hash] {
			continue
		}
		visited[current.hash] = true

		// Try each migration
		for _, m := range r.migrations {
			// Check if migration applies to current hash
			if m.SourceHash != "" && m.SourceHash != current.hash {
				continue
			}

			newPath := append([]Migration{}, current.path...)
			newPath = append(newPath, m)

			// Found target?
			if m.TargetHash == targetHash {
				return newPath
			}

			// Continue searching
			if !visited[m.TargetHash] {
				queue = append(queue, state{
					hash: m.TargetHash,
					path: newPath,
				})
			}
		}
	}

	return nil // No path found
}

// List returns all registered migrations.
func (r *Registry) List() []Migration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Migration, len(r.migrations))
	copy(result, r.migrations)
	return result
}

// Clear removes all migrations from the registry.
// Primarily useful for testing.
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.migrations = make([]Migration, 0)
}

// IdentityMigration creates a migration that updates schema version without transforming nodes.
// Useful when schema changes are additive (new optional fields) and nodes don't need modification.
func IdentityMigration(name, description, sourceHash, targetHash string) Migration {
	return Migration{
		Name:        name,
		Description: description,
		SourceHash:  sourceHash,
		TargetHash:  targetHash,
		Transform:   nil, // Identity transform
	}
}
