package auth

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
)

// BasicAuthGuard returns a Fiber handler enforcing HTTP Basic Authentication
// against the provided username and password using constant-time comparisons.
// If either the username or password is empty, access is denied for all requests.
func BasicAuthGuard(username, password string) fiber.Handler {
	return basicauth.New(basicauth.Config{
		Authorizer: func(u, p string, c fiber.Ctx) bool {
			if username == "" || password == "" {
				return false
			}
			uMatch := subtle.ConstantTimeCompare([]byte(u), []byte(username)) == 1
			pMatch := subtle.ConstantTimeCompare([]byte(p), []byte(password)) == 1
			return uMatch && pMatch
		},
		Realm: "Restricted",
	})
}
