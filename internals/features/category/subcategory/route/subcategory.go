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
	subcategoryCtrl := subcategoryController.NewSubcategoryController(db)
	subcategoryRoutes := api.Group("/subcategories")
	subcategoryRoutes.Get("/", subcategoryCtrl.GetSubcategories)
	subcategoryRoutes.Get("/:id", subcategoryCtrl.GetSubcategory)
	subcategoryRoutes.Get("/category/:category_id", subcategoryCtrl.GetSubcategoriesByCategory)
	subcategoryRoutes.Get("/with-category-themes/:difficulty_id", subcategoryCtrl.GetCategoryWithSubcategoryAndThemes) 
	subcategoryRoutes.Post("/", subcategoryCtrl.CreateSubcategory)
	subcategoryRoutes.Put("/:id", subcategoryCtrl.UpdateSubcategory)
	subcategoryRoutes.Delete("/:id", subcategoryCtrl.DeleteSubcategory)

	// 📰 Subcategory News Routes
	subcategoryNewsCtrl := subcategoryController.NewSubcategoryNewsController(db)
	subcategoryNewsRoutes := api.Group("/subcategories-news")
	subcategoryNewsRoutes.Get("/:subcategory_id", subcategoryNewsCtrl.GetAll)
	subcategoryNewsRoutes.Get("/:id", subcategoryNewsCtrl.GetByID)
	subcategoryNewsRoutes.Post("/", subcategoryNewsCtrl.Create)
	subcategoryNewsRoutes.Put("/:id", subcategoryNewsCtrl.Update)
	subcategoryNewsRoutes.Delete("/:id", subcategoryNewsCtrl.Delete)

	// ✅ User Subcategory Route
	userSubcategoryCtrl := subcategoryController.NewUserSubcategoryController(db)
	userSubcategoryRoutes := api.Group("/user-subcategory")
	userSubcategoryRoutes.Post("/", userSubcategoryCtrl.Create)
	userSubcategoryRoutes.Get("/:id", userSubcategoryCtrl.GetByUserId)
	userSubcategoryRoutes.Get("/category/:user_id/difficulty/:difficulty_id", userSubcategoryCtrl.GetWithProgressByParam)

}
