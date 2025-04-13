package controller

import (
	"errors"
	"log"

	"quiz-fiber/internals/features/quizzes/quizzes/model"
	"quiz-fiber/internals/features/quizzes/quizzes/services"

	unitModel "quiz-fiber/internals/features/category/units/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserQuizController struct {
	DB *gorm.DB
}

func NewUserQuizController(db *gorm.DB) *UserQuizController {
	return &UserQuizController{DB: db}
}

// POST user-quiz (create or update) + progress tracking
func (uc *UserQuizController) CreateOrUpdateUserQuiz(c *fiber.Ctx) error {
	log.Println("[INFO] Creating or updating user quiz progress")

	var input model.UserQuizzesModel
	if err := c.BodyParser(&input); err != nil {
		log.Println("[ERROR] Invalid input:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Cek apakah user_quiz sudah ada
	var existing model.UserQuizzesModel
	err := uc.DB.Where("user_id = ? AND quiz_id = ?", input.UserID, input.QuizID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := uc.DB.Create(&input).Error; err != nil {
			log.Println("[ERROR] Failed to create user quiz:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create user quiz"})
		}
		log.Printf("[SUCCESS] Created user_quiz for user_id=%s quiz_id=%d\n", input.UserID.String(), input.QuizID)
	} else if err != nil {
		log.Println("[ERROR] Failed to query user quiz:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user quiz"})
	} else {
		input.ID = existing.ID
		if err := uc.DB.Save(&input).Error; err != nil {
			log.Println("[ERROR] Failed to update user quiz:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update user quiz"})
		}
		log.Printf("[SUCCESS] Updated user_quiz for user_id=%s quiz_id=%d\n", input.UserID.String(), input.QuizID)
	}

	// Ambil quiz untuk ambil SectionID & UnitID
	var quiz model.QuizModel
	if err := uc.DB.First(&quiz, input.QuizID).Error; err != nil {
		log.Println("[ERROR] Quiz not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Quiz not found"})
	}

	var section model.SectionQuizzesModel
	if err := uc.DB.First(&section, quiz.SectionQuizID).Error; err != nil {
		log.Println("[ERROR] Section not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Section not found"})
	}

	// Dapatkan unit
	var unit unitModel.UnitModel
	if err := uc.DB.First(&unit, section.UnitID).Error; err != nil {
		log.Println("[ERROR] Failed to fetch UnitModel:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch related unit"})
	}

	// Update progres ke section dan unit
	_ = services.UpdateUserSectionIfQuizCompleted(uc.DB, input.UserID, section.ID, input.QuizID)
	_ = services.UpdateUserUnitIfSectionCompleted(uc.DB, input.UserID, section.UnitID, section.ID)

	return c.JSON(fiber.Map{
		"message": "User quiz progress saved and progress updated",
		"data":    input,
	})
}
