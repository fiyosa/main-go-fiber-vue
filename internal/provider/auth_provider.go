package provider

import (
	"go-fiber-svelte/internal/db"
	"go-fiber-svelte/internal/db/models"

	"gorm.io/gorm"
)

type AuthProvider struct {
	db *gorm.DB
}

func NewAuthProvider() *AuthProvider {
	return &AuthProvider{db: db.RUN}
}

func (p *AuthProvider) CheckPermission(userId int, permissionName string) bool {
	var count int64
	p.db.Table("user_has_roles").
		Joins("JOIN role_has_permissions ON role_has_permissions.role_id = user_has_roles.role_id").
		Joins("JOIN permissions ON permissions.id = role_has_permissions.permission_id").
		Where("user_has_roles.user_id = ? AND permissions.name = ? AND permissions.deleted_at IS NULL", userId, permissionName).
		Count(&count)
	return count > 0
}

func (p *AuthProvider) GetUserRoles(userId int) ([]models.Role, error) {
	var roles []models.Role
	err := p.db.Table("roles").
		Joins("JOIN user_has_roles ON user_has_roles.role_id = roles.id").
		Where("user_has_roles.user_id = ? AND roles.deleted_at IS NULL", userId).
		Find(&roles).Error
	return roles, err
}
