package category

import (
	"log"
	"quiz-fiber/internals/features/category/difficulty/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DifficultyController struct {
	DB *gorm.DB
}

func NewDifficultyController(db *gorm.DB) *DifficultyController {
	return &DifficultyController{DB: db}
}

// Get all difficulties
func (dc *DifficultyController) GetDifficulties(c *fiber.Ctx) error {
	var difficulties []model.DifficultyModel
	log.Println("[INFO] Received request to fetch all difficulties")

	if err := dc.DB.Find(&difficulties).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch difficulties: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[SUCCESS] Retrieved %d difficulties\n", len(difficulties))
	return c.JSON(difficulties)
}

// Get difficulty by ID
func (dc *DifficultyController) GetDifficulty(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching difficulty with ID: %s\n", id)

	var difficulty model.DifficultyModel
	if err := dc.DB.First(&difficulty, id).Error; err != nil {
		log.Printf("[ERROR] Difficulty with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{"error": "Difficulty not found"})
	}

	log.Printf("[SUCCESS] Retrieved difficulty: ID=%d, Name=%s\n", difficulty.ID, difficulty.Name)
	return c.JSON(difficulty)
}

// Create difficulty
func (dc *DifficultyController) CreateDifficulty(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create difficulty")
	difficulty := new(model.DifficultyModel)
	if err := c.BodyParser(difficulty); err != nil {
		log.Printf("[ERROR] Failed to parse JSON: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	log.Printf("[DEBUG] Parsed difficulty: %+v\n", difficulty)

	if err := dc.DB.Create(difficulty).Error; err != nil {
		log.Printf("[ERROR] Failed to insert difficulty: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[SUCCESS] Difficulty created: ID=%d, Name=%s\n", difficulty.ID, difficulty.Name)
	return c.Status(201).JSON(difficulty)
}

// Update difficulty
func (dc *DifficultyController) UpdateDifficulty(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating difficulty with ID: %s\n", id)

	var difficulty model.DifficultyModel
	if err := dc.DB.First(&difficulty, id).Error; err != nil {
		log.Printf("[ERROR] Difficulty with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{"error": "Difficulty not found"})
	}

	var input model.DifficultyModel
	if err := c.BodyParser(&input); err != nil {
		log.Printf("[ERROR] Invalid JSON input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := dc.DB.Model(&difficulty).Updates(input).Error; err != nil {
		log.Printf("[ERROR] Failed to update difficulty: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[SUCCESS] Difficulty with ID %s updated successfully\n", id)
	return c.JSON(difficulty)
}

// Delete difficulty
func (dc *DifficultyController) DeleteDifficulty(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting difficulty with ID: %s\n", id)

	if err := dc.DB.Delete(&model.DifficultyModel{}, id).Error; err != nil {
		log.Printf("[ERROR] Failed to delete difficulty: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[SUCCESS] Difficulty with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{"message": "Difficulty deleted"})
}
