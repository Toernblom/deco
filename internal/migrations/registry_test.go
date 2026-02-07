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
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()

	m := Migration{
		Name:        "test-migration",
		Description: "A test migration",
		SourceHash:  "abc123",
		TargetHash:  "def456",
	}

	if err := r.Register(m); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Registering duplicate should fail
	if err := r.Register(m); err == nil {
		t.Error("expected error when registering duplicate migration")
	}
}

func TestRegistry_Find(t *testing.T) {
	r := NewRegistry()

	m1 := Migration{
		Name:       "v1-to-v2",
		SourceHash: "v1hash",
		TargetHash: "v2hash",
	}
	m2 := Migration{
		Name:       "any-to-v3",
		SourceHash: "", // Matches any source
		TargetHash: "v3hash",
	}

	r.Register(m1)
	r.Register(m2)

	// Find specific migration
	found := r.Find("v1hash", "v2hash")
	if found == nil {
		t.Error("expected to find v1-to-v2 migration")
	}
	if found != nil && found.Name != "v1-to-v2" {
		t.Errorf("wrong migration found: %s", found.Name)
	}

	// Find migration with any source
	found = r.Find("anything", "v3hash")
	if found == nil {
		t.Error("expected to find any-to-v3 migration")
	}

	// No migration found
	found = r.Find("v1hash", "v4hash")
	if found != nil {
		t.Errorf("expected nil, got %s", found.Name)
	}
}

func TestRegistry_FindPath(t *testing.T) {
	r := NewRegistry()

	// Create a chain: v1 -> v2 -> v3
	r.Register(Migration{Name: "v1-v2", SourceHash: "v1", TargetHash: "v2"})
	r.Register(Migration{Name: "v2-v3", SourceHash: "v2", TargetHash: "v3"})

	// Direct path
	path := r.FindPath("v1", "v2")
	if len(path) != 1 {
		t.Errorf("expected 1-step path, got %d", len(path))
	}

	// Multi-step path
	path = r.FindPath("v1", "v3")
	if len(path) != 2 {
		t.Errorf("expected 2-step path, got %d", len(path))
	}

	// No path
	path = r.FindPath("v1", "v4")
	if path != nil {
		t.Error("expected nil path for unreachable target")
	}

	// Already at target
	path = r.FindPath("v2", "v2")
	if path != nil {
		t.Error("expected nil path when already at target")
	}
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()

	r.Register(Migration{Name: "m1", SourceHash: "a", TargetHash: "b"})
	r.Register(Migration{Name: "m2", SourceHash: "b", TargetHash: "c"})

	list := r.List()
	if len(list) != 2 {
		t.Errorf("expected 2 migrations, got %d", len(list))
	}
}

func TestRegistry_Clear(t *testing.T) {
	r := NewRegistry()

	r.Register(Migration{Name: "m1", SourceHash: "a", TargetHash: "b"})
	r.Clear()

	if len(r.List()) != 0 {
		t.Error("expected empty registry after Clear")
	}
}

func TestIdentityMigration(t *testing.T) {
	m := IdentityMigration("update-v2", "Update to v2", "v1", "v2")

	if m.Name != "update-v2" {
		t.Errorf("wrong name: %s", m.Name)
	}
	if m.Transform != nil {
		t.Error("identity migration should have nil transform")
	}
}

func TestMigration_Transform(t *testing.T) {
	// Create a migration that adds a tag
	addTagMigration := Migration{
		Name:       "add-migrated-tag",
		SourceHash: "v1",
		TargetHash: "v2",
		Transform: func(node domain.Node) (domain.Node, error) {
			result := node
			result.Tags = append(result.Tags, "migrated")
			return result, nil
		},
	}

	original := domain.Node{
		ID:    "test-1",
		Kind:  "item",
		Title: "Test",
		Tags:  []string{"existing"},
	}

	transformed, err := addTagMigration.Transform(original)
	if err != nil {
		t.Fatalf("transform failed: %v", err)
	}

	if len(transformed.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(transformed.Tags))
	}
	if transformed.Tags[1] != "migrated" {
		t.Errorf("expected 'migrated' tag, got %s", transformed.Tags[1])
	}

	// Original should be unchanged
	if len(original.Tags) != 1 {
		t.Error("original node was modified")
	}
}
