package routes

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/presentation/core/controllers"
)

// RegisterAuthRoutes mounts every /api/v1/auth/* endpoint on the controller.
func RegisterAuthRoutes(group fiber.Router, controller *controllers.AuthController) {
	auth := group.Group("/auth")
	auth.Post("/login", controller.Login)
	// Future endpoints live here: register, refresh, logout, me...
}
