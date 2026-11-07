package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type bearerTransport struct {
	token string
	base  http.RoundTripper
}

func (t bearerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	clone := r.Clone(r.Context())
	clone.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(clone)
}

func testHandlerConfig() HandlerConfig {
	return HandlerConfig{
		TokenHash:      hashToken("s3cret"),
		ActorUserID:    "user-1",
		Scopes:         []string{"mcp:read"},
		Version:        "test-1",
		MaxBodyBytes:   1 << 20,
		Stateless:      true,
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
	server := httptest.NewServer(NewHandler(testHandlerConfig()))
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
	cfg := testHandlerConfig()
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
	server := httptest.NewServer(NewHandler(testHandlerConfig()))
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
	cfg := testHandlerConfig()
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
