package routes

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/presentation/core/controllers"
	"go-fiber-vue/internal/app/presentation/core/middleware"
)

// RegisterAuthRoutes mounts every /api/v1/auth/* endpoint on the controller.
func RegisterAuthRoutes(group fiber.Router, controller *controllers.AuthController, tokens middleware.TokenValidator) {
	auth := group.Group("/auth")
	auth.Post("/login", controller.Login)
	auth.Post("/logout", controller.Logout)
	auth.Get("/me", middleware.Authenticate(tokens), controller.Me)
}
