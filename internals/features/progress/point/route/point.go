package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	pointController "quiz-fiber/internals/features/progress/point/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"
)

func UserPointRoutes(app *fiber.App, db *gorm.DB) {
	// 📌 Group /api dengan AuthMiddleware di sini langsung
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🧭 Group /user-point-logs di dalamnya
	userPointLogController := pointController.NewUserPointLogController(db)
	userPointRoutes := api.Group("/user-point-logs")

	userPointRoutes.Get("/:user_id", userPointLogController.GetByUserID)
	userPointRoutes.Post("/", userPointLogController.Create)
}
