CREATE TABLE IF NOT EXISTS mcp_audit_logs (
    id             BIGSERIAL PRIMARY KEY,
    tool_name      VARCHAR(128) NOT NULL,
    actor          VARCHAR(255) NULL,
    parameters     JSONB NULL,
    result_summary TEXT NULL,
    error          TEXT NULL,
    ip_address     VARCHAR(64) NULL,
    user_agent     VARCHAR(512) NULL,
    duration_ms    BIGINT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS mcp_audit_logs_created_at_idx ON mcp_audit_logs (created_at DESC);
CREATE INDEX IF NOT EXISTS mcp_audit_logs_tool_name_idx ON mcp_audit_logs (tool_name);
