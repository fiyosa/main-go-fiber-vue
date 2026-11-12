package controllers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/application/core/repository/auth"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	"go-fiber-vue/internal/app/presentation/core/middleware"
	"go-fiber-vue/internal/app/presentation/core/request"
	"go-fiber-vue/internal/app/presentation/core/response"
	"go-fiber-vue/internal/platform/lang"
)

// AuthController handles /api/v1/auth/* endpoints.
// All imports stay within: repository, request, response.
type AuthController struct {
	login  auth.LoginRepository
	logout auth.LogoutRepository
	// secure marks auth cookies Secure; it is disabled in development so the
	// cookies work over plain HTTP on localhost.
	secure bool
}

// NewAuthController wires the login and logout repositories (auth.LoginRepository,
// auth.LogoutRepository) and the cookie security policy resolved by the bootstrap.
func NewAuthController(login auth.LoginRepository, logout auth.LogoutRepository, secure bool) *AuthController {
	return &AuthController{login: login, logout: logout, secure: secure}
}

// Login handles POST /api/v1/auth/login.
func (c *AuthController) Login(fctx fiber.Ctx) error {
	var req request.LoginRequest
	if err := req.Bind(fctx); err != nil {
		return err
	}

	result, err := c.login.Login(
		fctx.Context(),
		req.Email,
		req.Password,
		auth.Meta{IP: fctx.IP(), UserAgent: fctx.Get("User-Agent")},
	)
	if err != nil {
		return err
	}

	auth.SetAuthCookies(
		fctx,
		result.Token,
		result.RefreshToken,
		time.Duration(result.ExpiresIn)*time.Second,
		time.Duration(result.RefreshExpiresIn)*time.Second,
		c.secure,
	)

	return fctx.Status(fiber.StatusOK).JSON(response.LoginSuccess(fctx))
}

// Me handles GET /api/v1/auth/me.
func (c *AuthController) Me(fctx fiber.Ctx) error {
	userID, ok := middleware.UserID(fctx)
	if !ok {
		return exceptions.Unauthorized("auth.unauthenticated", nil)
	}

	email, _ := fctx.Locals(middleware.LocalEmail).(string)

	operator := "Profile"
	if translator := lang.From(fctx); translator != nil {
		operator = translator.T(lang.Locale(fctx), "lang.profile", nil)
	}

	return response.Success(fctx, response.MeResource{
		ID:    userID,
		Email: email,
	}, "lang.retrieved", map[string]string{"operator": operator}, fiber.StatusOK)
}

// Logout handles POST /api/v1/auth/logout.
func (c *AuthController) Logout(fctx fiber.Ctx) error {
	token := middleware.AccessToken(fctx)
	auth.ClearAuthCookies(fctx, c.secure)

	if c.logout != nil {
		if err := c.logout.Logout(fctx.Context(), token); err != nil {
			slog.Warn("logout: session revoke failed", "error", err)
		}
	}

	return response.Success(fctx, nil, "lang.logged_out_successfully", nil, fiber.StatusOK)
}
