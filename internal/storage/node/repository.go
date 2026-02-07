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

package node

import "github.com/Toernblom/deco/internal/domain"

// Repository defines the interface for node persistence operations.
// Implementations handle loading and saving nodes to/from storage (filesystem, DB, etc.).
type Repository interface {
	// LoadAll loads all nodes from storage.
	LoadAll() ([]domain.Node, error)

	// Load retrieves a single node by its ID.
	// Returns the node if found, or an error if not found or on failure.
	Load(id string) (domain.Node, error)

	// Save persists a node to storage.
	// Creates a new node if it doesn't exist, updates if it does.
	Save(node domain.Node) error

	// Delete removes a node from storage by ID.
	// Returns an error if the node doesn't exist or on failure.
	Delete(id string) error

	// Exists checks if a node with the given ID exists in storage.
	Exists(id string) (bool, error)
}
