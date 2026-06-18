package policy_repository

import (
	"go-fiber-svelte/internal/db"
	"go-fiber-svelte/internal/db/models"
	"go-fiber-svelte/internal/helper"
	"go-fiber-svelte/internal/http/resources/policy_resource"
	"go-fiber-svelte/internal/lang"

	"github.com/gofiber/fiber/v2"
)

func PermissionListRepository(c *fiber.Ctx) error {
	database := db.RUN
	var permissions []models.Permission
	database.Find(&permissions)
	return c.JSON(helper.Res.SuccessData(policy_resource.PermissionToResource(permissions), lang.T.Convert(lang.T.Get().RETRIEVED_SUCCESSFULLY, map[string]any{"operator": "Permission"})))
}
