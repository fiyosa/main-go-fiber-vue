package policy_repository

import (
	"go-fiber-svelte/internal/db"
	"go-fiber-svelte/internal/db/models"
	"go-fiber-svelte/internal/helper"
	"go-fiber-svelte/internal/lang"

	"github.com/gofiber/fiber/v2"
)

func PermissionDestroyRepository(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.Res.Error("Invalid ID", nil))
	}
	database := db.RUN
	result := database.Delete(&models.Permission{}, id)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(helper.Res.Error(lang.T.Convert(lang.T.Get().NOT_FOUND, map[string]any{"operator": "Permission"}), nil))
	}
	return c.JSON(helper.Res.Success(lang.T.Convert(lang.T.Get().DELETED_SUCCESSFULLY, map[string]any{"operator": "Permission"})))
}
