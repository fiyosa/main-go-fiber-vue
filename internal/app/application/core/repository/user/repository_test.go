package user

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

func TestNormalizePaging(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		wantPage int
		wantSize int
	}{
		{name: "zero values use defaults", page: 0, pageSize: 0, wantPage: 1, wantSize: defaultPageSize},
		{name: "page size is capped", page: 2, pageSize: 500, wantPage: 2, wantSize: maxPageSize},
		{name: "negative page falls back", page: -3, pageSize: 10, wantPage: 1, wantSize: 10},
		{name: "negative size falls back", page: 1, pageSize: -1, wantPage: 1, wantSize: defaultPageSize},
		{name: "values pass through", page: 4, pageSize: 25, wantPage: 4, wantSize: 25},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			page, size := normalizePaging(tc.page, tc.pageSize)
			if page != tc.wantPage || size != tc.wantSize {
				t.Fatalf("normalizePaging(%d, %d) = (%d, %d), want (%d, %d)",
					tc.page, tc.pageSize, page, size, tc.wantPage, tc.wantSize)
			}
		})
	}
}

func TestSearchUsersWithoutDatabase(t *testing.T) {
	t.Run("nil db handle", func(t *testing.T) {
		repo := &gormRepository{db: nil}
		users, total, err := repo.SearchUsers(context.Background(), "someone@example.com", 1, 10)
		if err != nil {
			t.Fatalf("SearchUsers() error = %v, want nil", err)
		}
		if len(users) != 0 || total != 0 {
			t.Fatalf("SearchUsers() = %d users, total %d; want empty", len(users), total)
		}
	})

	t.Run("nil receiver", func(t *testing.T) {
		var repo *gormRepository
		users, total, err := repo.SearchUsers(context.Background(), "", 1, 10)
		if err != nil {
			t.Fatalf("SearchUsers() error = %v, want nil", err)
		}
		if len(users) != 0 || total != 0 {
			t.Fatalf("SearchUsers() = %d users, total %d; want empty", len(users), total)
		}
	})
}

func TestNewRepositoryReturnsNilWithoutDatabase(t *testing.T) {
	var db *gorm.DB
	if repo := NewRepository(db); repo != nil {
		t.Fatal("NewRepository(nil) must return nil so callers can degrade gracefully")
	}
}
