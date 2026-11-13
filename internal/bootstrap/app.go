package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"gorm.io/gorm"

	cronrepo "go-fiber-vue/internal/app/application/core/repository/cron"
	storagerepo "go-fiber-vue/internal/app/application/core/repository/storage"
	cronsvc "go-fiber-vue/internal/app/application/core/service/cron"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	"go-fiber-vue/internal/config"
	"go-fiber-vue/internal/platform/auth"
	"go-fiber-vue/internal/platform/database"
	"go-fiber-vue/internal/platform/lang"
	platformmiddleware "go-fiber-vue/internal/platform/middleware"
	"go-fiber-vue/internal/routes"
)

type Application struct {
	Fiber      *fiber.App
	DB         *gorm.DB
	Tokens     *auth.Service
	Translator *lang.Translator
	MCPHandler http.Handler
	// Storage is nil when STORAGE_* is not fully configured.
	Storage storagerepo.Repository
	// CronService is nil when cron is disabled.
	CronService *cronsvc.Service
}

// Router builds a Fiber app with every route registered but no runtime
// dependencies (no database, no translator). It exists so console commands
// such as `routes` can introspect the routing table from any working
// directory and without a live database.
func Router(cfg *config.Config) *fiber.App {
	deps := wireAPI(cfg, nil, nil, nil, nil)

	app := fiber.New(fiber.Config{
		AppName:       "app-api",
		CaseSensitive: true,
		ErrorHandler:  exceptions.Handler(cfg.IsDevelopment()),
	})
	app.Use(recover.New())
	app.Use(platformmiddleware.CORS(cfg.AppCORSOrigins))
	app.Use(platformmiddleware.Maintenance(cfg.AppMaintenance))

	routes.RegisterMCP(app, cfg.MCPPath, newMCPHandler(cfg, nil))
	routes.RegisterAPI(app, deps)
	routes.RegisterDocs(app, cfg.PublicDir, cfg.BasicAuthUsername, cfg.BasicAuthPassword)
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

	deps := wireAPI(cfg, db, sqlDB, tokens, tokens)

	var cronSvc *cronsvc.Service
	if cfg.CronEnabled && db != nil {
		cronRepo := cronrepo.NewRepository(db)
		cronSvc = cronsvc.NewService(cronsvc.Options{
			DB:        db,
			Repo:      cronRepo,
			SecretKey: cfg.CronSecretKey,
		})
		_ = cronSvc.Start(context.Background())
	}

	mcpHandler := newMCPHandler(cfg, sqlDB)

	app := fiber.New(fiber.Config{
		AppName:       "app-api",
		CaseSensitive: true,
		ErrorHandler:  exceptions.Handler(cfg.IsDevelopment()),
	})
	app.Use(recover.New())
	app.Use(platformmiddleware.CORS(cfg.AppCORSOrigins))
	app.Use(platformmiddleware.Locale(translator))
	app.Use(platformmiddleware.Maintenance(cfg.AppMaintenance))

	routes.RegisterMCP(app, cfg.MCPPath, mcpHandler)
	routes.RegisterAPI(app, deps)
	routes.RegisterDocs(app, cfg.PublicDir, cfg.BasicAuthUsername, cfg.BasicAuthPassword)
	routes.RegisterWeb(app, cfg.PublicDir)

	return &Application{
		Fiber:      app,
		DB:         db,
		Tokens:     tokens,
		Translator: translator,
		MCPHandler: mcpHandler,
		Storage:    newStorageRepository(cfg, db),
		CronService: cronSvc,
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
	candidates := []string{filepath.Join("lang")}
	if executable, err := os.Executable(); err == nil {
		base := filepath.Dir(executable)
		candidates = append(candidates, filepath.Join(base, "lang"))
	}
	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("lang directory not found (looked in %v)", candidates)
}
