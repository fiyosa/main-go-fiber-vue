package middleware

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/presentation/core/response"
)

// Maintenance returns a middleware that blocks all incoming requests with 503
// when enabled. It must be registered after Locale so that the error message
// is translated correctly.
func Maintenance(enabled bool) fiber.Handler {
	return func(c fiber.Ctx) error {
		if !enabled {
			return c.Next()
		}
		return response.Error(c, "common.maintenance", nil, fiber.StatusServiceUnavailable)
	}
}
