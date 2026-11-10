package routes

import (
	"path/filepath"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/platform/docs"
)

// RegisterDocs mounts the Scalar API reference UI and OpenAPI spec under /api/docs*.
// It must be called before RegisterWeb so the SPA fallback doesn't shadow it.
func RegisterDocs(app *fiber.App, publicDir string) {
	if app == nil {
		return
	}
	app.Get("/api/docs", docs.Handler)
	app.Get("/api/docs/openapi.json", docs.SpecHandler)

	// Serve specifically the local scalar esm bundle
	app.Get("/docs/scalar.js", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "application/javascript; charset=utf-8")
		filePath := filepath.Join(publicDir, "docs", "scalar.js")
		return c.SendFile(filePath)
	})
}
