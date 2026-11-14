package auth

import (
	"context"
	"log/slog"

	"go-fiber-vue/internal/app/application/core/repository/user"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	platformauth "go-fiber-vue/internal/platform/auth"
)

// ChangePasswordRepository is the single interface for
// POST /api/v1/auth/change-password. One file per API — no use-case layer.
type ChangePasswordRepository interface {
	ChangePassword(ctx context.Context, userID uint64, newPassword string) error
}

// gormChangePasswordRepository is the GORM implementation of
// ChangePasswordRepository. It updates the credential and revokes every active
// session of the same user.
type gormChangePasswordRepository struct {
	users user.Repository
	auths Repository
}

// NewChangePasswordRepository wires the change-password repository with its
// user and auth session repositories.
func NewChangePasswordRepository(users user.Repository, auths Repository) ChangePasswordRepository {
	return &gormChangePasswordRepository{users: users, auths: auths}
}

// ChangePassword implements ChangePasswordRepository.
func (r *gormChangePasswordRepository) ChangePassword(ctx context.Context, userID uint64, newPassword string) error {
	if r.users == nil || r.auths == nil {
		return exceptions.ServiceUnavailable(nil)
	}

	account, err := r.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if account == nil {
		return exceptions.NotFound(nil)
	}

	hashedPassword, err := platformauth.HashPassword(newPassword)
	if err != nil {
		return exceptions.Internal(err)
	}

	if err := r.users.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return err
	}

	// The password is already changed, so a failed revocation must not turn the
	// response into an error: the client is still told to log in again.
	if err := r.auths.RevokeAllForUser(ctx, userID); err != nil {
		slog.Warn("change password: session revoke failed", "user_id", userID, "error", err)
	}

	return nil
}
