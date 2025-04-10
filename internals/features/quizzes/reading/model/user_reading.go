package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quiz-fiber/internals/features/quizzes/reading/service"
)

type UserReading struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	ReadingID uint      `gorm:"not null" json:"reading_id"`
	UnitID    uint      `gorm:"not null" json:"unit_id"` // ➕ Tambahan
	Attempt   int       `gorm:"default:1;not null" json:"attempt"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserReading) TableName() string {
	return "user_readings"
}

func (u *UserReading) BeforeCreate(tx *gorm.DB) error {
	log.Println("👉 BeforeCreate dipanggil untuk:", u.UserID, u.ReadingID)

	var latestAttempt int
	err := tx.Table("user_readings").
		Select("COALESCE(MAX(attempt), 0)").
		Where("user_id = ? AND reading_id = ? ", u.UserID, u.ReadingID).
		Scan(&latestAttempt).Error

	if err != nil {
		return err
	}

	u.Attempt = latestAttempt + 1
	log.Println("🎯 Attempt terbaru:", u.Attempt)

	return nil
}

// Hook: Saat create → update user_units.is_reading = true
func (u *UserReading) AfterCreate(tx *gorm.DB) error {
	return service.UpdateUserUnitFromReading(tx, u.UserID, u.UnitID)
}

func (u *UserReading) AfterDelete(tx *gorm.DB) error {
	return service.CheckAndUnsetUserUnitReadingStatus(tx, u.UserID, u.UnitID)
}
