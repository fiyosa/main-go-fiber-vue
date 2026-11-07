package user

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

const (
	defaultPageSize = 25
	maxPageSize     = 50
)

// Repository reads user records.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	// SearchUsers returns users whose email contains the given fragment,
	// together with the total number of matches before pagination. Page is
	// 1-based and pageSize is capped at maxPageSize.
	SearchUsers(ctx context.Context, email string, page, pageSize int) ([]entity.User, int64, error)
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

// SearchUsers implements Repository. It returns an empty result when the
// repository has no database so callers can degrade gracefully.
func (r *gormRepository) SearchUsers(ctx context.Context, email string, page, pageSize int) ([]entity.User, int64, error) {
	if r == nil || r.db == nil {
		return []entity.User{}, 0, nil
	}

	page, pageSize = normalizePaging(page, pageSize)
	email = strings.TrimSpace(email)

	// Build the same filter twice: GORM's Count mutates the statement, so
	// reusing one query builder would leak the aggregate select into the page.
	filter := func(db *gorm.DB) *gorm.DB {
		query := db.WithContext(ctx).Model(&entity.User{})
		if email != "" {
			query = query.Where("email ILIKE ?", "%"+email+"%")
		}
		return query
	}

	var total int64
	if err := filter(r.db).Count(&total).Error; err != nil {
		return nil, 0, exceptions.FromDB(err)
	}

	var users []entity.User
	err := filter(r.db).
		Order("id ASC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&users).Error
	if err != nil {
		return nil, 0, exceptions.FromDB(err)
	}
	return users, total, nil
}

// normalizePaging clamps caller-supplied paging into the supported range.
func normalizePaging(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	switch {
	case pageSize < 1:
		pageSize = defaultPageSize
	case pageSize > maxPageSize:
		pageSize = maxPageSize
	}
	return page, pageSize
}
