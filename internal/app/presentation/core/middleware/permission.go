// Package middleware holds HTTP middleware that is specific to the API
// presentation layer: authentication from the access token and permission
// enforcement.
package middleware

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/application/core/repository/auth"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	platformauth "go-fiber-vue/internal/platform/auth"
	"go-fiber-vue/internal/platform/permission"
)

// Locals keys set on every authenticated request.
const (
	LocalUserID = "auth.user_id"
	LocalEmail  = "auth.email"
)

// TokenValidator verifies the access token and returns its claims.
type TokenValidator interface {
	ValidateAccessToken(raw string) (*platformauth.Claims, error)
}

// Authorizer answers whether a user holds a permission.
type Authorizer interface {
	Can(ctx context.Context, userID uint64, permission string) (bool, error)
}

// Permission returns the factory used at route registration:
//
//	perm := middleware.Permission(tokens, authorizer)
//	group.Get("/roles", perm("get-roles"), controller.Index)
//
// The returned handler authenticates the request, stores the caller in locals,
// then enforces the permission. It never panics when the token service or
// authorizer is missing; it answers 503 instead so console commands can still
// build the route table.
func Permission(tokens TokenValidator, authorizer Authorizer) func(string) fiber.Handler {
	return func(raw string) fiber.Handler {
		name := strings.TrimSpace(raw)
		if name == "" {
			panic("middleware: permission name must not be empty")
		}
		permission.Register(name)

		return func(c fiber.Ctx) error {
			if tokens == nil || authorizer == nil {
				return exceptions.ServiceUnavailable(nil)
			}

			userID, err := authenticate(c, tokens)
			if err != nil {
				return err
			}

			allowed, err := authorizer.Can(c.Context(), userID, name)
			if err != nil {
				return err
			}
			if !allowed {
				return exceptions.Forbidden(nil)
			}
			return c.Next()
		}
	}
}

// authenticate validates the access token from the cookie (preferred) or the
// Authorization bearer header and stores the caller in locals.
func authenticate(c fiber.Ctx, tokens TokenValidator) (uint64, error) {
	raw := AccessToken(c)
	if raw == "" {
		return 0, exceptions.Unauthorized("auth.unauthenticated", nil)
	}

	claims, err := tokens.ValidateAccessToken(raw)
	if err != nil || claims == nil {
		return 0, exceptions.Unauthorized("auth.invalid_token", nil)
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil || userID == 0 {
		return 0, exceptions.Unauthorized("auth.invalid_token", nil)
	}

	c.Locals(LocalUserID, userID)
	c.Locals(LocalEmail, claims.Email)
	return userID, nil
}

// AccessToken extracts the access token string from the cookie or Authorization header.
func AccessToken(c fiber.Ctx) string {
	if token := strings.TrimSpace(c.Cookies(auth.AccessCookieName)); token != "" {
		return token
	}

	header := strings.TrimSpace(c.Get("Authorization"))
	if len(header) > 7 && strings.EqualFold(header[:7], "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return ""
}

// UserID returns the caller id stored by Permission.
func UserID(c fiber.Ctx) (uint64, bool) {
	value, ok := c.Locals(LocalUserID).(uint64)
	return value, ok
}
