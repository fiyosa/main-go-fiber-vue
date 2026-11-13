package cron

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	cronrepo "go-fiber-vue/internal/app/application/core/repository/cron"
	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/platform/cronrunner"
	"go-fiber-vue/internal/platform/secretbox"
)

type Service struct {
	db        *gorm.DB
	repo      cronrepo.Repository
	secretKey string
	runner    *cronrunner.Runner
	scheduler *cron.Cron
	cronMap   map[uint64]cron.EntryID
	running   map[uint64]bool
	mu        sync.Mutex
}

type Options struct {
	DB        *gorm.DB
	Repo      cronrepo.Repository
	SecretKey string
}

func NewService(opts Options) *Service {
	return &Service{
		db:        opts.DB,
		repo:      opts.Repo,
		secretKey: opts.SecretKey,
		runner:    cronrunner.NewRunner(),
		cronMap:   make(map[uint64]cron.EntryID),
		running:   make(map[uint64]bool),
	}
}

func (s *Service) Start(ctx context.Context) error {
	if s.repo == nil || s.secretKey == "" {
		slog.Info("cron service disabled (repo or secret key missing)")
		return nil
	}

	s.scheduler = cron.New(cron.WithSeconds())
	s.scheduler.Start()

	crons, err := s.repo.List(ctx)
	if err != nil {
		slog.Warn("could not load crons on startup", "error", err)
		return nil
	}

	for _, c := range crons {
		if c.Enabled {
			s.scheduleCron(&c)
		}
	}

	slog.Info("cron scheduler started with active jobs", "count", len(s.cronMap))
	return nil
}

func (s *Service) Stop() {
	if s.scheduler != nil {
		s.scheduler.Stop()
	}
}

func (s *Service) scheduleCron(c *entity.Cron) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.scheduler == nil {
		return
	}
	cronID := c.ID
	scheduleExpr := c.Schedule

	// robfig/cron v3 parser: second, minute, hour, dom, month, dow.
	// Users may give the classic 5-field expression (minute hour dom month
	// dow) or the full 6-field one; pad the former with a zero second field.
	fullExpr := scheduleExpr
	fields := strings.Fields(scheduleExpr)
	if len(fields) == 5 {
		fullExpr = "0 " + scheduleExpr
	}

	// robfig/cron supports a CRON_TZ=<location> prefix, which lets every job
	// carry its own timezone instead of sharing the process default.
	if tz := strings.TrimSpace(c.Timezone); tz != "" {
		if _, err := time.LoadLocation(tz); err != nil {
			slog.Error("invalid cron timezone", "cron_id", cronID, "timezone", tz, "error", err)
			return
		}
		fullExpr = "CRON_TZ=" + tz + " " + fullExpr
	}

	entryID, err := s.scheduler.AddFunc(fullExpr, func() {
		s.ExecuteCron(context.Background(), cronID, "scheduled")
	})
	if err != nil {
		slog.Error("failed to schedule cron", "cron_id", cronID, "schedule", scheduleExpr, "error", err)
		return
	}

	s.cronMap[cronID] = entryID
}



func (s *Service) ReloadCron(ctx context.Context, cronID uint64) {
	if s.scheduler == nil {
		return
	}
	s.mu.Lock()
	if entryID, ok := s.cronMap[cronID]; ok {
		s.scheduler.Remove(entryID)
		delete(s.cronMap, cronID)
	}
	s.mu.Unlock()

	c, err := s.repo.Find(ctx, cronID)
	if err != nil || c == nil || !c.Enabled {
		return
	}

	s.scheduleCron(c)
}

func (s *Service) ExecuteCron(ctx context.Context, cronID uint64, triggerType string) (*entity.CronRun, error) {
	c, err := s.repo.Find(ctx, cronID)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	if c.OverlapPolicy == "skip" && s.running[cronID] {
		s.mu.Unlock()
		// Record skipped run
		run := &entity.CronRun{
			CronID:      cronID,
			Status:      "skipped",
			TriggerType: triggerType,
			StartedAt:   time.Now(),
			FinishedAt:  ptrTime(time.Now()),
			Output:      "skipped due to overlap_policy = skip",
		}
		_ = s.repo.CreateRun(ctx, run)
		return run, nil
	}
	s.running[cronID] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.running[cronID] = false
		s.mu.Unlock()
	}()

	startedAt := time.Now()
	run := &entity.CronRun{
		CronID:      cronID,
		Status:      "running",
		TriggerType: triggerType,
		StartedAt:   startedAt,
	}
	if err := s.repo.CreateRun(ctx, run); err != nil {
		return nil, err
	}

	// Decrypt env
	envMap := make(map[string]string)
	if c.EnvEncrypted != nil && *c.EnvEncrypted != "" {
		var cipherMap map[string]string
		if err := json.Unmarshal([]byte(*c.EnvEncrypted), &cipherMap); err == nil {
			for k, cipherVal := range cipherMap {
				if plain, err := secretbox.Decrypt(cipherVal, s.secretKey); err == nil {
					envMap[k] = plain
				} else {
					envMap[k] = "[decryption_failed]"
				}
			}
		}
	}

	event := map[string]interface{}{
		"cron":          c.Name,
		"scheduledTime": startedAt.Format(time.RFC3339),
		"type":          triggerType,
	}

	execCtx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	output, scriptErr, err := s.runner.Execute(execCtx, c.Script, event, envMap)
	finishedAt := time.Now()
	durationMS := finishedAt.Sub(startedAt).Milliseconds()

	status := "success"
	errMsg := ""
	if err != nil {
		status = "failed"
		errMsg = err.Error()
	} else if scriptErr != "" {
		status = "failed"
		errMsg = scriptErr
	}

	_ = s.repo.FinishRun(ctx, run.ID, status, finishedAt, durationMS, output, errMsg)

	run.Status = status
	run.FinishedAt = &finishedAt
	run.DurationMS = durationMS
	run.Output = output
	run.Error = errMsg

	return run, nil
}

func EncryptEnvMap(env map[string]string, secretKey string) (*string, error) {
	if len(env) == 0 {
		return nil, nil
	}
	cipherMap := make(map[string]string, len(env))
	for k, v := range env {
		cipher, err := secretbox.Encrypt(v, secretKey)
		if err != nil {
			return nil, err
		}
		cipherMap[k] = cipher
	}
	b, err := json.Marshal(cipherMap)
	if err != nil {
		return nil, err
	}
	s := string(b)
	return &s, nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
