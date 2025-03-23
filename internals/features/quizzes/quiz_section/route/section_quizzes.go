package category

import (
	sectionQuizzesController "quiz-fiber/internals/features/quizzes/quiz_section/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// QuizzesRoutes: Register semua routes terkait quizzes & evaluasi
func QuizzesRoutes(app *fiber.App, db *gorm.DB) {

	// 🔒 Middleware Auth diaktifkan untuk seluruh API /api/*
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🔥 Section Quizzes Routes
	sectionQuizzesController := sectionQuizzesController.NewSectionQuizController(db)
	sectionQuizzesRoutes := api.Group("/section-quizzes")
	sectionQuizzesRoutes.Get("/", sectionQuizzesController.GetSectionQuizzes)
	sectionQuizzesRoutes.Get("/:id", sectionQuizzesController.GetSectionQuiz)
	sectionQuizzesRoutes.Get("/unit/:unitId", sectionQuizzesController.GetSectionQuizzesByUnit)
	sectionQuizzesRoutes.Post("/", sectionQuizzesController.CreateSectionQuiz)
	sectionQuizzesRoutes.Put("/:id", sectionQuizzesController.UpdateSectionQuiz)
	sectionQuizzesRoutes.Delete("/:id", sectionQuizzesController.DeleteSectionQuiz)

}
