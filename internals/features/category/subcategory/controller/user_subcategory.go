package controller

import (
	"log"
	userSubcategoryModel "quiz-fiber/internals/features/category/subcategory/model"
	themesModel "quiz-fiber/internals/features/category/themes_or_levels/model"
	unitModel "quiz-fiber/internals/features/category/units/model"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserSubcategoryController struct {
	DB *gorm.DB
}

func NewUserSubcategoryController(db *gorm.DB) *UserSubcategoryController {
	return &UserSubcategoryController{DB: db}
}

func (ctrl *UserSubcategoryController) Create(c *fiber.Ctx) error {
	var input userSubcategoryModel.UserSubcategoryModel

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.SubcategoryID == 0 || input.UserID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UserID dan SubcategoryID tidak boleh kosong atau nol",
		})
	}

	input.CreatedAt = time.Now()
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Gagal simpan user_subcategory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan data user_subcategory",
		})
	}

	// Ambil semua themes_or_levels berdasarkan subcategory
	var themes []themesModel.ThemesOrLevelsModel
	if err := ctrl.DB.Where("subcategories_id = ?", input.SubcategoryID).Find(&themes).Error; err != nil {
		log.Println("[ERROR] Gagal ambil themes:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data themes yang terkait",
		})
	}

	for _, theme := range themes {
		// Buat user_themes_or_levels
		userTheme := themesModel.UserThemesOrLevelsModel{
			UserID:           input.UserID,
			ThemesOrLevelsID: theme.ID,
			CompleteUnit:     pq.Int64Array{},
			TotalUnit:        pq.Int64Array{},
			GradeResult:      0,
			CreatedAt:        time.Now(),
		}
		if err := ctrl.DB.Create(&userTheme).Error; err != nil {
			log.Printf("[ERROR] Gagal simpan user_theme untuk theme_id %d: %v", theme.ID, err)
		}

		// Ambil dan simpan semua unit
		var units []unitModel.UnitModel
		if err := ctrl.DB.Where("themes_or_level_id = ?", theme.ID).Find(&units).Error; err != nil {
			log.Println("[ERROR] Gagal ambil unit untuk theme:", theme.ID, err)
			continue
		}

		for _, unit := range units {
			userUnit := unitModel.UserUnitModel{
				UserID:                 input.UserID,
				UnitID:                 unit.ID,
				AttemptReading:         0,
				AttemptEvaluation:      0,
				CompleteSectionQuizzes: pq.Int64Array{},
				TotalSectionQuizzes:    pq.Int64Array{},
				GradeExam:              0,
				IsPassed:               false,
				GradeResult:            0,
				CreatedAt:              time.Now(),
				UpdatedAt:              time.Now(),
			}
			if err := ctrl.DB.Create(&userUnit).Error; err != nil {
				log.Printf("[ERROR] Gagal simpan user_unit untuk unit_id %d: %v", unit.ID, err)
			}
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "UserSubcategory, UserThemes, dan UserUnits berhasil dibuat",
		"data":    input,
	})
}
