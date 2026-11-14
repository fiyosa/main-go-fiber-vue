package auth

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// Repository persists auth sessions.
type Repository interface {
	CreateSession(ctx context.Context, session *entity.Auth) error
	RevokeSession(ctx context.Context, token string) error
	// RevokeAllForUser revokes every active session of a user. Used after a
	// credential change so the old sessions can no longer be used.
	RevokeAllForUser(ctx context.Context, userID uint64) error
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

// RevokeSession marks the session matching the given access token as revoked.
// A token that does not match any session is a no-op, so logout stays idempotent.
func (r *gormRepository) RevokeSession(ctx context.Context, token string) error {
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).
		Model(&entity.Auth{}).
		Where("token = ? AND revoked_at IS NULL", token).
		Update("revoked_at", now).Error
	return exceptions.FromDB(err)
}

// RevokeAllForUser implements Repository.
func (r *gormRepository) RevokeAllForUser(ctx context.Context, userID uint64) error {
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).
		Model(&entity.Auth{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
	return exceptions.FromDB(err)
}
