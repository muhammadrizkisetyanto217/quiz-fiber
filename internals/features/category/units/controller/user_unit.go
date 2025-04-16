package controller

import (
	"log"
	"net/http"
	"strconv"

	themesOrLevelsModel "quiz-fiber/internals/features/category/themes_or_levels/model"
	userModel "quiz-fiber/internals/features/category/units/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserUnitController struct {
	DB *gorm.DB
}

func NewUserUnitController(db *gorm.DB) *UserUnitController {
	return &UserUnitController{DB: db}
}

// GET /api/user-units/:user_id
func (ctrl *UserUnitController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var data []userModel.UserUnitModel
	if err := ctrl.DB.
		Preload("SectionProgress").
		Where("user_id = ?", userID).
		Find(&data).Error; err != nil {
		log.Println("[ERROR] Gagal ambil data user_unit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

func (ctrl *UserUnitController) GetUserUnitsByThemesOrLevelsAndUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	themesIDParam := c.Params("themes_or_levels_id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	themesID, err := strconv.Atoi(themesIDParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "themes_or_levels_id tidak valid",
		})
	}

	// Step 1: Ambil data user_themes_or_levels untuk user + theme
	var userTheme themesOrLevelsModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesID).
		First(&userTheme).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Data user_theme tidak ditemukan",
		})
	}

	// Step 2: Ambil list unit_id dari TotalUnit
	var unitIDs []int64
	for _, id := range userTheme.TotalUnit {
		unitIDs = append(unitIDs, id)
	}

	// Step 3: Ambil unit + section_quizzes
	var units []userModel.UnitModel
	if err := ctrl.DB.
		Preload("SectionQuizzes").
		Where("id IN ? AND themes_or_level_id = ?", unitIDs, themesID).
		Find(&units).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal ambil data unit dan quizzes",
		})
	}

	// Step 4: Ambil user_units untuk dapat progress tiap unit
	var userUnits []userModel.UserUnitModel
	if err := ctrl.DB.Where("user_id = ? AND unit_id IN ?", userID, unitIDs).
		Find(&userUnits).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal ambil data progress unit",
		})
	}

	// Mapping userUnit ke unit_id -> progress
	progressMap := make(map[uint]userModel.UserUnitModel)
	for _, u := range userUnits {
		progressMap[u.UnitID] = u
	}

	// Step 5: Build hasil response
	type ResponseUnit struct {
		userModel.UnitModel
		UserProgress userModel.UserUnitModel `json:"user_progress"`
	}

	var result []ResponseUnit
	for _, unit := range units {
		progress := progressMap[unit.ID]
		result = append(result, ResponseUnit{
			UnitModel:    unit,
			UserProgress: progress,
		})
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}
