package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"uniqueIndex;size:255;not null" json:"name"`
	Notes     string         `gorm:"type:text" json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Permissions []Permission `gorm:"many2many:role_has_permissions" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}
