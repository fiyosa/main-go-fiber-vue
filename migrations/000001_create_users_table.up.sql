CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    username      VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password      VARCHAR(255) NOT NULL,
    role          VARCHAR(125) NULL,
    division_id   BIGINT NULL,
    token_version INTEGER NOT NULL DEFAULT 0,
    last_login_at TIMESTAMPTZ NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_username_idx ON users (username);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);
