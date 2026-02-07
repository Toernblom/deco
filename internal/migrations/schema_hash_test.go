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

	"github.com/Toernblom/deco/internal/storage/config"
)

func TestComputeSchemaHash_Empty(t *testing.T) {
	cfg := config.Config{}
	hash := ComputeSchemaHash(cfg)
	if hash != "" {
		t.Errorf("expected empty hash for empty schema, got %q", hash)
	}
}

func TestComputeSchemaHash_Deterministic(t *testing.T) {
	cfg := config.Config{
		CustomBlockTypes: map[string]config.BlockTypeConfig{
			"dialogue": {RequiredFields: []string{"speaker", "text"}},
		},
		SchemaRules: map[string]config.SchemaRuleConfig{
			"character": {RequiredFields: []string{"name", "role"}},
		},
	}

	hash1 := ComputeSchemaHash(cfg)
	hash2 := ComputeSchemaHash(cfg)

	if hash1 != hash2 {
		t.Errorf("hash is not deterministic: %q != %q", hash1, hash2)
	}
	if hash1 == "" {
		t.Error("expected non-empty hash for non-empty schema")
	}
	if len(hash1) != 16 {
		t.Errorf("expected 16-char hash, got %d chars: %q", len(hash1), hash1)
	}
}

func TestComputeSchemaHash_SortedKeys(t *testing.T) {
	// Create config with keys in different order
	cfg1 := config.Config{
		CustomBlockTypes: map[string]config.BlockTypeConfig{
			"a": {RequiredFields: []string{"x", "y"}},
			"b": {RequiredFields: []string{"z"}},
		},
	}
	cfg2 := config.Config{
		CustomBlockTypes: map[string]config.BlockTypeConfig{
			"b": {RequiredFields: []string{"z"}},
			"a": {RequiredFields: []string{"y", "x"}}, // Different field order
		},
	}

	hash1 := ComputeSchemaHash(cfg1)
	hash2 := ComputeSchemaHash(cfg2)

	if hash1 != hash2 {
		t.Errorf("hash should be same regardless of key order: %q != %q", hash1, hash2)
	}
}

func TestComputeSchemaHash_DifferentSchemas(t *testing.T) {
	cfg1 := config.Config{
		SchemaRules: map[string]config.SchemaRuleConfig{
			"item": {RequiredFields: []string{"name"}},
		},
	}
	cfg2 := config.Config{
		SchemaRules: map[string]config.SchemaRuleConfig{
			"item": {RequiredFields: []string{"name", "cost"}},
		},
	}

	hash1 := ComputeSchemaHash(cfg1)
	hash2 := ComputeSchemaHash(cfg2)

	if hash1 == hash2 {
		t.Errorf("different schemas should produce different hashes: %q == %q", hash1, hash2)
	}
}

func TestComputeSchemaHash_CustomBlockOptionalFields(t *testing.T) {
	cfg1 := config.Config{
		CustomBlockTypes: map[string]config.BlockTypeConfig{
			"quest": {RequiredFields: []string{"name"}},
		},
	}
	cfg2 := config.Config{
		CustomBlockTypes: map[string]config.BlockTypeConfig{
			"quest": {RequiredFields: []string{"name"}, OptionalFields: []string{"reward"}},
		},
	}

	hash1 := ComputeSchemaHash(cfg1)
	hash2 := ComputeSchemaHash(cfg2)

	if hash1 == hash2 {
		t.Errorf("optional fields should affect schema hash: %q == %q", hash1, hash2)
	}
}

func TestSchemaVersionMatches(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.Config
		want bool
	}{
		{
			name: "empty schema no version",
			cfg:  config.Config{},
			want: true,
		},
		{
			name: "schema with matching version",
			cfg: func() config.Config {
				c := config.Config{
					SchemaRules: map[string]config.SchemaRuleConfig{
						"quest": {RequiredFields: []string{"objective"}},
					},
				}
				c.SchemaVersion = ComputeSchemaHash(c)
				return c
			}(),
			want: true,
		},
		{
			name: "schema with mismatched version",
			cfg: config.Config{
				SchemaRules: map[string]config.SchemaRuleConfig{
					"quest": {RequiredFields: []string{"objective"}},
				},
				SchemaVersion: "wronghash1234567",
			},
			want: false,
		},
		{
			name: "schema without version",
			cfg: config.Config{
				SchemaRules: map[string]config.SchemaRuleConfig{
					"quest": {RequiredFields: []string{"objective"}},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SchemaVersionMatches(tt.cfg)
			if got != tt.want {
				t.Errorf("SchemaVersionMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}
