package route

import (
	authController "quiz-fiber/internals/features/users/auth/controller"
	userController "quiz-fiber/internals/features/users/user/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes mengatur routing untuk user & user profile
func UserRoutes(app *fiber.App, db *gorm.DB) {

	// ✅ Group dengan Middleware Auth
	api := app.Group("/api", authController.AuthMiddleware(db))

	// 🔹 Users
	userCtrl := userController.NewUserController(db)
	userRoutes := api.Group("/users")
	userRoutes.Get("/", userCtrl.GetUsers)
	userRoutes.Get("/:id", userCtrl.GetUser)
	userRoutes.Put("/:id", userCtrl.UpdateUser)
	userRoutes.Delete("/:id", userCtrl.DeleteUser)

	// 🔹 Users Profile
	userProfileCtrl := userController.NewUsersProfileController(db)
	usersProfileRoutes := api.Group("/users-profiles")
	usersProfileRoutes.Get("/", userProfileCtrl.GetProfiles)
	usersProfileRoutes.Get("/:id", userProfileCtrl.GetProfile)
	usersProfileRoutes.Post("/", userProfileCtrl.CreateProfile)
	usersProfileRoutes.Put("/:id", userProfileCtrl.UpdateProfile)
	usersProfileRoutes.Delete("/:id", userProfileCtrl.DeleteProfile)
}
