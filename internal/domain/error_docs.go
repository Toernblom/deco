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

package domain

import (
	"fmt"
	"sort"
	"strings"
)

// ErrorDocsGenerator generates markdown documentation for error codes
type ErrorDocsGenerator struct {
	registry *ErrorCodeRegistry
}

// NewErrorDocsGenerator creates a new documentation generator
func NewErrorDocsGenerator(registry *ErrorCodeRegistry) *ErrorDocsGenerator {
	return &ErrorDocsGenerator{
		registry: registry,
	}
}

// GenerateMarkdown generates markdown documentation for all error codes
func (g *ErrorDocsGenerator) GenerateMarkdown() string {
	var b strings.Builder

	b.WriteString("# Deco Error Codes\n\n")
	b.WriteString("This document lists all error codes used by Deco.\n\n")

	// Get all categories
	categories := g.registry.Categories()
	sort.Strings(categories)

	// Generate documentation for each category
	for _, category := range categories {
		b.WriteString(g.generateCategorySection(category))
	}

	// Add footer
	b.WriteString("\n---\n\n")
	b.WriteString("*This documentation is automatically generated from the error code registry.*\n")

	return b.String()
}

// generateCategorySection generates markdown for a single category
func (g *ErrorDocsGenerator) generateCategorySection(category string) string {
	var b strings.Builder

	// Category header
	b.WriteString(fmt.Sprintf("## %s Errors\n\n", strings.Title(category)))

	// Get codes for this category
	codes := g.registry.ByCategory(category)

	// Sort by code
	sort.Slice(codes, func(i, j int) bool {
		return codes[i].Code < codes[j].Code
	})

	// Generate table
	b.WriteString("| Code | Message |\n")
	b.WriteString("|------|----------|\n")

	for _, code := range codes {
		// Skip reserved codes
		if strings.Contains(code.Message, "Reserved for future use") {
			continue
		}
		b.WriteString(fmt.Sprintf("| `%s` | %s |\n", code.Code, code.Message))
	}

	b.WriteString("\n")

	return b.String()
}

// GenerateDetailedMarkdown generates detailed markdown with examples
func (g *ErrorDocsGenerator) GenerateDetailedMarkdown() string {
	var b strings.Builder

	b.WriteString("# Deco Error Codes Reference\n\n")
	b.WriteString("This document provides detailed information about all error codes used by Deco.\n\n")
	b.WriteString("## Table of Contents\n\n")

	// Get all categories
	categories := g.registry.Categories()
	sort.Strings(categories)

	// TOC
	for _, category := range categories {
		anchor := strings.ToLower(category)
		b.WriteString(fmt.Sprintf("- [%s Errors](#%s-errors)\n", strings.Title(category), anchor))
	}
	b.WriteString("\n")

	// Generate detailed sections
	for _, category := range categories {
		b.WriteString(g.generateDetailedCategorySection(category))
	}

	// Add footer
	b.WriteString("\n---\n\n")
	b.WriteString("*This documentation is automatically generated from the error code registry.*\n")

	return b.String()
}

// generateDetailedCategorySection generates detailed markdown for a category
func (g *ErrorDocsGenerator) generateDetailedCategorySection(category string) string {
	var b strings.Builder

	// Category header
	b.WriteString(fmt.Sprintf("## %s Errors\n\n", strings.Title(category)))

	// Get codes for this category
	codes := g.registry.ByCategory(category)

	// Sort by code
	sort.Slice(codes, func(i, j int) bool {
		return codes[i].Code < codes[j].Code
	})

	// Generate detailed entries
	for _, code := range codes {
		// Skip reserved codes
		if strings.Contains(code.Message, "Reserved for future use") {
			continue
		}

		b.WriteString(fmt.Sprintf("### %s: %s\n\n", code.Code, code.Message))
		b.WriteString(fmt.Sprintf("**Category:** %s\n\n", category))

		// Add description based on the message
		b.WriteString(g.generateDescription(code))
		b.WriteString("\n")
	}

	return b.String()
}

// generateDescription generates a description for an error code
func (g *ErrorDocsGenerator) generateDescription(code ErrorCode) string {
	var b strings.Builder

	b.WriteString("**Description:**\n\n")

	// Generate contextual description based on category and message
	switch code.Category {
	case "schema":
		b.WriteString("This error occurs when there is an issue with the node schema or structure.\n\n")
	case "refs":
		b.WriteString("This error occurs when there is an issue with node references.\n\n")
	case "validation":
		b.WriteString("This error occurs during validation of node data.\n\n")
	case "io":
		b.WriteString("This error occurs during file I/O operations.\n\n")
	case "graph":
		b.WriteString("This error occurs when there is an issue with the node graph structure.\n\n")
	}

	return b.String()
}

// GenerateCodeList generates a simple list of all codes
func (g *ErrorDocsGenerator) GenerateCodeList() string {
	var b strings.Builder

	b.WriteString("# Error Code List\n\n")

	codes := g.registry.AllCodes()

	// Sort by code
	sort.Slice(codes, func(i, j int) bool {
		return codes[i].Code < codes[j].Code
	})

	for _, code := range codes {
		if !strings.Contains(code.Message, "Reserved for future use") {
			b.WriteString(fmt.Sprintf("- **%s** [%s]: %s\n", code.Code, code.Category, code.Message))
		}
	}

	return b.String()
}
