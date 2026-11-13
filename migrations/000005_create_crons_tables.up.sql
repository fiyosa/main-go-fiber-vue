CREATE TABLE IF NOT EXISTS crons (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    schedule       VARCHAR(100) NOT NULL,
    timezone       VARCHAR(100) NOT NULL DEFAULT 'Asia/Jakarta',
    script         TEXT NOT NULL,
    env_encrypted  JSONB NULL,
    overlap_policy VARCHAR(20) NOT NULL DEFAULT 'skip',
    enabled        BOOLEAN NOT NULL DEFAULT true,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS crons_name_idx ON crons (name);
CREATE INDEX IF NOT EXISTS crons_enabled_idx ON crons (enabled);

CREATE TABLE IF NOT EXISTS cron_runs (
    id           BIGSERIAL PRIMARY KEY,
    cron_id      BIGINT NOT NULL REFERENCES crons(id) ON DELETE CASCADE,
    status       VARCHAR(50) NOT NULL,
    trigger_type VARCHAR(50) NOT NULL,
    started_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at  TIMESTAMPTZ NULL,
    duration_ms  BIGINT NOT NULL DEFAULT 0,
    output       TEXT NULL,
    error        TEXT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS cron_runs_cron_id_idx ON cron_runs (cron_id);
CREATE INDEX IF NOT EXISTS cron_runs_started_at_idx ON cron_runs (started_at DESC);
