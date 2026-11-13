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
