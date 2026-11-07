package tools

import (
	"context"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// SystemConfig configures the built-in read-only system tools.
type SystemConfig struct {
	Version string
	Scopes  []string
	// Now is injectable so tests can pin the clock. Defaults to time.Now.
	Now func() time.Time
}

type pingOutput struct {
	Version string `json:"version" jsonschema:"server version"`
	Time    string `json:"time" jsonschema:"current server time in RFC3339"`
}

type whoamiOutput struct {
	Scopes  []string `json:"scopes" jsonschema:"granted scopes"`
	Version string   `json:"version" jsonschema:"server version"`
}

// RegisterSystem registers the ping and whoami tools on the server.
func RegisterSystem(server *sdkmcp.Server, cfg SystemConfig) {
	now := cfg.Now
	if now == nil {
		now = time.Now
	}

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "ping",
		Description: "Returns the server version and current time to verify connectivity.",
	}, func(context.Context, *sdkmcp.CallToolRequest, struct{}) (*sdkmcp.CallToolResult, pingOutput, error) {
		return nil, pingOutput{
			Version: cfg.Version,
			Time:    now().UTC().Format(time.RFC3339),
		}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "whoami",
		Description: "Returns the scopes and server version bound to this token.",
	}, func(context.Context, *sdkmcp.CallToolRequest, struct{}) (*sdkmcp.CallToolResult, whoamiOutput, error) {
		return nil, whoamiOutput{
			Scopes:  cfg.Scopes,
			Version: cfg.Version,
		}, nil
	})
}
