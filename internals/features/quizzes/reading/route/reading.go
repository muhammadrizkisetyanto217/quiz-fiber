package route

import (
	userController "quiz-fiber/internals/features/user/auth/controller"
	readingController "quiz-fiber/internals/features/quizzes/reading/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register routes untuk fitur Reading
func ReadingRoutes(app *fiber.App, db *gorm.DB) {
	// 🔒 Semua API reading dilindungi oleh Auth Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 📖 Reading Routes
	controller := readingController.NewReadingController(db)
	readingRoutes := api.Group("/readings")
	readingRoutes.Get("/", controller.GetReadings)                            // Ambil semua reading
	readingRoutes.Get("/:id", controller.GetReading)                         // Ambil satu reading
	readingRoutes.Get("/unit/:unitId", controller.GetReadingsByUnit)        // Ambil berdasarkan unit

	readingRoutes.Post("/", controller.CreateReading)                       // Tambah reading baru
	readingRoutes.Put("/:id", controller.UpdateReading)                     // Edit reading
	readingRoutes.Delete("/:id", controller.DeleteReading)                  // Hapus reading

	// 🧠 Tooltips integration
	readingRoutes.Get("/:id/with-tooltips", controller.GetReadingWithTooltips)      // Reading + Tooltips lengkap
	readingRoutes.Get("/:id/tooltips", controller.GetOnlyReadingTooltips)           // Hanya tooltips
	readingRoutes.Get("/:id/convert", controller.ConvertReadingWithTooltipsId)      // Tandai keyword + update DB
}
