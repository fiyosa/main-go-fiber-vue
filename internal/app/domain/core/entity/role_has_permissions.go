package entity

// RoleHasPermission links a role to a permission. The pair is the primary key.
type RoleHasPermission struct {
	RoleID       uint64 `gorm:"column:role_id;primaryKey"`
	PermissionID uint64 `gorm:"column:permission_id;primaryKey"`
}

func (RoleHasPermission) TableName() string {
	return "role_has_permissions"
}
