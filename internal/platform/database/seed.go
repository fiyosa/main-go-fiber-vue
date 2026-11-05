package database

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/platform/auth"
)

const (
	minSeedPasswordLength = 12

	// DefaultAdmin* define the account created by `console db seed`. Change the
	// password after seeding if the environment is shared.
	DefaultAdminEmail    = "admin@example.com"
	DefaultAdminPassword = "Admin#Local2026"
)

func SeedAdmin(ctx context.Context, db *gorm.DB, email, password string) error {
	if db == nil {
		return fmt.Errorf("database is not configured")
	}
	if strings.TrimSpace(password) == "" || len(password) < minSeedPasswordLength {
		return fmt.Errorf("seed admin password must be at least %d characters", minSeedPasswordLength)
	}

	var count int64
	if err := db.WithContext(ctx).Model(&entity.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash seed password: %w", err)
	}

	user := &entity.User{
		Email:    email,
		Password: hash,
	}
	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("create seed admin: %w", err)
	}
	return nil
}
