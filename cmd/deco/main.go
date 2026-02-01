package main

import (
	"fmt"
	"os"

	"github.com/Toernblom/deco/internal/cli"
)

func main() {
	root := cli.NewRootCommand()

	// Register subcommands
	root.AddCommand(cli.NewInitCommand())
	root.AddCommand(cli.NewCreateCommand())
	root.AddCommand(cli.NewRmCommand())
	root.AddCommand(cli.NewIssuesCommand())
	root.AddCommand(cli.NewValidateCommand())
	root.AddCommand(cli.NewListCommand())
	root.AddCommand(cli.NewShowCommand())
	root.AddCommand(cli.NewQueryCommand())
	root.AddCommand(cli.NewSetCommand())
	root.AddCommand(cli.NewAppendCommand())
	root.AddCommand(cli.NewUnsetCommand())
	root.AddCommand(cli.NewApplyCommand())
	root.AddCommand(cli.NewHistoryCommand())
	root.AddCommand(cli.NewMvCommand())
	root.AddCommand(cli.NewGraphCommand())
	root.AddCommand(cli.NewDiffCommand())

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
