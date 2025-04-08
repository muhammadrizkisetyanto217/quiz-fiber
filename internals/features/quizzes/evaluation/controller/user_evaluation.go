package controller


import (
	
	userEvaluationModel "quiz-fiber/internals/features/quizzes/evaluation/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserEvaluationController struct {
	DB *gorm.DB
}

func NewUserEvaluationController(db *gorm.DB) *UserEvaluationController {
	return &UserEvaluationController{DB: db}
}


// POST /api/user_evaluations
func (ctrl *UserEvaluationController) Create(c *fiber.Ctx) error {
	var input userEvaluationModel.UserEvaluationModel
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := ctrl.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user evaluation"})
	}

	return c.Status(fiber.StatusCreated).JSON(input)
}

// GET /api/user_evaluations/:user_id
func (ctrl *UserEvaluationController) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var evaluations []userEvaluationModel.UserEvaluationModel

	if err := ctrl.DB.Where("user_id = ?", userID).Find(&evaluations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get evaluations"})
	}

	return c.JSON(evaluations)
}