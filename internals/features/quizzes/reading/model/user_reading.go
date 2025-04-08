package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quiz-fiber/internals/features/quizzes/reading/service"
)

type UserReading struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"not null" json:"user_id"`
	ReadingID uint           `gorm:"not null" json:"reading_id"`
	IsReading bool           `gorm:"default:false" json:"is_reading"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserReading) TableName() string {
	return "user_readings"
}

// Hook: Saat create → update user_units.is_reading = true
func (u *UserReading) AfterCreate(tx *gorm.DB) error {
	return service.UpdateUserUnitFromReading(tx, u.UserID, u.ReadingID)
}

// Hook: Saat delete → jika tidak ada lagi reading, unset user_units.is_reading
func (u *UserReading) AfterDelete(tx *gorm.DB) error {
	return service.CheckAndUnsetUserUnitReadingStatus(tx, u.UserID, u.ReadingID)
}
