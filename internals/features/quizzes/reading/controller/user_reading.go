package controller

import (
	"log"
	UserReadingModel "quiz-fiber/internals/features/quizzes/reading/model"
	"quiz-fiber/internals/features/quizzes/reading/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserReadingController struct {
	DB *gorm.DB
}

func NewUserReadingController(db *gorm.DB) *UserReadingController {
	return &UserReadingController{DB: db}
}

// POST /user-readings
func (ctrl *UserReadingController) CreateUserReading(c *fiber.Ctx) error {
	var input UserReadingModel.UserReading
	body := c.Body()
	log.Println("[DEBUG] Raw request body:", string(body))

	// Parse body
	if err := c.BodyParser(&input); err != nil {
		log.Println("[ERROR] Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ Hitung Attempt ke-n
	var latestAttempt int
	err := ctrl.DB.Table("user_readings").
		Select("COALESCE(MAX(attempt), 0)").
		Where("user_id = ? AND reading_id = ?", input.UserID, input.ReadingID).
		Scan(&latestAttempt).Error
	if err != nil {
		log.Println("[ERROR] Failed to count latest attempt:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	input.Attempt = latestAttempt + 1

	// ✅ Simpan user reading
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user reading"})
	}

	// ✅ Panggil service untuk update user_unit
	if err := service.UpdateUserUnitFromReading(ctrl.DB, input.UserID, input.UnitID); err != nil {
		log.Println("[ERROR] Failed to update user unit from reading:", err)
	}

	log.Println("[SUCCESS] User reading created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User reading created successfully",
		"data":    input,
	})
}

// GET /user-readings
func (ctrl *UserReadingController) GetAllUserReading(c *fiber.Ctx) error {
	var readings []UserReadingModel.UserReading

	if err := ctrl.DB.Find(&readings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user readings",
		})
	}

	return c.JSON(readings)
}
