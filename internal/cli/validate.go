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
	"fmt"

	"github.com/Toernblom/deco/internal/cli/style"
	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/migrations"
	"github.com/Toernblom/deco/internal/services/validator"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type validateFlags struct {
	quiet     bool
	targetDir string
}

// NewValidateCommand creates the validate subcommand
func NewValidateCommand() *cobra.Command {
	flags := &validateFlags{}

	cmd := &cobra.Command{
		Use:   "validate [directory]",
		Short: "Validate all nodes in the project",
		Long: `Validate all nodes in the project against schema, references, constraints, and contracts.

Checks:
  - Schema: All required fields are present
  - References: All referenced nodes exist
  - Constraints: All CEL expressions evaluate to true
  - Contracts: Valid given/when/then structure, unique names, valid @node refs

Exit codes:
  0: All nodes are valid
  1: Validation errors found
  2: Schema version mismatch (run 'deco migrate')`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				flags.targetDir = args[0]
			} else {
				flags.targetDir = "."
			}
			return runValidate(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output (exit code only)")

	return cmd
}

func runValidate(flags *validateFlags) error {
	// Load config
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Check schema version before validation
	needsMigration, currentHash, expectedHash, err := migrations.NeedsMigration(flags.targetDir)
	if err != nil {
		return fmt.Errorf("failed to check schema version: %w", err)
	}
	if needsMigration {
		if !flags.quiet {
			fmt.Printf("%s %s\n", style.ErrorIcon(), style.Error.Sprint("Schema version mismatch"))
			fmt.Printf("  %s  %s\n", style.Muted.Sprint("Current:"), formatSchemaHash(currentHash))
			fmt.Printf("  %s %s\n", style.Muted.Sprint("Expected:"), formatSchemaHash(expectedHash))
			fmt.Printf("\n%s\n", style.Info.Sprint("Run 'deco migrate' to update nodes to the current schema."))
		}
		return NewExitError(ExitCodeSchemaMismatch, "schema version mismatch")
	}

	// Load all nodes
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	// Run validation with full config support (custom block types, schema rules, unknown field detection)
	orchestrator := validator.NewOrchestratorWithFullConfig(cfg.RequiredApprovals, cfg.CustomBlockTypes, cfg.SchemaRules)
	collector := orchestrator.ValidateAllWithDir(nodes, flags.targetDir)

	// Check if there are errors
	if !collector.HasErrors() {
		if !flags.quiet {
			fmt.Printf("%s All nodes are valid\n", style.SuccessIcon())
		}
		return nil
	}

	// Print errors unless quiet
	if !flags.quiet {
		errors := collector.Errors()
		fmt.Printf("%s Found %s validation error(s):\n\n", style.ErrorIcon(), style.Error.Sprint(collector.Count()))

		formatter := domain.NewErrorFormatter()
		formatter.SetColor(style.IsEnabled())

		for _, err := range errors {
			fmt.Println(formatter.Format(err))
		}
	}

	// Return exit error (message is for programmatic use, not printed again)
	return NewExitError(ExitCodeError, fmt.Sprintf("validation failed with %d error(s)", collector.Count()))
}

// formatSchemaHash formats a schema hash for display.
func formatSchemaHash(hash string) string {
	if hash == "" {
		return "(none)"
	}
	return hash
}
