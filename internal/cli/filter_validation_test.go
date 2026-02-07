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

package cli

import (
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/storage/config"
)

func TestValidateBlockType(t *testing.T) {
	t.Run("empty block type is valid", func(t *testing.T) {
		err := validateBlockType("", nil)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("built-in types are valid", func(t *testing.T) {
		for _, bt := range validBuiltInBlockTypes {
			err := validateBlockType(bt, nil)
			if err != nil {
				t.Errorf("Expected %q to be valid, got %v", bt, err)
			}
		}
	})

	t.Run("custom types are valid", func(t *testing.T) {
		custom := map[string]config.BlockTypeConfig{
			"building": {},
			"recipe":   {},
		}
		err := validateBlockType("building", custom)
		if err != nil {
			t.Errorf("Expected custom type 'building' to be valid, got %v", err)
		}
	})

	t.Run("unknown type returns error", func(t *testing.T) {
		err := validateBlockType("not_a_real_type", nil)
		if err == nil {
			t.Fatal("Expected error for unknown block type, got nil")
		}
		if !strings.Contains(err.Error(), "unknown block-type") {
			t.Errorf("Expected 'unknown block-type' in error, got %q", err.Error())
		}
		if !strings.Contains(err.Error(), "not_a_real_type") {
			t.Errorf("Expected 'not_a_real_type' in error, got %q", err.Error())
		}
	})

	t.Run("unknown type lists valid values", func(t *testing.T) {
		err := validateBlockType("xyz", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		for _, bt := range validBuiltInBlockTypes {
			if !strings.Contains(err.Error(), bt) {
				t.Errorf("Expected valid type %q in error message, got %q", bt, err.Error())
			}
		}
	})

	t.Run("unknown type with custom types lists all valid values", func(t *testing.T) {
		custom := map[string]config.BlockTypeConfig{
			"building": {},
		}
		err := validateBlockType("xyz", custom)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "building") {
			t.Errorf("Expected custom type 'building' in error, got %q", err.Error())
		}
	})

	t.Run("close typo suggests correction", func(t *testing.T) {
		err := validateBlockType("tabel", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "Did you mean") {
			t.Errorf("Expected did-you-mean suggestion, got %q", err.Error())
		}
	})
}
