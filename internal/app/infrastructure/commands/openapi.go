package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"go-fiber-vue/internal/bootstrap"
	"go-fiber-vue/internal/config"
	"go-fiber-vue/internal/platform/docs"
)

// OpenAPI returns a cobra.Command that generates or checks public/docs/openapi.json.
func OpenAPI() *cobra.Command {
	var outPath string
	var overridePath string
	var check bool

	cmd := &cobra.Command{
		Use:   "openapi",
		Short: "Generate OpenAPI 3.1 JSON spec from registered API routes",
		Long: "Generate OpenAPI 3.1 JSON spec from registered API routes.\n\n" +
			"Discovers API routes from the application router, excludes MCP and docs,\n" +
			"applies optional manual overrides, and writes to public/docs/openapi.json.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if outPath == "" {
				outPath = filepath.Join(cfg.PublicDir, "docs", "openapi.json")
			}

			if overridePath == "" {
				overridePath = filepath.Join("internal", "platform", "docs", "openapi.override.json")
			}

			app := bootstrap.Router(cfg)
			fiberRoutes := app.GetRoutes(true)

			var routes []docs.RouteInfo
			for _, r := range fiberRoutes {
				if docs.IsAPIRoute(r.Path, cfg.MCPPath) {
					routes = append(routes, docs.RouteInfo{
						Method: r.Method,
						Path:   r.Path,
					})
				}
			}

			var overrideBytes []byte
			if data, err := os.ReadFile(overridePath); err == nil {
				overrideBytes = data
			}

			spec, err := docs.BuildSpec(routes, overrideBytes)
			if err != nil {
				return fmt.Errorf("build spec: %w", err)
			}

			outData, err := json.MarshalIndent(spec, "", "  ")
			outData = append(outData, '\n')
			if err != nil {
				return fmt.Errorf("marshal spec: %w", err)
			}

			if check {
				existing, err := os.ReadFile(outPath)
				if err != nil {
					return fmt.Errorf("read existing spec %s: %w", outPath, err)
				}
				if string(existing) != string(outData) {
					return fmt.Errorf("openapi spec is out of sync; run 'console openapi' to update")
				}
				fmt.Println("openapi spec is up to date.")
				return nil
			}

			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return fmt.Errorf("mkdir output dir: %w", err)
			}

			if err := os.WriteFile(outPath, outData, 0644); err != nil {
				return fmt.Errorf("write spec %s: %w", outPath, err)
			}

			fmt.Printf("Successfully generated OpenAPI spec at %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output file path (default public/docs/openapi.json)")
	cmd.Flags().StringVar(&overridePath, "override", "", "manual JSON override file path (default docs/openapi.override.json)")
	cmd.Flags().BoolVar(&check, "check", false, "check if spec is up to date (exit 1 if not)")

	return cmd
}
