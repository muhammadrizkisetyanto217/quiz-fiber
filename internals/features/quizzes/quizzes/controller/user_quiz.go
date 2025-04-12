package controller

import (
	"errors"
	"log"

	"quiz-fiber/internals/features/quizzes/quizzes/model"
	"quiz-fiber/internals/features/quizzes/quizzes/services"

	categoryModel "quiz-fiber/internals/features/category/subcategory/model"
	themesOrLevelsModel "quiz-fiber/internals/features/category/themes_or_levels/model"
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

	// Dapatkan unit untuk themes_or_level_id
	var unit unitModel.UnitModel
	if err := uc.DB.First(&unit, section.UnitID).Error; err != nil {
		log.Println("[ERROR] Failed to fetch UnitModel:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch related unit"})
	}

	// Ambil themes_or_level untuk ambil subcategory_id
	var theme themesOrLevelsModel.ThemesOrLevelsModel
	if err := uc.DB.First(&theme, unit.ThemesOrLevelID).Error; err != nil {
		log.Println("[ERROR] Failed to fetch ThemesOrLevelsModel:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch related theme"})
	}

	// Ambil subcategory untuk ambil category_id
	var subcategory categoryModel.SubcategoryModel
	if err := uc.DB.First(&subcategory, theme.SubcategoriesID).Error; err != nil {
		log.Println("[ERROR] Failed to fetch SubcategoryModel:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch related subcategory"})
	}

	// Update progres ke section dan unit
	_ = services.UpdateUserSectionIfQuizCompleted(uc.DB, input.UserID, int(section.ID), int(input.QuizID))
	_ = services.UpdateUserUnitIfSectionCompleted(uc.DB, input.UserID, section.UnitID, section.ID)

	// _ = services.UpdateUserThemesOrLevelsIfUnitCompleted(uc.DB, input.UserID, int(unit.ID), int(unit.ThemesOrLevelID))
	// _ = services.UpdateUserSubcategoryIfThemeCompleted(uc.DB, input.UserID, int(theme.ID), int(theme.SubcategoriesID))
	// _ = services.UpdateUserCategoryIfSubcategoryCompleted(uc.DB, input.UserID, int(subcategory.ID), int(subcategory.CategoriesID))

	return c.JSON(fiber.Map{
		"message": "User quiz progress saved and progress updated",
		"data":    input,
	})
}
