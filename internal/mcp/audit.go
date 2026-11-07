package mcp

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// AuditEntry is one tool call recorded for compliance. Arguments are stored
// verbatim minus nothing: callers must not pass secrets through tool input.
type AuditEntry struct {
	Tool       string
	Actor      string
	Arguments  json.RawMessage
	Success    bool
	Error      string
	Duration   time.Duration
	OccurredAt time.Time
}

// AuditRecorder persists audit entries. A nil recorder means auditing is
// disabled; callers keep working without a database.
type AuditRecorder interface {
	RecordToolCall(ctx context.Context, entry AuditEntry) error
}

// NopAuditRecorder discards every entry. Use it when no database is available.
type NopAuditRecorder struct{}

// RecordToolCall implements AuditRecorder.
func (NopAuditRecorder) RecordToolCall(context.Context, AuditEntry) error { return nil }

// AuditMiddleware records exactly one entry per tools/call. Audit failures are
// logged but never turn a successful tool call into an error.
func AuditMiddleware(recorder AuditRecorder) sdkmcp.Middleware {
	return func(next sdkmcp.MethodHandler) sdkmcp.MethodHandler {
		return func(ctx context.Context, method string, req sdkmcp.Request) (sdkmcp.Result, error) {
			if recorder == nil || method != "tools/call" {
				return next(ctx, method, req)
			}

			occurredAt := time.Now()
			entry := AuditEntry{OccurredAt: occurredAt}
			if actor, ok := ActorFromContext(ctx); ok {
				entry.Actor = actor
			}
			if params, ok := req.GetParams().(*sdkmcp.CallToolParamsRaw); ok {
				entry.Tool = params.Name
				entry.Arguments = params.Arguments
			}

			result, err := next(ctx, method, req)

			entry.Duration = time.Since(occurredAt)
			switch {
			case err != nil:
				entry.Error = err.Error()
			case isToolError(result):
				entry.Error = "tool returned an error"
			default:
				entry.Success = true
			}

			if recordErr := recorder.RecordToolCall(ctx, entry); recordErr != nil {
				slog.Warn("mcp audit write failed", "tool", entry.Tool, "error", recordErr)
			}
			return result, err
		}
	}
}

func isToolError(result sdkmcp.Result) bool {
	callResult, ok := result.(*sdkmcp.CallToolResult)
	return ok && callResult.IsError
}
