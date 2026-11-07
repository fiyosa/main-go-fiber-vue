package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"gorm.io/gorm"
)

func TestNewGORMAuditRecorderWithoutDatabase(t *testing.T) {
	var db *gorm.DB
	recorder := NewGORMAuditRecorder(db)
	if recorder == nil {
		t.Fatal("NewGORMAuditRecorder(nil) returned nil; want a no-op recorder")
	}

	if err := recorder.RecordToolCall(context.Background(), AuditEntry{Tool: "ping"}); err != nil {
		t.Fatalf("no-op recorder returned error: %v", err)
	}
}

func TestGORMAuditRecorderNilReceiver(t *testing.T) {
	var recorder *gormAuditRecorder
	if err := recorder.RecordToolCall(context.Background(), AuditEntry{Tool: "ping"}); err != nil {
		t.Fatalf("nil recorder returned error: %v", err)
	}
}

func TestNullableHelpers(t *testing.T) {
	if got := nullableString(""); got != nil {
		t.Errorf("nullableString(\"\") = %v, want nil", got)
	}
	if got := nullableString("value"); got == nil || *got != "value" {
		t.Errorf("nullableString(\"value\") = %v, want pointer to value", got)
	}
	if got := nullableJSON(nil); got != nil {
		t.Errorf("nullableJSON(nil) = %v, want nil", got)
	}

	raw := json.RawMessage(`{"email":"a@b.c"}`)
	got := nullableJSON(raw)
	if got == nil || *got != string(raw) {
		t.Errorf("nullableJSON(raw) = %v, want %q", got, string(raw))
	}
}
