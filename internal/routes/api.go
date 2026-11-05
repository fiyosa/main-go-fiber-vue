package routes

import (
	"github.com/gofiber/fiber/v3"

	v1 "go-fiber-vue/internal/routes/v1"
)

func RegisterAPI(app *fiber.App, deps v1.Dependencies) {
	api := app.Group("/api")
	v1.Register(api.Group("/v1"), deps)
}
