package authorization

import (
	"context"
	"strings"

	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// Resolver is the read side of the authorization tables.
type Resolver interface {
	HasPermission(ctx context.Context, userID uint64, permission string) (bool, error)
}

// Authorizer answers permission questions. Access is decided only by the
// permission rows attached to the caller's roles: there is no admin bypass.
type Authorizer struct {
	resolver Resolver
}

// New returns an authorizer. A nil resolver is allowed for console commands
// that build the route table without a database.
func New(resolver Resolver) *Authorizer {
	return &Authorizer{resolver: resolver}
}

// Can reports whether the user may perform the permission. It returns a 503
// exception (not a panic) when the resolver is unavailable.
func (a *Authorizer) Can(ctx context.Context, userID uint64, permission string) (bool, error) {
	if a == nil || a.resolver == nil {
		return false, exceptions.ServiceUnavailable(nil)
	}
	permission = strings.TrimSpace(permission)
	if permission == "" {
		return false, exceptions.Internal(nil)
	}

	return a.resolver.HasPermission(ctx, userID, permission)
}
