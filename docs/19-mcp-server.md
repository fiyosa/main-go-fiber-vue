# 19. MCP Server

## 19.1 Decisions

| Area        | Decision                                                                 |
| ----------- | ------------------------------------------------------------------------ |
| SDK         | `github.com/modelcontextprotocol/go-sdk/mcp` (official SDK)              |
| Transport   | Streamable HTTP inside the existing Fiber process                        |
| Endpoint    | `/api/mcp` on the same port as the HTTP API                              |
| Auth        | Bearer token read from `MCP_BEARER_TOKEN` in `.env`; no OAuth            |
| Clients     | Transport-agnostic; any MCP client that supports remote HTTP + headers   |
| Initial use | Read-only tools; write tools only after the read-only surface is stable   |

There is no separate MCP process, port, or queue. MCP is a delivery adapter like the HTTP routes.

## 19.2 Transport and routing

- Build the MCP server with `mcp.NewServer` and register tools with `mcp.AddTool`.
- Expose it through the SDK's Streamable HTTP handler and mount it on Fiber with `adaptor.HTTPHandler`.
- Mount it before the SPA fallback so `/api/mcp` is not swallowed by `/*`.
- Prefer the stateless JSON response mode; it avoids sticky sessions and long-lived connections behind a load balancer.
- If a client uses server-sent events, exclude `/api/mcp` from the global Fiber write timeout and disable proxy buffering for that location.

Route order in `internal/routes`:

1. `/healthz`, `/readyz`
2. `/api/mcp`
3. `/api/v1/*`
4. `/build/*`
5. `/*` SPA fallback

Verify the exact constructor names of the pinned SDK version before coding; the SDK is still evolving.

## 19.3 Authentication

MCP uses one static bearer token loaded from `.env`, independent of the JWT access tokens used by the SPA.

- Header: `Authorization: Bearer <MCP_BEARER_TOKEN>`.
- Source: `MCP_BEARER_TOKEN`; it is required when `MCP_ENABLED=true`.
- Comparison: SHA-256 of the presented token compared with the configured hash using `subtle.ConstantTimeCompare`.
- Rotation: update `.env` and restart the process; there is no token table and no console command.
- Actor: `MCP_ACTOR_USER_ID` (optional) binds tool execution to that local user; tools that require an identity fail when it is unset.
- Scope: the static token grants `mcp:read`; `mcp:write` is reserved for mutating tools.
- Rate limit: reuse the PostgreSQL `rate_limits` table keyed by token fingerprint plus client IP.
- Audit: one `mcp_audit_logs` row per tool call.

Never commit `.env`, never log the token, and do not put it in a URL.

Auth wrapper:

```go
func (s *Server) WithAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, ok := bearerToken(r.Header.Get("Authorization"))
        if !ok {
            http.Error(w, "missing bearer token", http.StatusUnauthorized)
            return
        }
        sum := sha256.Sum256([]byte(token))
        if subtle.ConstantTimeCompare(sum[:], s.bearerHash[:]) != 1 {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r.WithContext(auth.WithActor(r.Context(), s.actor)))
    })
}
```

The wrapper runs on the `net/http` side, before the SDK handler. Do not rely on Fiber locals: the adaptor builds a new `http.Request`.

## 19.4 Tool model

- A tool is a typed Go function with input/output structs using `jsonschema` tags.
- Every tool executes as the configured `MCP_ACTOR_USER_ID` through the existing use case and `Authorizer`; never with super-admin bypass.
- Tools return structured content plus a short text summary.
- Limits: maximum page size 50, maximum range 90 days, per-token rate limit.
- No tool accepts raw SQL, a table name, a file path, a shell command, or a URL.
- Errors are returned as tool errors with a safe message; internal details go to logs only.
- Tool handlers are synchronous with a context deadline. Do not spawn a goroutine and return before the result exists.

## 19.5 Initial tools

| Tool           | Scope      | Behavior                                                        |
| -------------- | ---------- | --------------------------------------------------------------- |
| `ping`         | any valid  | Returns server version and current time; verifies connectivity  |
| `whoami`       | `mcp:read` | Returns the configured actor, scopes, and server version        |
| `search_users` | `mcp:read` | Paginated user search by email; max 50 rows                    |

Add domain tools as modules reach parity. The MCP layer must not duplicate business rules; it calls the same use cases as the HTTP handlers.

## 19.6 File layout

```text
internal/app/infrastructure/mcp/
├── server.go      # mcp.NewServer, tool registration, version
├── handler.go     # Streamable HTTP handler + auth wrapper
├── auth.go        # env bearer token + constant-time compare
├── audit.go       # writes mcp_audit_logs
└── tools/
    ├── system.go  # ping, whoami
    └── users.go   # search_users

internal/bootstrap/mcp.go   # builds the handler from Dependencies
internal/routes/mcp.go      # mounts /api/mcp through the adaptor
```

`internal/bootstrap/mcp.go` only wires dependencies and returns an `http.Handler`; it contains no tool logic.

## 19.7 Configuration

| Variable              | Default     | Purpose                                      |
| --------------------- | ----------- | -------------------------------------------- |
| `MCP_ENABLED`         | `true`      | Mount the endpoint                           |
| `MCP_PATH`            | `/api/mcp`  | Route path                                   |
| `MCP_BEARER_TOKEN`    | required    | Static bearer token; read from `.env`        |
| `MCP_ACTOR_USER_ID`   | empty       | Local user ID used as the tool actor          |
| `MCP_SCOPES`          | `mcp:read`  | Granted scopes                               |
| `MCP_MAX_BODY_BYTES`  | `1048576`   | Request size limit                           |
| `MCP_ALLOWED_ORIGINS` | empty       | Comma-separated `Origin` allowlist           |
| `MCP_SESSION_TIMEOUT` | `60s`       | Idle timeout for a session                   |
| `MCP_STATELESS`       | `true`      | Disable server-side session state            |

## 19.8 Client configuration

The exact JSON shape depends on the client; the values are the same everywhere.

```json
{
  "mcpServers": {
    "alma": {
      "type": "streamable-http",
      "url": "https://<domain>/api/mcp",
      "headers": {
        "Authorization": "Bearer <MCP_BEARER_TOKEN>"
      }
    }
  }
}
```

- A client must support remote HTTP transport and custom headers.
- A client that cannot send custom headers needs a local proxy that injects the `Authorization` header.
- Do not embed the token in a URL query string.
- Restart the client after changing `MCP_BEARER_TOKEN` on the server.

## 19.9 Testing

- Unit: env token loading, constant-time comparison, origin checks.
- Integration (testcontainers): set `MCP_BEARER_TOKEN`, connect with the SDK client, call `whoami` and `search_users`, assert the audit row and rate limit.
- Negative: missing token, malformed token, wrong token, missing actor where required, oversized body, disallowed origin.
- Contract: tool input/output JSON shape is stable; golden files for the tool schemas.
- Startup: `/api/mcp` is not registered when `MCP_ENABLED=false`, and startup fails when the token is empty.

## 19.10 Deployment

- The endpoint shares the API domain and port; only one reverse-proxy location is needed.
- TLS is terminated at the reverse proxy; the app listens on plain HTTP internally.
- Require HTTPS; reject plain HTTP at the proxy.
- For SSE clients, disable response buffering and raise the read timeout for `/api/mcp` only.
- Set `MCP_ALLOWED_ORIGINS` once the production domain is known.
- Rotate the token by updating `MCP_BEARER_TOKEN` in `.env` and restarting the process; there is no OAuth refresh flow.

The production domain is not fixed yet. All examples use `<domain>`; no hostname is hardcoded in code or docs.

## 19.11 Adding a domain tool

1. Create the tool function in `internal/app/infrastructure/mcp/tools/<domain>.go`.
2. Reuse the use case and authorization checks from `internal/app/Application/<domain>`.
3. Register it in `server.go` with the required scope.
4. Add an audit entry through the shared wrapper; do not write audit rows inside the tool.
5. Add integration and negative tests.
6. Document the tool in the table in 19.5.
