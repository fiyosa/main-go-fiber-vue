package config

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     APP_WEB_URL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}
}
