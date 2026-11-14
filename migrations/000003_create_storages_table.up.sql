CREATE TABLE IF NOT EXISTS storages (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
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
-- PostgreSQL does not index foreign key columns; both the cascade delete and
-- per-user listing filter on user_id.
CREATE INDEX IF NOT EXISTS storages_user_id_idx ON storages (user_id);
