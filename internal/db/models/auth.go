package models

import "time"

type Auth struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	Token     string    `gorm:"size:500;not null" json:"token"`
	Revoke    bool      `gorm:"default:false" json:"revoke"`
	IP        string    `gorm:"size:45" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Auth) TableName() string {
	return "auths"
}
