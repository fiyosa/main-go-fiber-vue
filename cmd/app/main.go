package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-fiber-vue/internal/bootstrap"
	"go-fiber-vue/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}
	bootstrap.SetupLogger(cfg.AppEnv)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.New(cfg)
	if err != nil {
		slog.Error("bootstrap application", "error", err)
		os.Exit(1)
	}

	defer func() {
		if app.CronService != nil {
			app.CronService.Stop()
		}
	}()

	if err := bootstrap.Run(ctx, app.Fiber, cfg); err != nil {
		slog.Error("run application", "error", err)
		os.Exit(1)
	}
}
