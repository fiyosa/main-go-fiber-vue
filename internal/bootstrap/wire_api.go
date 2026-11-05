package bootstrap

import (
	"database/sql"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/application/core/repository/auth"
	"go-fiber-vue/internal/app/application/core/repository/user"
	"go-fiber-vue/internal/app/presentation/core/controllers"
	v1 "go-fiber-vue/internal/routes/v1"
)

// wireAPI assembles every controller exposed by the HTTP API and returns the
// route layer contract. A nil database, nil sql handle, and nil token service
// are all valid, so console commands can build the route table without a live
// database.
func wireAPI(db *gorm.DB, sqlDB *sql.DB, tokens auth.TokenService) v1.Dependencies {
	loginRepo := auth.NewLoginRepository(
		user.NewRepository(db),
		auth.NewRepository(db),
		tokens,
	)

	var pinger controllers.Pinger
	if sqlDB != nil {
		pinger = sqlDB
	}

	return v1.Dependencies{
		AuthController:   controllers.NewAuthController(loginRepo),
		HealthController: controllers.NewHealthController(pinger),
	}
}
