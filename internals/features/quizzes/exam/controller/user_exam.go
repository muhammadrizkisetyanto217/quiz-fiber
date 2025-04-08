package controller


import (
	"net/http"

	"quiz-fiber/internals/features/quizzes/exam/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserExamController struct {
	DB *gorm.DB
}

func NewUserExamController(db *gorm.DB) *UserExamController {
	return &UserExamController{DB: db}
}

// Create user_exam
func (c *UserExamController) Create(ctx *fiber.Ctx) error {
	var payload model.UserExamModel

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validasi user_id
	if payload.UserID == uuid.Nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "UserID is required and must be a valid UUID",
		})
	}

	if err := c.DB.Create(&payload).Error; err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user exam record",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User exam record created successfully",
		"data":    payload,
	})
}

// Delete user_exam by ID
func (c *UserExamController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var exam model.UserExamModel
	if err := c.DB.First(&exam, id).Error; err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "User exam not found",
			"error":   err.Error(),
		})
	}

	if err := c.DB.Delete(&exam).Error; err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user exam",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "User exam deleted successfully",
	})
}

// Get all user_exams
func (c *UserExamController) GetAll(ctx *fiber.Ctx) error {
	var data []model.UserExamModel
	if err := c.DB.Find(&data).Error; err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve data",
			"error":   err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"data": data,
	})
}

// Get user_exam by ID
func (c *UserExamController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var data model.UserExamModel
	if err := c.DB.First(&data, id).Error; err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "User exam not found",
			"error":   err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"data": data,
	})
}
