package response

// RegisterResource is the data returned by POST /api/v1/auth/register.
type RegisterResource struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
}
