package models

type RoleHasPermission struct {
	RoleID       int `gorm:"primaryKey" json:"role_id"`
	PermissionID int `gorm:"primaryKey" json:"permission_id"`
}

func (RoleHasPermission) TableName() string {
	return "role_has_permissions"
}
