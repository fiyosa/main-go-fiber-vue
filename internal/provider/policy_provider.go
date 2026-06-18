package provider

import (
	"go-fiber-svelte/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func Policy(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id").(int)
		if !NewAuthProvider().CheckPermission(userId, permission) {
			return c.Status(fiber.StatusForbidden).JSON(helper.Res.Error("Forbidden", nil))
		}
		return c.Next()
	}
}
