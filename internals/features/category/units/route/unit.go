package route

import (
	unitController "quiz-fiber/internals/features/category/units/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register category-related routes
func CategoryRoutes(app *fiber.App, db *gorm.DB) {

	// 🔥 Proteksi seluruh kategori API dengan Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Unit Routes
	unitController := unitController.NewUnitController(db)
	unitRoutes := api.Group("/units")
	unitRoutes.Get("/", unitController.GetUnits)
	unitRoutes.Get("/:id", unitController.GetUnit)
	unitRoutes.Get("/themes-or-levels/:themesOrLevelId", unitController.GetUnitByThemesOrLevels)
	unitRoutes.Post("/", unitController.CreateUnit)
	unitRoutes.Put("/:id", unitController.UpdateUnit)
	unitRoutes.Delete("/:id", unitController.DeleteUnit)

}
