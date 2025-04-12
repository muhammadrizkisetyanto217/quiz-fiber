package service

import (
	"errors"
	userUnitModel "quiz-fiber/internals/features/category/units/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateUserUnitFromEvaluation(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	var userUnit userUnitModel.UserUnitModel
	result := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit)

	// Jika belum ada, buat dengan AttemptEvaluation = 1
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		userUnit = userUnitModel.UserUnitModel{
			UserID:            userID,
			UnitID:            unitID,
			AttemptEvaluation: 1,
		}
		return db.Create(&userUnit).Error
	} else if result.Error != nil {
		return result.Error
	}

	// Jika sudah ada, tambah AttemptEvaluation +1
	return db.Model(&userUnit).
		UpdateColumn("attempt_evaluation", gorm.Expr("attempt_evaluation + ?", 1)).Error
}


func CheckAndUnsetEvaluationStatus(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	var count int64
	err := db.Table("user_evaluations").
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Update("attempt_evaluation", 0).Error
	}

	return nil
}