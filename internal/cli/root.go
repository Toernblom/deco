package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version can be overridden at build time with -ldflags
var version = "0.1.0"

// Config holds global CLI configuration
type Config struct {
	ConfigPath string
	Verbose    bool
	Quiet      bool
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
		SilenceUsage: true, // Don't show usage on errors
	}

	// Global flags
	cmd.PersistentFlags().StringVarP(&globalConfig.ConfigPath, "config", "c", ".deco", "Path to deco project directory")
	cmd.PersistentFlags().BoolVar(&globalConfig.Verbose, "verbose", false, "Enable verbose output")
	cmd.PersistentFlags().BoolVarP(&globalConfig.Quiet, "quiet", "q", false, "Suppress non-error output")

	// Custom version template
	cmd.SetVersionTemplate(fmt.Sprintf("deco version %s\n", version))

	return cmd
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	return &globalConfig
}
