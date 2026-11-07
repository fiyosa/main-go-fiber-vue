package user

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewRepositoryReturnsNilWithoutDatabase(t *testing.T) {
	var db *gorm.DB
	if repo := NewRepository(db); repo != nil {
		t.Fatal("NewRepository(nil) must return nil so callers can degrade gracefully")
	}
}
