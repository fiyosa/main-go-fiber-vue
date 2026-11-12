package entity

import "time"

type Permission struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;size:125;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}
