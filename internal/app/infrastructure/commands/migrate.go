package commands

import (
	"github.com/spf13/cobra"

	"go-fiber-vue/internal/config"
	"go-fiber-vue/internal/platform/database"
)

// Migrate returns a cobra.Command that wraps all migration subcommands.
//
// Every subcommand resolves the configuration with config.LoadForConsole:
// migrations only need DatabaseDSN, so they must not fail on runtime secrets
// (HASHID_SALT, JWT_SECRET, CRON_SECRET_KEY, MCP_TOKEN) that the application
// itself validates at startup.
func Migrate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migrations",
	}
	cmd.AddCommand(migrateUp(), migrateStatus(), migrateDown())
	return cmd
}

func migrateUp() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadForConsole()
			if err != nil {
				return err
			}
			return database.MigrateUp(cfg.DatabaseDSN())
		},
	}
}

func migrateStatus() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current migration version",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadForConsole()
			if err != nil {
				return err
			}
			ver, dirty, err := database.MigrateStatus(cfg.DatabaseDSN())
			if err != nil {
				return err
			}
			cmd.Printf("version=%d dirty=%t\n", ver, dirty)
			return nil
		},
	}
}

func migrateDown() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Roll back the most recent migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadForConsole()
			if err != nil {
				return err
			}
			return database.MigrateDown(cfg.DatabaseDSN(), 1)
		},
	}
}
