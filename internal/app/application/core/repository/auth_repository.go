package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

type AuthRepository interface {
	CreateSession(ctx context.Context, auth *entity.Auth, lastLoginAt time.Time) error
}

type gormAuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	if db == nil {
		return nil
	}
	return &gormAuthRepository{db: db}
}

func (r *gormAuthRepository) CreateSession(ctx context.Context, auth *entity.Auth, lastLoginAt time.Time) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(auth).Error; err != nil {
			return err
		}
		return tx.Model(&entity.User{}).
			Where("id = ?", auth.UserID).
			Update("last_login_at", lastLoginAt).Error
	})
	return exceptions.FromDB(err)
}
