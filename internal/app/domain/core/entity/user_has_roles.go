package entity

// UserHasRole links a user to a role. The pair is the primary key; the row has
// no identity of its own. Both sides cascade because the link is meaningless
// once either endpoint is gone.
type UserHasRole struct {
	UserID uint64 `gorm:"column:user_id;primaryKey"`
	RoleID uint64 `gorm:"column:role_id;primaryKey"`
	User   *User  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role   *Role  `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserHasRole) TableName() string {
	return "user_has_roles"
}
