package mcp

import (
	"log/slog"

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
	ActorUserID string
	Logger      *slog.Logger
}

// NewMCPServer builds a server with the read-only system tools registered.
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
		Instructions: "Read-only tools for querying the application API.",
		Logger:       cfg.Logger,
	})

	tools.RegisterSystem(server, tools.SystemConfig{
		Version:     version,
		Scopes:      cfg.Scopes,
		ActorUserID: cfg.ActorUserID,
	})

	return server
}
