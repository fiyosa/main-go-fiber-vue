package provider

import (
	"go-fiber-svelte/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterMiddleware(app *fiber.App) {
	app.Use(middleware.HashMiddleware())
}
