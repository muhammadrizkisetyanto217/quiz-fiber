package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	levelController "quiz-fiber/internals/features/progress/level_rank/controller"
	userController "quiz-fiber/internals/features/users/auth/controller"
)

func LevelRequirementRoute(app *fiber.App, db *gorm.DB) {
	// 🔐 Group API dengan Auth Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Level Requirement Routes
	levelCtrl := levelController.NewLevelRequirementController(db)
	levelRoutes := api.Group("/level-requirements")
	levelRoutes.Get("/", levelCtrl.GetAll)
	levelRoutes.Get("/:id", levelCtrl.GetByID)
	levelRoutes.Post("/", levelCtrl.Create)
	levelRoutes.Put("/:id", levelCtrl.Update)
	levelRoutes.Delete("/:id", levelCtrl.Delete)

	// 🏆 Rank Requirement Routes
	rankCtrl := levelController.NewRankRequirementController(db)
	rankRoutes := api.Group("/rank-requirements")
	rankRoutes.Get("/", rankCtrl.GetAll)
	rankRoutes.Get("/:id", rankCtrl.GetByID)
	rankRoutes.Post("/", rankCtrl.Create)
	rankRoutes.Put("/:id", rankCtrl.Update)
	rankRoutes.Delete("/:id", rankCtrl.Delete)
}
