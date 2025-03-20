package category

import (
	"log"
	category "quiz-fiber/internals/features/category/subcategory/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SubcategoryController struct {
	DB *gorm.DB
}

func NewSubcategoryController(db *gorm.DB) *SubcategoryController {
	return &SubcategoryController{DB: db}
}

func (sc *SubcategoryController) GetSubcategories(c *fiber.Ctx) error {
	log.Println("Fetching all subcategories")
	var subcategories []category.SubcategoryModel
	if err := sc.DB.Find(&subcategories).Error; err != nil {
		log.Println("Error fetching subcategories:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch subcategories"})
	}
	return c.JSON(subcategories)
}

func (sc *SubcategoryController) GetSubcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Fetching subcategory with ID:", id)
	var subcategory category.SubcategoryModel
	if err := sc.DB.First(&subcategory, id).Error; err != nil {
		log.Println("Subcategory not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Subcategory not found"})
	}
	return c.JSON(subcategory)
}

func (sc *SubcategoryController) GetSubcategoriesByCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")
	log.Printf("[INFO] Fetching subcategories with category ID: %s\n", categoryID)

	var subcategories []category.SubcategoryModel
	if err := sc.DB.Where("categories_id = ?", categoryID).Find(&subcategories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch subcategories for category ID %s: %v\n", categoryID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch subcategories"})
	}
	log.Printf("[SUCCESS] Retrieved %d subcategories for category ID %s\n", len(subcategories), categoryID)
	return c.JSON(subcategories)
}

func (sc *SubcategoryController) CreateSubcategory(c *fiber.Ctx) error {
	log.Println("Creating a new subcategory")
	var subcategory category.SubcategoryModel
	if err := c.BodyParser(&subcategory); err != nil {
		log.Println("Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := sc.DB.Create(&subcategory).Error; err != nil {
		log.Println("Error creating subcategory:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create subcategory"})
	}
	return c.Status(201).JSON(subcategory)
}

func (sc *SubcategoryController) UpdateSubcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Updating subcategory with ID:", id)
	var subcategory category.SubcategoryModel
	if err := sc.DB.First(&subcategory, id).Error; err != nil {
		log.Println("Subcategory not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Subcategory not found"})
	}

	if err := c.BodyParser(&subcategory); err != nil {
		log.Println("Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := sc.DB.Save(&subcategory).Error; err != nil {
		log.Println("Error updating subcategory:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update subcategory"})
	}

	return c.JSON(subcategory)
}

func (sc *SubcategoryController) DeleteSubcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Deleting subcategory with ID:", id)
	if err := sc.DB.Delete(&category.SubcategoryModel{}, id).Error; err != nil {
		log.Println("Error deleting subcategory:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete subcategory"})
	}
	return c.JSON(fiber.Map{"message": "Subcategory deleted successfully"})
}
