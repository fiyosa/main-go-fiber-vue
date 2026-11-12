package routes

import (
	"path/filepath"

	"github.com/gofiber/fiber/v3"

	platformauth "go-fiber-vue/internal/platform/auth"
	"go-fiber-vue/internal/platform/docs"
)

// RegisterDocs mounts the Scalar API reference UI and OpenAPI spec under /api/docs*
// guarded by HTTP Basic Authentication.
// It must be called before RegisterWeb so the SPA fallback doesn't shadow it.
func RegisterDocs(app *fiber.App, publicDir string, username, password string) {
	if app == nil {
		return
	}
	guard := platformauth.BasicAuthGuard(username, password)
	app.Get("/api/docs", guard, docs.Handler)
	app.Get("/api/docs/openapi.json", guard, docs.SpecHandler)

	// Serve specifically the local scalar esm bundle
	app.Get("/docs/scalar.js", guard, func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "application/javascript; charset=utf-8")
		filePath := filepath.Join(publicDir, "docs", "scalar.js")
		return c.SendFile(filePath)
	})
}
