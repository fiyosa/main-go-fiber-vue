package v1

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/presentation/core/controllers"
	coreroutes "go-fiber-vue/internal/app/presentation/core/routes"
)

type Dependencies struct {
	AuthController   *controllers.AuthController
	HealthController *controllers.HealthController
}

func Register(group fiber.Router, deps Dependencies) {
	health := group.Group("/")
	health.Get("healthz", deps.HealthController.Health)
	health.Get("readyz", deps.HealthController.Ready)

	coreroutes.RegisterAuthRoutes(group, deps.AuthController)
}
