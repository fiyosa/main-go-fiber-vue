CREATE TABLE IF NOT EXISTS roles (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(125) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS roles_name_idx ON roles (name);

CREATE TABLE IF NOT EXISTS permissions (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(125) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS permissions_name_idx ON permissions (name);

CREATE TABLE IF NOT EXISTS user_has_roles (
    user_id BIGINT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS user_has_roles_role_idx ON user_has_roles (role_id);

CREATE TABLE IF NOT EXISTS role_has_permissions (
    role_id       BIGINT NOT NULL REFERENCES roles (id) ON UPDATE CASCADE ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX IF NOT EXISTS role_has_permissions_permission_idx ON role_has_permissions (permission_id);
