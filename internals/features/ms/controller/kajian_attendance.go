package controller

import (
	"log"
	"net/http"
	"quiz-fiber/internals/features/ms/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KajianAttendanceController struct {
	DB *gorm.DB
}

func NewKajianAttendanceController(db *gorm.DB) *KajianAttendanceController {
	return &KajianAttendanceController{DB: db}
}

// GET /api/kajian-attendance
func (ctrl *KajianAttendanceController) GetAll(c *fiber.Ctx) error {
	var attendances []model.KajianAttendance
	if err := ctrl.DB.Find(&attendances).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data kajian attendance",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data",
		"data":    attendances,
	})
}

// GET /api/kajian-attendance/user/:user_id
func (ctrl *KajianAttendanceController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "user_id tidak valid"})
	}

	var attendances []model.KajianAttendance
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&attendances).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data untuk user ini",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data berdasarkan user",
		"data":    attendances,
	})
}

// POST /api/kajian-attendance
func (ctrl *KajianAttendanceController) Create(c *fiber.Ctx) error {
	var input model.KajianAttendance

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Permintaan tidak valid",
		})
	}

	// Set waktu akses jika tidak dikirim
	if input.AccessTime.IsZero() {
		input.AccessTime = time.Now()
	}

	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Gagal menyimpan data kajian attendance:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan data",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Data kajian attendance berhasil disimpan",
		"data":    input,
	})
}
