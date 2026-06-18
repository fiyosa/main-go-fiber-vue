package policy_request

type PermissionStoreRequest struct {
	Name  string `json:"name" validate:"required"`
	Notes string `json:"notes"`
}
