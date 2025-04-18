package service

import (
	"log"
	"time"

	userPointLog "quiz-fiber/internals/features/progress/point/model"
	updateUserProgressTotalService "quiz-fiber/internals/features/progress/progress/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddPointFromReading(db *gorm.DB, userID uuid.UUID, readingID uint, attempt int) error {
	log.Println("[SERVICE] AddPointFromReading - userID:", userID, "readingID:", readingID, "attempt:", attempt)

	// Hitung poin berdasarkan attempt
	var point int
	switch attempt {
	case 1:
		point = 10
	case 2:
		point = 20
	default:
		point = 5
	}

	const sourceTypeReading = 0 // reading = 0 (int)

	pointLog := userPointLog.UserPointLog{
		UserID:     userID,
		Points:     point,
		SourceType: sourceTypeReading,
		SourceID:   int(readingID),
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&pointLog).Error; err != nil {
		log.Println("[ERROR] Gagal insert user_point_log:", err)
		return err
	}

	log.Printf("[SUCCESS] Poin reading attempt %d ditambahkan: %d poin", attempt, point)

	// ✅ Tambahkan update total poin ke user_progress
	if err := updateUserProgressTotalService.UpdateUserProgressTotal(db, userID); err != nil {
		log.Println("[WARNING] Gagal update user_progress:", err)
	}
	
	return nil
}
