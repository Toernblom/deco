package cli

import (
	"fmt"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/Toernblom/deco/internal/storage/node"
	"github.com/spf13/cobra"
)

type createFlags struct {
	kind      string
	title     string
	status    string
	force     bool
	targetDir string
}

// NewCreateCommand creates the create subcommand
func NewCreateCommand() *cobra.Command {
	flags := &createFlags{}

	cmd := &cobra.Command{
		Use:   "create <id>",
		Short: "Create a new node with required fields",
		Long: `Create a new node with the specified ID and scaffold required fields.

The ID can include path segments for organization (e.g., systems/combat).
Node files are stored in .deco/nodes/ directory.

Examples:
  deco create sword-001
  deco create systems/combat --kind mechanic --title "Combat System"
  deco create items/weapons/axe -k item -t "Battle Axe"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(args[0], flags)
		},
	}

	cmd.Flags().StringVarP(&flags.kind, "kind", "k", "system", "Node kind (system, mechanic, feature, etc.)")
	cmd.Flags().StringVarP(&flags.title, "title", "t", "", "Node title (defaults to ID)")
	cmd.Flags().StringVarP(&flags.status, "status", "s", "draft", "Node status (draft, approved, deprecated)")
	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, "Overwrite existing node")
	cmd.Flags().StringVarP(&flags.targetDir, "dir", "d", ".", "Project directory")

	return cmd
}

func runCreate(id string, flags *createFlags) error {
	// Verify project exists
	configRepo := config.NewYAMLRepository(flags.targetDir)
	_, err := configRepo.Load()
	if err != nil {
		return fmt.Errorf(".deco directory not found or invalid: %w", err)
	}

	nodeRepo := node.NewYAMLRepository(flags.targetDir)

	// Check if node already exists
	exists, err := nodeRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check node existence: %w", err)
	}
	if exists && !flags.force {
		return fmt.Errorf("node %q already exists (use --force to overwrite)", id)
	}

	// Derive title from ID if not provided
	title := flags.title
	if title == "" {
		title = deriveTitle(id)
	}

	// Create the node
	newNode := domain.Node{
		ID:      id,
		Kind:    flags.kind,
		Version: 1,
		Status:  flags.status,
		Title:   title,
	}

	// Save the node
	if err := nodeRepo.Save(newNode); err != nil {
		return fmt.Errorf("failed to save node: %w", err)
	}

	fmt.Printf("Created node: %s\n", id)
	return nil
}

// deriveTitle creates a human-readable title from the node ID
func deriveTitle(id string) string {
	// Get the last segment of the path
	parts := strings.Split(id, "/")
	name := parts[len(parts)-1]

	// Replace hyphens and underscores with spaces
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Capitalize first letter of each word
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return strings.Join(words, " ")
}
