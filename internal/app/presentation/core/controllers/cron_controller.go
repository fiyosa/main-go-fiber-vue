package controllers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"

	cronrepo "go-fiber-vue/internal/app/application/core/repository/cron"
	cronsvc "go-fiber-vue/internal/app/application/core/service/cron"
	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	"go-fiber-vue/internal/app/presentation/core/request"
	"go-fiber-vue/internal/app/presentation/core/response"
	"go-fiber-vue/internal/platform/hashid"
	"go-fiber-vue/internal/platform/secretbox"
)

type CronController struct {
	repo      cronrepo.Repository
	svc       *cronsvc.Service
	ids       *hashid.Encoder
	secretKey string
}

func NewCronController(repo cronrepo.Repository, svc *cronsvc.Service, ids *hashid.Encoder, secretKey string) *CronController {
	return &CronController{
		repo:      repo,
		svc:       svc,
		ids:       ids,
		secretKey: secretKey,
	}
}

type cronResource struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Schedule      string            `json:"schedule"`
	Timezone      string            `json:"timezone"`
	Script        string            `json:"script"`
	Env           map[string]string `json:"env,omitempty"`
	OverlapPolicy string            `json:"overlap_policy"`
	Enabled       bool              `json:"enabled"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type cronRunResource struct {
	ID          string     `json:"id"`
	CronID      string     `json:"cron_id"`
	Status      string     `json:"status"`
	TriggerType string     `json:"trigger_type"`
	StartedAt   time.Time  `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	DurationMS  int64      `json:"duration_ms"`
	Output      string     `json:"output,omitempty"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// 1. GET /api/v1/crons
func (c *CronController) Index(fctx fiber.Ctx) error {
	if c == nil || c.repo == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	rows, err := c.repo.List(fctx.Context())
	if err != nil {
		return err
	}

	data := make([]cronResource, 0, len(rows))
	for _, row := range rows {
		id, err := c.ids.Encode(row.ID)
		if err != nil {
			return exceptions.Internal(err)
		}
		data = append(data, cronResource{
			ID:            id,
			Name:          row.Name,
			Schedule:      row.Schedule,
			Timezone:      row.Timezone,
			Script:        row.Script,
			OverlapPolicy: row.OverlapPolicy,
			Enabled:       row.Enabled,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
	}

	return response.Success(fctx, data, "lang.retrieved", nil, fiber.StatusOK)
}

// 2. POST /api/v1/crons
func (c *CronController) Store(fctx fiber.Ctx) error {
	if c == nil || c.repo == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	var req request.CronCreateRequest
	if err := fctx.Bind().JSON(&req); err != nil {
		return exceptions.BadRequest(err)
	}

	envEncrypted, err := cronsvc.EncryptEnvMap(req.Env, c.secretKey)
	if err != nil {
		return exceptions.Internal(err)
	}

	tz := req.Timezone
	if tz == "" {
		tz = "Asia/Jakarta"
	}
	policy := req.OverlapPolicy
	if policy == "" {
		policy = "skip"
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	cronEntity := &entity.Cron{
		Name:          req.Name,
		Schedule:      req.Schedule,
		Timezone:      tz,
		Script:        req.Script,
		EnvEncrypted:  envEncrypted,
		OverlapPolicy: policy,
		Enabled:       enabled,
	}

	if err := c.repo.Create(fctx.Context(), cronEntity); err != nil {
		return err
	}

	if c.svc != nil {
		c.svc.ReloadCron(fctx.Context(), cronEntity.ID)
	}

	id, err := c.ids.Encode(cronEntity.ID)
	if err != nil {
		return exceptions.Internal(err)
	}

	res := cronResource{
		ID:            id,
		Name:          cronEntity.Name,
		Schedule:      cronEntity.Schedule,
		Timezone:      cronEntity.Timezone,
		Script:        cronEntity.Script,
		Env:           req.Env,
		OverlapPolicy: cronEntity.OverlapPolicy,
		Enabled:       cronEntity.Enabled,
		CreatedAt:     cronEntity.CreatedAt,
		UpdatedAt:     cronEntity.UpdatedAt,
	}

	return response.Success(fctx, res, "lang.created", nil, fiber.StatusCreated)
}

// 3. GET /api/v1/crons/:id
func (c *CronController) Show(fctx fiber.Ctx) error {
	if c == nil || c.repo == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	dbID, err := c.ids.Decode(fctx.Params("id"))
	if err != nil {
		return exceptions.NotFound(err)
	}

	row, err := c.repo.Find(fctx.Context(), dbID)
	if err != nil {
		return err
	}

	id, err := c.ids.Encode(row.ID)
	if err != nil {
		return exceptions.Internal(err)
	}

	var env map[string]string
	if row.EnvEncrypted != nil && *row.EnvEncrypted != "" {
		var cipherMap map[string]string
		if err := json.Unmarshal([]byte(*row.EnvEncrypted), &cipherMap); err == nil {
			env = make(map[string]string, len(cipherMap))
			for k, cipherVal := range cipherMap {
				if plain, err := secretbox.Decrypt(cipherVal, c.secretKey); err == nil {
					env[k] = plain
				} else {
					env[k] = "[decryption_failed]"
				}
			}
		}
	}

	res := cronResource{
		ID:            id,
		Name:          row.Name,
		Schedule:      row.Schedule,
		Timezone:      row.Timezone,
		Script:        row.Script,
		Env:           env,
		OverlapPolicy: row.OverlapPolicy,
		Enabled:       row.Enabled,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}

	return response.Success(fctx, res, "lang.retrieved", nil, fiber.StatusOK)
}

// 4. PUT /api/v1/crons/:id
func (c *CronController) Update(fctx fiber.Ctx) error {
	if c == nil || c.repo == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	dbID, err := c.ids.Decode(fctx.Params("id"))
	if err != nil {
		return exceptions.NotFound(err)
	}

	var req request.CronUpdateRequest
	if err := fctx.Bind().JSON(&req); err != nil {
		return exceptions.BadRequest(err)
	}

	envEncrypted, err := cronsvc.EncryptEnvMap(req.Env, c.secretKey)
	if err != nil {
		return exceptions.Internal(err)
	}

	tz := req.Timezone
	if tz == "" {
		tz = "Asia/Jakarta"
	}
	policy := req.OverlapPolicy
	if policy == "" {
		policy = "skip"
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	cronEntity := &entity.Cron{
		ID:            dbID,
		Name:          req.Name,
		Schedule:      req.Schedule,
		Timezone:      tz,
		Script:        req.Script,
		EnvEncrypted:  envEncrypted,
		OverlapPolicy: policy,
		Enabled:       enabled,
	}

	if err := c.repo.Update(fctx.Context(), cronEntity); err != nil {
		return err
	}

	if c.svc != nil {
		c.svc.ReloadCron(fctx.Context(), dbID)
	}

	id, err := c.ids.Encode(dbID)
	if err != nil {
		return exceptions.Internal(err)
	}

	res := cronResource{
		ID:            id,
		Name:          cronEntity.Name,
		Schedule:      cronEntity.Schedule,
		Timezone:      cronEntity.Timezone,
		Script:        cronEntity.Script,
		Env:           req.Env,
		OverlapPolicy: cronEntity.OverlapPolicy,
		Enabled:       cronEntity.Enabled,
		UpdatedAt:     time.Now(),
	}

	return response.Success(fctx, res, "lang.updated", nil, fiber.StatusOK)
}

// 5. DELETE /api/v1/crons/:id
func (c *CronController) Delete(fctx fiber.Ctx) error {
	if c == nil || c.repo == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	dbID, err := c.ids.Decode(fctx.Params("id"))
	if err != nil {
		return exceptions.NotFound(err)
	}

	if err := c.repo.Delete(fctx.Context(), dbID); err != nil {
		return err
	}

	if c.svc != nil {
		c.svc.ReloadCron(fctx.Context(), dbID)
	}

	return response.Success(fctx, nil, "lang.deleted", nil, fiber.StatusOK)
}

// 6. POST /api/v1/crons/:id/toggle
func (c *CronController) Toggle(fctx fiber.Ctx) error {
	if c == nil || c.repo == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	dbID, err := c.ids.Decode(fctx.Params("id"))
	if err != nil {
		return exceptions.NotFound(err)
	}

	row, err := c.repo.Find(fctx.Context(), dbID)
	if err != nil {
		return err
	}

	newStatus := !row.Enabled
	if err := c.repo.SetEnabled(fctx.Context(), dbID, newStatus); err != nil {
		return err
	}

	if c.svc != nil {
		c.svc.ReloadCron(fctx.Context(), dbID)
	}

	id, err := c.ids.Encode(dbID)
	if err != nil {
		return exceptions.Internal(err)
	}

	res := cronResource{
		ID:            id,
		Name:          row.Name,
		Schedule:      row.Schedule,
		Timezone:      row.Timezone,
		Script:        row.Script,
		OverlapPolicy: row.OverlapPolicy,
		Enabled:       newStatus,
		UpdatedAt:     time.Now(),
	}

	return response.Success(fctx, res, "lang.updated", nil, fiber.StatusOK)
}

// 7. POST /api/v1/crons/:id/trigger
func (c *CronController) Trigger(fctx fiber.Ctx) error {
	if c == nil || c.svc == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	dbID, err := c.ids.Decode(fctx.Params("id"))
	if err != nil {
		return exceptions.NotFound(err)
	}

	run, err := c.svc.ExecuteCron(fctx.Context(), dbID, "manual")
	if err != nil {
		return err
	}

	runID, err := c.ids.Encode(run.ID)
	if err != nil {
		return exceptions.Internal(err)
	}
	cronIDStr, err := c.ids.Encode(run.CronID)
	if err != nil {
		return exceptions.Internal(err)
	}

	res := cronRunResource{
		ID:          runID,
		CronID:      cronIDStr,
		Status:      run.Status,
		TriggerType: run.TriggerType,
		StartedAt:   run.StartedAt,
		FinishedAt:  run.FinishedAt,
		DurationMS:  run.DurationMS,
		Output:      run.Output,
		Error:       run.Error,
		CreatedAt:   run.CreatedAt,
	}

	return response.Success(fctx, res, "lang.triggered", nil, fiber.StatusOK)
}

// 8. GET /api/v1/crons/:id/runs
func (c *CronController) ListRuns(fctx fiber.Ctx) error {
	if c == nil || c.repo == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	dbID, err := c.ids.Decode(fctx.Params("id"))
	if err != nil {
		return exceptions.NotFound(err)
	}

	runs, err := c.repo.ListRuns(fctx.Context(), dbID, 100)
	if err != nil {
		return err
	}

	data := make([]cronRunResource, 0, len(runs))
	for _, run := range runs {
		runID, err := c.ids.Encode(run.ID)
		if err != nil {
			return exceptions.Internal(err)
		}
		cronIDStr, err := c.ids.Encode(run.CronID)
		if err != nil {
			return exceptions.Internal(err)
		}
		data = append(data, cronRunResource{
			ID:          runID,
			CronID:      cronIDStr,
			Status:      run.Status,
			TriggerType: run.TriggerType,
			StartedAt:   run.StartedAt,
			FinishedAt:  run.FinishedAt,
			DurationMS:  run.DurationMS,
			Output:      run.Output,
			Error:       run.Error,
			CreatedAt:   run.CreatedAt,
		})
	}

	return response.Success(fctx, data, "lang.retrieved", nil, fiber.StatusOK)
}

// 9. GET /api/v1/crons/runs/:run_id (and GET /api/v1/crons/runs/:id)
func (c *CronController) ShowRun(fctx fiber.Ctx) error {
	if c == nil || c.repo == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	paramID := fctx.Params("run_id")
	if paramID == "" {
		paramID = fctx.Params("id")
	}

	dbID, err := c.ids.Decode(paramID)
	if err != nil {
		return exceptions.NotFound(err)
	}

	run, err := c.repo.FindRun(fctx.Context(), dbID)
	if err != nil {
		return err
	}

	runID, err := c.ids.Encode(run.ID)
	if err != nil {
		return exceptions.Internal(err)
	}
	cronIDStr, err := c.ids.Encode(run.CronID)
	if err != nil {
		return exceptions.Internal(err)
	}

	res := cronRunResource{
		ID:          runID,
		CronID:      cronIDStr,
		Status:      run.Status,
		TriggerType: run.TriggerType,
		StartedAt:   run.StartedAt,
		FinishedAt:  run.FinishedAt,
		DurationMS:  run.DurationMS,
		Output:      run.Output,
		Error:       run.Error,
		CreatedAt:   run.CreatedAt,
	}

	return response.Success(fctx, res, "lang.retrieved", nil, fiber.StatusOK)
}

// 10. GET /api/v1/crons/runs (List all runs across all crons)
func (c *CronController) ListAllRuns(fctx fiber.Ctx) error {
	if c == nil || c.repo == nil || c.ids == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	runs, err := c.repo.ListRuns(fctx.Context(), 0, 100)
	if err != nil {
		return err
	}

	data := make([]cronRunResource, 0, len(runs))
	for _, run := range runs {
		runID, err := c.ids.Encode(run.ID)
		if err != nil {
			return exceptions.Internal(err)
		}
		cronIDStr, err := c.ids.Encode(run.CronID)
		if err != nil {
			return exceptions.Internal(err)
		}
		data = append(data, cronRunResource{
			ID:          runID,
			CronID:      cronIDStr,
			Status:      run.Status,
			TriggerType: run.TriggerType,
			StartedAt:   run.StartedAt,
			FinishedAt:  run.FinishedAt,
			DurationMS:  run.DurationMS,
			Output:      run.Output,
			Error:       run.Error,
			CreatedAt:   run.CreatedAt,
		})
	}

	return response.Success(fctx, data, "lang.retrieved", nil, fiber.StatusOK)
}
