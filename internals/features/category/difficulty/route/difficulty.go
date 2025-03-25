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
	difficultyCtrl := difficultyController.NewDifficultyController(db)
	difficultyRoutes := api.Group("/difficulties")
	difficultyRoutes.Get("/", difficultyCtrl.GetDifficulties)
	difficultyRoutes.Get("/:id", difficultyCtrl.GetDifficulty)
	difficultyRoutes.Post("/", difficultyCtrl.CreateDifficulty)
	difficultyRoutes.Put("/:id", difficultyCtrl.UpdateDifficulty)
	difficultyRoutes.Delete("/:id", difficultyCtrl.DeleteDifficulty)

	// 📰 Difficulty News Routes
	difficultyNewsCtrl := difficultyController.NewDifficultyNewsController(db)
	difficultyNewsRoutes := api.Group("/difficulties-news")
	difficultyNewsRoutes.Get("/:difficulty_id", difficultyNewsCtrl.GetNewsByDifficulty)
	difficultyNewsRoutes.Get("/:id", difficultyNewsCtrl.GetNewsByID)
	difficultyNewsRoutes.Post("/", difficultyNewsCtrl.CreateNews)
	difficultyNewsRoutes.Put("/:id", difficultyNewsCtrl.UpdateNews)
	difficultyNewsRoutes.Delete("/:id", difficultyNewsCtrl.DeleteNews)
}
