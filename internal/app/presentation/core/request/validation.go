package request

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/platform/lang"
)

// validationFields converts validator errors into translated per-field messages.
// Attribute names use a lang.<field> label when one exists.
func validationFields(c fiber.Ctx, err error) map[string][]string {
	translator := lang.From(c)
	locale := lang.Locale(c)

	fields := make(map[string][]string)
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			label := fieldLabel(translator, locale, field)
			message := validationMessage(translator, locale, fieldError, label)
			fields[field] = append(fields[field], message)
		}
	}
	if len(fields) == 0 {
		fields["request"] = append(fields["request"], commonMessage(translator, locale, "common.bad_request"))
	}
	return fields
}

func fieldLabel(translator *lang.Translator, locale, field string) string {
	key := "lang." + field
	if translator != nil && translator.Has(locale, key) {
		return translator.T(locale, key, nil)
	}
	return field
}

func validationMessage(translator *lang.Translator, locale string, fieldError validator.FieldError, label string) string {
	key := "validation." + fieldError.Tag()
	if fieldError.Kind() == reflect.String && translator != nil {
		if nested := key + ".string"; translator.Has(locale, nested) {
			key = nested
		}
	}
	if translator == nil || !translator.Has(locale, key) {
		return commonMessage(translator, locale, "common.validation_failed")
	}
	return translator.T(locale, key, map[string]string{
		"attribute": label,
		"min":       fieldError.Param(),
		"max":       fieldError.Param(),
		"values":    fieldError.Param(),
	})
}

func commonMessage(translator *lang.Translator, locale, key string) string {
	if translator == nil {
		return key
	}
	return translator.T(locale, key, nil)
}
