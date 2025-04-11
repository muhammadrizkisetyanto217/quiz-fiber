package service

import (
	"errors"

	userUnitModel "quiz-fiber/internals/features/category/units/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//////////////////////////////////////////////////////////
// === BAGIAN UNTUK USER READING ===
//////////////////////////////////////////////////////////

func UpdateUserUnitFromReading(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	var userUnit userUnitModel.UserUnitModel
	result := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit)

	// Jika belum ada, buat dengan AttemptReading = 1
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		userUnit = userUnitModel.UserUnitModel{
			UserID:         userID,
			UnitID:         unitID,
			AttemptReading: 1,
		}
		return db.Create(&userUnit).Error
	} else if result.Error != nil {
		return result.Error
	}

	// Jika sudah ada, tambah +1
	return db.Model(&userUnit).
		UpdateColumn("attempt_reading", gorm.Expr("attempt_reading + ?", 1)).Error
}



func CheckAndUnsetUserUnitReadingStatus(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	var count int64
	err := db.Table("user_readings").
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Update("attempt_reading", 0).Error
	}

	return nil
}
