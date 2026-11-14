package controllers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/application/core/repository/auth"
	"go-fiber-vue/internal/app/application/core/repository/role"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	"go-fiber-vue/internal/app/presentation/core/middleware"
	"go-fiber-vue/internal/app/presentation/core/request"
	"go-fiber-vue/internal/app/presentation/core/response"
	"go-fiber-vue/internal/platform/lang"
)

// AuthController handles /api/v1/auth/* endpoints.
// All imports stay within: repository, request, response.
type AuthController struct {
	login          auth.LoginRepository
	logout         auth.LogoutRepository
	register       auth.RegisterRepository
	changePassword auth.ChangePasswordRepository
	roles          role.Repository
	// secure marks auth cookies Secure; it is disabled in development so the
	// cookies work over plain HTTP on localhost.
	secure bool
}

// NewAuthController wires the login, logout, register, and change-password
// repositories (auth.LoginRepository, auth.LogoutRepository,
// auth.RegisterRepository, auth.ChangePasswordRepository), the role reader
// used by Me, and the cookie security policy resolved by the bootstrap.
func NewAuthController(login auth.LoginRepository, logout auth.LogoutRepository, register auth.RegisterRepository, changePassword auth.ChangePasswordRepository, roles role.Repository, secure bool) *AuthController {
	return &AuthController{login: login, logout: logout, register: register, changePassword: changePassword, roles: roles, secure: secure}
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

	roleNames := []string{}
	if c.roles != nil {
		names, err := c.roles.NamesByUserID(fctx.Context(), userID)
		if err != nil {
			return err
		}
		roleNames = names
	}

	operator := "Profile"
	if translator := lang.From(fctx); translator != nil {
		operator = translator.T(lang.Locale(fctx), "lang.profile", nil)
	}

	return response.Success(fctx, response.MeResource{
		ID:    userID,
		Email: email,
		Roles: roleNames,
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

// Register handles POST /api/v1/auth/register.
func (c *AuthController) Register(fctx fiber.Ctx) error {
	var req request.RegisterRequest
	if err := req.Bind(fctx); err != nil {
		return err
	}

	if c.register == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	result, err := c.register.Register(fctx.Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	operator := "User"
	if translator := lang.From(fctx); translator != nil {
		operator = translator.T(lang.Locale(fctx), "lang.user", nil)
	}

	return response.Success(fctx, response.RegisterResource{
		ID:    result.ID,
		Email: result.Email,
	}, "lang.created", map[string]string{"operator": operator}, fiber.StatusCreated)
}

// ChangePassword handles POST /api/v1/auth/change-password. The caller must be
// authenticated; on success every session of the user is revoked and the client
// cookies are cleared so the next request requires a fresh login.
func (c *AuthController) ChangePassword(fctx fiber.Ctx) error {
	userID, ok := middleware.UserID(fctx)
	if !ok {
		return exceptions.Unauthorized("auth.unauthenticated", nil)
	}

	var req request.ChangePasswordRequest
	if err := req.Bind(fctx); err != nil {
		return err
	}

	if c.changePassword == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	if err := c.changePassword.ChangePassword(fctx.Context(), userID, req.NewPassword); err != nil {
		return err
	}

	auth.ClearAuthCookies(fctx, c.secure)

	return response.Success(fctx, nil, "lang.password_changed", nil, fiber.StatusOK)
}
