package service

import (
	"errors"
	"log"
	"time"
	"quiz-fiber/internals/features/progress/progress/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



func UpdateUserProgressTotal(db *gorm.DB, userID uuid.UUID) error {
	var total int64
	err := db.Table("user_point_logs").
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(points), 0)").
		Scan(&total).Error
	if err != nil {
		log.Println("[ERROR] Gagal hitung total poin:", err)
		return err
	}

	var progress model.UserProgress
	if err := db.Where("user_id = ?", userID).First(&progress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Buat baru
			progress = model.UserProgress{
				UserID:      userID,
				TotalPoints: int(total),
				LastUpdated: time.Now(),
			}
			return db.Create(&progress).Error
		}
		return err
	}

	// Update existing
	progress.TotalPoints = int(total)
	progress.LastUpdated = time.Now()
	return db.Save(&progress).Error
}
