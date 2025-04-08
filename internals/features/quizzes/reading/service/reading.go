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

func UpdateUserUnitFromReading(db *gorm.DB, userID uuid.UUID, readingID uint) error {
	var unitID uint

	err := db.Table("readings").
		Select("unit_id").
		Where("id = ?", readingID).
		Scan(&unitID).Error

	if err != nil {
		return err
	}
	if unitID == 0 {
		return errors.New("unit_id not found for reading_id")
	}

	var userUnit userUnitModel.UserUnitModel
	result := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		userUnit = userUnitModel.UserUnitModel{
			UserID:    userID,
			UnitID:    unitID,
			IsReading: true,
		}
		return db.Create(&userUnit).Error
	} else if result.Error != nil {
		return result.Error
	}

	return db.Model(&userUnit).Update("is_reading", true).Error
}

func CheckAndUnsetUserUnitReadingStatus(db *gorm.DB, userID uuid.UUID, readingID uint) error {
	var unitID uint
	err := db.Table("readings").
		Select("unit_id").
		Where("id = ?", readingID).
		Scan(&unitID).Error
	if err != nil || unitID == 0 {
		return err
	}

	var count int64
	err = db.Table("user_readings").
		Joins("JOIN readings ON readings.id = user_readings.reading_id").
		Where("user_readings.user_id = ? AND readings.unit_id = ?", userID, unitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Update("is_reading", false).Error
	}

	return nil
}
