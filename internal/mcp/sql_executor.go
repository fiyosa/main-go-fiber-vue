package mcp

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"go-fiber-vue/internal/mcp/tools"
)

// sqlDBExecutor runs raw statements through database/sql. It deliberately keeps
// no allowlist: the run_sql tool is meant for full database access.
type sqlDBExecutor struct {
	db *sql.DB
}

// NewSQLExecutor wraps db for the run_sql tool. It returns nil when db is nil so
// the tool reports "database unavailable" instead of panicking.
func NewSQLExecutor(db *sql.DB) tools.DBExecutor {
	if db == nil {
		return nil
	}
	return &sqlDBExecutor{db: db}
}

// Execute implements tools.DBExecutor.
func (e *sqlDBExecutor) Execute(ctx context.Context, statement string, maxRows int) (tools.SQLResult, error) {
	if isRowReturning(statement) {
		return e.query(ctx, statement, maxRows)
	}

	res, err := e.db.ExecContext(ctx, statement)
	if err != nil {
		return tools.SQLResult{}, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		// Some PostgreSQL statements do not report an affected-row count.
		affected = 0
	}
	return tools.SQLResult{RowsAffected: affected}, nil
}

func (e *sqlDBExecutor) query(ctx context.Context, statement string, maxRows int) (tools.SQLResult, error) {
	rows, err := e.db.QueryContext(ctx, statement)
	if err != nil {
		return tools.SQLResult{}, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return tools.SQLResult{}, err
	}

	result := tools.SQLResult{
		Columns: columns,
		Rows:    make([]map[string]any, 0, 16),
	}

	for rows.Next() {
		if len(result.Rows) >= maxRows {
			result.Truncated = true
			break
		}

		values := make([]any, len(columns))
		targets := make([]any, len(columns))
		for i := range values {
			targets[i] = &values[i]
		}
		if err := rows.Scan(targets...); err != nil {
			return tools.SQLResult{}, err
		}

		record := make(map[string]any, len(columns))
		for i, column := range columns {
			record[column] = normalizeValue(values[i])
		}
		result.Rows = append(result.Rows, record)
	}
	if err := rows.Err(); err != nil {
		return tools.SQLResult{}, err
	}

	result.RowCount = len(result.Rows)
	return result, nil
}

// isRowReturning reports whether a statement should be executed as a query.
// Statements with a RETURNING clause are treated as queries even when they are
// INSERT/UPDATE/DELETE.
func isRowReturning(statement string) bool {
	cleaned := stripLeadingComments(statement)
	if cleaned == "" {
		return false
	}

	keyword := strings.ToUpper(firstWord(cleaned))
	switch keyword {
	case "SELECT", "WITH", "SHOW", "EXPLAIN", "TABLE", "VALUES":
		return true
	}

	return strings.Contains(" "+strings.ToUpper(statement)+" ", " RETURNING ")
}

// stripLeadingComments removes whitespace, line comments, and block comments
// from the start of a statement so the first keyword can be read.
func stripLeadingComments(statement string) string {
	s := strings.TrimSpace(statement)
	for {
		switch {
		case strings.HasPrefix(s, "--"):
			index := strings.IndexByte(s, '\n')
			if index < 0 {
				return ""
			}
			s = strings.TrimSpace(s[index+1:])
		case strings.HasPrefix(s, "/*"):
			index := strings.Index(s, "*/")
			if index < 0 {
				return ""
			}
			s = strings.TrimSpace(s[index+2:])
		default:
			return s
		}
	}
}

func firstWord(statement string) string {
	for i, r := range statement {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '(' {
			return statement[:i]
		}
	}
	return statement
}

// normalizeValue converts database/sql scan values into JSON-friendly values.
func normalizeValue(value any) any {
	switch typed := value.(type) {
	case []byte:
		return string(typed)
	case time.Time:
		return typed.UTC().Format(time.RFC3339Nano)
	default:
		return value
	}
}
