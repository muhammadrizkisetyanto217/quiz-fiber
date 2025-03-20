package routes

import (
	// Add this line.
	userRoute "quiz-fiber/internals/features/user/user/route"
	authRoute "quiz-fiber/internals/features/user/auth/route"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register routes
func SetupRoutes(app *fiber.App, db *gorm.DB) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Fiber & Supabase PostgreSQL connected successfully 🚀")
	})

	userRoute.UserRoutes(app, db)
	authRoute.AuthRoutes(app, db)

}
