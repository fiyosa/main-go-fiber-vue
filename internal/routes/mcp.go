package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
)

// RegisterMCP mounts the MCP handler on every method the Streamable HTTP
// transport uses. It is a no-op when the handler is nil so a disabled MCP server
// leaves /api/mcp to the regular API 404 path.
//
// RegisterMCP must run before RegisterWeb: the SPA fallback only skips /api/*,
// but mounting first keeps the intent explicit and safe if the fallback changes.
func RegisterMCP(app *fiber.App, path string, handler http.Handler) {
	if app == nil || handler == nil || path == "" {
		return
	}
	app.All(path, adaptor.HTTPHandler(handler))
}
