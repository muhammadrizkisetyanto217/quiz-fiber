package routes

import (
	// Add this line.
	categoryRoute "quiz-fiber/internals/features/category/category/route"
	difficultyRoute "quiz-fiber/internals/features/category/difficulty/route"
	subcategoryRoute "quiz-fiber/internals/features/category/subcategory/route"
	themesOrLevles "quiz-fiber/internals/features/category/themes_or_levels/route"
	units "quiz-fiber/internals/features/category/units/route"
	authRoute "quiz-fiber/internals/features/user/auth/route"
	userRoute "quiz-fiber/internals/features/user/user/route"

	// Quizzes
	evaluationRoute "quiz-fiber/internals/features/quizzes/evaluation/route"
	examRoute "quiz-fiber/internals/features/quizzes/exam/route"
	questionRoute "quiz-fiber/internals/features/quizzes/question/route"
	SectionQuizzesRoutes "quiz-fiber/internals/features/quizzes/quizzes/route"
	readingRoute "quiz-fiber/internals/features/quizzes/reading/route"

	// Utils
	tooltipRoute "quiz-fiber/internals/features/utils/tooltip/route"

	// Progress
	levelRankRoute "quiz-fiber/internals/features/progress/level_rank/route"
	pointRoutes "quiz-fiber/internals/features/progress/point/route"

	// Donation
	donationRoutes "quiz-fiber/internals/features/donation/donations/routes"

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
