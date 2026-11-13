package cron

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	if db == nil {
		return nil
	}
	return &gormRepository{db: db}
}

func (r *gormRepository) List(ctx context.Context) ([]entity.Cron, error) {
	var crons []entity.Cron
	err := r.db.WithContext(ctx).Order("name ASC").Find(&crons).Error
	if err != nil {
		return nil, exceptions.FromDB(err)
	}
	return crons, nil
}

func (r *gormRepository) Find(ctx context.Context, id uint64) (*entity.Cron, error) {
	var cron entity.Cron
	err := r.db.WithContext(ctx).First(&cron, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NotFound(err)
		}
		return nil, exceptions.FromDB(err)
	}
	return &cron, nil
}

func (r *gormRepository) Create(ctx context.Context, cron *entity.Cron) error {
	err := r.db.WithContext(ctx).Create(cron).Error
	if err != nil {
		return exceptions.FromDB(err)
	}
	return nil
}

func (r *gormRepository) Update(ctx context.Context, cron *entity.Cron) error {
	err := r.db.WithContext(ctx).Model(cron).Updates(map[string]interface{}{
		"name":           cron.Name,
		"schedule":       cron.Schedule,
		"timezone":       cron.Timezone,
		"script":         cron.Script,
		"env_encrypted":  cron.EnvEncrypted,
		"overlap_policy": cron.OverlapPolicy,
		"enabled":        cron.Enabled,
		"updated_at":     time.Now(),
	}).Error
	if err != nil {
		return exceptions.FromDB(err)
	}
	return nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint64) error {
	err := r.db.WithContext(ctx).Delete(&entity.Cron{}, id).Error
	if err != nil {
		return exceptions.FromDB(err)
	}
	return nil
}

func (r *gormRepository) SetEnabled(ctx context.Context, id uint64, enabled bool) error {
	err := r.db.WithContext(ctx).Model(&entity.Cron{}).Where("id = ?", id).Updates(map[string]interface{}{
		"enabled":    enabled,
		"updated_at": time.Now(),
	}).Error
	if err != nil {
		return exceptions.FromDB(err)
	}
	return nil
}

func (r *gormRepository) ListRuns(ctx context.Context, cronID uint64, limit int) ([]entity.CronRun, error) {
	var runs []entity.CronRun
	query := r.db.WithContext(ctx).Order("started_at DESC")
	if cronID > 0 {
		query = query.Where("cron_id = ?", cronID)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&runs).Error
	if err != nil {
		return nil, exceptions.FromDB(err)
	}
	return runs, nil
}

func (r *gormRepository) FindRun(ctx context.Context, runID uint64) (*entity.CronRun, error) {
	var run entity.CronRun
	err := r.db.WithContext(ctx).First(&run, runID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NotFound(err)
		}
		return nil, exceptions.FromDB(err)
	}
	return &run, nil
}

func (r *gormRepository) CreateRun(ctx context.Context, run *entity.CronRun) error {
	err := r.db.WithContext(ctx).Create(run).Error
	if err != nil {
		return exceptions.FromDB(err)
	}
	return nil
}

func (r *gormRepository) FinishRun(ctx context.Context, runID uint64, status string, finishedAt time.Time, durationMS int64, output string, errMsg string) error {
	err := r.db.WithContext(ctx).Model(&entity.CronRun{}).Where("id = ?", runID).Updates(map[string]interface{}{
		"status":      status,
		"finished_at": finishedAt,
		"duration_ms": durationMS,
		"output":      output,
		"error":       errMsg,
		"updated_at":  time.Now(),
	}).Error
	if err != nil {
		return exceptions.FromDB(err)
	}
	return nil
}
