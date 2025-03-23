package main

import (
	"log"
	"os"
	"quiz-fiber/internals/configs"
	"quiz-fiber/internals/database"
	"quiz-fiber/internals/middleware"
	"quiz-fiber/internals/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	configs.LoadEnv()

	// ✅ Cek JWT_SECRET wajib ada
	jwtSecret := configs.GetEnv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET is not set! Please check your .env or Railway environment settings.")
	}

	log.Println("✅ JWT_SECRET:", configs.JWTSecret)
	log.Println("🧪 os.Getenv JWT_SECRET:", os.Getenv("JWT_SECRET"))

	// ✅ Koneksi ke database
	database.ConnectDB()

	app := fiber.New()

	// ✅ Setup Middleware
	middleware.SetupMiddleware(app)

	// ✅ Setup Routes
	routes.SetupRoutes(app, database.DB)

	log.Fatal(app.Listen(":8080"))
}
