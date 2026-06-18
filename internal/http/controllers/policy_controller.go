package controllers

import (
	"go-fiber-svelte/internal/http/repositories/policy_repository"

	"github.com/gofiber/fiber/v2"
)

func RoleList(c *fiber.Ctx) error {
	return policy_repository.RoleListRepository(c)
}

func PermissionList(c *fiber.Ctx) error {
	return policy_repository.PermissionListRepository(c)
}

func PermissionStore(c *fiber.Ctx) error {
	return policy_repository.PermissionStoreRepository(c)
}

func PermissionDestroy(c *fiber.Ctx) error {
	return policy_repository.PermissionDestroyRepository(c)
}

func PolicyOpenAPIPaths() map[string]any {
	security := []map[string]any{{"bearer": []string{}}}
	return map[string]any{
		"/api/policy/role": map[string]any{
			"get": map[string]any{
				"summary":     "List roles",
				"description": "Get all roles",
				"tags":        []string{"Policy"},
				"security":    security,
				"responses": map[string]any{
					"200": map[string]any{"description": "Role list"},
				},
			},
		},
		"/api/policy/permission": map[string]any{
			"get": map[string]any{
				"summary":     "List permissions",
				"description": "Get all permissions",
				"tags":        []string{"Policy"},
				"security":    security,
				"responses": map[string]any{
					"200": map[string]any{"description": "Permission list"},
				},
			},
			"post": map[string]any{
				"summary":     "Create permission",
				"description": "Add a new permission",
				"tags":        []string{"Policy"},
				"security":    security,
				"responses": map[string]any{
					"201": map[string]any{"description": "Permission created"},
				},
			},
		},
		"/api/policy/permission/{id}": map[string]any{
			"delete": map[string]any{
				"summary":     "Delete permission",
				"description": "Remove a permission by ID",
				"tags":        []string{"Policy"},
				"security":    security,
				"parameters": []map[string]any{
					{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Permission deleted"},
				},
			},
		},
	}
}
