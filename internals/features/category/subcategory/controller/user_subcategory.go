package controller

import (
	"log"
	subcategoryModel "quiz-fiber/internals/features/category/subcategory/model"
	themesModel "quiz-fiber/internals/features/category/themes_or_levels/model"
	unitModel "quiz-fiber/internals/features/category/units/model"

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
		userUnits = append(userUnits, unitModel.UserUnitModel{
			UserID:                 input.UserID,
			UnitID:                 unit.ID,
			AttemptReading:         0,
			AttemptEvaluation:      datatypes.JSON([]byte(`{"attempt":0,"grade_evaluation":0}`)),
			CompleteSectionQuizzes: datatypes.JSON([]byte(`[]`)),
			TotalSectionQuizzes:    pq.Int64Array{},
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

	// Ambil kategori berdasarkan difficulty_id
	var categories []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	if err := ctrl.DB.Table("categories").
		Select("id, name").
		Where("difficulty_id = ?", difficultyID).
		Scan(&categories).Error; err != nil {
		log.Println("[ERROR] Gagal ambil kategori:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil kategori"})
	}

	var categoryIDs []uint
	for _, cat := range categories {
		categoryIDs = append(categoryIDs, cat.ID)
	}

	// Subcategory dengan progres user
	type SubcategoryWithProgress struct {
		ID                     uint           `json:"id"`
		Name                   string         `json:"name"`
		CategoryID             uint           `json:"category_id"`
		GradeResult            int            `json:"grade_result"`
		CompleteThemesOrLevels datatypes.JSON `json:"complete_themes_or_levels"`
		TotalThemesOrLevels    pq.Int64Array  `json:"total_themes_or_levels"`
	}

	var subcategories []SubcategoryWithProgress
	if err := ctrl.DB.Table("subcategories").
		Select(`
			subcategories.id,
			subcategories.name,
			subcategories.categories_id AS category_id,
			COALESCE(user_subcategory.grade_result, 0) AS grade_result,
			COALESCE(user_subcategory.complete_themes_or_levels, '{}') AS complete_themes_or_levels,
			COALESCE(user_subcategory.total_themes_or_levels, '{}') AS total_themes_or_levels
		`).
		Joins("LEFT JOIN user_subcategory ON user_subcategory.subcategory_id = subcategories.id AND user_subcategory.user_id = ?", userID).
		Where("subcategories.categories_id IN ?", categoryIDs).
		Scan(&subcategories).Error; err != nil {
		log.Println("[ERROR] Gagal ambil subkategori:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil subkategori"})
	}

	// Ambil subcategory_id langsung tanpa helper
	var subcategoryIDs []uint
	for _, sub := range subcategories {
		subcategoryIDs = append(subcategoryIDs, sub.ID)
	}

	// Ambil themes_or_levels + user progress
	type ThemeWithProgress struct {
		ID            uint           `json:"id"`
		Name          string         `json:"name"`
		SubcategoryID uint           `json:"subcategory_id"`
		GradeResult   int            `json:"grade_result"`
		CompleteUnit  datatypes.JSON `json:"complete_unit"`
	}

	var themes []ThemeWithProgress
	if err := ctrl.DB.Table("themes_or_levels").
		Select(`
			themes_or_levels.id,
			themes_or_levels.name,
			themes_or_levels.subcategories_id AS subcategory_id,
			COALESCE(user_themes_or_levels.grade_result, 0) AS grade_result,
			COALESCE(user_themes_or_levels.complete_unit, '{}') AS complete_unit
		`).
		Joins("LEFT JOIN user_themes_or_levels ON user_themes_or_levels.themes_or_levels_id = themes_or_levels.id AND user_themes_or_levels.user_id = ?", userID).
		Where("themes_or_levels.subcategories_id IN ?", subcategoryIDs).
		Scan(&themes).Error; err != nil {
		log.Println("[ERROR] Gagal ambil themes:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data themes_or_levels"})
	}

	// Gabungkan themes ke masing-masing subcategory
	themesMap := make(map[uint][]ThemeWithProgress)
	for _, t := range themes {
		themesMap[t.SubcategoryID] = append(themesMap[t.SubcategoryID], t)
	}

	// Gabungkan subcategory ke masing-masing kategori
	type SubcategoryWithThemes struct {
		ID                     uint                `json:"id"`
		Name                   string              `json:"name"`
		CategoryID             uint                `json:"category_id"`
		GradeResult            int                 `json:"grade_result"`
		CompleteThemesOrLevels datatypes.JSON      `json:"complete_themes_or_levels"`
		TotalThemesOrLevels    pq.Int64Array       `json:"total_themes_or_levels"`
		ThemesOrLevels         []ThemeWithProgress `json:"themes_or_levels"`
	}

	categoryMap := make(map[uint][]SubcategoryWithThemes)
	for _, sub := range subcategories {
		categoryMap[sub.CategoryID] = append(categoryMap[sub.CategoryID], SubcategoryWithThemes{
			ID:                     sub.ID,
			Name:                   sub.Name,
			CategoryID:             sub.CategoryID,
			GradeResult:            sub.GradeResult,
			CompleteThemesOrLevels: sub.CompleteThemesOrLevels,
			TotalThemesOrLevels:    sub.TotalThemesOrLevels,
			ThemesOrLevels:         themesMap[sub.ID],
		})
	}

	// Final response structure
	type CategoryResponse struct {
		ID            uint                    `json:"id"`
		Name          string                  `json:"name"`
		Subcategories []SubcategoryWithThemes `json:"subcategories"`
	}

	var result []CategoryResponse
	for _, cat := range categories {
		result = append(result, CategoryResponse{
			ID:            cat.ID,
			Name:          cat.Name,
			Subcategories: categoryMap[cat.ID],
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil ambil data lengkap",
		"data":    result,
	})
}
