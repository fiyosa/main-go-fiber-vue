package bootstrap

import (
	"database/sql"
	"net/http"
	"strings"

	"go-fiber-vue/internal/config"
	mcpserver "go-fiber-vue/internal/mcp"
)

// newMCPHandler builds the MCP HTTP handler from configuration and an optional
// database handle. It returns nil when MCP is disabled so callers can skip route
// registration entirely; a nil database is valid and only makes run_sql return a
// safe "database unavailable" tool error.
func newMCPHandler(cfg *config.Config, sqlDB *sql.DB) http.Handler {
	if cfg == nil || !cfg.MCPEnabled {
		return nil
	}

	return mcpserver.NewHandler(mcpserver.HandlerConfig{
		TokenHash:      cfg.MCPTokenHash(),
		ActorUserID:    cfg.MCPActorUserID,
		Scopes:         splitCSV(cfg.MCPScopes),
		Version:        mcpserver.DefaultVersion,
		MaxBodyBytes:   cfg.MCPMaxBodyBytes,
		AllowedOrigins: splitCSV(cfg.MCPAllowedOrigins),
		SessionTimeout: cfg.MCPSessionTimeout,
		Stateless:      cfg.MCPStateless,
		SQLExecutor:    mcpserver.NewSQLExecutor(sqlDB),
		SQLMaxRows:     cfg.MCPSQLMaxRows,
		SQLTimeout:     cfg.MCPSQLTimeout,
	})
}

// splitCSV splits a comma-separated setting, trimming spaces and dropping empty
// entries.
func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
