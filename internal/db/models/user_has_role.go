package models

type UserHasRole struct {
	UserID int `gorm:"primaryKey" json:"user_id"`
	RoleID int `gorm:"primaryKey" json:"role_id"`
}

func (UserHasRole) TableName() string {
	return "user_has_roles"
}
