package database

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-fiber-vue/internal/app/domain/core/entity"
)

// SeedRBAC writes the default admin role and every permission registered by the
// route table, then attaches all permissions to the admin role. It is idempotent:
// running it twice leaves exactly the same rows. Permission rows are the only
// source of access; there is no admin bypass.
//
// User accounts are intentionally not seeded. Create them manually and attach a
// role through user_has_roles when an account is needed.
func SeedRBAC(ctx context.Context, db *gorm.DB, permissions []string) error {
	if db == nil {
		return fmt.Errorf("database is not configured")
	}

	names := normalizePermissions(permissions)
	if len(names) == 0 {
		return fmt.Errorf("no permission registered; build the route table before seeding")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, name := range names {
			if err := tx.Exec(
				`INSERT INTO permissions (name) VALUES (?) ON CONFLICT (name) DO NOTHING`,
				name,
			).Error; err != nil {
				return fmt.Errorf("seed permission %q: %w", name, err)
			}
		}

		if err := tx.Exec(
			`INSERT INTO roles (name) VALUES (?) ON CONFLICT (name) DO NOTHING`,
			entity.AdminRole,
		).Error; err != nil {
			return fmt.Errorf("seed role %q: %w", entity.AdminRole, err)
		}

		// Grant every registered permission to the admin role.
		var role entity.Role
		if err := tx.Where("name = ?", entity.AdminRole).First(&role).Error; err != nil {
			return fmt.Errorf("find admin role: %w", err)
		}

		var perms []entity.Permission
		if err := tx.Find(&perms).Error; err != nil {
			return fmt.Errorf("find permissions: %w", err)
		}

		for _, p := range perms {
			rp := entity.RoleHasPermission{RoleID: role.ID, PermissionID: p.ID}
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&rp).Error; err != nil {
				return fmt.Errorf("seed role permission: %w", err)
			}
		}
		return nil
	})
}

// normalizePermissions trims, deduplicates, and sorts the names so the seed is
// stable regardless of route registration order.
func normalizePermissions(permissions []string) []string {
	seen := make(map[string]struct{}, len(permissions))
	names := make([]string, 0, len(permissions))
	for _, name := range permissions {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
