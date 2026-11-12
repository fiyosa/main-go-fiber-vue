package controllers

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/application/core/repository/role"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	"go-fiber-vue/internal/app/presentation/core/response"
	"go-fiber-vue/internal/platform/hashid"
	"go-fiber-vue/internal/platform/lang"
)

// RoleController handles /api/v1/roles* endpoints.
type RoleController struct {
	roles role.Repository
	ids   *hashid.Encoder
}

// NewRoleController wires the role repository and the id encoder.
func NewRoleController(roles role.Repository, ids *hashid.Encoder) *RoleController {
	return &RoleController{roles: roles, ids: ids}
}

type permissionResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type roleResource struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Permissions []permissionResource `json:"permissions"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// Index handles GET /api/v1/roles. It returns every role with its permissions,
// ordered by name. Numeric ids are never exposed: every id is a hashid.
func (c *RoleController) Index(fctx fiber.Ctx) error {
	if c == nil || c.roles == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	rows, err := c.roles.List(fctx.Context())
	if err != nil {
		return err
	}

	data := make([]roleResource, 0, len(rows))
	for _, row := range rows {
		id, err := c.ids.Encode(row.ID)
		if err != nil {
			return exceptions.Internal(err)
		}

		permissions := make([]permissionResource, 0, len(row.Permissions))
		for _, permission := range row.Permissions {
			permissionID, err := c.ids.Encode(permission.ID)
			if err != nil {
				return exceptions.Internal(err)
			}
			permissions = append(permissions, permissionResource{
				ID:   permissionID,
				Name: permission.Name,
			})
		}

		data = append(data, roleResource{
			ID:          id,
			Name:        row.Name,
			Permissions: permissions,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	operator := "Role"
	if translator := lang.From(fctx); translator != nil {
		operator = translator.T(lang.Locale(fctx), "lang.role", nil)
	}

	return response.Success(fctx, data, "lang.retrieved", map[string]string{"operator": operator}, fiber.StatusOK)
}
