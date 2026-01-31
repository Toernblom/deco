package cli

import (
	"fmt"

	"github.com/Toernblom/deco/internal/services/patcher"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type appendFlags struct {
	quiet     bool
	targetDir string
	nodeID    string
	path      string
	value     string
}

// NewAppendCommand creates the append subcommand
func NewAppendCommand() *cobra.Command {
	flags := &appendFlags{}

	cmd := &cobra.Command{
		Use:   "append <node-id> <path> <value> [directory]",
		Short: "Append a value to an array field on a node",
		Long: `Append a value to an array field on a node.

The path should point to an array field. The value will be appended to the end
of the array. This command errors if the path does not point to an array field.

Examples:
  deco append sword-001 tags legendary
  deco append hero-001 tags combat
  deco append quest-001 tags story

The version number is automatically incremented after a successful append.`,
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			flags.path = args[1]
			flags.value = args[2]
			if len(args) > 3 {
				flags.targetDir = args[3]
			} else {
				flags.targetDir = "."
			}
			return runAppend(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runAppend(flags *appendFlags) error {
	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Apply the append
	p := patcher.New()
	err = p.Append(&n, flags.path, flags.value)
	if err != nil {
		return fmt.Errorf("failed to append: %w", err)
	}

	// Increment version
	n.Version++

	// Save the node
	err = nodeRepo.Save(n)
	if err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	if !flags.quiet {
		fmt.Printf("Appended %q to %s.%s (version %d)\n", flags.value, flags.nodeID, flags.path, n.Version)
	}

	return nil
}
