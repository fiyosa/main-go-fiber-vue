package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// renderHelp prints a Laravel-artisan style command list for the root
// command, and defers to cobra's default help for every subcommand.
func renderHelp(cmd *cobra.Command, _ []string) {
	root := cmd.Root()
	if cmd != root {
		// Detach the custom help function so cobra renders the default
		// per-command help, then restore it afterwards.
		root.SetHelpFunc(nil)
		defer root.SetHelpFunc(renderHelp)
		_ = cmd.Help()
		return
	}
	renderHome(root)
}

// renderHome writes the grouped command index.
func renderHome(root *cobra.Command) {
	fmt.Printf("%s %s\n\n", consoleName, consoleVersion)
	fmt.Println("Usage:")
	fmt.Printf("  %s [command] [flags]\n\n", root.Name())

	var leaves []*cobra.Command
	var groups []*cobra.Command
	for _, c := range root.Commands() {
		if c.Hidden {
			continue
		}
		// Completion is a grouped command but reads better as a single entry.
		if c.HasSubCommands() && c.Name() != "completion" {
			groups = append(groups, c)
			continue
		}
		leaves = append(leaves, c)
	}

	// Size the name column so every description lines up.
	width := 0
	for _, c := range leaves {
		if l := len(c.Name()); l > width {
			width = l
		}
	}
	for _, g := range groups {
		for _, c := range g.Commands() {
			if c.Hidden {
				continue
			}
			if l := len(c.Name()); l > width {
				width = l
			}
		}
	}
	width += 4

	fmt.Println("Available commands:")
	for _, c := range leaves {
		fmt.Printf("  %-*s %s\n", width, c.Name(), c.Short)
	}

	for _, g := range groups {
		fmt.Printf("\n %s\n", g.Name())
		for _, c := range g.Commands() {
			if c.Hidden {
				continue
			}
			fmt.Printf("  %-*s %s\n", width, c.Name(), c.Short)
		}
	}

	fmt.Printf("\nRun '%s [command] --help' for more information on a command.\n", root.Name())
}
