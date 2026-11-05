package entity

import "time"

type User struct {
	ID           uint64     `gorm:"column:id;primaryKey"`
	Name         string     `gorm:"column:name;size:255;not null"`
	Username     string     `gorm:"column:username;size:100;not null;uniqueIndex"`
	Email        string     `gorm:"column:email;size:255;not null;uniqueIndex"`
	Password     string     `gorm:"column:password;size:255;not null"`
	Role         *string    `gorm:"column:role;size:125"`
	DivisionID   *uint64    `gorm:"column:division_id"`
	TokenVersion int        `gorm:"column:token_version;not null;default:0"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}
