package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Toernblom/deco/internal/cli"
)

func main() {
	root := cli.NewRootCommand()

	// Register subcommands
	root.AddCommand(cli.NewInitCommand())
	root.AddCommand(cli.NewIssuesCommand())
	root.AddCommand(cli.NewValidateCommand())
	root.AddCommand(cli.NewListCommand())
	root.AddCommand(cli.NewShowCommand())
	root.AddCommand(cli.NewQueryCommand())
	root.AddCommand(cli.NewHistoryCommand())
	root.AddCommand(cli.NewGraphCommand())
	root.AddCommand(cli.NewDiffCommand())
	root.AddCommand(cli.NewStatsCommand())
	root.AddCommand(cli.NewReviewCommand())
	root.AddCommand(cli.NewSyncCommand())
	root.AddCommand(cli.NewMigrateCommand())

	if err := root.Execute(); err != nil {
		// Check for ExitError with custom exit code
		var exitErr *cli.ExitError
		if errors.As(err, &exitErr) {
			fmt.Fprintf(os.Stderr, "Error: %v\n", exitErr.Message)
			os.Exit(exitErr.Code)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
