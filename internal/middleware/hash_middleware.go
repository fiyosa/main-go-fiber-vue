package middleware

import (
	"go-fiber-svelte/internal/lib"

	"github.com/gofiber/fiber/v2"
)

func HashMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		for _, param := range c.Route().Params {
			val := c.Params(param)
			if val == "" {
				continue
			}
			decoded, err := lib.Hash.DecodeId(val)
			if err == nil {
				c.Locals("decoded_"+param, decoded)
			}
		}
		return c.Next()
	}
}
