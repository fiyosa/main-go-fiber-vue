package entity

import "time"

// BypassRole is the only role that skips every permission check. It is seeded
// by the console and must not be created through the API.
const BypassRole = "admin"

type Role struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;size:125;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Role) TableName() string {
	return "roles"
}
