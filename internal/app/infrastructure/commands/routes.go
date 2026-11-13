package commands

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/cobra"

	"go-fiber-vue/internal/bootstrap"
	"go-fiber-vue/internal/config"
)

// apiPrefix is the mount point of every API route (see routes.RegisterAPI and
// routes.RegisterDocs). The MCP endpoint is excluded separately because it is
// only listed by --mcp.
const apiPrefix = "/api"

// routeScope selects which slice of the routing table the command prints.
type routeScope int

const (
	scopeAPI routeScope = iota // default: only /api/v1 routes
	scopeMCP                   // only the MCP endpoint
	scopeAll                   // every route, middleware included
)

// Routes returns a cobra.Command that lists registered HTTP routes.
func Routes() *cobra.Command {
	var all bool
	var mcpOnly bool

	cmd := &cobra.Command{
		Use:   "routes",
		Short: "List registered routes",
		Long: "List registered routes.\n\n" +
			"By default only API routes are shown. Use --mcp to list only the MCP\n" +
			"endpoint, or --all to include every route (middleware, docs, web).",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if all && mcpOnly {
				return fmt.Errorf("--all and --mcp cannot be combined")
			}

			cfg, err := config.LoadForConsole()
			if err != nil {
				return err
			}

			scope := scopeAPI
			switch {
			case all:
				scope = scopeAll
			case mcpOnly:
				scope = scopeMCP
			}

			listRoutes(bootstrap.Router(cfg), scope, cfg.MCPPath)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "include every route (middleware, mcp, docs, web)")
	cmd.Flags().BoolVarP(&mcpOnly, "mcp", "m", false, "show only the MCP endpoint")
	return cmd
}

// listRoutes prints a Laravel-style route table for the requested scope.
func listRoutes(app *fiber.App, scope routeScope, mcpPath string) {
	fmt.Printf("%-8s %s\n", "Method", "URI")
	fmt.Println(strings.Repeat("-", 60))

	// GetRoutes' filterUse flag skips middleware routes when true.
	includeMiddleware := scope == scopeAll

	count := 0
	for _, r := range app.GetRoutes(!includeMiddleware) {
		if !inScope(r.Path, scope, mcpPath) {
			continue
		}
		fmt.Printf("%-8s %s\n", r.Method, r.Path)
		count++
	}

	if count == 0 {
		fmt.Println("no routes registered")
	}
}

// inScope reports whether path belongs to the selected route scope.
func inScope(path string, scope routeScope, mcpPath string) bool {
	switch scope {
	case scopeAll:
		return true
	case scopeMCP:
		return mcpPath != "" && path == mcpPath
	default:
		if path == mcpPath {
			return false
		}
		return path == apiPrefix || strings.HasPrefix(path, apiPrefix+"/")
	}
}
