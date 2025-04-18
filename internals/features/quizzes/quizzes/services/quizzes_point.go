package services

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userPointLog "quiz-fiber/internals/features/progress/point/model"
	updateUserProgressTotalService "quiz-fiber/internals/features/progress/progress/service"
)

func AddPointFromQuiz(db *gorm.DB, userID uuid.UUID, quizID uint, attempt int) error {
	log.Println("[SERVICE] AddPointFromQuiz - userID:", userID, "quizID:", quizID, "attempt:", attempt)

	// Hitung poin berdasarkan attempt
	var point int
	switch attempt {
	case 1:
		point = 20
	case 2:
		point = 40
	default:
		point = 10
	}

	const sourceTypeQuiz = 1 // ✅ quiz = 1

	pointLog := userPointLog.UserPointLog{
		UserID:     userID,
		Points:     point,
		SourceType: sourceTypeQuiz,
		SourceID:   int(quizID),
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&pointLog).Error; err != nil {
		log.Println("[ERROR] Gagal insert user_point_log (quiz):", err)
		return err
	}

	// ✅ Tambahkan update total poin ke user_progress
	if err := updateUserProgressTotalService.UpdateUserProgressTotal(db, userID); err != nil {
		log.Println("[WARNING] Gagal update user_progress:", err)
	}

	log.Printf("[SUCCESS] Poin quiz attempt %d ditambahkan: %d poin", attempt, point)
	return nil
}
