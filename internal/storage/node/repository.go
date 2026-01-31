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
