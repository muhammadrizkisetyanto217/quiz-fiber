package ro

import (
	themes_or_levelsController "quiz-fiber/internals/features/category/themes_or_levels/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register category-related routes
func CategoryRoutes(app *fiber.App, db *gorm.DB) {

	// 🔥 Proteksi seluruh kategori API dengan Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Themes or Levels Routes
	themeOrLevelController := themes_or_levelsController.NewThemeOrLevelController(db)
	themeOrLevelRoutes := api.Group("/themes-or-levels")
	themeOrLevelRoutes.Get("/", themeOrLevelController.GetThemeOrLevels)
	themeOrLevelRoutes.Get("/:id", themeOrLevelController.GetThemeOrLevelById)
	themeOrLevelRoutes.Post("/", themeOrLevelController.CreateThemeOrLevel)
	themeOrLevelRoutes.Get("/subcategories/:subcategory_id", themeOrLevelController.GetThemesOrLevelsBySubcategory)
	themeOrLevelRoutes.Put("/:id", themeOrLevelController.UpdateThemeOrLevel)
	themeOrLevelRoutes.Delete("/:id", themeOrLevelController.DeleteThemeOrLevel)

}
