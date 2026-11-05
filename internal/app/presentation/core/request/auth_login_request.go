package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// validate is created once and reused across every request type.
var validate = validator.New()

// LoginRequest carries the POST /api/v1/auth/login payload.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Bind decodes and validates the JSON body. On failure it returns a typed 422
// exception carrying translated per-field messages; the controller returns it
// to Fiber's global error handler unchanged. On success the struct is populated.
func (r *LoginRequest) Bind(fctx fiber.Ctx) error {
	if err := fctx.Bind().JSON(r); err != nil {
		return exceptions.Validation(validationFields(fctx, err))
	}
	if err := validate.Struct(r); err != nil {
		return exceptions.Validation(validationFields(fctx, err))
	}
	return nil
}
