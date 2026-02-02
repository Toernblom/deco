package cli

import (
	"fmt"

	"github.com/Toernblom/deco/internal/migrations"
	"github.com/spf13/cobra"
)

type migrateFlags struct {
	dryRun    bool
	noBackup  bool
	quiet     bool
	targetDir string
}

// NewMigrateCommand creates the migrate subcommand.
func NewMigrateCommand() *cobra.Command {
	flags := &migrateFlags{}

	cmd := &cobra.Command{
		Use:   "migrate [directory]",
		Short: "Migrate nodes to current schema version",
		Long: `Migrate all nodes to match the current schema configuration.

When schema configuration (custom_block_types, schema_rules) changes,
nodes may need to be updated. This command:

1. Creates a backup (unless --no-backup)
2. Applies registered migration transforms to each node
3. Updates the schema_version in config
4. Logs migration to audit history

Exit codes:
  0: Migration successful (or no migration needed)
  1: Error occurred`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				flags.targetDir = args[0]
			} else {
				flags.targetDir = "."
			}
			return runMigrate(flags)
		},
	}

	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Show what would change without making changes")
	cmd.Flags().BoolVar(&flags.noBackup, "no-backup", false, "Skip creating backup before migration")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress non-essential output")

	return cmd
}

func runMigrate(flags *migrateFlags) error {
	// Check if migration is needed
	needs, currentHash, expectedHash, err := migrations.NeedsMigration(flags.targetDir)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if !needs {
		if !flags.quiet {
			fmt.Println("✓ Schema is up to date, no migration needed")
		}
		return nil
	}

	// Show what will happen
	if !flags.quiet {
		if flags.dryRun {
			fmt.Println("Dry run - showing what would change:")
		}
		fmt.Printf("Schema version: %s -> %s\n", formatHash(currentHash), formatHash(expectedHash))
	}

	// Execute migration
	executor := migrations.NewExecutor(migrations.ExecutorOptions{
		DryRun:    flags.dryRun,
		NoBackup:  flags.noBackup,
		Quiet:     flags.quiet,
		TargetDir: flags.targetDir,
	}, nil)

	result, err := executor.Execute()
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Report results
	if !flags.quiet {
		if result.DryRun {
			fmt.Printf("\nWould process %d node(s), modify %d node(s)\n",
				result.NodesProcessed, result.NodesModified)
		} else {
			if result.BackupDir != "" {
				fmt.Printf("Backup created: %s\n", result.BackupDir)
			}
			fmt.Printf("✓ Migrated %d node(s) (%d modified)\n",
				result.NodesProcessed, result.NodesModified)
		}
	}

	return nil
}

// formatHash formats a schema hash for display.
func formatHash(hash string) string {
	if hash == "" {
		return "(none)"
	}
	return hash
}
