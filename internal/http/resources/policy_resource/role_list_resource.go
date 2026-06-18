package policy_resource

import (
	"go-fiber-svelte/internal/db/models"
	"go-fiber-svelte/internal/lib"
)

type RoleListResource struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Notes       string               `json:"notes"`
	Permissions []PermissionResource `json:"permissions"`
}

func RoleListToResource(roles []models.Role) []RoleListResource {
	var result []RoleListResource
	for _, role := range roles {
		encodedID, _ := lib.Hash.EncodeId(role.ID)
		var perms []PermissionResource
		for _, p := range role.Permissions {
			encodedPermID, _ := lib.Hash.EncodeId(p.ID)
			perms = append(perms, PermissionResource{
				ID:    encodedPermID,
				Name:  p.Name,
				Notes: p.Notes,
			})
		}
		result = append(result, RoleListResource{
			ID:          encodedID,
			Name:        role.Name,
			Notes:       role.Notes,
			Permissions: perms,
		})
	}
	return result
}
