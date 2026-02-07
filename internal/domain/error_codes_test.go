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

package domain_test

import (
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestErrorCodeRegistry_Lookup(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()

	// Test looking up a valid code
	code, exists := registry.Lookup("E001")
	if !exists {
		t.Error("Expected E001 to exist in registry")
	}
	if code.Code != "E001" {
		t.Errorf("Expected code 'E001', got %q", code.Code)
	}
	if code.Category == "" {
		t.Error("Expected code to have a category")
	}
	if code.Message == "" {
		t.Error("Expected code to have a message")
	}

	// Test looking up non-existent code
	_, exists = registry.Lookup("E999")
	if exists {
		t.Error("Expected E999 to not exist in registry")
	}
}

func TestErrorCodeRegistry_CategoryRanges(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()

	tests := []struct {
		category string
		minCode  string
		maxCode  string
	}{
		{"schema", "E001", "E019"},
		{"refs", "E020", "E039"},
		{"validation", "E040", "E059"},
		{"io", "E060", "E079"},
		{"graph", "E080", "E099"},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			// Check minimum code in range
			code, exists := registry.Lookup(tt.minCode)
			if !exists {
				t.Errorf("Expected %s to exist in registry", tt.minCode)
			}
			if code.Category != tt.category {
				t.Errorf("Expected %s to be in category %q, got %q", tt.minCode, tt.category, code.Category)
			}

			// Check maximum code in range
			code, exists = registry.Lookup(tt.maxCode)
			if !exists {
				t.Errorf("Expected %s to exist in registry", tt.maxCode)
			}
			if code.Category != tt.category {
				t.Errorf("Expected %s to be in category %q, got %q", tt.maxCode, tt.category, code.Category)
			}
		})
	}
}

func TestErrorCodeRegistry_Uniqueness(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()

	seen := make(map[string]bool)
	codes := registry.AllCodes()

	for _, code := range codes {
		if seen[code.Code] {
			t.Errorf("Duplicate code found: %s", code.Code)
		}
		seen[code.Code] = true
	}
}

func TestErrorCodeRegistry_AllCodes(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()

	codes := registry.AllCodes()

	if len(codes) == 0 {
		t.Error("Expected registry to have at least one code")
	}

	// Verify all codes have required fields
	for _, code := range codes {
		if code.Code == "" {
			t.Error("Found code with empty Code field")
		}
		if code.Category == "" {
			t.Errorf("Code %s has empty Category field", code.Code)
		}
		if code.Message == "" {
			t.Errorf("Code %s has empty Message field", code.Code)
		}
	}
}

func TestErrorCode_Type(t *testing.T) {
	code := domain.ErrorCode{
		Code:     "E001",
		Category: "schema",
		Message:  "Invalid schema",
	}

	if code.Code != "E001" {
		t.Errorf("Expected Code 'E001', got %q", code.Code)
	}
	if code.Category != "schema" {
		t.Errorf("Expected Category 'schema', got %q", code.Category)
	}
	if code.Message != "Invalid schema" {
		t.Errorf("Expected Message 'Invalid schema', got %q", code.Message)
	}
}

func TestErrorCodeRegistry_ByCategory(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()

	// Test getting codes by category
	schemaCodes := registry.ByCategory("schema")
	if len(schemaCodes) == 0 {
		t.Error("Expected at least one schema error code")
	}

	// Verify all returned codes are in the correct category
	for _, code := range schemaCodes {
		if code.Category != "schema" {
			t.Errorf("Expected code %s to be in 'schema' category, got %q", code.Code, code.Category)
		}
	}

	// Test non-existent category
	nonExistent := registry.ByCategory("nonexistent")
	if len(nonExistent) != 0 {
		t.Error("Expected no codes for non-existent category")
	}
}

func TestErrorCodeRegistry_Categories(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()

	categories := registry.Categories()

	if len(categories) == 0 {
		t.Error("Expected at least one category")
	}

	// Check that all expected categories exist
	expectedCategories := []string{"schema", "refs", "validation", "io", "graph"}
	for _, expected := range expectedCategories {
		found := false
		for _, cat := range categories {
			if cat == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected category %q to exist", expected)
		}
	}
}
