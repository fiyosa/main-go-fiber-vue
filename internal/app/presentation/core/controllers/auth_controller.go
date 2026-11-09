package controllers

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/application/core/repository/auth"
	"go-fiber-vue/internal/app/presentation/core/request"
	"go-fiber-vue/internal/app/presentation/core/response"
)

// AuthController handles /api/v1/auth/* endpoints.
// All imports stay within: repository, request, response.
type AuthController struct {
	login auth.LoginRepository
	// secure marks auth cookies Secure; it is disabled in development so the
	// cookies work over plain HTTP on localhost.
	secure bool
}

// NewAuthController wires the login repository (auth.LoginRepository) and the
// cookie security policy resolved by the bootstrap.
func NewAuthController(login auth.LoginRepository, secure bool) *AuthController {
	return &AuthController{login: login, secure: secure}
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
