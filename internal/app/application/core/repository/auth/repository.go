package auth

import (
	"context"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// Repository persists auth sessions.
type Repository interface {
	CreateSession(ctx context.Context, session *entity.Auth) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository returns a GORM-backed auth repository, or nil when db is nil.
func NewRepository(db *gorm.DB) Repository {
	if db == nil {
		return nil
	}
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateSession(ctx context.Context, session *entity.Auth) error {
	err := r.db.WithContext(ctx).Create(session).Error
	return exceptions.FromDB(err)
}
