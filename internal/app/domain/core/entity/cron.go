package entity

import "time"

type Cron struct {
	ID            uint64    `gorm:"column:id;primaryKey"`
	Name          string    `gorm:"column:name;size:255;not null;uniqueIndex"`
	Schedule      string    `gorm:"column:schedule;size:100;not null"`
	Timezone      string    `gorm:"column:timezone;size:100;not null;default:Asia/Jakarta"`
	Script        string    `gorm:"column:script;type:text;not null"`
	EnvEncrypted  *string   `gorm:"column:env_encrypted;type:jsonb"`
	OverlapPolicy string    `gorm:"column:overlap_policy;size:20;not null;default:skip"`
	Enabled       bool      `gorm:"column:enabled;not null;default:true"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (Cron) TableName() string {
	return "crons"
}

type CronRun struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	CronID      uint64     `gorm:"column:cron_id;not null;index"`
	Cron        *Cron      `gorm:"foreignKey:CronID"`
	Status      string     `gorm:"column:status;size:50;not null"`
	TriggerType string     `gorm:"column:trigger_type;size:50;not null"`
	StartedAt   time.Time  `gorm:"column:started_at;not null;index"`
	FinishedAt  *time.Time `gorm:"column:finished_at"`
	DurationMS  int64      `gorm:"column:duration_ms;not null;default:0"`
	Output      string     `gorm:"column:output;type:text"`
	Error       string     `gorm:"column:error;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (CronRun) TableName() string {
	return "cron_runs"
}
