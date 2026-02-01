package cli

import (
	"encoding/json"
	"fmt"
	"os"
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

type applyFlags struct {
	quiet     bool
	dryRun    bool
	targetDir string
	nodeID    string
	patchFile string
}

// NewApplyCommand creates the apply subcommand
func NewApplyCommand() *cobra.Command {
	flags := &applyFlags{}

	cmd := &cobra.Command{
		Use:   "apply <node-id> <patch-file> [directory]",
		Short: "Apply a batch of patch operations to a node",
		Long: `Apply a batch of patch operations from a JSON file to a node.

The patch file should contain a JSON array of operations:
  [
    {"op": "set", "path": "title", "value": "New Title"},
    {"op": "append", "path": "tags", "value": "new-tag"},
    {"op": "unset", "path": "summary"}
  ]

Supported operations:
  - set:    Set a field value
  - append: Append to an array field
  - unset:  Remove a field or array element

If any operation fails, all changes are rolled back (transactional).

Examples:
  deco apply sword-001 patch.json
  deco apply sword-001 changes.json --dry-run
  deco apply hero-001 updates.json /path/to/project

The version number is automatically incremented after a successful apply.`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			flags.patchFile = args[1]
			if len(args) > 2 {
				flags.targetDir = args[2]
			} else {
				flags.targetDir = "."
			}
			return runApply(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")
	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Validate patch without applying")

	return cmd
}

func runApply(flags *applyFlags) error {
	// Load config to verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Read and parse patch file
	data, err := os.ReadFile(flags.patchFile)
	if err != nil {
		return fmt.Errorf("failed to read patch file: %w", err)
	}

	var operations []patcher.PatchOperation
	if err := json.Unmarshal(data, &operations); err != nil {
		return fmt.Errorf("failed to parse patch file: %w", err)
	}

	// Load the node
	nodeRepo := node.NewYAMLRepository(flags.targetDir)
	n, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Capture before values for history
	beforeValues := make(map[string]interface{})
	for _, op := range operations {
		beforeValues[op.Path] = getFieldValueApply(&n, op.Path)
	}

	// Apply the patch operations
	p := patcher.New()
	err = p.Apply(&n, operations)
	if err != nil {
		return fmt.Errorf("failed to apply patch: %w", err)
	}

	if flags.dryRun {
		if !flags.quiet {
			fmt.Printf("Dry run: %d operation(s) would be applied to %s\n", len(operations), flags.nodeID)
		}
		return nil
	}

	// Capture after values for history
	afterValues := make(map[string]interface{})
	for _, op := range operations {
		if op.Op != "unset" {
			afterValues[op.Path] = getFieldValueApply(&n, op.Path)
		}
	}

	// Increment version
	n.Version++

	// Save the node
	err = nodeRepo.Save(n)
	if err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log apply operation in history
	if err := logApplyOperation(flags.targetDir, n.ID, beforeValues, afterValues); err != nil {
		fmt.Printf("Warning: failed to log apply operation: %v\n", err)
	}

	if !flags.quiet {
		fmt.Printf("Applied %d operation(s) to %s (version %d)\n", len(operations), flags.nodeID, n.Version)
	}

	return nil
}

// getFieldValueApply extracts a field value from a node using reflection
func getFieldValueApply(n *domain.Node, path string) interface{} {
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

			field := v.FieldByName(capitalizeFirstApply(fieldName))
			if !field.IsValid() || field.Kind() != reflect.Slice {
				return nil
			}
			if index < 0 || index >= field.Len() {
				return nil
			}
			v = field.Index(index)
			continue
		}

		field := v.FieldByName(capitalizeFirstApply(part))
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

func capitalizeFirstApply(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// logApplyOperation adds an apply entry to the history log
func logApplyOperation(targetDir, nodeID string, beforeValues, afterValues map[string]interface{}) error {
	historyRepo := history.NewYAMLRepository(targetDir)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	entry := domain.AuditEntry{
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Operation: "update",
		User:      username,
		Before:    beforeValues,
		After:     afterValues,
	}

	return historyRepo.Append(entry)
}
