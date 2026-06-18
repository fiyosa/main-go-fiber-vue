package seed

import (
	"go-fiber-svelte/internal/db"
	"go-fiber-svelte/internal/db/models"
	"go-fiber-svelte/internal/lib"

	"gorm.io/gorm"
)

func Run() {
	database := db.RUN

	createRoles(database)
	createPermissions(database)
	createAdminUser(database)
}

func createRoles(database *gorm.DB) {
	roles := []models.Role{
		{Name: "Super Admin", Notes: "Full access"},
		{Name: "Admin", Notes: "Administrative access"},
		{Name: "User", Notes: "Regular user"},
	}
	for _, role := range roles {
		database.FirstOrCreate(&role, models.Role{Name: role.Name})
	}
}

func createPermissions(database *gorm.DB) {
	permissions := []models.Permission{
		{Name: "user.view", Notes: "View users"},
		{Name: "user.create", Notes: "Create users"},
		{Name: "user.update", Notes: "Update users"},
		{Name: "user.delete", Notes: "Delete users"},
		{Name: "role.view", Notes: "View roles"},
		{Name: "role.create", Notes: "Create roles"},
		{Name: "role.update", Notes: "Update roles"},
		{Name: "role.delete", Notes: "Delete roles"},
		{Name: "permission.view", Notes: "View permissions"},
		{Name: "permission.create", Notes: "Create permissions"},
		{Name: "permission.delete", Notes: "Delete permissions"},
	}
	for _, perm := range permissions {
		database.FirstOrCreate(&perm, models.Permission{Name: perm.Name})
	}
}

func createAdminUser(database *gorm.DB) {
	password, _ := lib.Hash.Create("admin123")
	var existing models.User
	result := database.Where("email = ?", "admin@example.com").First(&existing)
	if result.Error != nil {
		admin := models.User{
			Email:    "admin@example.com",
			Username: "admin",
			Password: password,
		}
		database.Create(&admin)
		var superAdminRole models.Role
		database.Where("name = ?", "Super Admin").First(&superAdminRole)
		database.Create(&models.UserHasRole{
			UserID: admin.ID,
			RoleID: superAdminRole.ID,
		})
	}
}
