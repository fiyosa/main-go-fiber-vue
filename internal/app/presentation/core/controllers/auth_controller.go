package controllers

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/application/core/repository/auth"
	"go-fiber-vue/internal/app/presentation/core/request"
	"go-fiber-vue/internal/app/presentation/core/response"
)

// AuthController handles /api/v1/auth/* endpoints.
// All imports stay within: repository, request, response.
type AuthController struct {
	login auth.LoginRepository
}

// NewAuthController wires the login repository (auth.LoginRepository).
func NewAuthController(login auth.LoginRepository) *AuthController {
	return &AuthController{login: login}
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

	return fctx.Status(fiber.StatusOK).JSON(response.LoginSuccess(
		fctx,
		result.Token,
		result.RefreshToken,
		result.ExpiresIn,
	))
}
