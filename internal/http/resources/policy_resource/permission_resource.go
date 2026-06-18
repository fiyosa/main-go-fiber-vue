package policy_resource

import (
	"go-fiber-svelte/internal/db/models"
	"go-fiber-svelte/internal/lib"
)

type PermissionResource struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Notes string `json:"notes"`
}

func PermissionToResource(permissions []models.Permission) []PermissionResource {
	var result []PermissionResource
	for _, p := range permissions {
		encodedID, _ := lib.Hash.EncodeId(p.ID)
		result = append(result, PermissionResource{
			ID:    encodedID,
			Name:  p.Name,
			Notes: p.Notes,
		})
	}
	return result
}
