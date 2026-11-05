package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"go-fiber-vue/internal/config"
	"go-fiber-vue/internal/platform/database"
)

// Seed returns a cobra.Command that seeds the initial admin user.
func Seed() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed the initial admin user",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
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

			if err := database.SeedAdmin(context.Background(), db,
				database.DefaultAdminEmail,
				database.DefaultAdminPassword,
			); err != nil {
				return err
			}
			fmt.Printf("seed complete: email=%s\n", database.DefaultAdminEmail)
			return nil
		},
	}
}
