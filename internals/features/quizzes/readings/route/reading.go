package route

import (
	readingController "quiz-fiber/internals/features/quizzes/readings/controller"
	userController "quiz-fiber/internals/features/users/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register routes untuk fitur Reading
func ReadingRoutes(app *fiber.App, db *gorm.DB) {
	// 🔒 Semua API reading dilindungi oleh Auth Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 📖 Reading Routes
	readingCtrl := readingController.NewReadingController(db)
	readingRoutes := api.Group("/readings")
	readingRoutes.Get("/", readingCtrl.GetReadings)                   // Ambil semua reading
	readingRoutes.Get("/:id", readingCtrl.GetReading)                 // Ambil satu reading
	readingRoutes.Get("/unit/:unitId", readingCtrl.GetReadingsByUnit) // Ambil berdasarkan unit

	readingRoutes.Post("/", readingCtrl.CreateReading)      // Tambah reading baru
	readingRoutes.Put("/:id", readingCtrl.UpdateReading)    // Edit reading
	readingRoutes.Delete("/:id", readingCtrl.DeleteReading) // Hapus reading

	// 🧠 Tooltips integration
	readingRoutes.Get("/:id/with-tooltips", readingCtrl.GetReadingWithTooltips) // Reading + Tooltips lengkap
	readingRoutes.Get("/:id/tooltips", readingCtrl.GetOnlyReadingTooltips)      // Hanya tooltips
	readingRoutes.Get("/:id/convert", readingCtrl.ConvertReadingWithTooltipsId) // Tandai keyword + update DB

	// User Reading Routes
	userReadingCtrl := readingController.NewUserReadingController(db)
	userReadingRoutes := api.Group("/user-readings")
	userReadingRoutes.Post("/", userReadingCtrl.CreateUserReading)       // Ambil semua user_reading
	userReadingRoutes.Get("/:id", userReadingCtrl.GetAllUserReading)     // Ambil satu user_reading
	userReadingRoutes.Get("/user/:user_id", userReadingCtrl.GetByUserID) // Ambil berdasarkan user_id

}
