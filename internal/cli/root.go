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

	"github.com/Toernblom/deco/internal/cli/style"
	"github.com/spf13/cobra"
)

// version can be overridden at build time with -ldflags
var version = "0.9.0"

// Config holds global CLI configuration
type Config struct {
	ConfigPath string
	Verbose    bool
	Quiet      bool
	Color      string
}

var globalConfig Config

// NewRootCommand creates and returns the root cobra command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deco",
		Short: "Manage game design documents as structured YAML",
		Long: `deco is a CLI tool for managing game design documents.

It provides structured YAML files with validation, references,
and dependency tracking for game design artifacts.`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Show help when no subcommand is provided
			return cmd.Help()
		},
		SilenceUsage:  true, // Don't show usage on errors
		SilenceErrors: true, // main.go handles error output
	}

	// Global flags
	cmd.PersistentFlags().StringVarP(&globalConfig.ConfigPath, "config", "c", ".deco", "Path to deco project directory")
	cmd.PersistentFlags().BoolVar(&globalConfig.Verbose, "verbose", false, "Enable verbose output")
	cmd.PersistentFlags().BoolVarP(&globalConfig.Quiet, "quiet", "q", false, "Suppress non-error output")
	cmd.PersistentFlags().StringVar(&globalConfig.Color, "color", "auto", "Color output: auto, always, never")

	// Initialize style system based on color flag
	cobra.OnInitialize(func() {
		style.SetMode(style.ParseColorMode(globalConfig.Color))
	})

	// Custom version template
	cmd.SetVersionTemplate(fmt.Sprintf("deco version %s\n", version))

	return cmd
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	return &globalConfig
}
