package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	consoleName    = "App Console"
	consoleVersion = "1.0"
)

// Root returns the top-level console command that all subcommands attach to.
func Root() *cobra.Command {
	root := &cobra.Command{
		Use:     "console",
		Short:   "App one-off commands",
		Long:    "App one-off commands.",
		Version: consoleVersion,
	}
	root.SetVersionTemplate(fmt.Sprintf("%s %s\n", consoleName, consoleVersion))
	root.SetHelpFunc(renderHelp)
	root.AddCommand(Migrate(), DB(), Routes(), OpenAPI())
	return root
}

// DB returns a cobra.Command that groups all database subcommands (e.g. db seed).
func DB() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database maintenance",
	}
	cmd.AddCommand(Seed(), UserCreate())
	return cmd
}

// Execute runs the console application and exits on error.
func Execute() {
	if err := Root().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
