package response

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/platform/lang"
)

// Envelope is the canonical JSON response shape. Message is always present.
type Envelope struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// Success writes a 2xx envelope using a translation key.
func Success(c fiber.Ctx, data any, key string, params map[string]string, code int) error {
	return c.Status(code).JSON(Envelope{
		Success: true,
		Data:    data,
		Message: translate(c, key, params),
	})
}

// Error writes an error envelope using a translation key.
func Error(c fiber.Ctx, key string, params map[string]string, code int) error {
	return c.Status(code).JSON(Envelope{
		Success: false,
		Message: translate(c, key, params),
	})
}

// ErrorDebug writes an error envelope with debug detail. Only used in debug mode.
func ErrorDebug(c fiber.Ctx, key string, code int, detail string) error {
	return c.Status(code).JSON(struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Debug   string `json:"debug,omitempty"`
	}{
		Success: false,
		Message: translate(c, key, nil),
		Debug:   detail,
	})
}

// ValidationError writes a 422 envelope with translated field messages.
func ValidationError(c fiber.Ctx, fields map[string][]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(struct {
		Success bool                `json:"success"`
		Message string              `json:"message"`
		Errors  map[string][]string `json:"errors"`
	}{
		Success: false,
		Message: translate(c, "common.validation_failed", nil),
		Errors:  fields,
	})
}

// LoginSuccess builds a login response with a translated message.
func LoginSuccess(c fiber.Ctx, token, refreshToken string, expiresIn int64) LoginResponse {
	return LoginResponse{
		Success:      true,
		Message:      translate(c, "lang.logged_in_successfully", nil),
		Token:        token,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
	}
}

func translate(c fiber.Ctx, key string, params map[string]string) string {
	translator := lang.From(c)
	if translator == nil {
		return key
	}
	return translator.T(lang.Locale(c), key, params)
}
