package mcp

import (
	"context"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// gormAuditRecorder persists audit entries in mcp_audit_logs.
type gormAuditRecorder struct {
	db *gorm.DB
}

// NewGORMAuditRecorder returns a recorder backed by the given database. It
// degrades to a no-op when db is nil so the MCP server still runs without a
// database.
func NewGORMAuditRecorder(db *gorm.DB) AuditRecorder {
	if db == nil {
		return NopAuditRecorder{}
	}
	return &gormAuditRecorder{db: db}
}

// RecordToolCall implements AuditRecorder.
func (r *gormAuditRecorder) RecordToolCall(ctx context.Context, entry AuditEntry) error {
	if r == nil || r.db == nil {
		return nil
	}

	row := &entity.MCPAuditLog{
		ToolName:   entry.Tool,
		Actor:      nullableString(entry.Actor),
		Parameters: nullableJSON(entry.Arguments),
		Error:      nullableString(entry.Error),
		DurationMS: entry.Duration.Milliseconds(),
		CreatedAt:  entry.OccurredAt.UTC(),
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return exceptions.FromDB(err)
	}
	return nil
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func nullableJSON(raw []byte) *string {
	if len(raw) == 0 {
		return nil
	}
	value := string(raw)
	return &value
}
