package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Toernblom/deco/internal/services/patcher"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type setFlags struct {
	quiet     bool
	targetDir string
	nodeID    string
	path      string
	value     string
}

// NewSetCommand creates the set subcommand
func NewSetCommand() *cobra.Command {
	flags := &setFlags{}

	cmd := &cobra.Command{
		Use:   "set <node-id> <path> <value> [directory]",
		Short: "Set a field value on a node",
		Long: `Set a field value on a node.

The path can be a simple field name or use dot notation for nested fields.
Array indices are supported using bracket notation.

Examples:
  deco set sword-001 title "Golden Sword"
  deco set sword-001 status published
  deco set sword-001 summary "A legendary golden sword"
  deco set sword-001 tags[0] legendary

The version number is automatically incremented after a successful set.`,
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
			return runSet(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runSet(flags *setFlags) error {
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

	// Parse the value to appropriate type
	value := parseValue(flags.value)

	// Apply the patch
	p := patcher.New()
	err = p.Set(&n, flags.path, value)
	if err != nil {
		return fmt.Errorf("failed to set field: %w", err)
	}

	// Increment version
	n.Version++

	// Save the node
	err = nodeRepo.Save(n)
	if err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	if !flags.quiet {
		fmt.Printf("Updated %s.%s = %v (version %d)\n", flags.nodeID, flags.path, value, n.Version)
	}

	return nil
}

// parseValue attempts to parse a string value into the appropriate Go type.
// Order of attempts: int, float, bool, JSON array/object, string.
func parseValue(s string) interface{} {
	// Try int
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}

	// Try float (only if contains decimal point to avoid converting "123" to 123.0)
	if strings.Contains(s, ".") {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}

	// Try bool
	if strings.EqualFold(s, "true") {
		return true
	}
	if strings.EqualFold(s, "false") {
		return false
	}

	// Try JSON array or object
	if (strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) ||
		(strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")) {
		var jsonVal interface{}
		if err := json.Unmarshal([]byte(s), &jsonVal); err == nil {
			return jsonVal
		}
	}

	// Default to string
	return s
}
