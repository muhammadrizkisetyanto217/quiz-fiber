package route

import (
	"quiz-fiber/internals/features/ms/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func KajianAttendanceRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", userController.AuthMiddleware(db))

	kajianCtrl := controller.NewKajianAttendanceController(db)
	kajianRoutes := api.Group("/kajian-attendances")

	kajianRoutes.Post("/", kajianCtrl.Create)                  // Create attendance
	kajianRoutes.Get("/", kajianCtrl.GetAll)                   // Get all attendance
	kajianRoutes.Get("/user/:user_id", kajianCtrl.GetByUserID) // Get by user
}
