package route

import (
	"quiz-fiber/internals/features/donations/donations/controller"
	userController "quiz-fiber/internals/features/users/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DonationRoutes(app *fiber.App, db *gorm.DB) {
	// 🔒 Lindungi semua API donasi dengan middleware auth
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎁 Donasi Routes
	donationCtrl := controller.NewDonationController(db)
	donationRoutes := api.Group("/donations")
	donationRoutes.Post("/", donationCtrl.CreateDonation)                   // Buat donasi + Snap token
	donationRoutes.Get("/", donationCtrl.GetAllDonations)                   // Semua donasi
	donationRoutes.Get("/user/:user_id", donationCtrl.GetDonationsByUserID) // Donasi per user

	// ✅ Webhook tanpa middleware auth, inject db manual
	app.Post("/api/donations/notification", func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return donationCtrl.HandleMidtransNotification(c)
	})
}