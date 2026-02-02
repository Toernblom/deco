package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/patcher"
	"github.com/Toernblom/deco/internal/services/validator"
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
	// Load config to get validation settings
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
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

	// Capture before values for diff and history
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

	// Capture after values for diff and history
	afterValues := make(map[string]interface{})
	for _, op := range operations {
		afterValues[op.Path] = getFieldValueApply(&n, op.Path)
	}

	// Validate the resulting node before save
	orchestrator := validator.NewOrchestratorWithFullConfig(cfg.RequiredApprovals, cfg.CustomBlockTypes, cfg.SchemaRules)
	collector := orchestrator.ValidateNode(&n)
	validationPassed := !collector.HasErrors()

	if flags.dryRun {
		if !flags.quiet {
			printApplyDiff(flags.nodeID, beforeValues, afterValues, operations)
			printValidationResult(validationPassed, collector)
			if validationPassed {
				fmt.Println("\nRun without --dry-run to apply.")
			}
		}
		return nil
	}

	// Abort if validation fails
	if !validationPassed {
		fmt.Printf("Patch would create invalid node %s:\n\n", flags.nodeID)
		for _, err := range collector.Errors() {
			fmt.Println("  " + err.Error())
		}
		return fmt.Errorf("patch rejected: validation failed with %d error(s)", collector.Count())
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

// printApplyDiff shows field-level changes for dry-run
func printApplyDiff(nodeID string, before, after map[string]interface{}, operations []patcher.PatchOperation) {
	fmt.Printf("Proposed changes to %s:\n", nodeID)

	// Collect all paths
	paths := make(map[string]bool)
	for _, op := range operations {
		paths[op.Path] = true
	}

	// Sort paths for consistent output
	var sortedPaths []string
	for path := range paths {
		sortedPaths = append(sortedPaths, path)
	}
	sort.Strings(sortedPaths)

	for _, path := range sortedPaths {
		beforeVal, hasBefore := before[path]
		afterVal, hasAfter := after[path]

		// Determine the operation type
		opType := ""
		for _, op := range operations {
			if op.Path == path {
				opType = op.Op
				break
			}
		}

		switch opType {
		case "unset":
			fmt.Printf("  - %s: %s\n", path, formatApplyValue(beforeVal))
		case "set":
			if hasBefore && beforeVal != nil && !isZeroValue(beforeVal) {
				fmt.Printf("  %s: %s → %s\n", path, formatApplyValue(beforeVal), formatApplyValue(afterVal))
			} else {
				fmt.Printf("  + %s: %s\n", path, formatApplyValue(afterVal))
			}
		case "append":
			if hasBefore && hasAfter {
				fmt.Printf("  %s: %s → %s\n", path, formatApplyValue(beforeVal), formatApplyValue(afterVal))
			} else {
				fmt.Printf("  + %s: %s\n", path, formatApplyValue(afterVal))
			}
		default:
			if hasBefore && hasAfter {
				fmt.Printf("  %s: %s → %s\n", path, formatApplyValue(beforeVal), formatApplyValue(afterVal))
			}
		}
	}
}

// printValidationResult shows validation status
func printValidationResult(passed bool, collector interface{ Count() int; Errors() []domain.DecoError }) {
	fmt.Println()
	if passed {
		fmt.Println("Validation: ✓ Valid")
	} else {
		fmt.Printf("Validation: ✗ Invalid (%d error(s))\n", collector.Count())
		for _, err := range collector.Errors() {
			fmt.Println("  " + err.Summary)
		}
	}
}

// formatApplyValue formats a value for diff display
func formatApplyValue(v interface{}) string {
	if v == nil {
		return "(nil)"
	}
	switch val := v.(type) {
	case []interface{}:
		if len(val) == 0 {
			return "[]"
		}
		var items []string
		for _, item := range val {
			items = append(items, fmt.Sprintf("%v", item))
		}
		return "[" + strings.Join(items, ", ") + "]"
	case []string:
		if len(val) == 0 {
			return "[]"
		}
		return "[" + strings.Join(val, ", ") + "]"
	case string:
		return fmt.Sprintf("%q", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// isZeroValue checks if a value is considered "empty"
func isZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case []interface{}:
		return len(val) == 0
	case []string:
		return len(val) == 0
	case int:
		return val == 0
	case float64:
		return val == 0
	default:
		return false
	}
}
