package route

import (

	questionController "quiz-fiber/internals/features/quizzes/question/controller"

	userController "quiz-fiber/internals/features/user/auth/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


// Register quiz-related routes
func QuestionRoutes(app *fiber.App, db *gorm.DB) {

	// 🔒 Middleware Auth diaktifkan untuk seluruh API /api/*
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Quiz Question Routes
	questionController := questionController.NewQuestionController(db)
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
	
}