CREATE TABLE IF NOT EXISTS storages (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL,
    key           VARCHAR(32) NOT NULL,
    path          TEXT NOT NULL,
    bucket        VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    mime_type     VARCHAR(255) NOT NULL,
    size_bytes    BIGINT NOT NULL,
    meta          JSONB NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS storages_key_idx ON storages (key);
