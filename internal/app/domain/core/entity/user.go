package entity

import "time"

type User struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	Email     string    `gorm:"column:email;size:255;not null;uniqueIndex"`
	Password  string    `gorm:"column:password;size:255;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}
