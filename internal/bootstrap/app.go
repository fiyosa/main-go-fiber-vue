package bootstrap

import (
	"go-fiber-svelte/internal/config"
	"go-fiber-svelte/internal/db"
	"go-fiber-svelte/internal/provider"
	"go-fiber-svelte/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func Init() {
	config.InitConfigApp()
	config.I18n()
	db.Init()
}

func RegisterApp(app *fiber.App) {
	provider.RegisterMiddleware(app)
	routes.RegisterAPI(app)
	routes.RegisterWeb(app)
}
