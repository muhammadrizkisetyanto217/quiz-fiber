package route

import (
	"quiz-fiber/internals/features/quizzes/questions/controller"

	userController "quiz-fiber/internals/features/users/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register quiz-related routes
func QuestionRoutes(app *fiber.App, db *gorm.DB) {

	// 🔒 Middleware Auth diaktifkan untuk seluruh API /api/*
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Quiz Question Routes
	questionController := controller.NewQuestionController(db)
	questionRoutes := api.Group("/question")
	questionRoutes.Get("/", questionController.GetQuestions)
	questionRoutes.Get("/:id", questionController.GetQuestion)
	questionRoutes.Get("/:quizId/questionsQuiz", questionController.GetQuestionsByQuizID)
	questionRoutes.Get("/:evaluationId/questionsEvaluation/", questionController.GetQuestionsByEvaluationID)
	questionRoutes.Get("/:examId/questionsExam/", questionController.GetQuestionsByExamID)

	questionRoutes.Post("/", questionController.CreateQuestion)
	questionRoutes.Put("/:id", questionController.UpdateQuestion)
	questionRoutes.Delete("/:id", questionController.DeleteQuestion)

	questionRoutes.Get(":id/questionTooltips", questionController.GetQuestionWithTooltips)
	questionRoutes.Get(":id/questionTooltips/:tooltipId", questionController.GetOnlyQuestionTooltips)
	questionRoutes.Get("/:id/questionTooltipsMarked", questionController.GetQuestionWithTooltipsMarked)

	// Quiz Saved Routes
	questionSavedController := controller.NewQuestionSavedController(db)
	questionSavedRoutes := api.Group("/question-saved")

	questionSavedRoutes.Post("/", questionSavedController.Create)
	questionSavedRoutes.Get("/user/:user_id", questionSavedController.GetByUserID)
	questionSavedRoutes.Get("/question_saved_with_question/:user_id", questionSavedController.GetByUserIDWithQuestions)
	questionSavedRoutes.Delete("/user/:id", questionSavedController.Delete)

	// Quiz Mistake Routes
	questionMistakeController := controller.NewQuestionMistakeController(db)
	questionMistakeRoutes := api.Group("/question-mistakes")
	questionMistakeRoutes.Post("/", questionMistakeController.Create)
	questionMistakeRoutes.Get("/user/:user_id", questionMistakeController.GetByUserID)
	questionMistakeRoutes.Delete("/:id", questionMistakeController.Delete)

	// User Question Routes
	userQuestionController := controller.NewUserQuestionController(db)
	userQuestionRoutes := api.Group("/user-questions")
	userQuestionRoutes.Post("/", userQuestionController.Create)
	userQuestionRoutes.Get("/user/:user_id", userQuestionController.GetByUserID)
	userQuestionRoutes.Get("/user/:user_id/question/:question_id", userQuestionController.GetByUserIDAndQuestionID)

}
