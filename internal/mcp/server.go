package mcp

import (
	"log/slog"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"go-fiber-vue/internal/mcp/tools"
)

const (
	// ServerName is the programmatic name advertised during initialization.
	ServerName = "app-mcp"
	// DefaultVersion is used when no build version is supplied.
	DefaultVersion = "0.1.0"
)

// ServerConfig configures an MCP server instance.
type ServerConfig struct {
	Version     string
	Scopes      []string
	SQLExecutor tools.DBExecutor
	SQLMaxRows  int
	SQLTimeout  time.Duration
	Logger      *slog.Logger
}

// NewMCPServer builds a server with the system and SQL tools registered.
func NewMCPServer(cfg ServerConfig) *sdkmcp.Server {
	version := cfg.Version
	if version == "" {
		version = DefaultVersion
	}

	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    ServerName,
		Title:   "Application MCP",
		Version: version,
	}, &sdkmcp.ServerOptions{
		Instructions: "Tools for inspecting and modifying the application database over MCP.",
		Logger:       cfg.Logger,
	})

	tools.RegisterSystem(server, tools.SystemConfig{
		Version: version,
		Scopes:  cfg.Scopes,
	})

	tools.RegisterSQL(server, tools.SQLConfig{
		Executor: cfg.SQLExecutor,
		MaxRows:  cfg.SQLMaxRows,
		Timeout:  cfg.SQLTimeout,
	})

	return server
}
