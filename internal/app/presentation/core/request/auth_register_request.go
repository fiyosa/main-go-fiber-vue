package request

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
)

// RegisterRequest carries the POST /api/v1/auth/register payload.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// Bind decodes and validates the JSON body.
func (r *RegisterRequest) Bind(fctx fiber.Ctx) error {
	if err := fctx.Bind().JSON(r); err != nil {
		return exceptions.Validation(validationFields(fctx, err))
	}
	if err := validate.Struct(r); err != nil {
		return exceptions.Validation(validationFields(fctx, err))
	}
	return nil
}
