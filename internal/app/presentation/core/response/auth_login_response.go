package response

// LoginResponse is the JSON body returned by POST /api/v1/auth/login.
// Credentials are delivered through HTTP cookies, never in the body.
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
