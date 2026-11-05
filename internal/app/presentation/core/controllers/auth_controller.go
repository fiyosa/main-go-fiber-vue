package controllers

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/application/core/repository"
	"go-fiber-vue/internal/app/presentation/core/request"
	"go-fiber-vue/internal/app/presentation/core/response"
)

// AuthController handles /api/v1/auth/* endpoints.
// All imports stay within: repository, request, response.
type AuthController struct {
	login repository.AuthLoginRepository
}

// NewAuthController wires the login repository (repository.AuthLoginRepository).
func NewAuthController(login repository.AuthLoginRepository) *AuthController {
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
		req.Username,
		req.Password,
		repository.Meta{IP: fctx.IP(), UserAgent: fctx.Get("User-Agent")},
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
