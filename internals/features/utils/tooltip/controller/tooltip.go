package utils

import (
	"log"
	"time"

	"quiz-fiber/internals/database"
	"quiz-fiber/internals/features/utils/tooltip/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// TooltipsController menangani semua operasi terkait tooltips
type TooltipsController struct {
	DB *gorm.DB
}

// NewTooltipsController membuat instance baru dari TooltipsController
func NewTooltipsController(db *gorm.DB) *TooltipsController {
	return &TooltipsController{DB: db}
}

func (tc *TooltipsController) GetTooltips(c *fiber.Ctx) error {
	log.Println("Fetching tooltips for given keywords")

	var request struct {
		Keywords []string `json:"keywords"`
	}

	// Parsing request body
	if err := c.BodyParser(&request); err != nil {
		log.Println("Error parsing request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Inisialisasi array untuk menyimpan ID tooltips
	var tooltipIDs []uint

	// Loop pencarian keyword dalam database
	for _, keyword := range request.Keywords {
		var tooltip model.Tooltip
		if err := database.DB.Select("id").Where("keyword = ?", keyword).First(&tooltip).Error; err == nil {
			tooltipIDs = append(tooltipIDs, tooltip.ID)
		}
	}

	// Mengembalikan array ID tooltips
	return c.JSON(fiber.Map{
		"tooltips_id": tooltipIDs,
	})
}

// InsertTooltip menangani permintaan untuk menambahkan tooltips baru
func (tc *TooltipsController) InsertTooltip(c *fiber.Ctx) error {
	log.Println("Inserting new tooltip")

	var request struct {
		Keyword          string `json:"keyword"`
		DescriptionShort string `json:"description_short"`
		DescriptionLong  string `json:"description_long"`
	}

	// Parsing request body
	if err := c.BodyParser(&request); err != nil {
		log.Println("Error parsing request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Cek apakah keyword sudah ada di database
	var existingTooltip model.Tooltip
	if err := tc.DB.Where("keyword = ?", request.Keyword).First(&existingTooltip).Error; err == nil {
		log.Println("Keyword already exists:", request.Keyword)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Keyword already exists"})
	}

	// Insert data baru
	newTooltip := model.Tooltip{
		Keyword:          request.Keyword,
		DescriptionShort: request.DescriptionShort,
		DescriptionLong:  request.DescriptionLong,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := tc.DB.Create(&newTooltip).Error; err != nil {
		log.Println("Error inserting tooltip:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert tooltip"})
	}

	return c.JSON(fiber.Map{"message": "Tooltip added successfully", "data": newTooltip})
}

// GetAllTooltips menangani permintaan untuk mendapatkan semua data tooltips
func (tc *TooltipsController) GetAllTooltips(c *fiber.Ctx) error {
	log.Println("Fetching all tooltips")

	var tooltips []model.Tooltip

	// Ambil semua data dari database
	if err := tc.DB.Find(&tooltips).Error; err != nil {
		log.Println("Error fetching tooltips:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tooltips"})
	}

	return c.JSON(tooltips)
}
