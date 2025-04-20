package route

import (
	"quiz-fiber/internals/features/donation/donations/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DonationRoutes(app *fiber.App, db *gorm.DB) {
	// 🔒 Lindungi semua API donasi dengan middleware auth
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎁 Donasi Routes
	donationCtrl := controller.NewDonationController(db)
	donationRoutes := api.Group("/donations")
	donationRoutes.Post("/", donationCtrl.CreateDonation) // Buat donasi + Snap token

	// 🌐 Webhook dari Midtrans → tidak perlu pakai middleware auth
	app.Post("/api/donations/notification", controller.HandleMidtransNotification)
}
