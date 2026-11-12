package entity

// UserHasRole links a user to a role. The pair is the primary key; the row has
// no identity of its own.
type UserHasRole struct {
	UserID uint64 `gorm:"column:user_id;primaryKey"`
	RoleID uint64 `gorm:"column:role_id;primaryKey"`
}

func (UserHasRole) TableName() string {
	return "user_has_roles"
}
