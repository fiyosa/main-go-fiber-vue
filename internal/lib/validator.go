package lib

import (
	"encoding/json"
	"go-fiber-svelte/internal/config"
	"go-fiber-svelte/internal/lang"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validate validate

type validate struct{}

type ValidationError struct {
	Message string
	Errors  map[string]string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (validate) Check(c *fiber.Ctx, input any) error {
	if err := c.BodyParser(input); err != nil {
		return generateError(err)
	}

	if err := config.Validate.Struct(input); err != nil {
		return generateError(err)
	}

	return nil
}

func generateError(err error) *ValidationError {
	newErrors := map[string]string{}
	msg := "Invalid data"

	switch v := err.(type) {
	case *json.UnmarshalTypeError:
		field := strings.ToLower(v.Field)
		newErrors[field] = "Json binding error: " + field + " type error"

	case validator.ValidationErrors:
		for _, e := range v {
			field := strings.ToLower(e.Field())
			newErrors[field] = strings.ToLower(e.Translate(config.Translator))
		}

	default:
		if v != nil {
			msg = v.Error()
		} else {
			msg = lang.T.Get().SOMETHING_WENT_WRONG
		}
	}

	return &ValidationError{Message: msg, Errors: newErrors}
}
