package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-fiber-vue/internal/config"
)

func mcpConfig() *config.Config {
	return &config.Config{
		MCPEnabled:        true,
		MCPPath:           "/api/mcp",
		MCPBearerToken:    "s3cret",
		MCPActorUserID:    "user-1",
		MCPScopes:         "mcp:read, mcp:write",
		MCPMaxBodyBytes:   1 << 20,
		MCPSessionTimeout: 0,
		MCPStateless:      true,
		MCPSQLMaxRows:     10,
		MCPSQLTimeout:     0,
	}
}

func TestNewMCPHandlerDisabled(t *testing.T) {
	cfg := mcpConfig()
	cfg.MCPEnabled = false

	if handler := newMCPHandler(cfg, nil); handler != nil {
		t.Fatal("newMCPHandler must return nil when MCP_ENABLED=false")
	}
}

func TestNewMCPHandlerNilConfig(t *testing.T) {
	if handler := newMCPHandler(nil, nil); handler != nil {
		t.Fatal("newMCPHandler must return nil without config")
	}
}

func TestNewMCPHandlerRequiresBearerToken(t *testing.T) {
	handler := newMCPHandler(mcpConfig(), nil)
	if handler == nil {
		t.Fatal("newMCPHandler returned nil while MCP is enabled")
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/mcp", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestNewMCPHandlerAcceptsConfiguredToken(t *testing.T) {
	handler := newMCPHandler(mcpConfig(), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/mcp", strings.NewReader("{}"))
	req.Header.Set("Authorization", "Bearer s3cret")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code == http.StatusUnauthorized {
		t.Fatal("the configured bearer token was rejected")
	}
}

func TestSplitCSV(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{input: "", want: nil},
		{input: "  ", want: nil},
		{input: "mcp:read", want: []string{"mcp:read"}},
		{input: "mcp:read, mcp:write", want: []string{"mcp:read", "mcp:write"}},
		{input: " a ,, b ", want: []string{"a", "b"}},
	}

	for _, tc := range tests {
		got := splitCSV(tc.input)
		if len(got) != len(tc.want) {
			t.Fatalf("splitCSV(%q) = %v, want %v", tc.input, got, tc.want)
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Fatalf("splitCSV(%q) = %v, want %v", tc.input, got, tc.want)
			}
		}
	}
}
