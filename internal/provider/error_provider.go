package provider

import (
	"errors"

	"go-fiber-svelte/internal/helper"
	"go-fiber-svelte/internal/lib"

	"github.com/gofiber/fiber/v2"
)

func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var ve *lib.ValidationError
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(helper.Res.Error(ve.Message, ve.Errors))
		}
		var fe *fiber.Error
		if errors.As(err, &fe) {
			return c.Status(fe.Code).JSON(helper.Res.Error(fe.Message, nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(helper.Res.Error(err.Error(), nil))
	}
}
