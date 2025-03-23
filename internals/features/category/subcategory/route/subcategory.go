package route

import (

	subcategoryController "quiz-fiber/internals/features/category/subcategory/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register category-related routes
func CategoryRoutes(app *fiber.App, db *gorm.DB) {

	// 🔥 Proteksi seluruh kategori API dengan Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Subcategory Routes
	subcategoryController := subcategoryController.NewSubcategoryController(db)
	subcategoryRoutes := api.Group("/subcategories")
	subcategoryRoutes.Get("/", subcategoryController.GetSubcategories)
	subcategoryRoutes.Get("/:id", subcategoryController.GetSubcategory)
	subcategoryRoutes.Get("/category/:category_id", subcategoryController.GetSubcategoriesByCategory)
	subcategoryRoutes.Post("/", subcategoryController.CreateSubcategory)
	subcategoryRoutes.Put("/:id", subcategoryController.UpdateSubcategory)
	subcategoryRoutes.Delete("/:id", subcategoryController.DeleteSubcategory)

}
