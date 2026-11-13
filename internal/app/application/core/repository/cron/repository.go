package cron

import (
	"context"
	"time"

	"go-fiber-vue/internal/app/domain/core/entity"
)

type Repository interface {
	List(ctx context.Context) ([]entity.Cron, error)
	Find(ctx context.Context, id uint64) (*entity.Cron, error)
	Create(ctx context.Context, cron *entity.Cron) error
	Update(ctx context.Context, cron *entity.Cron) error
	Delete(ctx context.Context, id uint64) error
	SetEnabled(ctx context.Context, id uint64, enabled bool) error

	ListRuns(ctx context.Context, cronID uint64, limit int) ([]entity.CronRun, error)
	FindRun(ctx context.Context, runID uint64) (*entity.CronRun, error)
	CreateRun(ctx context.Context, run *entity.CronRun) error
	FinishRun(ctx context.Context, runID uint64, status string, finishedAt time.Time, durationMS int64, output string, errMsg string) error
}
