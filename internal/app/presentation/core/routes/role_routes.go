package routes

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/presentation/core/controllers"
)

// RegisterRoleRoutes mounts every /api/v1/roles endpoint.
func RegisterRoleRoutes(group fiber.Router, perm func(string) fiber.Handler, controller *controllers.RoleController) {
	group.Get("/roles", perm("get-roles"), controller.Index)
}
