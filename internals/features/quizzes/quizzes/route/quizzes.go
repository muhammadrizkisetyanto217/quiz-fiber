package route

import (
	quizzesController "quiz-fiber/internals/features/quizzes/quizzes/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// QuizzesRoutes: Register semua routes terkait quizzes & evaluasi
func SectionQuizzesRoutes(app *fiber.App, db *gorm.DB) {

	// 🔒 Middleware Auth diaktifkan untuk seluruh API /api/*
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🔥 Section Quizzes Routes
	sectionQuizzesController := quizzesController.NewSectionQuizController(db)
	sectionQuizzesRoutes := api.Group("/section-quizzes")
	sectionQuizzesRoutes.Get("/", sectionQuizzesController.GetSectionQuizzes)
	sectionQuizzesRoutes.Get("/:id", sectionQuizzesController.GetSectionQuiz)
	sectionQuizzesRoutes.Get("/unit/:unitId", sectionQuizzesController.GetSectionQuizzesByUnit)
	sectionQuizzesRoutes.Post("/", sectionQuizzesController.CreateSectionQuiz)
	sectionQuizzesRoutes.Put("/:id", sectionQuizzesController.UpdateSectionQuiz)
	sectionQuizzesRoutes.Delete("/:id", sectionQuizzesController.DeleteSectionQuiz)

	// 🧠 Quiz Routes
	quizController := quizzesController.NewQuizController(db)
	quizRoutes := api.Group("/quizzes")
	quizRoutes.Get("/", quizController.GetQuizzes)
	quizRoutes.Get("/:id", quizController.GetQuiz)
	quizRoutes.Get("/section/:sectionId", quizController.GetQuizzesBySection)
	quizRoutes.Post("/", quizController.CreateQuiz)
	quizRoutes.Put("/:id", quizController.UpdateQuiz)
	quizRoutes.Delete("/:id", quizController.DeleteQuiz)


	// 🧑‍🎓 User Quiz Routes
	userQuizController := quizzesController.NewUserQuizController(db)
	userQuizRoutes := api.Group("/user-quizzes")
	userQuizRoutes.Post("/", userQuizController.CreateOrUpdateUserQuiz)
}
