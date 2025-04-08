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
	categoryCtrl := categoryController.NewCategoryController(db)
	categoryRoutes := api.Group("/categories")
	categoryRoutes.Get("/", categoryCtrl.GetCategories)
	categoryRoutes.Get("/:id", categoryCtrl.GetCategory)
	categoryRoutes.Get("/difficulty/:difficulty_id", categoryCtrl.GetCategoriesByDifficulty)
	categoryRoutes.Post("/", categoryCtrl.CreateCategory)
	categoryRoutes.Put("/:id", categoryCtrl.UpdateCategory)
	categoryRoutes.Delete("/:id", categoryCtrl.DeleteCategory)

	// 📰 Category News Routes
	categoryNewsCtrl := categoryController.NewCategoryNewsController(db)
	categoryNewsRoutes := api.Group("/categories-news")
	categoryNewsRoutes.Get("/:category_id", categoryNewsCtrl.GetAll)
	categoryNewsRoutes.Get("/:id", categoryNewsCtrl.GetByID)
	categoryNewsRoutes.Post("/", categoryNewsCtrl.Create)
	categoryNewsRoutes.Put("/:id", categoryNewsCtrl.Update)
	categoryNewsRoutes.Delete("/:id", categoryNewsCtrl.Delete)

}
