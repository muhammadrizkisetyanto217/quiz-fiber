package controller

import (
	"encoding/json"
	"fmt"
	"log"
	categoryModel "quiz-fiber/internals/features/lessons/categories/model"
	subcategoryModel "quiz-fiber/internals/features/lessons/subcategory/model"
	themesModel "quiz-fiber/internals/features/lessons/themes_or_levels/model"
	unitModel "quiz-fiber/internals/features/lessons/units/model"
	sectionQuizzesModel "quiz-fiber/internals/features/quizzes/quizzes/model"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserSubcategoryController struct {
	DB *gorm.DB
}

func NewUserSubcategoryController(db *gorm.DB) *UserSubcategoryController {
	return &UserSubcategoryController{DB: db}
}

func (ctrl *UserSubcategoryController) Create(c *fiber.Ctx) error {
	var input subcategoryModel.UserSubcategoryModel

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

	// Mulai transaction
	tx := ctrl.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memulai transaksi database",
		})
	}

	// Ambil data subcategory
	var subcategory subcategoryModel.SubcategoryModel
	if err := tx.First(&subcategory, input.SubcategoryID).Error; err != nil {
		tx.Rollback()
		log.Println("[ERROR] Gagal ambil subcategory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data subcategory",
		})
	}

	// Set total themes dari subcategory
	input.TotalThemesOrLevels = subcategory.TotalThemesOrLevels
	input.CreatedAt = time.Now()

	// Simpan user_subcategory
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		log.Println("[ERROR] Gagal simpan user_subcategory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan data user_subcategory",
		})
	}

	// Ambil semua themes_or_levels berdasarkan subcategory
	var themes []themesModel.ThemesOrLevelsModel
	if err := tx.Where("subcategories_id = ?", input.SubcategoryID).Find(&themes).Error; err != nil {
		tx.Rollback()
		log.Println("[ERROR] Gagal ambil themes:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data themes yang terkait",
		})
	}

	// Siapkan data user_themes dan ambil semua theme.ID untuk query units
	var themeIDs []uint
	var userThemes []themesModel.UserThemesOrLevelsModel
	for _, theme := range themes {
		themeIDs = append(themeIDs, theme.ID)
		userThemes = append(userThemes, themesModel.UserThemesOrLevelsModel{
			UserID:           input.UserID,
			ThemesOrLevelsID: theme.ID,
			CompleteUnit:     datatypes.JSONMap{},
			TotalUnit:        theme.TotalUnit,
			GradeResult:      0,
			CreatedAt:        time.Now(),
		})
	}

	// Batch insert userThemes
	if len(userThemes) > 0 {
		if err := tx.CreateInBatches(&userThemes, 100).Error; err != nil {
			tx.Rollback()
			log.Println("[ERROR] Gagal simpan user_themes batch:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menyimpan data user_themes",
			})
		}
	}

	// Ambil semua units yang terkait dengan themes
	var units []unitModel.UnitModel
	if err := tx.Where("themes_or_level_id IN ?", themeIDs).Find(&units).Error; err != nil {
		tx.Rollback()
		log.Println("[ERROR] Gagal ambil semua units:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data unit",
		})
	}

	// Buat data userUnits untuk batch insert
	var userUnits []unitModel.UserUnitModel
	now := time.Now()

	for _, unit := range units {
		// Ambil semua ID dari section_quizzes berdasarkan unit_id
		var sectionQuizIDs []int64
		if err := tx.Model(&sectionQuizzesModel.SectionQuizzesModel{}).
			Where("unit_id = ?", unit.ID).
			Pluck("id", &sectionQuizIDs).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] Gagal ambil section_quizzes untuk unit_id %d: %v", unit.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Gagal mengambil section_quizzes untuk unit_id %d", unit.ID),
			})
		}

		// Tambahkan ke list user_unit
		userUnits = append(userUnits, unitModel.UserUnitModel{
			UserID:                 input.UserID,
			UnitID:                 unit.ID,
			AttemptReading:         0,
			AttemptEvaluation:      datatypes.JSON([]byte(`{"attempt":0,"grade_evaluation":0}`)),
			CompleteSectionQuizzes: datatypes.JSON([]byte(`[]`)),
			TotalSectionQuizzes:    pq.Int64Array(sectionQuizIDs),
			GradeExam:              0,
			IsPassed:               false,
			GradeResult:            0,
			CreatedAt:              now,
			UpdatedAt:              now,
		})
	}

	if len(userUnits) > 0 {
		if err := tx.CreateInBatches(&userUnits, 100).Error; err != nil {
			tx.Rollback()
			log.Println("[ERROR] Gagal simpan user_units batch:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menyimpan data user_units",
			})
		}
	}

	// Commit transaksi jika semua berhasil
	if err := tx.Commit().Error; err != nil {
		log.Println("[ERROR] Commit transaksi gagal:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal commit transaksi database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "UserSubcategory, UserThemes, dan UserUnits berhasil dibuat",
		"data":    input,
	})
}

func (ctrl *UserSubcategoryController) GetByUserId(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validasi UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID user tidak valid",
		})
	}

	var userSub subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.Where("user_id = ?", userID).First(&userSub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Data tidak ditemukan",
			})
		}
		log.Println("[ERROR] Gagal ambil user_subcategory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": userSub,
	})
}

func (ctrl *UserSubcategoryController) GetWithProgressByParam(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	difficultyID := c.Params("difficulty_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id tidak valid"})
	}

	if difficultyID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "difficulty_id wajib diisi"})
	}

	// Step 1: Ambil semua kategori + subkategori + themes
	var categories []categoryModel.CategoryModel
	if err := ctrl.DB.
		Where("difficulty_id = ?", difficultyID).
		Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "active").Preload("ThemesOrLevels")
		}).
		Find(&categories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil kategori"})
	}

	// Step 2: Ambil progres user_subcategory
	var userSubcat []subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&userSubcat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil progress user_subcategory"})
	}
	userSubcatMap := map[uint]subcategoryModel.UserSubcategoryModel{}
	for _, us := range userSubcat {
		userSubcatMap[uint(us.SubcategoryID)] = us
	}

	// Step 3: Ambil progres user_themes_or_levels
	var userThemes []themesModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&userThemes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil progress user_themes_or_levels"})
	}
	userThemeMap := map[uint]themesModel.UserThemesOrLevelsModel{}
	for _, ut := range userThemes {
		userThemeMap[ut.ThemesOrLevelsID] = ut
	}

	// Step 4: Build response akhir
	type ThemeWithProgress struct {
		ID               uint           `json:"id"`
		Name             string         `json:"name"`
		Status           string         `json:"status"`
		DescriptionShort string         `json:"description_short"`
		DescriptionLong  string         `json:"description_long"`
		TotalUnit        pq.Int64Array  `json:"total_unit"`
		ImageURL         string         `json:"image_url"`
		UpdateNews       datatypes.JSON `json:"update_news"`
		CreatedAt        time.Time      `json:"created_at"`
		UpdatedAt        *time.Time     `json:"updated_at"`
		SubcategoriesID  uint           `json:"subcategories_id"`
		GradeResult      int            `json:"grade_result"`
		CompleteUnit     datatypes.JSON `json:"complete_unit"`
	}

	type SubcategoryWithProgress struct {
		ID                     uint                `json:"id"`
		Name                   string              `json:"name"`
		Status                 string              `json:"status"`
		DescriptionLong        string              `json:"description_long"`
		TotalThemesOrLevels    pq.Int64Array       `json:"total_themes_or_levels"`
		ImageURL               string              `json:"image_url"`
		UpdateNews             datatypes.JSON      `json:"update_news"`
		CreatedAt              time.Time           `json:"created_at"`
		UpdatedAt              *time.Time          `json:"updated_at"`
		CategoriesID           uint                `json:"categories_id"`
		GradeResult            int                 `json:"grade_result"`
		CompleteThemesOrLevels any                 `json:"complete_themes_or_levels"`
		UserSubcategoryID      uint                `json:"user_subcategory_id"`
		UserID                 uuid.UUID           `json:"user_id"`
		ThemesOrLevels         []ThemeWithProgress `json:"themes_or_levels"`
	}

	type CategoryWithSubcat struct {
		ID                 uint                      `json:"id"`
		Name               string                    `json:"name"`
		Status             string                    `json:"status"`
		DescriptionShort   string                    `json:"description_short"`
		DescriptionLong    string                    `json:"description_long"`
		TotalSubcategories pq.Int64Array             `json:"total_subcategories"`
		ImageURL           string                    `json:"image_url"`
		UpdateNews         datatypes.JSON            `json:"update_news"`
		DifficultyID       uint                      `json:"difficulty_id"`
		CreatedAt          time.Time                 `json:"created_at"`
		UpdatedAt          *time.Time                `json:"updated_at"`
		Subcategories      []SubcategoryWithProgress `json:"subcategories"`
	}

	var result []CategoryWithSubcat
	totalGrade := 0
	totalCount := 0

	for _, cat := range categories {
		var subcatWithProgress []SubcategoryWithProgress

		for _, sub := range cat.Subcategories {
			us, ok := userSubcatMap[sub.ID]
			if !ok {
				us = subcategoryModel.UserSubcategoryModel{}
			}

			var themes []ThemeWithProgress
			for _, theme := range sub.ThemesOrLevels {
				userTheme := userThemeMap[theme.ID]
				rawJSON, _ := json.Marshal(userTheme.CompleteUnit)

				themes = append(themes, ThemeWithProgress{
					ID:               theme.ID,
					Name:             theme.Name,
					Status:           theme.Status,
					DescriptionShort: theme.DescriptionShort,
					DescriptionLong:  theme.DescriptionLong,
					TotalUnit:        theme.TotalUnit,
					ImageURL:         theme.ImageURL,
					UpdateNews:       theme.UpdateNews,
					CreatedAt:        theme.CreatedAt,
					UpdatedAt:        theme.UpdatedAt,
					SubcategoriesID:  uint(theme.SubcategoriesID),
					GradeResult:      userTheme.GradeResult,
					CompleteUnit:     datatypes.JSON(rawJSON),
				})
				if userTheme.GradeResult > 0 {
					totalGrade += userTheme.GradeResult
					totalCount++
				}
			}

			if us.GradeResult > 0 {
				totalGrade += us.GradeResult
				totalCount++
			}

			subcatWithProgress = append(subcatWithProgress, SubcategoryWithProgress{
				ID:                     sub.ID,
				Name:                   sub.Name,
				Status:                 sub.Status,
				DescriptionLong:        sub.DescriptionLong,
				TotalThemesOrLevels:    sub.TotalThemesOrLevels,
				ImageURL:               sub.ImageURL,
				UpdateNews:             sub.UpdateNews,
				CreatedAt:              sub.CreatedAt,
				UpdatedAt:              sub.UpdatedAt,
				CategoriesID:           sub.CategoriesID,
				GradeResult:            us.GradeResult,
				CompleteThemesOrLevels: us.CompleteThemesOrLevels,
				UserSubcategoryID:      us.ID,
				UserID:                 us.UserID,
				ThemesOrLevels:         themes,
			})
		}

		result = append(result, CategoryWithSubcat{
			ID:                 cat.ID,
			Name:               cat.Name,
			Status:             cat.Status,
			DescriptionShort:   cat.DescriptionShort,
			DescriptionLong:    cat.DescriptionLong,
			TotalSubcategories: cat.TotalSubcategories,
			ImageURL:           cat.ImageURL,
			UpdateNews:         cat.UpdateNews,
			DifficultyID:       cat.DifficultyID,
			CreatedAt:          cat.CreatedAt,
			// UpdatedAt:          cat.UpdatedAt,
			Subcategories: subcatWithProgress,
		})
	}

	type CombinedProgress struct {
		UserID       uuid.UUID `json:"user_id"`
		AverageGrade int       `json:"average_grade"`
		DataCount    int       `json:"data_count"`
	}
	combined := CombinedProgress{
		UserID:       userID,
		AverageGrade: 0,
		DataCount:    totalCount,
	}
	if totalCount > 0 {
		combined.AverageGrade = totalGrade / totalCount
	}

	return c.JSON(fiber.Map{
		"message":       "Berhasil ambil data lengkap",
		"data":          result,
		"user_progress": combined,
	})
}
