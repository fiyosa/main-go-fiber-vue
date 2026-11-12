package v1

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/presentation/core/controllers"
	"go-fiber-vue/internal/app/presentation/core/middleware"
	coreroutes "go-fiber-vue/internal/app/presentation/core/routes"
	"go-fiber-vue/internal/platform/permission"
)

// Dependencies is everything the /api/v1 route table needs. Tokens and
// Authorizer may be nil: the route table is still built so console commands can
// inspect it, and protected routes answer 503 instead of panicking.
type Dependencies struct {
	AuthController   *controllers.AuthController
	HealthController *controllers.HealthController
	RoleController   *controllers.RoleController

	Tokens     middleware.TokenValidator
	Authorizer middleware.Authorizer
}

// Register mounts every /api/v1 route. It resets the permission registry first
// so a repeated build in the same process starts from an empty list.
func Register(group fiber.Router, deps Dependencies) {
	permission.Reset()

	perm := middleware.Permission(deps.Tokens, deps.Authorizer)

	health := group.Group("/")
	health.Get("healthz", deps.HealthController.Health)
	health.Get("readyz", deps.HealthController.Ready)

	coreroutes.RegisterAuthRoutes(group, deps.AuthController, deps.Tokens)
	coreroutes.RegisterRoleRoutes(group, perm, deps.RoleController)
}
