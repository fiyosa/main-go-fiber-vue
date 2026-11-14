package entity

import "time"

type Auth struct {
	ID             uint64     `gorm:"column:id;primaryKey"`
	UserID         uint64     `gorm:"column:user_id;not null"`
	User           *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Token          string     `gorm:"column:token;type:text;not null"`
	RefreshToken   string     `gorm:"column:refresh_token;size:64;not null"`
	ImpersonatorID *uint64    `gorm:"column:impersonator_id"`
	Impersonator   *User      `gorm:"foreignKey:ImpersonatorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Meta           *string    `gorm:"column:meta;type:jsonb"`
	RevokedAt      *time.Time `gorm:"column:revoked_at"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
}

func (Auth) TableName() string {
	return "auths"
}
