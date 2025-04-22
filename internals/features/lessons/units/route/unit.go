package route

import (
	unitController "quiz-fiber/internals/features/lessons/units/controller"
	userController "quiz-fiber/internals/features/users/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register category-related routes
func CategoryRoutes(app *fiber.App, db *gorm.DB) {

	// 🔥 Proteksi seluruh kategori API dengan Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Unit Routes
	unitCtrl := unitController.NewUnitController(db)
	unitRoutes := api.Group("/units")
	unitRoutes.Get("/", unitCtrl.GetUnits)
	unitRoutes.Get("/:id", unitCtrl.GetUnit)
	unitRoutes.Get("/themes-or-levels/:themesOrLevelId", unitCtrl.GetUnitByThemesOrLevels)
	unitRoutes.Post("/", unitCtrl.CreateUnit)
	unitRoutes.Put("/:id", unitCtrl.UpdateUnit)
	unitRoutes.Delete("/:id", unitCtrl.DeleteUnit)

	// 📰 Unit News Routes
	unitNewsCtrl := unitController.NewUnitNewsController(db)
	unitNewsRoutes := api.Group("/units-news")
	unitNewsRoutes.Get("/", unitNewsCtrl.GetAll)
	unitNewsRoutes.Get("/:id", unitNewsCtrl.GetByID)
	unitNewsRoutes.Post("/", unitNewsCtrl.Create)
	unitNewsRoutes.Put("/:id", unitNewsCtrl.Update)
	unitNewsRoutes.Delete("/:id", unitNewsCtrl.Delete)

	// User Unit
	userUnitCtrl := unitController.NewUserUnitController(db)
	userUnitRoutes := api.Group("/user-units")
	userUnitRoutes.Get("/:user_id", userUnitCtrl.GetByUserID)
	userUnitRoutes.Get("/:user_id/themes-or-levels/:themes_or_levels_id", userUnitCtrl.GetUserUnitsByThemesOrLevelsAndUserID)

}
