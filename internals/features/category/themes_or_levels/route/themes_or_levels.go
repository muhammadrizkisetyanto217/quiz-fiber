package route

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
	themeOrLevelCtrl := themes_or_levelsController.NewThemeOrLevelController(db)
	themeOrLevelRoutes := api.Group("/themes-or-levels")
	themeOrLevelRoutes.Get("/", themeOrLevelCtrl.GetThemeOrLevels)
	themeOrLevelRoutes.Get("/:id", themeOrLevelCtrl.GetThemeOrLevelById)
	themeOrLevelRoutes.Post("/", themeOrLevelCtrl.CreateThemeOrLevel)
	themeOrLevelRoutes.Get("/subcategories/:subcategory_id", themeOrLevelCtrl.GetThemesOrLevelsBySubcategory)
	themeOrLevelRoutes.Put("/:id", themeOrLevelCtrl.UpdateThemeOrLevel)
	themeOrLevelRoutes.Delete("/:id", themeOrLevelCtrl.DeleteThemeOrLevel)

	// 📰 Themes or Levels News RoutesthemeOrLevelCtrl
	themesNewsCtrl := themes_or_levelsController.NewThemesOrLevelsNewsController(db)
	themesNewsRoutes := api.Group("/themes-or-levels-news")
	themesNewsRoutes.Get("/", themesNewsCtrl.GetAll)
	themesNewsRoutes.Get("/:id", themesNewsCtrl.GetByID)
	themesNewsRoutes.Post("/", themesNewsCtrl.Create)
	themesNewsRoutes.Put("/:id", themesNewsCtrl.Update)
	themesNewsRoutes.Delete("/:id", themesNewsCtrl.Delete)
}
