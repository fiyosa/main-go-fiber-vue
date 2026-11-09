package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// CORS allows the configured browser origins to call the API cross-origin. The
// setting is a comma-separated list; an empty list disables CORS entirely by
// returning a pass-through handler, because Fiber treats an empty allowlist as
// "allow every origin" and we never want that by accident.
func CORS(origins string) fiber.Handler {
	allowed := make([]string, 0)
	for _, origin := range strings.Split(origins, ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			allowed = append(allowed, trimmed)
		}
	}

	if len(allowed) == 0 {
		return func(c fiber.Ctx) error { return c.Next() }
	}

	return cors.New(cors.Config{
		AllowOrigins: allowed,
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           600,
	})
}
