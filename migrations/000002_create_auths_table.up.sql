CREATE TABLE IF NOT EXISTS auths (
    id                 BIGSERIAL PRIMARY KEY,
    user_id            BIGINT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    token              TEXT NOT NULL,
    refresh_token      CHAR(64) NOT NULL,
    impersonator_id    BIGINT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE SET NULL,
    meta               JSONB NULL,
    revoked_at         TIMESTAMPTZ NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS auths_token_idx ON auths (token);
CREATE UNIQUE INDEX IF NOT EXISTS auths_refresh_token_idx ON auths (refresh_token);
CREATE INDEX IF NOT EXISTS auths_user_revoked_idx ON auths (user_id, revoked_at);
CREATE INDEX IF NOT EXISTS auths_impersonator_idx ON auths (impersonator_id)
    WHERE impersonator_id IS NOT NULL;
