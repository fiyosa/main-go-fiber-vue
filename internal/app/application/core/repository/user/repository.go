package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// Repository reads user records.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository returns a GORM-backed user repository, or nil when db is nil.
func NewRepository(db *gorm.DB) Repository {
	if db == nil {
		return nil
	}
	return &gormRepository{db: db}
}

func (r *gormRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var found entity.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&found).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, exceptions.FromDB(err)
	}
	return &found, nil
}
