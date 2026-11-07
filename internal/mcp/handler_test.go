package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type fakeAudit struct {
	mu      sync.Mutex
	entries []AuditEntry
}

func (f *fakeAudit) RecordToolCall(_ context.Context, entry AuditEntry) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entries = append(f.entries, entry)
	return nil
}

func (f *fakeAudit) recorded() []AuditEntry {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]AuditEntry(nil), f.entries...)
}

type bearerTransport struct {
	token string
	base  http.RoundTripper
}

func (t bearerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	clone := r.Clone(r.Context())
	clone.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(clone)
}

func testHandlerConfig(audit AuditRecorder) HandlerConfig {
	return HandlerConfig{
		TokenHash:      hashToken("s3cret"),
		ActorUserID:    "user-1",
		Scopes:         []string{"mcp:read"},
		Version:        "test-1",
		MaxBodyBytes:   1 << 20,
		Stateless:      true,
		Audit:          audit,
		SessionTimeout: time.Minute,
	}
}

func connectClient(t *testing.T, url, token string) *sdkmcp.ClientSession {
	t.Helper()

	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "test-client", Version: "0.0.0"}, nil)
	transport := &sdkmcp.StreamableClientTransport{
		Endpoint:             url,
		HTTPClient:           &http.Client{Transport: bearerTransport{token: token, base: http.DefaultTransport}},
		DisableStandaloneSSE: true,
	}
	session, err := client.Connect(context.Background(), transport, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })
	return session
}

func decodeStructured[T any](t *testing.T, result *sdkmcp.CallToolResult) T {
	t.Helper()
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal structured content: %v", err)
	}
	return out
}

func TestHandlerSystemTools(t *testing.T) {
	server := httptest.NewServer(NewHandler(testHandlerConfig(&fakeAudit{})))
	defer server.Close()

	session := connectClient(t, server.URL, "s3cret")
	ctx := context.Background()

	pingResult, err := session.CallTool(ctx, &sdkmcp.CallToolParams{Name: "ping"})
	if err != nil {
		t.Fatalf("ping: %v", err)
	}
	if pingResult.IsError {
		t.Fatalf("ping returned tool error: %v", pingResult.Content)
	}
	ping := decodeStructured[struct {
		Version string `json:"version"`
		Time    string `json:"time"`
	}](t, pingResult)
	if ping.Version != "test-1" {
		t.Errorf("ping version = %q, want %q", ping.Version, "test-1")
	}
	if _, err := time.Parse(time.RFC3339, ping.Time); err != nil {
		t.Errorf("ping time %q is not RFC3339: %v", ping.Time, err)
	}

	whoamiResult, err := session.CallTool(ctx, &sdkmcp.CallToolParams{Name: "whoami"})
	if err != nil {
		t.Fatalf("whoami: %v", err)
	}
	if whoamiResult.IsError {
		t.Fatalf("whoami returned tool error: %v", whoamiResult.Content)
	}
	whoami := decodeStructured[struct {
		Actor   string   `json:"actor"`
		Scopes  []string `json:"scopes"`
		Version string   `json:"version"`
	}](t, whoamiResult)
	if whoami.Actor != "user-1" {
		t.Errorf("whoami actor = %q, want %q", whoami.Actor, "user-1")
	}
	if len(whoami.Scopes) != 1 || whoami.Scopes[0] != "mcp:read" {
		t.Errorf("whoami scopes = %v, want [mcp:read]", whoami.Scopes)
	}
}

func TestWhoamiWithoutActorReturnsToolError(t *testing.T) {
	cfg := testHandlerConfig(&fakeAudit{})
	cfg.ActorUserID = ""

	server := httptest.NewServer(NewHandler(cfg))
	defer server.Close()

	session := connectClient(t, server.URL, "s3cret")
	result, err := session.CallTool(context.Background(), &sdkmcp.CallToolParams{Name: "whoami"})
	if err != nil {
		t.Fatalf("whoami: %v", err)
	}
	if !result.IsError {
		t.Fatal("whoami without an actor should be a tool error")
	}
}

func TestHandlerRequiresBearerToken(t *testing.T) {
	server := httptest.NewServer(NewHandler(testHandlerConfig(&fakeAudit{})))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestHandlerRejectsOversizedBody(t *testing.T) {
	cfg := testHandlerConfig(&fakeAudit{})
	cfg.MaxBodyBytes = 32

	server := httptest.NewServer(NewHandler(cfg))
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL,
		strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"ping","padding":"`+strings.Repeat("a", 4096)+`"}`))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer s3cret")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusRequestEntityTooLarge)
	}
}

func TestAuditMiddlewareRecordsToolCalls(t *testing.T) {
	audit := &fakeAudit{}
	server := httptest.NewServer(NewHandler(testHandlerConfig(audit)))
	defer server.Close()

	session := connectClient(t, server.URL, "s3cret")
	if _, err := session.CallTool(context.Background(), &sdkmcp.CallToolParams{Name: "ping"}); err != nil {
		t.Fatalf("ping: %v", err)
	}

	entries := audit.recorded()
	if len(entries) != 1 {
		t.Fatalf("audit entries = %d, want 1", len(entries))
	}
	if entries[0].Tool != "ping" {
		t.Errorf("audited tool = %q, want %q", entries[0].Tool, "ping")
	}
	if !entries[0].Success {
		t.Errorf("audited call marked failed: %q", entries[0].Error)
	}
	if entries[0].OccurredAt.IsZero() {
		t.Error("audit entry has zero timestamp")
	}
}
