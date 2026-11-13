package request

type CronCreateRequest struct {
	Name          string            `json:"name" validate:"required"`
	Schedule      string            `json:"schedule" validate:"required"`
	Timezone      string            `json:"timezone"`
	Script        string            `json:"script" validate:"required"`
	Env           map[string]string `json:"env"`
	OverlapPolicy string            `json:"overlap_policy" validate:"omitempty,oneof=skip allow"`
	Enabled       *bool             `json:"enabled"`
}

type CronUpdateRequest struct {
	Name          string            `json:"name" validate:"required"`
	Schedule      string            `json:"schedule" validate:"required"`
	Timezone      string            `json:"timezone"`
	Script        string            `json:"script" validate:"required"`
	Env           map[string]string `json:"env"`
	OverlapPolicy string            `json:"overlap_policy" validate:"omitempty,oneof=skip allow"`
	Enabled       *bool             `json:"enabled"`
}
