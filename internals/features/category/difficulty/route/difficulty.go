package route

import (
	difficultyController "quiz-fiber/internals/features/category/difficulty/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register category-related routes
func CategoryRoutes(app *fiber.App, db *gorm.DB) {

	// 🔥 Proteksi seluruh kategori API dengan Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Difficulty Routes
	difficultyController := difficultyController.NewDifficultyController(db)
	difficultyRoutes := api.Group("/difficulties")
	difficultyRoutes.Get("/", difficultyController.GetDifficulties)
	difficultyRoutes.Get("/:id", difficultyController.GetDifficulty)
	difficultyRoutes.Post("/", difficultyController.CreateDifficulty)
	difficultyRoutes.Put("/:id", difficultyController.UpdateDifficulty)
	difficultyRoutes.Delete("/:id", difficultyController.DeleteDifficulty)

}
