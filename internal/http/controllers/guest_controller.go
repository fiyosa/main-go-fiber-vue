package controllers

import (
	"go-fiber-svelte/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func Ping(c *fiber.Ctx) error {
	return c.JSON(helper.Res.Success("pong"))
}

func GuestOpenAPIPaths() map[string]any {
	return map[string]any{
		"/api/guest/ping": map[string]any{
			"get": map[string]any{
				"summary":     "Ping",
				"description": "Health check endpoint",
				"tags":        []string{"Guest"},
				"responses": map[string]any{
					"200": map[string]any{"description": "Pong"},
				},
			},
		},
	}
}
