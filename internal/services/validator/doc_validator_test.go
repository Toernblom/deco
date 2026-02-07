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

package validator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/errors"
	"github.com/Toernblom/deco/internal/services/validator"
)

func TestDocValidator_NodeLevelDocs_FileExists(t *testing.T) {
	dir := t.TempDir()
	// Create a .md file
	mdPath := filepath.Join(dir, "narrative.md")
	os.WriteFile(mdPath, []byte("# Chapter 1\nThe protagonist enters the ancient temple."), 0644)

	node := &domain.Node{
		ID: "stories/chapter-1",
		Docs: []domain.DocRef{
			{Path: "narrative.md"},
		},
	}

	collector := errors.NewCollector()
	dv := validator.NewDocValidator()
	dv.ValidateNodeDocs(node, dir, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors, got: %v", collector.Errors())
	}
}

func TestDocValidator_NodeLevelDocs_FileMissing(t *testing.T) {
	dir := t.TempDir()

	node := &domain.Node{
		ID: "stories/chapter-1",
		Docs: []domain.DocRef{
			{Path: "missing.md"},
		},
	}

	collector := errors.NewCollector()
	dv := validator.NewDocValidator()
	dv.ValidateNodeDocs(node, dir, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for missing file")
	}
	errs := collector.Errors()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Code != "E055" {
		t.Errorf("expected E055, got %s", errs[0].Code)
	}
}

func TestDocValidator_NodeLevelDocs_KeywordsPresent(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "narrative.md")
	os.WriteFile(mdPath, []byte("The protagonist enters the ancient temple and discovers a betrayal."), 0644)

	node := &domain.Node{
		ID: "stories/chapter-1",
		Docs: []domain.DocRef{
			{
				Path:     "narrative.md",
				Keywords: []string{"protagonist", "ancient temple", "betrayal"},
			},
		},
	}

	collector := errors.NewCollector()
	dv := validator.NewDocValidator()
	dv.ValidateNodeDocs(node, dir, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors, got: %v", collector.Errors())
	}
}

func TestDocValidator_NodeLevelDocs_KeywordMissing(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "narrative.md")
	os.WriteFile(mdPath, []byte("The hero enters the castle."), 0644)

	node := &domain.Node{
		ID: "stories/chapter-1",
		Docs: []domain.DocRef{
			{
				Path:     "narrative.md",
				Keywords: []string{"protagonist", "betrayal"},
			},
		},
	}

	collector := errors.NewCollector()
	dv := validator.NewDocValidator()
	dv.ValidateNodeDocs(node, dir, collector)

	if !collector.HasErrors() {
		t.Fatal("expected errors for missing keywords")
	}
	errs := collector.Errors()
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors (one per missing keyword), got %d", len(errs))
	}
	for _, err := range errs {
		if err.Code != "E056" {
			t.Errorf("expected E056, got %s", err.Code)
		}
	}
}

func TestDocValidator_NodeLevelDocs_KeywordsCaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "narrative.md")
	os.WriteFile(mdPath, []byte("The PROTAGONIST found the Ancient Temple."), 0644)

	node := &domain.Node{
		ID: "stories/chapter-1",
		Docs: []domain.DocRef{
			{
				Path:     "narrative.md",
				Keywords: []string{"protagonist", "ancient temple"},
			},
		},
	}

	collector := errors.NewCollector()
	dv := validator.NewDocValidator()
	dv.ValidateNodeDocs(node, dir, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors (case-insensitive), got: %v", collector.Errors())
	}
}

func TestDocValidator_NodeLevelDocs_EmptyKeywords(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "narrative.md")
	os.WriteFile(mdPath, []byte("Some content."), 0644)

	node := &domain.Node{
		ID: "stories/chapter-1",
		Docs: []domain.DocRef{
			{Path: "narrative.md", Keywords: []string{}},
		},
	}

	collector := errors.NewCollector()
	dv := validator.NewDocValidator()
	dv.ValidateNodeDocs(node, dir, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors with empty keywords, got: %v", collector.Errors())
	}
}

func TestDocValidator_NodeLevelDocs_NilNode(t *testing.T) {
	dv := validator.NewDocValidator()
	collector := errors.NewCollector()
	dv.ValidateNodeDocs(nil, "/tmp", collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for nil node")
	}
}

func TestDocValidator_NodeLevelDocs_NoDocs(t *testing.T) {
	node := &domain.Node{ID: "test/node"}
	dv := validator.NewDocValidator()
	collector := errors.NewCollector()
	dv.ValidateNodeDocs(node, "/tmp", collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for node without docs")
	}
}

func TestDocValidator_BlockDocs_FileExists(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "scene.md")
	os.WriteFile(mdPath, []byte("A storm rages outside."), 0644)

	block := domain.Block{
		Type: "doc",
		Data: map[string]interface{}{
			"path": "scene.md",
		},
	}

	dv := validator.NewDocValidator()
	collector := errors.NewCollector()
	dv.ValidateDocBlock(&block, "test/node", "Narrative", 0, dir, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors, got: %v", collector.Errors())
	}
}

func TestDocValidator_BlockDocs_FileMissing(t *testing.T) {
	dir := t.TempDir()

	block := domain.Block{
		Type: "doc",
		Data: map[string]interface{}{
			"path": "missing.md",
		},
	}

	dv := validator.NewDocValidator()
	collector := errors.NewCollector()
	dv.ValidateDocBlock(&block, "test/node", "Narrative", 0, dir, collector)

	if !collector.HasErrors() {
		t.Fatal("expected error for missing file")
	}
	errs := collector.Errors()
	if errs[0].Code != "E055" {
		t.Errorf("expected E055, got %s", errs[0].Code)
	}
}

func TestDocValidator_BlockDocs_KeywordMissing(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "scene.md")
	os.WriteFile(mdPath, []byte("A calm day in the village."), 0644)

	block := domain.Block{
		Type: "doc",
		Data: map[string]interface{}{
			"path":     "scene.md",
			"keywords": []interface{}{"storm", "mysterious stranger"},
		},
	}

	dv := validator.NewDocValidator()
	collector := errors.NewCollector()
	dv.ValidateDocBlock(&block, "test/node", "Narrative", 0, dir, collector)

	if !collector.HasErrors() {
		t.Fatal("expected errors for missing keywords")
	}
	errs := collector.Errors()
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
	for _, err := range errs {
		if err.Code != "E056" {
			t.Errorf("expected E056, got %s", err.Code)
		}
	}
}

func TestDocValidator_BlockDocs_NoPath(t *testing.T) {
	block := domain.Block{
		Type: "doc",
		Data: map[string]interface{}{},
	}

	dv := validator.NewDocValidator()
	collector := errors.NewCollector()
	dv.ValidateDocBlock(&block, "test/node", "Narrative", 0, "/tmp", collector)

	// No path means nothing to validate for doc content (block validator handles required field)
	if collector.HasErrors() {
		t.Errorf("expected no errors when path missing (block validator handles that)")
	}
}

func TestDocValidator_BlockDocs_SubdirectoryPath(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "narratives")
	os.MkdirAll(subDir, 0755)
	mdPath := filepath.Join(subDir, "chapter-1.md")
	os.WriteFile(mdPath, []byte("The story begins."), 0644)

	block := domain.Block{
		Type: "doc",
		Data: map[string]interface{}{
			"path": "narratives/chapter-1.md",
		},
	}

	dv := validator.NewDocValidator()
	collector := errors.NewCollector()
	dv.ValidateDocBlock(&block, "test/node", "Narrative", 0, dir, collector)

	if collector.HasErrors() {
		t.Errorf("expected no errors for subdirectory path, got: %v", collector.Errors())
	}
}

func TestDocValidator_MultipleDocs(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "intro.md"), []byte("The story protagonist begins a journey."), 0644)
	// note: outro.md intentionally not created

	node := &domain.Node{
		ID: "stories/arc",
		Docs: []domain.DocRef{
			{Path: "intro.md", Keywords: []string{"protagonist"}},
			{Path: "outro.md", Keywords: []string{"resolution"}},
		},
	}

	dv := validator.NewDocValidator()
	collector := errors.NewCollector()
	dv.ValidateNodeDocs(node, dir, collector)

	errs := collector.Errors()
	// Should have E055 for missing outro.md (keyword check skipped for missing files)
	foundE055 := false
	for _, err := range errs {
		if err.Code == "E055" {
			foundE055 = true
		}
	}
	if !foundE055 {
		t.Error("expected E055 for missing outro.md")
	}
}
