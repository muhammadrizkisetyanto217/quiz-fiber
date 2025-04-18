package service

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userPointLog "quiz-fiber/internals/features/progress/point/model"
	updateUserProgressTotalService "quiz-fiber/internals/features/progress/progress/service"
)

func AddPointFromExam(db *gorm.DB, userID uuid.UUID, examID uint, attempt int) error {
	log.Println("[SERVICE] AddPointFromExam - userID:", userID, "examID:", examID, "attempt:", attempt)

	// Hitung poin berdasarkan attempt ke-n
	var point int
	switch attempt {
	case 1:
		point = 20
	case 2:
		point = 40
	default:
		point = 10
	}

	const sourceTypeExam = 3 // ✅ exam = 3

	pointLog := userPointLog.UserPointLog{
		UserID:     userID,
		Points:     point,
		SourceType: sourceTypeExam,
		SourceID:   int(examID),
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&pointLog).Error; err != nil {
		log.Println("[ERROR] Gagal insert user_point_log (exam):", err)
		return err
	}

	// ✅ Tambahkan update total poin ke user_progress
	if err := updateUserProgressTotalService.UpdateUserProgressTotal(db, userID); err != nil {
		log.Println("[WARNING] Gagal update user_progress:", err)
	}

	log.Printf("[SUCCESS] Poin exam attempt %d ditambahkan: %d poin", attempt, point)
	return nil
}
