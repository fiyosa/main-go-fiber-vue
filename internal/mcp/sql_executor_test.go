package mcp

import (
	"testing"
	"time"
)

func TestIsRowReturning(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want bool
	}{
		{name: "select", sql: "SELECT 1", want: true},
		{name: "select with leading comment", sql: "-- note\nSELECT 1", want: true},
		{name: "block comment", sql: "/* note */ select 1", want: true},
		{name: "with", sql: "WITH x AS (SELECT 1) SELECT * FROM x", want: true},
		{name: "show", sql: "SHOW server_version", want: true},
		{name: "explain", sql: "EXPLAIN SELECT 1", want: true},
		{name: "table", sql: "TABLE users", want: true},
		{name: "insert", sql: "INSERT INTO users (email) VALUES ('a')", want: false},
		{name: "insert returning", sql: "INSERT INTO users (email) VALUES ('a') RETURNING id", want: true},
		{name: "update", sql: "UPDATE users SET email = 'a'", want: false},
		{name: "update returning", sql: "UPDATE users SET email = 'a' RETURNING id", want: true},
		{name: "delete", sql: "DELETE FROM users", want: false},
		{name: "create table", sql: "CREATE TABLE t (id int)", want: false},
		{name: "drop table", sql: "DROP TABLE t", want: false},
		{name: "alter table", sql: "ALTER TABLE users ADD COLUMN x int", want: false},
		{name: "empty", sql: "   ", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isRowReturning(tc.sql); got != tc.want {
				t.Fatalf("isRowReturning(%q) = %v, want %v", tc.sql, got, tc.want)
			}
		})
	}
}

func TestNormalizeValue(t *testing.T) {
	t.Run("bytes become strings", func(t *testing.T) {
		if got := normalizeValue([]byte("hello")); got != "hello" {
			t.Fatalf("normalizeValue([]byte) = %v, want hello", got)
		}
	})

	t.Run("time becomes RFC3339Nano", func(t *testing.T) {
		stamp := time.Date(2026, 11, 7, 12, 30, 0, 0, time.UTC)
		if got := normalizeValue(stamp); got != stamp.Format(time.RFC3339Nano) {
			t.Fatalf("normalizeValue(time) = %v", got)
		}
	})

	t.Run("scalars pass through", func(t *testing.T) {
		values := []any{int64(1), float64(1.5), true, nil, "text"}
		for _, value := range values {
			if got := normalizeValue(value); got != value {
				t.Fatalf("normalizeValue(%v) = %v, want unchanged", value, got)
			}
		}
	})
}
