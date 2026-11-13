package role

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// Permission is the read view of a permission row. IDs stay numeric here and
// are hashed by the presentation layer.
type Permission struct {
	ID   uint64
	Name string
}

// Role is a role with its attached permissions, ordered by permission name.
type Role struct {
	ID          uint64
	Name        string
	Permissions []Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository reads authorization data. Every method is a single query so the
// middleware stays free of SQL.
type Repository interface {
	// List returns every role ordered by name, each with its permissions.
	List(ctx context.Context) ([]Role, error)
	// HasPermission reports whether any of the user's roles grants permission.
	HasPermission(ctx context.Context, userID uint64, permission string) (bool, error)
	// NamesByUserID returns every role name held by the user, ordered by name.
	NamesByUserID(ctx context.Context, userID uint64) ([]string, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository returns a GORM-backed role repository, or nil when db is nil.
func NewRepository(db *gorm.DB) Repository {
	if db == nil {
		return nil
	}
	return &gormRepository{db: db}
}

func (r *gormRepository) List(ctx context.Context) ([]Role, error) {
	var rows []struct {
		ID        uint64
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	err := r.db.WithContext(ctx).
		Table("roles").
		Select("id, name, created_at, updated_at").
		Order("name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, exceptions.FromDB(err)
	}
	if len(rows) == 0 {
		return []Role{}, nil
	}

	var permissionRows []struct {
		RoleID uint64
		ID     uint64
		Name   string
	}
	err = r.db.WithContext(ctx).
		Table("role_has_permissions AS rp").
		Select("rp.role_id, p.id, p.name").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Order("p.name ASC").
		Scan(&permissionRows).Error
	if err != nil {
		return nil, exceptions.FromDB(err)
	}

	byRole := make(map[uint64][]Permission, len(rows))
	for _, row := range permissionRows {
		byRole[row.RoleID] = append(byRole[row.RoleID], Permission{ID: row.ID, Name: row.Name})
	}

	roles := make([]Role, 0, len(rows))
	for _, row := range rows {
		permissions := byRole[row.ID]
		if permissions == nil {
			permissions = []Permission{}
		}
		roles = append(roles, Role{
			ID:          row.ID,
			Name:        row.Name,
			Permissions: permissions,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}
	return roles, nil
}

const hasPermissionQuery = `
SELECT EXISTS (
    SELECT 1
    FROM user_has_roles ur
    JOIN role_has_permissions rp ON rp.role_id = ur.role_id
    JOIN permissions p ON p.id = rp.permission_id
    WHERE ur.user_id = ? AND p.name = ?
)`

const namesByUserIDQuery = `
SELECT r.name
FROM user_has_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = ?
ORDER BY r.name ASC`

func (r *gormRepository) HasPermission(ctx context.Context, userID uint64, permission string) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).
		Raw(hasPermissionQuery, userID, permission).
		Scan(&exists).Error
	if err != nil {
		return false, exceptions.FromDB(err)
	}
	return exists, nil
}

func (r *gormRepository) NamesByUserID(ctx context.Context, userID uint64) ([]string, error) {
	var rows []struct{ Name string }
	err := r.db.WithContext(ctx).
		Raw(namesByUserIDQuery, userID).
		Scan(&rows).Error
	if err != nil {
		return nil, exceptions.FromDB(err)
	}
	names := make([]string, 0, len(rows))
	for _, row := range rows {
		names = append(names, row.Name)
	}
	return names, nil
}
