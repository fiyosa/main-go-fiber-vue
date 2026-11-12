package auth

import (
	"context"
)

// LogoutRepository is the single interface for POST /api/v1/auth/logout.
// One file per API — no use-case layer.
type LogoutRepository interface {
	// Logout revokes the session bound to the access token. An empty token is
	// a no-op so the endpoint stays idempotent.
	Logout(ctx context.Context, accessToken string) error
}

// gormLogoutRepository is the GORM implementation of LogoutRepository.
type gormLogoutRepository struct {
	auths Repository
}

// NewLogoutRepository wires the logout repository with its GORM implementation.
// A nil session repository is allowed: logout then only clears the client
// cookies so console commands can still build the route table.
func NewLogoutRepository(auths Repository) LogoutRepository {
	return &gormLogoutRepository{auths: auths}
}

// Logout implements LogoutRepository.
func (r *gormLogoutRepository) Logout(ctx context.Context, accessToken string) error {
	if accessToken == "" || r.auths == nil {
		return nil
	}
	return r.auths.RevokeSession(ctx, accessToken)
}
