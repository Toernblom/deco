package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type initFlags struct {
	force     bool
	targetDir string
}

// NewInitCommand creates the init subcommand
func NewInitCommand() *cobra.Command {
	flags := &initFlags{}

	cmd := &cobra.Command{
		Use:   "init [directory]",
		Short: "Initialize a new deco project",
		Long: `Initialize a new deco project in the current or specified directory.

Creates the .deco directory structure with:
  - config.yaml: Project configuration
  - nodes/: Directory for design document YAML files`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				flags.targetDir = args[0]
			} else {
				flags.targetDir = "."
			}
			return runInit(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, "Reinitialize existing project")

	return cmd
}

func runInit(flags *initFlags) error {
	decoDir := filepath.Join(flags.targetDir, ".deco")

	// Check if project already exists
	if _, err := os.Stat(decoDir); err == nil {
		if !flags.force {
			return fmt.Errorf("project already initialized in .deco/ (use --force to reinitialize)")
		}
		// Remove existing directory if force flag is set
		if err := os.RemoveAll(decoDir); err != nil {
			return fmt.Errorf("failed to remove existing .deco directory: %w", err)
		}
	}

	// Create .deco directory
	if err := os.MkdirAll(decoDir, 0755); err != nil {
		return fmt.Errorf("failed to create .deco directory: %w", err)
	}

	// Create nodes directory
	nodesDir := filepath.Join(decoDir, "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		return fmt.Errorf("failed to create nodes directory: %w", err)
	}

	// Create default config.yaml
	if err := createDefaultConfig(decoDir, flags.targetDir); err != nil {
		return fmt.Errorf("failed to create config.yaml: %w", err)
	}

	fmt.Printf("Initialized deco project in %s\n", decoDir)
	return nil
}

func createDefaultConfig(decoDir, targetDir string) error {
	configPath := filepath.Join(decoDir, "config.yaml")

	// Get project name from target directory
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		absTarget = targetDir
	}

	// Create default configuration
	defaultConfig := config.Config{
		Version:     1,
		ProjectName: filepath.Base(absTarget),
		NodesPath:   filepath.Join(decoDir, "nodes"),
		HistoryPath: filepath.Join(decoDir, "history.jsonl"),
		Custom:      make(map[string]interface{}),
	}

	// Marshal to YAML
	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
