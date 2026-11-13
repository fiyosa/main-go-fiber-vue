package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/bootstrap"
	"go-fiber-vue/internal/config"
	"go-fiber-vue/internal/platform/database"
	"go-fiber-vue/internal/platform/permission"
)

// Seed returns a cobra.Command that seeds roles and permissions.
func Seed() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed roles and permissions",
		Long: "Seed the default admin role and every permission registered by the API\n" +
			"route table, then attach all permissions to the admin role. Safe to run\n" +
			"repeatedly; no user account is created.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			// Migrate up first to ensure tables exist
			if err := database.MigrateUp(cfg.DatabaseDSN()); err != nil {
				return fmt.Errorf("migrate up: %w", err)
			}

			// Building the route table registers every permission name, so the
			// seed list always matches the code.
			bootstrap.Router(cfg)
			names := permission.Names()
			if len(names) == 0 {
				return fmt.Errorf("no permission registered; check the route table")
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

			if err := database.SeedRBAC(context.Background(), db, names); err != nil {
				return err
			}

			fmt.Printf("seed complete: role=%s permissions=%d\n", entity.AdminRole, len(names))
			for _, name := range names {
				fmt.Printf("  - %s\n", name)
			}
			return nil
		},
	}
}
