package auth

import (
	"context"
	"strings"
	"time"

	"go-fiber-vue/internal/app/application/core/repository/user"
	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	platformauth "go-fiber-vue/internal/platform/auth"
)

// RegisterResult holds the data returned after a successful registration.
type RegisterResult struct {
	ID    uint64
	Email string
}

// RegisterRepository is the single interface for POST /api/v1/auth/register.
// One file per API — no use-case layer.
type RegisterRepository interface {
	Register(ctx context.Context, email, password string) (*RegisterResult, error)
}

type gormRegisterRepository struct {
	users user.Repository
}

// NewRegisterRepository wires the register repository with its GORM user repository.
func NewRegisterRepository(users user.Repository) RegisterRepository {
	return &gormRegisterRepository{users: users}
}

// Register implements RegisterRepository.
func (r *gormRegisterRepository) Register(ctx context.Context, email, password string) (*RegisterResult, error) {
	if r.users == nil {
		return nil, exceptions.ServiceUnavailable(nil)
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	existing, err := r.users.FindByEmail(ctx, normalizedEmail)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, exceptions.New("auth.email_taken", 409, "lang.email_already_taken", nil)
	}

	hashedPassword, err := platformauth.HashPassword(password)
	if err != nil {
		return nil, exceptions.Internal(err)
	}

	now := time.Now().UTC()
	account := &entity.User{
		Email:     normalizedEmail,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := r.users.Create(ctx, account, "user"); err != nil {
		return nil, err
	}

	return &RegisterResult{
		ID:    account.ID,
		Email: account.Email,
	}, nil
}
