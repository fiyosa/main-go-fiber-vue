package docs

import (
	"os"

	"github.com/gofiber/fiber/v3"
)

func readFile(filename string) ([]byte, error) {
	paths := []string{
		filename,
		"../" + filename,
		"../../" + filename,
		"../../../" + filename,
	}
	for _, p := range paths {
		if data, err := os.ReadFile(p); err == nil {
			return data, nil
		}
	}
	return nil, os.ErrNotExist
}

// Handler serves the Scalar API reference HTML UI from public/docs.html.
func Handler(c fiber.Ctx) error {
	body, err := readFile("public/docs/scalar.html")
	if err != nil {
		return fiber.ErrNotFound
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	c.Set(fiber.HeaderCacheControl, "no-cache")
	return c.Send(body)
}

// SpecHandler serves the public OpenAPI 3.1 JSON document from public/docs/openapi.json.
func SpecHandler(c fiber.Ctx) error {
	body, err := readFile("public/docs/openapi.json")
	if err != nil {
		return fiber.ErrNotFound
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.Set(fiber.HeaderCacheControl, "private, max-age=3600")
	return c.Send(body)
}
