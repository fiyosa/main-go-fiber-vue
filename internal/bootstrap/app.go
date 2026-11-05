package bootstrap

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"gorm.io/gorm"

	"go-fiber-vue/internal/app/application/core/repository"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	"go-fiber-vue/internal/app/presentation/core/controllers"
	"go-fiber-vue/internal/config"
	"go-fiber-vue/internal/platform/auth"
	"go-fiber-vue/internal/platform/database"
	platformmiddleware "go-fiber-vue/internal/platform/http/middleware"
	"go-fiber-vue/internal/platform/lang"
	"go-fiber-vue/internal/routes"
	v1 "go-fiber-vue/internal/routes/v1"
)

type Application struct {
	Fiber      *fiber.App
	DB         *gorm.DB
	Tokens     *auth.Service
	Translator *lang.Translator
}

// Router builds a Fiber app with every route registered but no runtime
// dependencies (no database, no translator). It exists so console commands
// such as `routes` can introspect the routing table from any working
// directory and without a live database.
func Router(cfg *config.Config) *fiber.App {
	loginRepo := repository.NewAuthLoginRepository(nil, nil, nil, nil)
	authController := controllers.NewAuthController(loginRepo)
	healthController := controllers.NewHealthController(nil)

	app := fiber.New(fiber.Config{
		AppName:       "app-api",
		CaseSensitive: true,
		ErrorHandler:  exceptions.Handler(cfg.IsDevelopment()),
	})
	app.Use(recover.New())
	app.Use(platformmiddleware.Maintenance(cfg.AppMaintenance))

	routes.RegisterAPI(app, v1.Dependencies{
		AuthController:   authController,
		HealthController: healthController,
	})
	routes.RegisterWeb(app, cfg.PublicDir)

	return app
}

func New(cfg *config.Config) (*Application, error) {
	translator, err := newTranslator(cfg)
	if err != nil {
		return nil, fmt.Errorf("init translator: %w", err)
	}

	db, err := database.Open(cfg.DatabaseDSN(), cfg.DBLogSQL)
	if err != nil {
		slog.Warn("database unavailable; starting without it", "error", err)
		db = nil
	}

	var sqlDB *sql.DB
	if db != nil {
		sqlDB, err = db.DB()
		if err != nil {
			slog.Warn("database handle unavailable", "error", err)
			sqlDB = nil
		}
	}

	tokens, err := auth.NewService(auth.Config{
		Secret:     cfg.JWTSecret,
		AccessTTL:  cfg.JWTAccessTTL,
		RefreshTTL: cfg.JWTRefreshTTL,
	})
	if err != nil {
		return nil, fmt.Errorf("init token service: %w", err)
	}

	var users repository.UserRepository
	var auths repository.AuthRepository
	if db != nil {
		users = repository.NewUserRepository(db)
		auths = repository.NewAuthRepository(db)
	}

	loginRepo := repository.NewAuthLoginRepository(db, users, auths, tokens)
	authController := controllers.NewAuthController(loginRepo)

	var pinger controllers.Pinger
	if sqlDB != nil {
		pinger = sqlDB
	}
	healthController := controllers.NewHealthController(pinger)

	app := fiber.New(fiber.Config{
		AppName:       "app-api",
		CaseSensitive: true,
		ErrorHandler:  exceptions.Handler(cfg.IsDevelopment()),
	})
	app.Use(recover.New())
	app.Use(platformmiddleware.Locale(translator))
	app.Use(platformmiddleware.Maintenance(cfg.AppMaintenance))

	routes.RegisterAPI(app, v1.Dependencies{
		AuthController:   authController,
		HealthController: healthController,
	})
	routes.RegisterWeb(app, cfg.PublicDir)

	return &Application{
		Fiber:      app,
		DB:         db,
		Tokens:     tokens,
		Translator: translator,
	}, nil
}

func newTranslator(cfg *config.Config) (*lang.Translator, error) {
	dir, err := localeDir()
	if err != nil {
		return nil, err
	}
	return lang.Load(dir, cfg.AppLocale)
}

func localeDir() (string, error) {
	candidates := []string{filepath.Join("internal", "lang")}
	if executable, err := os.Executable(); err == nil {
		base := filepath.Dir(executable)
		candidates = append(candidates,
			filepath.Join(base, "internal", "lang"),
			filepath.Join(base, "lang"),
		)
	}
	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("lang directory not found (looked in %v)", candidates)
}
