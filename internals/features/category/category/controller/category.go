package controller

import (
	"log"

	"quiz-fiber/internals/features/category/category/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{DB: db}
}

// Get all categories
func (cc *CategoryController) GetCategories(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all categories")
	var categories []model.CategoryModel
	if err := cc.DB.Find(&categories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch categories: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch categories"})
	}
	log.Printf("[SUCCESS] Retrieved %d categories\n", len(categories))
	return c.JSON(fiber.Map{
		"message": "All categories fetched successfully",
		"total":   len(categories),
		"data":    categories,
	})
}

// Get single category
func (cc *CategoryController) GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching category with ID: %s\n", id)

	var category model.CategoryModel
	if err := cc.DB.First(&category, id).Error; err != nil {
		log.Printf("[ERROR] Category with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}
	log.Printf("[SUCCESS] Retrieved category: ID=%s, Name=%s\n", id, category.Name)
	return c.JSON(fiber.Map{
		"message": "Category fetched successfully",
		"data":    category,
	})
}

// Get categories by difficulty_id
func (cc *CategoryController) GetCategoriesByDifficulty(c *fiber.Ctx) error {
	difficultyID := c.Params("difficulty_id")
	log.Printf("[INFO] Fetching categories with difficulty ID: %s\n", difficultyID)

	var categories []model.CategoryModel
	if err := cc.DB.Where("difficulty_id = ?", difficultyID).Find(&categories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch categories for difficulty ID %s: %v\n", difficultyID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch categories"})
	}
	log.Printf("[SUCCESS] Retrieved %d categories for difficulty ID %s\n", len(categories), difficultyID)
	return c.JSON(fiber.Map{
		"message": "Categories fetched successfully by difficulty",
		"total":   len(categories),
		"data":    categories,
	})
}

// CreateCategory menangani permintaan untuk membuat satu atau banyak kategori
func (cc *CategoryController) CreateCategory(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create category")

	var singleCategory model.CategoryModel       // Untuk input tunggal
	var multipleCategories []model.CategoryModel // Untuk input array / banyak data

	// 💡 Coba parsing body request sebagai array terlebih dahulu
	if err := c.BodyParser(&multipleCategories); err == nil && len(multipleCategories) > 0 {
		// ✅ Jika berhasil dan array tidak kosong → buat banyak kategori sekaligus
		if err := cc.DB.Create(&multipleCategories).Error; err != nil {
			log.Printf("[ERROR] Failed to create multiple categories: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create multiple categories"})
		}
		log.Printf("[SUCCESS] %d categories created\n", len(multipleCategories))
		return c.Status(201).JSON(fiber.Map{
			"message": "Multiple categories created successfully",
			"data":    multipleCategories,
		})
	}

	// 💡 Jika gagal parsing sebagai array, coba parsing sebagai satu objek (single category)
	if err := c.BodyParser(&singleCategory); err != nil {
		log.Printf("[ERROR] Invalid input for single category: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input format"})
	}

	// ✅ Simpan satu kategori
	if err := cc.DB.Create(&singleCategory).Error; err != nil {
		log.Printf("[ERROR] Failed to create single category: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create category"})
	}

	log.Printf("[SUCCESS] Category created: ID=%d, Name=%s\n", singleCategory.ID, singleCategory.Name)
	return c.Status(201).JSON(fiber.Map{
		"message": "Category created successfully",
		"data":    singleCategory,
	})
}

// Update category
func (cc *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating category with ID: %s\n", id)

	var category model.CategoryModel
	if err := cc.DB.First(&category, id).Error; err != nil {
		log.Printf("[ERROR] Category with ID %s not found\n", id)
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	if err := c.BodyParser(&category); err != nil {
		log.Printf("[ERROR] Invalid input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := cc.DB.Save(&category).Error; err != nil {
		log.Printf("[ERROR] Failed to update category: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update category"})
	}
	log.Printf("[SUCCESS] Category updated: ID=%s, Name=%s\n", id, category.Name)
	return c.JSON(fiber.Map{
		"message": "Category updated successfully",
		"data":    category,
	})
}

// Delete category
func (cc *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting category with ID: %s\n", id)

	if err := cc.DB.Delete(&model.CategoryModel{}, id).Error; err != nil {
		log.Printf("[ERROR] Failed to delete category: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete category"})
	}
	log.Printf("[SUCCESS] Category with ID %s deleted successfully\n", id)
	return c.JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
