package auth_resource

import (
	"go-fiber-svelte/internal/db/models"
	"go-fiber-svelte/internal/lib"
)

type UserResource struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	Username  string         `json:"username"`
	Roles     []RoleResource `json:"roles,omitempty"`
	CreatedAt string         `json:"created_at"`
}

type RoleResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func UserToResource(user models.User) UserResource {
	encodedID, _ := lib.Hash.EncodeId(user.ID)
	res := UserResource{
		ID:        encodedID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	for _, role := range user.Roles {
		encodedRoleID, _ := lib.Hash.EncodeId(role.ID)
		res.Roles = append(res.Roles, RoleResource{
			ID:   encodedRoleID,
			Name: role.Name,
		})
	}
	return res
}
