package entity

// RoleHasPermission links a role to a permission. The pair is the primary key.
// Both sides cascade: a grant cannot outlive its role or its permission.
type RoleHasPermission struct {
	RoleID       uint64      `gorm:"column:role_id;primaryKey"`
	PermissionID uint64      `gorm:"column:permission_id;primaryKey"`
	Role         *Role       `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permission   *Permission `gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RoleHasPermission) TableName() string {
	return "role_has_permissions"
}
