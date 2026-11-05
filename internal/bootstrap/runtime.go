package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/config"
)

func Run(ctx context.Context, app *fiber.App, cfg *config.Config) error {
	errCh := make(chan error, 1)

	go func() {
		addr := ":" + cfg.AppPort
		slog.Info("http server starting", "addr", addr)
		errCh <- app.Listen(addr)
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}
		return nil
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("serve http: %w", err)
		}
		return nil
	}
}
