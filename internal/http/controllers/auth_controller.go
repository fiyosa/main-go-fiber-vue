package controllers

import (
	"go-fiber-svelte/internal/http/repositories/auth_repository"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	return auth_repository.LoginRepository(c)
}

func Logout(c *fiber.Ctx) error {
	return auth_repository.LogoutRepository(c)
}

func User(c *fiber.Ctx) error {
	return auth_repository.UserRepository(c)
}

func AuthOpenAPIPaths() map[string]any {
	security := []map[string]any{{"bearer": []string{}}}
	return map[string]any{
		"/api/auth/login": map[string]any{
			"post": map[string]any{
				"summary":     "Login",
				"description": "Authenticate user and return JWT token",
				"tags":        []string{"Auth"},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"email":    map[string]any{"type": "string", "format": "email"},
									"password": map[string]any{"type": "string", "format": "password"},
								},
							},
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Login successful"},
					"401": map[string]any{"description": "Invalid credentials"},
				},
			},
		},
		"/api/auth/logout": map[string]any{
			"delete": map[string]any{
				"summary":     "Logout",
				"description": "Invalidate current session",
				"tags":        []string{"Auth"},
				"security":    security,
				"responses": map[string]any{
					"200": map[string]any{"description": "Logout successful"},
				},
			},
		},
		"/api/auth/user": map[string]any{
			"get": map[string]any{
				"summary":     "Get current user",
				"description": "Return authenticated user profile",
				"tags":        []string{"Auth"},
				"security":    security,
				"responses": map[string]any{
					"200": map[string]any{"description": "User data"},
				},
			},
		},
	}
}
