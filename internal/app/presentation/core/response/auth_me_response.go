package response

// MeResource is the identity payload returned by GET /api/v1/auth/me.
// The id is the caller's numeric user id; email comes from the access token;
// roles are read from the user_has_roles join.
type MeResource struct {
	ID    uint64   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}
