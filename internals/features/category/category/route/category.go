package category

import (
	categoryController "quiz-fiber/internals/features/category/category/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register category-related routes
func CategoryRoutes(app *fiber.App, db *gorm.DB) {

	// 🔥 Proteksi seluruh kategori API dengan Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Category Routes
	categoryController := categoryController.NewCategoryController(db)
	categoryRoutes := api.Group("/categories")
	categoryRoutes.Get("/", categoryController.GetCategories)
	categoryRoutes.Get("/:id", categoryController.GetCategory)
	categoryRoutes.Get("/difficulty/:difficulty_id", categoryController.GetCategoriesByDifficulty)
	categoryRoutes.Post("/", categoryController.CreateCategory)
	categoryRoutes.Put("/:id", categoryController.UpdateCategory)
	categoryRoutes.Delete("/:id", categoryController.DeleteCategory)

}
