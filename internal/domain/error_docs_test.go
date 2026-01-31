package domain_test

import (
	"strings"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
)

func TestErrorDocsGenerator_GenerateMarkdown(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()
	generator := domain.NewErrorDocsGenerator(registry)

	markdown := generator.GenerateMarkdown()

	if markdown == "" {
		t.Error("Expected non-empty markdown output")
	}

	// Should contain header
	if !strings.Contains(markdown, "# Deco Error Codes") {
		t.Error("Markdown should contain main header")
	}

	// Should contain category sections
	categories := []string{"Schema", "Refs", "Validation", "Io", "Graph"}
	for _, category := range categories {
		if !strings.Contains(markdown, category+" Errors") {
			t.Errorf("Markdown should contain %s section", category)
		}
	}

	// Should contain error codes
	if !strings.Contains(markdown, "E001") {
		t.Error("Markdown should contain error codes")
	}

	// Should contain table format
	if !strings.Contains(markdown, "| Code | Message |") {
		t.Error("Markdown should contain table header")
	}
}

func TestErrorDocsGenerator_GenerateDetailedMarkdown(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()
	generator := domain.NewErrorDocsGenerator(registry)

	markdown := generator.GenerateDetailedMarkdown()

	if markdown == "" {
		t.Error("Expected non-empty markdown output")
	}

	// Should contain header
	if !strings.Contains(markdown, "# Deco Error Codes Reference") {
		t.Error("Markdown should contain main header")
	}

	// Should contain table of contents
	if !strings.Contains(markdown, "## Table of Contents") {
		t.Error("Markdown should contain table of contents")
	}

	// Should contain detailed entries
	if !strings.Contains(markdown, "### E001:") {
		t.Error("Markdown should contain detailed error entries")
	}

	// Should contain descriptions
	if !strings.Contains(markdown, "**Description:**") {
		t.Error("Markdown should contain descriptions")
	}
}

func TestErrorDocsGenerator_GenerateCodeList(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()
	generator := domain.NewErrorDocsGenerator(registry)

	list := generator.GenerateCodeList()

	if list == "" {
		t.Error("Expected non-empty code list output")
	}

	// Should contain header
	if !strings.Contains(list, "# Error Code List") {
		t.Error("List should contain header")
	}

	// Should contain error codes in list format
	if !strings.Contains(list, "- **E001**") {
		t.Error("List should contain error codes in list format")
	}

	// Should contain categories
	if !strings.Contains(list, "[schema]") {
		t.Error("List should contain category tags")
	}
}

func TestErrorDocsGenerator_NoReservedCodes(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()
	generator := domain.NewErrorDocsGenerator(registry)

	markdown := generator.GenerateMarkdown()

	// Should not include reserved codes in output
	if strings.Contains(markdown, "Reserved for future use") {
		t.Error("Markdown should not include reserved error codes")
	}
}

func TestErrorDocsGenerator_CategoryOrder(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()
	generator := domain.NewErrorDocsGenerator(registry)

	markdown := generator.GenerateMarkdown()

	// Categories should appear in the output
	schemaPos := strings.Index(markdown, "Schema Errors")
	refsPos := strings.Index(markdown, "Refs Errors")

	if schemaPos == -1 || refsPos == -1 {
		t.Error("Expected both Schema and Refs categories in output")
	}
}

func TestErrorDocsGenerator_TableFormat(t *testing.T) {
	registry := domain.NewErrorCodeRegistry()
	generator := domain.NewErrorDocsGenerator(registry)

	markdown := generator.GenerateMarkdown()

	// Should have table headers
	if !strings.Contains(markdown, "|------") {
		t.Error("Markdown should contain table separators")
	}

	// Should have code cells with backticks
	if !strings.Contains(markdown, "`E001`") {
		t.Error("Markdown should contain code cells with backticks")
	}
}
