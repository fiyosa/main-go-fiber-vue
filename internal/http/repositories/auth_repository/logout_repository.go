package auth_repository

import (
	"go-fiber-svelte/internal/db"
	"go-fiber-svelte/internal/db/models"
	"go-fiber-svelte/internal/helper"
	"go-fiber-svelte/internal/lang"

	"github.com/gofiber/fiber/v2"
)

func LogoutRepository(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(int)
	database := db.RUN
	database.Model(&models.Auth{}).
		Where("user_id = ? AND revoke = ?", userId, false).
		Update("revoke", true)
	c.ClearCookie("token")
	return c.JSON(helper.Res.Success(lang.T.Convert(lang.T.Get().SAVED_SUCCESSFULLY, map[string]any{"operator": lang.T.Get().LOGOUT})))
}
