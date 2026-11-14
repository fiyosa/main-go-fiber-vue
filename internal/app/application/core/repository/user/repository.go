package user

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// Repository reads user records.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	// FindByID returns a user by primary key, or nil when it does not exist.
	FindByID(ctx context.Context, id uint64) (*entity.User, error)
	// UpdatePassword replaces the password hash of an existing user.
	UpdatePassword(ctx context.Context, id uint64, passwordHash string) error
	// Create inserts a new user and, when roleName is not empty and a role with
	// that name exists, links the user to the role. Both writes run in one
	// transaction so a user can never be created without its default role.
	Create(ctx context.Context, account *entity.User, roleName string) error
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

func (r *gormRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	var found entity.User
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&found).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, exceptions.FromDB(err)
	}
	return &found, nil
}

// UpdatePassword implements Repository. A missing row is not treated as an
// error here; the change-password use case resolves the user first.
func (r *gormRepository) UpdatePassword(ctx context.Context, id uint64, passwordHash string) error {
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"password":   passwordHash,
			"updated_at": time.Now().UTC(),
		}).Error
	return exceptions.FromDB(err)
}

// Create implements Repository. The user insert and the default role assignment
// share one transaction; an existing role link is left untouched on conflict.
func (r *gormRepository) Create(ctx context.Context, account *entity.User, roleName string) error {
	return exceptions.FromDB(r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(account).Error; err != nil {
			return err
		}
		if roleName == "" {
			return nil
		}
		const linkUserHasRole = `
INSERT INTO user_has_roles (user_id, role_id)
SELECT ?, id FROM roles WHERE name = ?
ON CONFLICT DO NOTHING`
		return tx.Exec(linkUserHasRole, account.ID, roleName).Error
	}))
}
