package routes

import (
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func newWebTestApp(t *testing.T) *fiber.App {
	t.Helper()

	publicDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(publicDir, "build"), 0o755); err != nil {
		t.Fatalf("create build dir: %v", err)
	}

	app := fiber.New()
	RegisterWeb(app, publicDir)

	return app
}

func request(t *testing.T, app *fiber.App, path string) (int, string, string) {
	t.Helper()

	res, err := app.Test(httptest.NewRequest("GET", path, nil))
	if err != nil {
		t.Fatalf("request %s: %v", path, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	return res.StatusCode, res.Header.Get("Content-Type"), string(body)
}

func TestViteFaviconIsServed(t *testing.T) {
	publicDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(publicDir, "build"), 0o755); err != nil {
		t.Fatalf("create build dir: %v", err)
	}

	const svg = `<svg xmlns="http://www.w3.org/2000/svg"></svg>`
	if err := os.WriteFile(filepath.Join(publicDir, "vite.svg"), []byte(svg), 0o644); err != nil {
		t.Fatalf("write vite.svg: %v", err)
	}

	app := fiber.New()
	RegisterWeb(app, publicDir)

	status, contentType, body := request(t, app, "/vite.svg")
	if status != 200 {
		t.Fatalf("status = %d, want 200", status)
	}
	if !strings.Contains(contentType, "image/svg+xml") {
		t.Fatalf("content type = %q, want image/svg+xml", contentType)
	}
	if body != svg {
		t.Fatalf("body = %q, want %q", body, svg)
	}
}

func TestUnknownPathFallsBackToSPA(t *testing.T) {
	app := newWebTestApp(t)

	status, contentType, body := request(t, app, "/admin/login")
	if status != 200 {
		t.Fatalf("status = %d, want 200", status)
	}
	if !strings.Contains(contentType, "text/html") {
		t.Fatalf("content type = %q, want text/html", contentType)
	}
	if !strings.Contains(body, `<div id="app"></div>`) {
		t.Fatalf("body does not contain the SPA mount point")
	}
}

func TestUnknownAPIPathIsNotFound(t *testing.T) {
	app := newWebTestApp(t)

	status, _, _ := request(t, app, "/api/v1/nope")
	if status != 404 {
		t.Fatalf("status = %d, want 404", status)
	}
}
