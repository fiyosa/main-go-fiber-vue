package entity

import "time"

// MCPAuditLog is one recorded MCP tool call. Parameters are stored as raw JSON
// text so the row can be inspected without a schema migration per tool.
type MCPAuditLog struct {
	ID            uint64    `gorm:"column:id;primaryKey"`
	ToolName      string    `gorm:"column:tool_name;size:128;not null"`
	Actor         *string   `gorm:"column:actor;size:255"`
	Parameters    *string   `gorm:"column:parameters;type:jsonb"`
	ResultSummary *string   `gorm:"column:result_summary;type:text"`
	Error         *string   `gorm:"column:error;type:text"`
	IPAddress     *string   `gorm:"column:ip_address;size:64"`
	UserAgent     *string   `gorm:"column:user_agent;size:512"`
	DurationMS    int64     `gorm:"column:duration_ms;not null;default:0"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
}

// TableName maps the entity to its migration-managed table.
func (MCPAuditLog) TableName() string {
	return "mcp_audit_logs"
}
