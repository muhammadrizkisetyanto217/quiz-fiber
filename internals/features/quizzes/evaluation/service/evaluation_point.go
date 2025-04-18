package service

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userPointLog "quiz-fiber/internals/features/progress/point/model"
	updateUserProgressTotalService "quiz-fiber/internals/features/progress/progress/service"
)

func AddPointFromEvaluation(db *gorm.DB, userID uuid.UUID, evaluationID uint, attempt int) error {
	log.Println("[SERVICE] AddPointFromEvaluation - userID:", userID, "evaluationID:", evaluationID, "attempt:", attempt)

	// Hitung poin berdasarkan attempt
	var point int
	switch attempt {
	case 1:
		point = 25
	case 2:
		point = 15
	default:
		point = 10
	}

	const sourceTypeEvaluation = 2

	pointLog := userPointLog.UserPointLog{
		UserID:     userID,
		Points:     point,
		SourceType: sourceTypeEvaluation,
		SourceID:   int(evaluationID),
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&pointLog).Error; err != nil {
		log.Println("[ERROR] Gagal insert user_point_log (evaluation):", err)
		return err
	}

	// ✅ Tambahkan update total poin ke user_progress
	if err := updateUserProgressTotalService.UpdateUserProgressTotal(db, userID); err != nil {
		log.Println("[WARNING] Gagal update user_progress:", err)
	}

	log.Printf("[SUCCESS] Poin evaluation attempt %d ditambahkan: %d poin", attempt, point)
	return nil
}
