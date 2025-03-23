package quizzes

import (
	"log"

	"quiz-fiber/internals/features/quizzes/quiz_section/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SectionQuizController struct {
	DB *gorm.DB
}

func NewSectionQuizController(db *gorm.DB) *SectionQuizController {
	return &SectionQuizController{DB: db}
}

func (sqc *SectionQuizController) GetSectionQuizzes(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all section quizzes")
	var quizzes []model.SectionQuizzesModel
	if err := sqc.DB.Find(&quizzes).Error; err != nil {
		log.Println("[ERROR] Failed to fetch section quizzes:", err)
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Failed to fetch section quizzes"})
	}
	log.Printf("[SUCCESS] Retrieved %d section quizzes\n", len(quizzes))
	return c.JSON(fiber.Map{"status": true, "message": "Section quizzes fetched successfully", "data": quizzes})
}

func (sqc *SectionQuizController) GetSectionQuiz(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching section quiz with ID: %s\n", id)
	var quiz model.SectionQuizzesModel
	if err := sqc.DB.First(&quiz, id).Error; err != nil {
		log.Println("[ERROR] Section quiz not found:", err)
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Section quiz not found"})
	}
	return c.JSON(fiber.Map{"status": true, "message": "Section quiz fetched successfully", "data": quiz})
}

func (sqc *SectionQuizController) GetSectionQuizzesByUnit(c *fiber.Ctx) error {
	unitID := c.Params("unitId")
	log.Printf("[INFO] Fetching section quizzes for unit_id: %s\n", unitID)

	var sectionQuizzes []model.SectionQuizzesModel
	if err := sqc.DB.Where("unit_id = ?", unitID).Find(&sectionQuizzes).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch section quizzes for unit_id %s: %v\n", unitID, err)
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Failed to fetch section quizzes by unit ID"})
	}

	log.Printf("[SUCCESS] Retrieved %d section quizzes for unit_id %s\n", len(sectionQuizzes), unitID)
	return c.JSON(fiber.Map{"status": true, "message": "Section quizzes fetched by unit ID successfully", "data": sectionQuizzes})
}

func (sqc *SectionQuizController) CreateSectionQuiz(c *fiber.Ctx) error {
	log.Println("[INFO] Creating a new section quiz")
	var quiz model.SectionQuizzesModel
	if err := c.BodyParser(&quiz); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Invalid request"})
	}
	if err := sqc.DB.Create(&quiz).Error; err != nil {
		log.Println("[ERROR] Failed to create section quiz:", err)
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Failed to create section quiz"})
	}
	log.Printf("[SUCCESS] Section quiz created with ID: %d\n", quiz.ID)
	return c.Status(201).JSON(fiber.Map{"status": true, "message": "Section quiz created successfully", "data": quiz})
}

func (sqc *SectionQuizController) UpdateSectionQuiz(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating section quiz with ID: %s\n", id)

	var quiz model.SectionQuizzesModel
	if err := sqc.DB.First(&quiz, id).Error; err != nil {
		log.Println("[ERROR] Section quiz not found:", err)
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Section quiz not found"})
	}

	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Invalid request"})
	}

	if err := sqc.DB.Model(&quiz).Updates(requestData).Error; err != nil {
		log.Println("[ERROR] Failed to update section quiz:", err)
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Failed to update section quiz"})
	}

	log.Printf("[SUCCESS] Section quiz with ID %s updated\n", id)
	return c.JSON(fiber.Map{"status": true, "message": "Section quiz updated successfully", "data": quiz})
}

func (sqc *SectionQuizController) DeleteSectionQuiz(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting section quiz with ID: %s\n", id)
	if err := sqc.DB.Delete(&model.SectionQuizzesModel{}, id).Error; err != nil {
		log.Println("[ERROR] Failed to delete section quiz:", err)
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Failed to delete section quiz"})
	}
	log.Printf("[SUCCESS] Section quiz with ID %s deleted\n", id)
	return c.JSON(fiber.Map{"status": true, "message": "Section quiz deleted successfully"})
}
