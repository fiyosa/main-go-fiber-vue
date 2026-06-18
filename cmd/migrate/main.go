package main

import (
	"fmt"
	"os"

	"go-fiber-svelte/internal/config"
	"go-fiber-svelte/internal/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	config.InitConfigApp()

	db, err := gorm.Open(postgres.Open(config.DB_URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("Failed to connect database:", err)
		os.Exit(1)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.UserDetail{},
		&models.Auth{},
		&models.Role{},
		&models.Permission{},
		&models.RoleHasPermission{},
		&models.UserHasRole{},
	)
	if err != nil {
		fmt.Println("Migration failed:", err)
		os.Exit(1)
	}

	fmt.Println("Migration completed successfully")
}
