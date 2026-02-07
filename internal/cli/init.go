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
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Toernblom/deco/internal/migrations"
	"github.com/Toernblom/deco/internal/storage/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type initFlags struct {
	force     bool
	template  string
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
  - nodes/: Directory for design document YAML files

Use --template to start with pre-built content:
  - game-design: Game design document with combat mechanics and controls
  - api-spec: API specification with endpoints and auth schemas`,
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
	cmd.Flags().StringVar(&flags.template, "template", "", "Project template (game-design, api-spec)")

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

	if flags.template != "" {
		// Validate template name
		valid := false
		for _, t := range availableTemplates {
			if t == flags.template {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unknown template %q (available: %s)", flags.template, strings.Join(availableTemplates, ", "))
		}

		if err := applyTemplate(flags.template, decoDir, flags.targetDir); err != nil {
			return fmt.Errorf("failed to apply template: %w", err)
		}

		// Compute and set schema version so validate passes out of the box
		if err := setSchemaVersion(decoDir); err != nil {
			// Non-fatal: template still works, user can run deco migrate
			fmt.Printf("Warning: could not set schema version: %v\n", err)
		}
	} else {
		// Create default config.yaml
		if err := createDefaultConfig(decoDir, flags.targetDir); err != nil {
			return fmt.Errorf("failed to create config.yaml: %w", err)
		}
	}

	fmt.Printf("Initialized deco project in %s\n", decoDir)
	if flags.template != "" {
		fmt.Printf("Template: %s\n", flags.template)
	}
	return nil
}

func applyTemplate(templateName, decoDir, targetDir string) error {
	templateRoot := "templates/" + templateName

	// Get project name
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		absTarget = targetDir
	}
	projectName := filepath.Base(absTarget)

	return fs.WalkDir(embeddedTemplates, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path from template root
		relPath, err := filepath.Rel(templateRoot, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(decoDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := embeddedTemplates.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		// Replace template variables
		content := strings.ReplaceAll(string(data), "{{PROJECT_NAME}}", projectName)

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(destPath, []byte(content), 0644)
	})
}

func setSchemaVersion(decoDir string) error {
	// Load the config that the template just wrote
	// decoDir is the .deco directory, parent is the project root
	projectRoot := filepath.Dir(decoDir)
	configRepo := config.NewYAMLRepository(projectRoot)
	cfg, err := configRepo.Load()
	if err != nil {
		return err
	}

	hash := migrations.ComputeSchemaHash(cfg)
	if hash == "" {
		return nil // No schema constraints, no version needed
	}

	cfg.SchemaVersion = hash
	return configRepo.Save(cfg)
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
