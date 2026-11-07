package mcp

import (
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"testing"
)

func hashToken(token string) [32]byte { return sha256.Sum256([]byte(token)) }

func TestRequireBearer(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   int
	}{
		{name: "missing header", header: "", want: http.StatusUnauthorized},
		{name: "wrong scheme", header: "Basic abc", want: http.StatusUnauthorized},
		{name: "empty token", header: "Bearer   ", want: http.StatusUnauthorized},
		{name: "wrong token", header: "Bearer nope", want: http.StatusUnauthorized},
		{name: "correct token", header: "Bearer s3cret", want: http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			handler := RequireBearer(next, hashToken("s3cret"), "user-1")

			req := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d (body %q)", rec.Code, tc.want, rec.Body.String())
			}
		})
	}
}

func TestRequireBearerInjectsActor(t *testing.T) {
	var got string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actor, ok := ActorFromContext(r.Context())
		if !ok {
			t.Error("actor missing from context")
		}
		got = actor
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
	req.Header.Set("Authorization", "Bearer s3cret")
	RequireBearer(next, hashToken("s3cret"), "user-1").ServeHTTP(httptest.NewRecorder(), req)

	if got != "user-1" {
		t.Fatalf("actor = %q, want %q", got, "user-1")
	}
}

func TestRequireOrigin(t *testing.T) {
	t.Run("no allowlist permits any origin", func(t *testing.T) {
		handler := RequireOrigin(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), nil)

		req := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
		req.Header.Set("Origin", "https://example.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("allowed origin passes and gets CORS headers", func(t *testing.T) {
		handler := RequireOrigin(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), []string{"https://app.example.com"})

		req := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
		req.Header.Set("Origin", "https://app.example.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
			t.Errorf("Access-Control-Allow-Origin = %q", got)
		}
	})

	t.Run("disallowed origin is rejected", func(t *testing.T) {
		handler := RequireOrigin(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), []string{"https://app.example.com"})

		req := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
		req.Header.Set("Origin", "https://evil.example.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
		}
	})

	t.Run("preflight does not reach the handler", func(t *testing.T) {
		called := false
		handler := RequireOrigin(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}), []string{"https://app.example.com"})

		req := httptest.NewRequest(http.MethodOptions, "/api/mcp", nil)
		req.Header.Set("Origin", "https://app.example.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
		}
		if called {
			t.Fatal("preflight reached the wrapped handler")
		}
	})
}
