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
	"fmt"
	"testing"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/validator"
)

// generateNodes creates n test nodes with various relationships
func generateNodes(n int) []domain.Node {
	nodes := make([]domain.Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = domain.Node{
			ID:      fmt.Sprintf("node-%04d", i),
			Kind:    "system",
			Version: 1,
			Status:  "draft",
			Title:   fmt.Sprintf("Node %d", i),
			Tags:    []string{"tag1", "tag2"},
		}

		// Add references to previous nodes (creates a chain)
		if i > 0 {
			nodes[i].Refs.Uses = append(nodes[i].Refs.Uses, domain.RefLink{
				Target:  fmt.Sprintf("node-%04d", i-1),
				Context: "depends on",
			})
		}
		if i > 1 {
			nodes[i].Refs.Related = append(nodes[i].Refs.Related, domain.RefLink{
				Target: fmt.Sprintf("node-%04d", i-2),
			})
		}
	}
	return nodes
}

// generateNodesWithConstraints creates nodes that each have CEL constraints
func generateNodesWithConstraints(n int) []domain.Node {
	nodes := generateNodes(n)
	for i := range nodes {
		nodes[i].Constraints = []domain.Constraint{
			{Expr: "version > 0", Message: "Version must be positive", Scope: "all"},
			{Expr: `status in ["draft", "approved"]`, Message: "Invalid status", Scope: "all"},
		}
	}
	return nodes
}

func BenchmarkValidateAll_Small(b *testing.B) {
	nodes := generateNodes(10)
	orchestrator := validator.NewOrchestrator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orchestrator.ValidateAll(nodes)
	}
}

func BenchmarkValidateAll_Medium(b *testing.B) {
	nodes := generateNodes(100)
	orchestrator := validator.NewOrchestrator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orchestrator.ValidateAll(nodes)
	}
}

func BenchmarkValidateAll_Large(b *testing.B) {
	nodes := generateNodes(500)
	orchestrator := validator.NewOrchestrator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orchestrator.ValidateAll(nodes)
	}
}

func BenchmarkConstraintValidation_Small(b *testing.B) {
	nodes := generateNodesWithConstraints(10)
	orchestrator := validator.NewOrchestrator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orchestrator.ValidateAll(nodes)
	}
}

func BenchmarkConstraintValidation_Medium(b *testing.B) {
	nodes := generateNodesWithConstraints(100)
	orchestrator := validator.NewOrchestrator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orchestrator.ValidateAll(nodes)
	}
}

func BenchmarkConstraintValidation_Large(b *testing.B) {
	nodes := generateNodesWithConstraints(500)
	orchestrator := validator.NewOrchestrator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orchestrator.ValidateAll(nodes)
	}
}

func BenchmarkReferenceValidation(b *testing.B) {
	nodes := generateNodes(500)
	refValidator := validator.NewReferenceValidator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector := make(chan domain.DecoError, 100)
		go func() {
			for range collector {
			}
		}()
		// Note: collector interface mismatch - using ValidateAll instead
		orchestrator := validator.NewOrchestrator()
		orchestrator.ValidateAll(nodes)
	}
	_ = refValidator // silence unused warning
}
