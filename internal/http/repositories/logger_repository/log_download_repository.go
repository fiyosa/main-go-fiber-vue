package logger_repository

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func LogDownloadRepository(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" || strings.Contains(filename, "..") {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid filename")
	}

	return c.SendFile("./logs/" + filename)
}
