package routes

import (
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"

	mcpserver "go-fiber-vue/internal/mcp"
)

func newMCPTestApp(t *testing.T, handler http.Handler) *fiber.App {
	t.Helper()

	publicDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(publicDir, "build"), 0o755); err != nil {
		t.Fatalf("create build dir: %v", err)
	}

	app := fiber.New()
	RegisterMCP(app, "/api/mcp", handler)
	RegisterWeb(app, publicDir)

	return app
}

func mcpHandlerForTest() http.Handler {
	return mcpserver.NewHandler(mcpserver.HandlerConfig{
		TokenHash:      sha256.Sum256([]byte("s3cret")),
		ActorUserID:    "user-1",
		Scopes:         []string{"mcp:read"},
		Version:        "test-1",
		MaxBodyBytes:   1 << 20,
		SessionTimeout: 0,
		Stateless:      true,
	})
}

func TestMCPRouteRejectsMissingToken(t *testing.T) {
	app := newMCPTestApp(t, mcpHandlerForTest())

	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			res, err := app.Test(httptest.NewRequest(method, "/api/mcp", nil))
			if err != nil {
				t.Fatalf("request: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusUnauthorized)
			}
		})
	}
}

func TestMCPRouteAcceptsConfiguredToken(t *testing.T) {
	app := newMCPTestApp(t, mcpHandlerForTest())

	req := httptest.NewRequest(http.MethodPost, "/api/mcp", strings.NewReader("{}"))
	req.Header.Set("Authorization", "Bearer s3cret")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		t.Fatal("the configured bearer token was rejected")
	}
}

func TestMCPRouteIsSkippedWhenDisabled(t *testing.T) {
	app := newMCPTestApp(t, nil)

	status, contentType, _ := request(t, app, "/api/mcp")
	if status != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", status, http.StatusNotFound)
	}
	if strings.Contains(contentType, "text/html") {
		t.Fatal("/api/mcp must not be swallowed by the SPA fallback")
	}
}

func TestMCPRouteKeepsSPAFallback(t *testing.T) {
	app := newMCPTestApp(t, mcpHandlerForTest())

	status, contentType, body := request(t, app, "/admin/login")
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if !strings.Contains(contentType, "text/html") {
		t.Fatalf("content type = %q, want text/html", contentType)
	}
	if !strings.Contains(body, `<div id="app"></div>`) {
		t.Fatal("SPA fallback body is missing the mount point")
	}
}

func TestMCPRouteKeepsUnknownAPINotFound(t *testing.T) {
	app := newMCPTestApp(t, mcpHandlerForTest())

	status, _, _ := request(t, app, "/api/v1/nope")
	if status != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", status)
	}
}
