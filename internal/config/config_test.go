package config

import (
	"crypto/sha256"
	"os"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// clearMCPEnv unsets every MCP_* variable so the struct tag defaults are
// exercised regardless of the developer's shell, then restores them.
func clearMCPEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"MCP_ENABLED",
		"MCP_PATH",
		"MCP_TOKEN",
		"MCP_SCOPES",
		"MCP_MAX_BODY_BYTES",
		"MCP_ALLOWED_ORIGINS",
		"MCP_SESSION_TIMEOUT",
		"MCP_STATELESS",
	} {
		previous, had := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
		t.Cleanup(func() {
			if had {
				_ = os.Setenv(key, previous)
			}
		})
	}
}

func TestMCPDefaults(t *testing.T) {
	clearMCPEnv(t)

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		t.Fatalf("process config: %v", err)
	}

	if !cfg.MCPEnabled {
		t.Errorf("MCPEnabled = false, want true")
	}
	if cfg.MCPPath != "/api/mcp" {
		t.Errorf("MCPPath = %q, want %q", cfg.MCPPath, "/api/mcp")
	}
	if cfg.MCPToken != "" {
		t.Errorf("MCPToken = %q, want empty", cfg.MCPToken)
	}
	if cfg.MCPScopes != "mcp:read" {
		t.Errorf("MCPScopes = %q, want %q", cfg.MCPScopes, "mcp:read")
	}
	if cfg.MCPMaxBodyBytes != 1048576 {
		t.Errorf("MCPMaxBodyBytes = %d, want %d", cfg.MCPMaxBodyBytes, 1048576)
	}
	if cfg.MCPAllowedOrigins != "" {
		t.Errorf("MCPAllowedOrigins = %q, want empty", cfg.MCPAllowedOrigins)
	}
	if cfg.MCPSessionTimeout != 60*time.Second {
		t.Errorf("MCPSessionTimeout = %s, want 60s", cfg.MCPSessionTimeout)
	}
	if !cfg.MCPStateless {
		t.Errorf("MCPStateless = false, want true")
	}
}

func TestMCPValidate(t *testing.T) {
	t.Run("disabled without token is allowed", func(t *testing.T) {
		cfg := Config{MCPEnabled: false}
		if err := cfg.validateMCP(); err != nil {
			t.Fatalf("validateMCP() = %v, want nil", err)
		}
	})

	t.Run("enabled without token fails", func(t *testing.T) {
		cfg := Config{MCPEnabled: true}
		if err := cfg.validateMCP(); err == nil {
			t.Fatal("validateMCP() = nil, want error")
		}
	})

	t.Run("enabled with blank token fails", func(t *testing.T) {
		cfg := Config{MCPEnabled: true, MCPToken: "   "}
		if err := cfg.validateMCP(); err == nil {
			t.Fatal("validateMCP() = nil, want error")
		}
	})

	t.Run("enabled with token passes", func(t *testing.T) {
		cfg := Config{MCPEnabled: true, MCPToken: "s3cret"}
		if err := cfg.validateMCP(); err != nil {
			t.Fatalf("validateMCP() = %v, want nil", err)
		}
	})
}

func TestMCPTokenHash(t *testing.T) {
	cfg := Config{MCPToken: "s3cret"}
	want := sha256.Sum256([]byte("s3cret"))

	if got := cfg.MCPTokenHash(); got != want {
		t.Errorf("MCPTokenHash() = %x, want %x", got, want)
	}
}
