package response

// LoginResponse is the JSON body returned by POST /api/v1/auth/login.
type LoginResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Token        string `json:"token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
