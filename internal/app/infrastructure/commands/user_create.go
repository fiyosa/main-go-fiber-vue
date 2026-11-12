package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/config"
	"go-fiber-vue/internal/platform/auth"
	"go-fiber-vue/internal/platform/database"
)

// UserCreate returns a cobra.Command that creates or updates a user and assigns a role.
func UserCreate() *cobra.Command {
	var email string
	var password string
	var roleName string

	cmd := &cobra.Command{
		Use:   "user:create",
		Short: "Create or update a user and assign a role",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			// Migrate up first to ensure tables exist
			if err := database.MigrateUp(cfg.DatabaseDSN()); err != nil {
				return fmt.Errorf("migrate up: %w", err)
			}

			db, err := database.Open(cfg.DatabaseDSN(), cfg.DBLogSQL)
			if err != nil {
				return err
			}
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			defer sqlDB.Close()

			ctx := context.Background()

			hashedPassword, err := auth.HashPassword(password)
			if err != nil {
				return fmt.Errorf("hash password: %w", err)
			}

			var user entity.User
			err = db.WithContext(ctx).Where("email = ?", email).First(&user).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				user = entity.User{
					Email:    email,
					Password: hashedPassword,
				}
				if err := db.WithContext(ctx).Create(&user).Error; err != nil {
					return fmt.Errorf("create user: %w", err)
				}
				fmt.Printf("Created user: %s (id=%d)\n", email, user.ID)
			} else if err != nil {
				return fmt.Errorf("find user: %w", err)
			} else {
				user.Password = hashedPassword
				if err := db.WithContext(ctx).Save(&user).Error; err != nil {
					return fmt.Errorf("update user: %w", err)
				}
				fmt.Printf("Updated user: %s (id=%d)\n", email, user.ID)
			}

			if roleName != "" {
				var role entity.Role
				if err := db.WithContext(ctx).Where("name = ?", roleName).First(&role).Error; err != nil {
					return fmt.Errorf("find role %q: %w (did you run 'console db seed'?)", roleName, err)
				}

				userRole := entity.UserHasRole{UserID: user.ID, RoleID: role.ID}
				err := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&userRole).Error
				if err != nil {
					return fmt.Errorf("assign role: %w", err)
				}
				fmt.Printf("Assigned role %q to user %s\n", roleName, email)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&email, "email", "fiyosa97@gmail.com", "User email")
	cmd.Flags().StringVar(&password, "password", "password123", "User password")
	cmd.Flags().StringVar(&roleName, "role", "admin", "Role name to assign")

	return cmd
}
