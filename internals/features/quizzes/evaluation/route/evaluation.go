package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	evaluationController "quiz-fiber/internals/features/quizzes/evaluation/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"
)

func EvaluationRoute(app *fiber.App, db *gorm.DB) {

	// 🔒 Semua API reading dilindungi oleh Auth Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🏆 Evaluation Routes
	evaluationCtrl := evaluationController.NewEvaluationController(db)
	evaluationRoutes := api.Group("/evaluations")
	evaluationRoutes.Get("/", evaluationCtrl.GetEvaluations)
	evaluationRoutes.Get("/:id", evaluationCtrl.GetEvaluation)
	evaluationRoutes.Get("/unit/:unitId", evaluationCtrl.GetEvaluationsByUnitID)
	evaluationRoutes.Post("/", evaluationCtrl.CreateEvaluation)
	evaluationRoutes.Put("/:id", evaluationCtrl.UpdateEvaluation)
	evaluationRoutes.Delete("/:id", evaluationCtrl.DeleteEvaluation)

	// 🧠 User Evaluation Routes
	userEvaluationController := evaluationController.NewUserEvaluationController(db)
	userEvaluationRoutes := api.Group("/user-evaluations")
	userEvaluationRoutes.Post("/", userEvaluationController.Create)
	userEvaluationRoutes.Get("/:user_id", userEvaluationController.GetByUserID) // Ambil semua user_reading

}
