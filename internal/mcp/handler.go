package mcp

import (
	"log/slog"
	"net/http"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"go-fiber-vue/internal/mcp/tools"
)

// HandlerConfig configures the net/http delivery layer. It is transport-level
// only: tool behavior lives in the tools package.
type HandlerConfig struct {
	TokenHash      [32]byte
	ActorUserID    string
	Scopes         []string
	Version        string
	MaxBodyBytes   int64
	AllowedOrigins []string
	SessionTimeout time.Duration
	Stateless      bool
	SQLExecutor    tools.DBExecutor
	SQLMaxRows     int
	SQLTimeout     time.Duration
	Logger         *slog.Logger
}

// NewHandler builds the Streamable HTTP handler wrapped with origin checks and
// bearer authentication. The result is a plain net/http handler so Fiber can
// mount it through the adaptor without leaking Fiber types into this package.
//
// Auth runs on the net/http side because the adaptor builds a fresh request;
// Fiber locals are not available downstream.
func NewHandler(cfg HandlerConfig) http.Handler {
	server := NewMCPServer(ServerConfig{
		Version:     cfg.Version,
		Scopes:      cfg.Scopes,
		ActorUserID: cfg.ActorUserID,
		SQLExecutor: cfg.SQLExecutor,
		SQLMaxRows:  cfg.SQLMaxRows,
		SQLTimeout:  cfg.SQLTimeout,
		Logger:      cfg.Logger,
	})

	streamable := sdkmcp.NewStreamableHTTPHandler(
		func(*http.Request) *sdkmcp.Server { return server },
		&sdkmcp.StreamableHTTPOptions{
			Stateless:           cfg.Stateless,
			JSONResponse:        true,
			Logger:              cfg.Logger,
			SessionTimeout:      cfg.SessionTimeout,
			MaxRequestBodyBytes: cfg.MaxBodyBytes,
		},
	)

	var handler http.Handler = streamable
	handler = RequireBearer(handler, cfg.TokenHash, cfg.ActorUserID)
	handler = RequireOrigin(handler, cfg.AllowedOrigins)
	return handler
}
