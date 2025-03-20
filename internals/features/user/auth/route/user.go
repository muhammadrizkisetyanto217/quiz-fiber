package route

import (
	"quiz-fiber/internals/features/user/auth/controller"
	// "quiz-fiber/internals/features/user/user/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes mengatur routing untuk authentication dan user
// SetupRoutes mengatur routing untuk user tanpa middleware
func AuthRoutes(app *fiber.App, db *gorm.DB) {

	//* Dengan constructor
	authController := controller.NewAuthController(db)

	// // 🔥 Setup AuthController
	auth := app.Group("/auth")
	auth.Post("/register", authController.Register) // ✅ Register user baru
	auth.Post("/login", authController.Login)       // ✅ Login user

	auth.Post("/forgot-password/check", authController.CheckSecurityAnswer) // validasi email dan jawaban keamanan
	auth.Post("/forgot-password/reset", authController.ResetPassword)       // reset password setelah validasi berhasil

	// 🔥 Setup AuthController with middleware
	protectedRoutes := app.Group("/api/auth", controller.AuthMiddleware(db))
	protectedRoutes.Post("/logout", authController.Logout)                  // ✅ Logout User
	protectedRoutes.Post("/change-password", authController.ChangePassword) // ✅ Ganti Password User

	googleAuthController := controller.NewGoogleAuthController(db)

	// Auth routes group
	authGoogle := app.Group("/auth")

	// Regular auth routes
	authGoogle.Post("/register", authController.Register)
	authGoogle.Post("/login", authController.Login)
	authGoogle.Post("/logout", authController.Logout)
	authGoogle.Post("/check-security", authController.CheckSecurityAnswer)
	authGoogle.Post("/reset-password", authController.ResetPassword)

	// Google auth routes
	authGoogle.Get("/google", googleAuthController.GoogleLogin)
	authGoogle.Get("/google/callback", googleAuthController.GoogleCallback)

	// Protected routes
	protected := app.Group("/user")
	protected.Use(controller.AuthMiddleware(db))
	protected.Post("/change-password", authController.ChangePassword)
}
