package entity

import "time"

// AdminRole is the default role seeded with every permission.
// Access is determined purely by permission rows; there is no hardcoded bypass.
const AdminRole = "admin"

type Role struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;size:125;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Role) TableName() string {
	return "roles"
}
