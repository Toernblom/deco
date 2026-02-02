package cli

import (
	"fmt"
	"os"
	"os/user"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/services/validator"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/history"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type rewriteFlags struct {
	quiet     bool
	dryRun    bool
	targetDir string
	nodeID    string
	inputFile string
}

// NewRewriteCommand creates the rewrite subcommand
func NewRewriteCommand() *cobra.Command {
	flags := &rewriteFlags{}

	cmd := &cobra.Command{
		Use:   "rewrite <node-id> <file.yaml> [directory]",
		Short: "Replace entire node content from a YAML file",
		Long: `Replace a node's content entirely with content from a YAML file.

Unlike 'apply' which performs surgical patches, 'rewrite' replaces the
entire node content. This is useful when an AI rewrites a complete node.

The input file must:
  - Be valid YAML
  - Include all required node fields (id, title, status, kind, version)
  - The ID in the file must match the target node ID

The node is validated before writing. If validation fails, no changes are made.

Examples:
  deco rewrite sword-001 new-sword.yaml
  deco rewrite sword-001 new-sword.yaml --dry-run  # Show diff without applying
  deco rewrite hero-001 hero-v2.yaml /path/to/project`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.nodeID = args[0]
			flags.inputFile = args[1]
			if len(args) > 2 {
				flags.targetDir = args[2]
			} else {
				flags.targetDir = "."
			}
			return runRewrite(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output")
	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Show diff without applying changes")

	return cmd
}

func runRewrite(flags *rewriteFlags) error {
	// Load config to get validation settings
	configRepo := config.NewYAMLRepository(flags.targetDir)
	cfg, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	// Read and parse input file
	data, err := os.ReadFile(flags.inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	var newNode domain.Node
	if err := yaml.Unmarshal(data, &newNode); err != nil {
		return fmt.Errorf("failed to parse input YAML: %w", err)
	}

	// Validate that the new node ID matches the target
	if newNode.ID != flags.nodeID {
		return fmt.Errorf("node ID in file (%q) does not match target (%q)", newNode.ID, flags.nodeID)
	}

	// Load the existing node
	nodeRepo := node.NewYAMLRepository(config.ResolveNodesPath(cfg, flags.targetDir))
	oldNode, err := nodeRepo.Load(flags.nodeID)
	if err != nil {
		return fmt.Errorf("node %q not found: %w", flags.nodeID, err)
	}

	// Preserve source file from original node
	newNode.SourceFile = oldNode.SourceFile

	// Validate the new node
	orchestrator := validator.NewOrchestratorWithFullConfig(cfg.RequiredApprovals, cfg.CustomBlockTypes, cfg.SchemaRules)
	collector := orchestrator.ValidateNode(&newNode)
	validationPassed := !collector.HasErrors()

	if flags.dryRun {
		if !flags.quiet {
			printRewriteDiff(flags.nodeID, &oldNode, &newNode)
			printValidationResult(validationPassed, collector)
			if validationPassed {
				fmt.Println("\nRun without --dry-run to apply.")
			}
		}
		return nil
	}

	// Abort if validation fails
	if !validationPassed {
		fmt.Printf("Rewrite would create invalid node %s:\n\n", flags.nodeID)
		for _, err := range collector.Errors() {
			fmt.Println("  " + err.Error())
		}
		return fmt.Errorf("rewrite rejected: validation failed with %d error(s)", collector.Count())
	}

	// Save the new node
	err = nodeRepo.Save(newNode)
	if err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	// Log rewrite operation in history
	historyPath := config.ResolveHistoryPath(cfg, flags.targetDir)
	if err := logRewriteOperation(historyPath, &oldNode, &newNode); err != nil {
		fmt.Printf("Warning: failed to log rewrite operation: %v\n", err)
	}

	if !flags.quiet {
		fmt.Printf("Rewrote %s (version %d → %d)\n", flags.nodeID, oldNode.Version, newNode.Version)
	}

	return nil
}

// printRewriteDiff shows a comprehensive diff between old and new nodes
func printRewriteDiff(nodeID string, oldNode, newNode *domain.Node) {
	fmt.Printf("Proposed rewrite of %s:\n", nodeID)
	fmt.Println(strings.Repeat("-", 50))

	// Compare core fields
	diffField("id", oldNode.ID, newNode.ID)
	diffField("kind", oldNode.Kind, newNode.Kind)
	diffField("version", oldNode.Version, newNode.Version)
	diffField("status", oldNode.Status, newNode.Status)
	diffField("title", oldNode.Title, newNode.Title)
	diffField("summary", oldNode.Summary, newNode.Summary)

	// Compare tags
	diffSlice("tags", oldNode.Tags, newNode.Tags)

	// Compare LLM context
	diffField("llm_context", oldNode.LLMContext, newNode.LLMContext)

	// Note about complex fields
	if !reflect.DeepEqual(oldNode.Content, newNode.Content) {
		fmt.Println("  content: (changed)")
	}
	if !reflect.DeepEqual(oldNode.Refs, newNode.Refs) {
		fmt.Println("  refs: (changed)")
	}
	if !reflect.DeepEqual(oldNode.Issues, newNode.Issues) {
		fmt.Println("  issues: (changed)")
	}
	if !reflect.DeepEqual(oldNode.Contracts, newNode.Contracts) {
		fmt.Println("  contracts: (changed)")
	}
	if !reflect.DeepEqual(oldNode.Constraints, newNode.Constraints) {
		fmt.Println("  constraints: (changed)")
	}
	if !reflect.DeepEqual(oldNode.Glossary, newNode.Glossary) {
		fmt.Println("  glossary: (changed)")
	}
	if !reflect.DeepEqual(oldNode.Reviewers, newNode.Reviewers) {
		fmt.Println("  reviewers: (changed)")
	}
	if !reflect.DeepEqual(oldNode.Custom, newNode.Custom) {
		fmt.Println("  custom: (changed)")
	}
}

// diffField prints a diff for a single field if changed
func diffField(name string, oldVal, newVal interface{}) {
	if reflect.DeepEqual(oldVal, newVal) {
		return
	}

	oldStr := formatRewriteValue(oldVal)
	newStr := formatRewriteValue(newVal)

	if isEmptyValue(oldVal) {
		fmt.Printf("  + %s: %s\n", name, newStr)
	} else if isEmptyValue(newVal) {
		fmt.Printf("  - %s: %s\n", name, oldStr)
	} else {
		fmt.Printf("  %s: %s → %s\n", name, oldStr, newStr)
	}
}

// diffSlice prints a diff for slice fields
func diffSlice(name string, oldSlice, newSlice []string) {
	if reflect.DeepEqual(oldSlice, newSlice) {
		return
	}

	oldStr := formatStringSlice(oldSlice)
	newStr := formatStringSlice(newSlice)

	if len(oldSlice) == 0 {
		fmt.Printf("  + %s: %s\n", name, newStr)
	} else if len(newSlice) == 0 {
		fmt.Printf("  - %s: %s\n", name, oldStr)
	} else {
		fmt.Printf("  %s: %s → %s\n", name, oldStr, newStr)
	}
}

// formatRewriteValue formats a value for display
func formatRewriteValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		if val == "" {
			return "(empty)"
		}
		if len(val) > 50 {
			return fmt.Sprintf("%q...", val[:47])
		}
		return fmt.Sprintf("%q", val)
	case int:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatStringSlice formats a string slice for display
func formatStringSlice(s []string) string {
	if len(s) == 0 {
		return "[]"
	}
	return "[" + strings.Join(s, ", ") + "]"
}

// isEmptyValue checks if a value is empty/zero
func isEmptyValue(v interface{}) bool {
	switch val := v.(type) {
	case string:
		return val == ""
	case int:
		return val == 0
	case nil:
		return true
	default:
		return false
	}
}

// logRewriteOperation adds a rewrite entry to the history log
func logRewriteOperation(historyPath string, oldNode, newNode *domain.Node) error {
	historyRepo := history.NewYAMLRepository(historyPath)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	// Capture key field changes
	before := make(map[string]interface{})
	after := make(map[string]interface{})

	// Track changed fields
	fields := []string{"version", "status", "title", "summary", "tags"}
	for _, field := range fields {
		oldVal := getNodeFieldValue(oldNode, field)
		newVal := getNodeFieldValue(newNode, field)
		if !reflect.DeepEqual(oldVal, newVal) {
			before[field] = oldVal
			after[field] = newVal
		}
	}

	entry := domain.AuditEntry{
		Timestamp:   time.Now(),
		NodeID:      newNode.ID,
		Operation:   "rewrite",
		User:        username,
		ContentHash: ComputeContentHash(*newNode),
		Before:      before,
		After:       after,
	}

	return historyRepo.Append(entry)
}

// getNodeFieldValue extracts a field value from a node
func getNodeFieldValue(n *domain.Node, field string) interface{} {
	switch field {
	case "id":
		return n.ID
	case "kind":
		return n.Kind
	case "version":
		return n.Version
	case "status":
		return n.Status
	case "title":
		return n.Title
	case "summary":
		return n.Summary
	case "tags":
		// Return a copy to avoid slice aliasing
		tags := make([]string, len(n.Tags))
		copy(tags, n.Tags)
		return tags
	case "llm_context":
		return n.LLMContext
	default:
		return nil
	}
}

// nodeFieldNames returns field names for diff output
func nodeFieldNames() []string {
	return []string{
		"id", "kind", "version", "status", "title", "summary",
		"tags", "llm_context", "content", "refs", "issues",
		"contracts", "constraints", "glossary", "reviewers", "custom",
	}
}

// sortedKeys returns sorted keys from a map
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
