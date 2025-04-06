package routes

import (
	// Add this line.
	userRoute "quiz-fiber/internals/features/user/user/route"
	authRoute "quiz-fiber/internals/features/user/auth/route"
	difficultyRoute "quiz-fiber/internals/features/category/difficulty/route"
	categoryRoute "quiz-fiber/internals/features/category/category/route"
	subcategoryRoute "quiz-fiber/internals/features/category/subcategory/route"
	themesOrLevles "quiz-fiber/internals/features/category/themes_or_levels/route"
	units "quiz-fiber/internals/features/category/units/route"

	// Quizzes
	SectionQuizzesRoutes "quiz-fiber/internals/features/quizzes/quizzes/route"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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



}
