package commands

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/cobra"

	"go-fiber-vue/internal/bootstrap"
	"go-fiber-vue/internal/config"
)

// Routes returns a cobra.Command that lists all registered HTTP routes.
// Middleware routes are hidden by default; pass --all to include them.
func Routes() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "routes",
		Short: "List all registered routes",
		Long:  "List all registered routes. Middleware routes are hidden by default; use --all to include them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			listRoutes(bootstrap.Router(cfg), all)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "include middleware routes")
	return cmd
}

// listRoutes prints a Laravel-style route table.
func listRoutes(app *fiber.App, includeMiddleware bool) {
	fmt.Printf("%-8s %s\n", "Method", "URI")
	fmt.Println(strings.Repeat("-", 60))

	// GetRoutes' filterUse flag skips middleware routes when true.
	for _, r := range app.GetRoutes(!includeMiddleware) {
		fmt.Printf("%-8s %s\n", r.Method, r.Path)
	}
}
