package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupMiddleware mengatur middleware untuk aplikasi Fiber
func SetupMiddleware(app *fiber.App) {
	app.Use(logger.New()) // Logging request

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Mengizinkan semua origin (bisa diatur lebih ketat)
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
}
