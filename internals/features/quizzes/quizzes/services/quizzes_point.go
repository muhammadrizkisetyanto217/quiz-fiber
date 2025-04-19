package services

import (
	addUserPointService "quiz-fiber/internals/features/progress/point/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddPointFromQuiz(db *gorm.DB, userID uuid.UUID, quizID uint, attempt int) error {
	var point int
	switch attempt {
	case 1:
		point = 20
	case 2:
		point = 40
	default:
		point = 10
	}
	return addUserPointService.AddUserPointLogAndUpdateProgress(db, userID, 1, int(quizID), point)
}
