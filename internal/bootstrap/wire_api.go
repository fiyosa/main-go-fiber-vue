package bootstrap

import (
	"database/sql"
	"log/slog"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/application/core/authorization"
	"go-fiber-vue/internal/app/application/core/repository/auth"
	"go-fiber-vue/internal/app/application/core/repository/role"
	"go-fiber-vue/internal/app/application/core/repository/user"
	"go-fiber-vue/internal/app/presentation/core/controllers"
	"go-fiber-vue/internal/app/presentation/core/middleware"
	"go-fiber-vue/internal/config"
	"go-fiber-vue/internal/platform/hashid"
	v1 "go-fiber-vue/internal/routes/v1"
)

// wireAPI assembles every controller exposed by the HTTP API and returns the
// route layer contract. A nil database, nil sql handle, and nil token service
// are all valid, so console commands can build the route table without a live
// database.
func wireAPI(cfg *config.Config, db *gorm.DB, sqlDB *sql.DB, tokens auth.TokenService, validator middleware.TokenValidator) v1.Dependencies {
	loginRepo := auth.NewLoginRepository(
		user.NewRepository(db),
		auth.NewRepository(db),
		tokens,
	)

	var pinger controllers.Pinger
	if sqlDB != nil {
		pinger = sqlDB
	}

	roleRepo := role.NewRepository(db)

	return v1.Dependencies{
		AuthController:   controllers.NewAuthController(loginRepo, !cfg.IsDevelopment()),
		HealthController: controllers.NewHealthController(pinger),
		RoleController:   controllers.NewRoleController(roleRepo, newHashidEncoder(cfg)),
		Tokens:           validator,
		Authorizer:       authorization.New(roleRepo),
	}
}

// newHashidEncoder builds the id encoder. A missing or invalid salt disables
// hashing instead of crashing: the route table stays inspectable and the
// affected handlers answer 503.
func newHashidEncoder(cfg *config.Config) *hashid.Encoder {
	encoder, err := hashid.New(cfg.HashidSaltOrDefault(), cfg.HashidMinLength)
	if err != nil {
		slog.Warn("hashid disabled", "error", err)
		return nil
	}
	return encoder
}
