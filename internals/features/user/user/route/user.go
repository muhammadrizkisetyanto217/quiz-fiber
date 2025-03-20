package route

import (
	userController "quiz-fiber/internals/features/user/user/controller"
	authController "quiz-fiber/internals/features/user/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes mengatur routing untuk authentication dan user
// SetupRoutes mengatur routing untuk user tanpa middleware
func UserRoutes(app *fiber.App, db *gorm.DB) {

	// 🔥 Setup UserController (dengan middleware untuk proteksi API)
	userController := userController.NewUserController(db)
	userRoutes := app.Group("/api/users", authController.AuthMiddleware(db)) // ✅ Proteksi semua user route
	userRoutes.Get("/", userController.GetUsers)                         // ✅ Get semua users (Hanya Admin)
	userRoutes.Get("/:id", userController.GetUser)                       // ✅ Get satu user berdasarkan ID
	userRoutes.Put("/:id", userController.UpdateUser)                    // ✅ Update user
	userRoutes.Delete("/:id", userController.DeleteUser)                 // ✅ Hapus user
}
