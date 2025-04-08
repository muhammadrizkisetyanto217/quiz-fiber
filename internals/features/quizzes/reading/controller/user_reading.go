package controller

import (
	"log"
	UserReadingModel "quiz-fiber/internals/features/quizzes/reading/model"

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

	// Debug: cek body request yang masuk
	body := c.Body()
	log.Println("[DEBUG] Raw request body:", string(body))

	// Parse JSON ke struct
	if err := c.BodyParser(&input); err != nil {
		log.Println("[ERROR] Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Debug: lihat hasil parsing ke struct
	log.Printf("[DEBUG] Parsed UserReading input: %+v\n", input)

	// Simpan ke database
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user reading",
		})
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
