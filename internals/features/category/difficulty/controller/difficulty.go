package category

import (
	"fmt"
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
	return c.JSON(fiber.Map{
		"message": "All difficulties fetched successfully",
		"total":   len(difficulties),
		"data":    difficulties,
	})
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
	return c.JSON(fiber.Map{
		"message": "Difficulty fetched successfully",
		"data":    difficulty,
	})
}

// CreateDifficulty menangani request untuk membuat satu atau banyak data difficulty
// CreateDifficulty menangani input satu atau banyak difficulty
func (dc *DifficultyController) CreateDifficulty(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create difficulty")

	var single model.DifficultyModel
	var multiple []model.DifficultyModel

	// 🧠 Coba parse sebagai array terlebih dahulu
	if err := c.BodyParser(&multiple); err == nil && len(multiple) > 0 {
		log.Printf("[DEBUG] Parsed %d difficulties as array\n", len(multiple))

		// ✅ Validasi setiap data di array
		for i, d := range multiple {
			if d.Name == "" {
				return c.Status(400).JSON(fiber.Map{
					"error": "Name is required in each difficulty",
					"index": i,
				})
			}
		}

		// ✅ Simpan jika semua valid
		if err := dc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Failed to insert multiple difficulties: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create difficulties"})
		}

		log.Printf("[SUCCESS] %d difficulties created successfully\n", len(multiple))
		return c.Status(201).JSON(fiber.Map{
			"message": "Multiple difficulties created successfully",
			"data":    multiple,
		})
	}

	// 🔁 Jika bukan array, parse sebagai satu objek
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single difficulty input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.Printf("[DEBUG] Parsed single difficulty: %+v\n", single)

	// ✅ Validasi name
	if single.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}

	// ✅ Simpan ke database
	if err := dc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Failed to insert single difficulty: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create difficulty"})
	}

	log.Printf("[SUCCESS] Difficulty created: ID=%d, Name=%s\n", single.ID, single.Name)
	return c.Status(201).JSON(fiber.Map{
		"message": "Difficulty created successfully",
		"data":    single,
	})
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
	return c.JSON(fiber.Map{
		"message": "Difficulty updated successfully",
		"data":    difficulty,
	})
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
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Difficulty with ID %s deleted successfully", id),
	})
}
