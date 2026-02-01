package cli

import (
	"fmt"

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
  1: Validation errors found`,
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
	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load all nodes
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	nodes, err := nodeRepo.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	// Run validation (including unknown field detection)
	orchestrator := validator.NewOrchestrator()
	collector := orchestrator.ValidateAllWithDir(nodes, flags.targetDir)

	// Check if there are errors
	if !collector.HasErrors() {
		if !flags.quiet {
			fmt.Println("✓ All nodes are valid")
		}
		return nil
	}

	// Print errors unless quiet
	if !flags.quiet {
		errors := collector.Errors()
		fmt.Printf("✗ Found %d validation error(s):\n\n", collector.Count())
		for _, err := range errors {
			fmt.Println(err.Error())
		}
	}

	// Return error to trigger exit code 1
	return fmt.Errorf("validation failed with %d error(s)", collector.Count())
}
