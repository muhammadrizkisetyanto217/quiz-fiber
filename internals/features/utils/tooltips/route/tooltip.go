package route

import (
	tooltipController "quiz-fiber/internals/features/utils/tooltips/controller"
	userController "quiz-fiber/internals/features/users/auth/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func TooltipRoute(app *fiber.App, db *gorm.DB) {

	// 🔥 Proteksi seluruh kategori API dengan Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 🎯 Tooltip Routes
	tooltipCtrl := tooltipController.NewTooltipsController(db)
	tooltipRoutes := api.Group("/tooltip")
	tooltipRoutes.Get("/", tooltipCtrl.GetAllTooltips)
	tooltipRoutes.Post("/get-tooltips-id", tooltipCtrl.GetTooltipsID)
	tooltipRoutes.Post("/create-tooltips", tooltipCtrl.CreateTooltip)
}