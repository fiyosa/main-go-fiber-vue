package middleware

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// Authenticate returns a handler that verifies the access token and sets caller locals
// without enforcing any permission. It responds with 401 Unauthorized when unauthenticated
// or invalid, and 503 ServiceUnavailable when the token service is nil.
func Authenticate(tokens TokenValidator) fiber.Handler {
	return func(c fiber.Ctx) error {
		if tokens == nil {
			return exceptions.ServiceUnavailable(nil)
		}

		if _, err := authenticate(c, tokens); err != nil {
			return err
		}

		return c.Next()
	}
}
