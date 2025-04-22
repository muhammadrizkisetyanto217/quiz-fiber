package routes

import (
	// Add this line.
	categoryRoute "quiz-fiber/internals/features/lessons/categories/route"
	difficultyRoute "quiz-fiber/internals/features/lessons/difficulty/route"
	subcategoryRoute "quiz-fiber/internals/features/lessons/subcategory/route"
	themesOrLevles "quiz-fiber/internals/features/lessons/themes_or_levels/route"
	units "quiz-fiber/internals/features/lessons/units/route"
	authRoute "quiz-fiber/internals/features/users/auth/route"
	userRoute "quiz-fiber/internals/features/users/user/route"

	// Quizzes
	evaluationRoute "quiz-fiber/internals/features/quizzes/evaluations/route"
	examRoute "quiz-fiber/internals/features/quizzes/exams/route"
	questionRoute "quiz-fiber/internals/features/quizzes/questions/route"
	SectionQuizzesRoutes "quiz-fiber/internals/features/quizzes/quizzes/route"
	readingRoute "quiz-fiber/internals/features/quizzes/readings/route"

	// Utils
	tooltipRoute "quiz-fiber/internals/features/utils/tooltips/route"

	// Progress
	levelRankRoute "quiz-fiber/internals/features/progress/level_rank/route"
	pointRoutes "quiz-fiber/internals/features/progress/points/route"

	// Donation
	donationRoutes "quiz-fiber/internals/features/donations/donations/routes"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"


	// Kajian Attendance
	kajianAttendanceRoute "quiz-fiber/internals/features/ms/route"
)

// Register routes
func SetupRoutes(app *fiber.App, db *gorm.DB) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Fiber & Supabase PostgreSQL connected successfully 🚀")
	})

	userRoute.UserRoutes(app, db)
	authRoute.AuthRoutes(app, db)

	//* Category
	difficultyRoute.CategoryRoutes(app, db)
	categoryRoute.CategoryRoutes(app, db)
	subcategoryRoute.CategoryRoutes(app, db)
	themesOrLevles.CategoryRoutes(app, db)
	units.CategoryRoutes(app, db)

	//* Quizzes
	SectionQuizzesRoutes.SectionQuizzesRoutes(app, db)
	questionRoute.QuestionRoutes(app, db)
	readingRoute.ReadingRoutes(app, db)
	evaluationRoute.EvaluationRoute(app, db)
	examRoute.ExamRoute(app, db)

	//* Utils
	tooltipRoute.TooltipRoute(app, db)

	//* Progress
	pointRoutes.UserPointRoutes(app, db)
	levelRankRoute.LevelRequirementRoute(app, db)

	//* Donation
	donationRoutes.DonationRoutes(app, db)


	//* Kajian Attendance
	kajianAttendanceRoute.KajianAttendanceRoutes(app, db)

}
