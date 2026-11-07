package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	// DefaultSQLMaxRows bounds how many rows a SELECT may return per call.
	DefaultSQLMaxRows = 1000
	// DefaultSQLTimeout bounds how long a single statement may run.
	DefaultSQLTimeout = 30 * time.Second
)

// SQLResult is the structured output of the run_sql tool.
type SQLResult struct {
	Columns      []string         `json:"columns,omitempty" jsonschema:"column names for row-returning statements"`
	Rows         []map[string]any `json:"rows,omitempty" jsonschema:"returned rows keyed by column name"`
	RowCount     int              `json:"rowCount" jsonschema:"number of rows returned"`
	RowsAffected int64            `json:"rowsAffected" jsonschema:"rows affected by a non-returning statement"`
	Truncated    bool             `json:"truncated" jsonschema:"true when the result was cut at the row limit"`
	DurationMS   int64            `json:"durationMs" jsonschema:"execution time in milliseconds"`
}

// DBExecutor runs one raw SQL statement. Implementations decide how to split
// row-returning statements from statements that only report a row count.
type DBExecutor interface {
	Execute(ctx context.Context, statement string, maxRows int) (SQLResult, error)
}

// SQLConfig configures the run_sql tool.
type SQLConfig struct {
	Executor DBExecutor
	MaxRows  int
	Timeout  time.Duration
}

type sqlInput struct {
	SQL string `json:"sql" jsonschema:"raw PostgreSQL statement to execute; any statement is allowed"`
}

// RegisterSQL registers the run_sql tool. The tool is intentionally unrestricted
// at the SQL level: the bearer token owner can run DDL and DML. It only enforces
// a non-empty statement, an execution timeout, and a row cap.
func RegisterSQL(server *sdkmcp.Server, cfg SQLConfig) {
	maxRows := cfg.MaxRows
	if maxRows <= 0 {
		maxRows = DefaultSQLMaxRows
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = DefaultSQLTimeout
	}

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "run_sql",
		Description: "Executes a raw SQL statement against the application PostgreSQL database and returns rows or affected-row counts.",
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, in sqlInput) (*sdkmcp.CallToolResult, SQLResult, error) {
		statement := strings.TrimSpace(in.SQL)
		if statement == "" {
			return nil, SQLResult{}, fmt.Errorf("sql is required")
		}
		if cfg.Executor == nil {
			return nil, SQLResult{}, fmt.Errorf("database unavailable")
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		started := time.Now()
		result, err := cfg.Executor.Execute(ctx, statement, maxRows)
		if err != nil {
			return nil, SQLResult{}, fmt.Errorf("execute sql: %w", err)
		}
		result.DurationMS = time.Since(started).Milliseconds()
		return nil, result, nil
	})
}
