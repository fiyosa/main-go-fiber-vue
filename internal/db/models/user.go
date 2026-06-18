package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Username  string         `gorm:"uniqueIndex;size:255;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UserDetails *UserDetail   `gorm:"foreignKey:UserID" json:"user_details,omitempty"`
	Auths       []Auth        `gorm:"foreignKey:UserID" json:"auths,omitempty"`
	Roles       []Role        `gorm:"many2many:user_has_roles" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "users"
}
