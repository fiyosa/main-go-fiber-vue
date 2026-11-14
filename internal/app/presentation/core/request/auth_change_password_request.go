package request

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	"go-fiber-vue/internal/platform/lang"
)

// ChangePasswordRequest carries the POST /api/v1/auth/change-password payload.
type ChangePasswordRequest struct {
	NewPassword     string `json:"new-password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm-password" validate:"required"`
}

// Bind decodes and validates the JSON body. The validator reports lowercased
// Go field names (newpassword, confirmpassword), so field keys are remapped to
// the JSON names before the response is written.
func (r *ChangePasswordRequest) Bind(fctx fiber.Ctx) error {
	if err := fctx.Bind().JSON(r); err != nil {
		return exceptions.Validation(changePasswordFields(validationFields(fctx, err)))
	}
	if err := validate.Struct(r); err != nil {
		return exceptions.Validation(changePasswordFields(validationFields(fctx, err)))
	}
	if r.NewPassword != r.ConfirmPassword {
		translator := lang.From(fctx)
		locale := lang.Locale(fctx)
		message := "lang.password_mismatch"
		if translator != nil {
			message = translator.T(locale, "lang.password_mismatch", nil)
		}
		return exceptions.Validation(map[string][]string{
			"confirm-password": {message},
		})
	}
	return nil
}

// changePasswordFields renames the lowercased validator field keys to the JSON
// names used by this endpoint.
func changePasswordFields(fields map[string][]string) map[string][]string {
	renames := map[string]string{
		"newpassword":     "new-password",
		"confirmpassword": "confirm-password",
	}
	for from, to := range renames {
		messages, ok := fields[from]
		if !ok {
			continue
		}
		delete(fields, from)
		fields[to] = append(fields[to], messages...)
	}
	return fields
}
