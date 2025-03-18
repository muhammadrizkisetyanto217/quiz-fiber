package main

import (
	"log"
	"quiz-fiber/internals/configs"
	"quiz-fiber/internals/database"
	"quiz-fiber/internals/middleware"
	"quiz-fiber/internals/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	configs.LoadEnv()

	// ✅ Koneksi ke database
	database.ConnectDB()

	app := fiber.New()

	// ✅ Setup Middleware
	middleware.SetupMiddleware(app)

	// ✅ Setup Routes
	routes.SetupRoutes(app, database.DB)

	log.Fatal(app.Listen(":3000"))
}
