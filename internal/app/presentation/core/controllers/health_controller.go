package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/presentation/core/response"
)

type Pinger interface {
	PingContext(ctx context.Context) error
}

type HealthController struct {
	db Pinger
}

func NewHealthController(db Pinger) *HealthController {
	return &HealthController{db: db}
}

func (c *HealthController) Health(fctx fiber.Ctx) error {
	return response.Success(fctx, map[string]string{"status": "ok"}, "common.success", nil, fiber.StatusOK)
}

func (c *HealthController) Ready(fctx fiber.Ctx) error {
	if c.db == nil {
		return response.Error(fctx, "common.try_again_later", nil, fiber.StatusServiceUnavailable)
	}
	ctx, cancel := context.WithTimeout(fctx.Context(), 2*time.Second)
	defer cancel()
	if err := c.db.PingContext(ctx); err != nil {
		return response.Error(fctx, "common.try_again_later", nil, fiber.StatusServiceUnavailable)
	}
	return response.Success(fctx, map[string]string{"status": "ready"}, "common.success", nil, fiber.StatusOK)
}
