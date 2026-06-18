package routes

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterWeb(app *fiber.App) {
	app.Static("/", "public")
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("public/index.html")
	})
}
