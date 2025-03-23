package category

import (
	"log"
	"quiz-fiber/internals/features/category/units/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UnitController struct {
	DB *gorm.DB
}

func NewUnitController(db *gorm.DB) *UnitController {
	return &UnitController{DB: db}
}

// GET all units
func (uc *UnitController) GetUnits(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all units")
	var units []model.UnitModel

	if err := uc.DB.Preload("SectionQuizzes").Find(&units).Error; err != nil {
		log.Println("[ERROR] Failed to fetch units:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch units"})
	}

	log.Printf("[SUCCESS] Retrieved %d units\n", len(units))
	return c.JSON(fiber.Map{
		"message": "All units fetched successfully",
		"total":   len(units),
		"data":    units,
	})
}

// GET single unit by ID
func (uc *UnitController) GetUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Fetching unit with ID:", id)

	var unit model.UnitModel
	if err := uc.DB.Preload("SectionQuizzes").First(&unit, id).Error; err != nil {
		log.Println("[ERROR] Unit not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Unit not found"})
	}

	log.Printf("[SUCCESS] Unit found: ID=%s\n", id)
	return c.JSON(fiber.Map{
		"message": "Unit fetched successfully",
		"data":    unit,
	})
}

// GET units by themes_or_level_id
func (uc *UnitController) GetUnitByThemesOrLevels(c *fiber.Ctx) error {
	themesOrLevelID := c.Params("themesOrLevelId")
	log.Printf("[INFO] Fetching units with themes_or_level_id: %s\n", themesOrLevelID)

	var units []model.UnitModel
	if err := uc.DB.Preload("SectionQuizzes").
		Where("themes_or_level_id = ?", themesOrLevelID).
		Find(&units).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch units for themes_or_level_id %s: %v\n", themesOrLevelID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch units"})
	}

	log.Printf("[SUCCESS] Retrieved %d units for themes_or_level_id %s\n", len(units), themesOrLevelID)
	return c.JSON(fiber.Map{
		"message": "Units fetched successfully by themes_or_level",
		"total":   len(units),
		"data":    units,
	})
}

// CreateUnit menangani input satu atau banyak unit
func (uc *UnitController) CreateUnit(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create unit")

	var single model.UnitModel
	var multiple []model.UnitModel

	// 🧠 Coba parse sebagai array terlebih dahulu
	if err := c.BodyParser(&multiple); err == nil && len(multiple) > 0 {
		log.Printf("[DEBUG] Parsed %d units as array\n", len(multiple))

		// ✅ Validasi setiap data jika diperlukan
		for i, unit := range multiple {
			if unit.ThemesOrLevelID == 0 || unit.Name == "" {
				return c.Status(400).JSON(fiber.Map{
					"error":      "Each unit must have a valid themes_or_level_id and name",
					"index":      i,
					"unit_input": unit,
				})
			}
		}

		// ✅ Simpan ke database
		if err := uc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Failed to insert multiple units: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create units"})
		}

		log.Printf("[SUCCESS] %d units created successfully\n", len(multiple))
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Multiple units created successfully",
			"data":    multiple,
		})
	}

	// 🔁 Jika bukan array, parse sebagai satu objek
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single unit input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.Printf("[DEBUG] Parsed single unit: %+v\n", single)

	// ✅ Validasi minimal
	if single.ThemesOrLevelID == 0 || single.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "themes_or_level_id and name are required"})
	}

	if err := uc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Failed to insert single unit: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create unit"})
	}

	log.Printf("[SUCCESS] Unit created: ID=%d, Name=%s\n", single.ID, single.Name)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Unit created successfully",
		"data":    single,
	})
}

// UPDATE unit
func (uc *UnitController) UpdateUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Updating unit with ID:", id)

	var unit model.UnitModel
	if err := uc.DB.First(&unit, id).Error; err != nil {
		log.Println("[ERROR] Unit not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Unit not found"})
	}

	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := uc.DB.Model(&unit).Updates(requestData).Error; err != nil {
		log.Println("[ERROR] Failed to update unit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update unit"})
	}

	log.Printf("[SUCCESS] Unit updated: ID=%s\n", id)
	return c.JSON(fiber.Map{
		"message": "Unit updated successfully",
		"data":    unit,
	})
}

// DELETE unit
func (uc *UnitController) DeleteUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting unit with ID:", id)

	if err := uc.DB.Delete(&model.UnitModel{}, id).Error; err != nil {
		log.Println("[ERROR] Failed to delete unit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete unit"})
	}

	log.Printf("[SUCCESS] Unit with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": "Unit deleted successfully",
	})
}
