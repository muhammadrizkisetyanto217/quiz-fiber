package service

import (
	"log"
	userLogPoint "quiz-fiber/internals/features/progress/point/model"
	userProgress "quiz-fiber/internals/features/progress/progress/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddUserPointLogAndUpdateProgress(db *gorm.DB, userID uuid.UUID, sourceType int, sourceID int, points int) error {
	log.Printf("[SERVICE] AddUserPointLogAndUpdateProgress - userID: %s sourceType: %d sourceID: %d point: %d",
		userID.String(), sourceType, sourceID, points)

	// 1. Simpan ke log
	logEntry := userLogPoint.UserPointLog{
		UserID:     userID,
		Points:     points,
		SourceType: sourceType,
		SourceID:   sourceID,
		CreatedAt:  time.Now(),
	}
	if err := db.Create(&logEntry).Error; err != nil {
		log.Println("[ERROR] Gagal insert user_point_log:", err)
		return err
	}

	// 2. Tambah poin secara efisien ke user_progress
	if err := db.Model(&userProgress.UserProgress{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"total_points": gorm.Expr("total_points + ?", points),
			"last_updated": time.Now(),
		}).Error; err != nil {
		log.Println("[ERROR] Gagal update user_progress:", err)
		return err
	}

	log.Printf("[SUCCESS] Poin ditambahkan: %d poin", points)
	return nil
}
