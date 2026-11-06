package routes

import (
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"

	"go-fiber-vue/internal/app/presentation/core/view"
)

// RegisterWeb serves the built SPA assets and the SPA shell fallback.
// It must be called after RegisterAPI so the fallback never shadows API routes.
func RegisterWeb(app *fiber.App, publicDir string) {
	buildDir := filepath.Join(publicDir, "build")

	// /build/* holds the Vite output: index.js, index.css and assets/*.
	app.Use("/build", static.New(buildDir, static.Config{
		MaxAge:   3600,
		Compress: true,
		NotFoundHandler: func(c fiber.Ctx) error {
			return fiber.ErrNotFound
		},
	}))

	// The favicon lives at the root of public/ and is linked as /vite.svg.
	app.Get("/vite.svg", func(c fiber.Ctx) error {
		return c.SendFile(filepath.Join(publicDir, "vite.svg"))
	})

	// Every other non-API path renders the SPA shell so the frontend router
	// can take over. Unknown API paths keep their JSON 404.
	app.Get("/*", func(c fiber.Ctx) error {
		path := c.Path()
		switch {
		case strings.HasPrefix(path, "/api"):
			return c.Next()
		case strings.HasPrefix(path, "/build"):
			return fiber.ErrNotFound
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		c.Set(fiber.HeaderCacheControl, "no-cache")
		return c.SendString(view.SPA)
	})
}
