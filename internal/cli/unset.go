package cli

import (
	"fmt"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/patcher"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type unsetFlags struct {
	quiet     bool
	targetDir string
	nodeID    string
	path      string
}

// NewUnsetCommand creates the unset subcommand
func NewUnsetCommand() *cobra.Command {
	flags := &unsetFlags{}

	cmd := &cobra.Command{
		Use:   "unset <node-id> <path> [directory]",
		Short: "Remove a field value from a node",
		Long: `Remove a field value or array element from a node.

The path can be a simple field name or use bracket notation for array elements.
Required fields (id, kind, version, status, title) cannot be unset.

Examples:
  deco unset sword-001 summary          # Remove summary field
  deco unset sword-001 tags             # Remove all tags
  deco unset sword-001 tags[0]          # Remove first tag

The version number is automatically incremented after a successful unset.`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			flags.path = args[1]
			if len(args) > 2 {
				flags.targetDir = args[2]
			} else {
				flags.targetDir = "."
			}
			return runUnset(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")

	return cmd
}

func runUnset(flags *unsetFlags) error {
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

	// Capture old value for history
	oldValue := getFieldValueUnset(&n, flags.path)

	// Apply the unset
	p := patcher.New()
	err = p.Unset(&n, flags.path)
	if err != nil {
		return fmt.Errorf("failed to unset: %w", err)
	}

	// Increment version
	n.Version++

	// Save the node
	err = nodeRepo.Save(n)
	if err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log unset operation in history
	if err := logUnsetOperation(flags.targetDir, n.ID, flags.path, oldValue); err != nil {
		fmt.Printf("Warning: failed to log unset operation: %v\n", err)
	}

	if !flags.quiet {
		fmt.Printf("Unset %s.%s (version %d)\n", flags.nodeID, flags.path, n.Version)
	}

	return nil
}

// getFieldValueUnset extracts a field value from a node using reflection
func getFieldValueUnset(n *domain.Node, path string) interface{} {
	parts := strings.Split(path, ".")
	v := reflect.ValueOf(n).Elem()

	for _, part := range parts {
		// Handle array index notation
		if idx := strings.Index(part, "["); idx != -1 {
			fieldName := part[:idx]
			endIdx := strings.Index(part, "]")
			if endIdx == -1 {
				return nil
			}
			indexStr := part[idx+1 : endIdx]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil
			}

			field := v.FieldByName(capitalizeFirstUnset(fieldName))
			if !field.IsValid() || field.Kind() != reflect.Slice {
				return nil
			}
			if index < 0 || index >= field.Len() {
				return nil
			}
			v = field.Index(index)
			continue
		}

		field := v.FieldByName(capitalizeFirstUnset(part))
		if !field.IsValid() {
			return nil
		}
		v = field
	}

	if v.IsValid() && v.CanInterface() {
		return v.Interface()
	}
	return nil
}

func capitalizeFirstUnset(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// logUnsetOperation adds an unset entry to the history log
func logUnsetOperation(targetDir, nodeID, path string, oldValue interface{}) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Operation: "unset",
		User:      username,
		Before:    map[string]interface{}{path: oldValue},
	}

	return historyRepo.Append(entry)
}
