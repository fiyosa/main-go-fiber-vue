package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

type UserRepository interface {
	FindByUsernameOrEmail(ctx context.Context, value string) (*entity.User, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	if db == nil {
		return nil
	}
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) FindByUsernameOrEmail(ctx context.Context, value string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("username = ? OR email = ?", value, value).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, exceptions.FromDB(err)
	}
	return &user, nil
}
