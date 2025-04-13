package service

import (
	"log"
	userUnitModel "quiz-fiber/internals/features/category/units/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateUserUnitFromEvaluation(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	// Langsung update, tanpa create jika tidak ditemukan
	result := db.Model(&userUnitModel.UserUnitModel{}).
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		UpdateColumn("attempt_evaluation", gorm.Expr("attempt_evaluation + 1"))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// Safety: log jika tidak ditemukan (tidak membuat)
		log.Printf("[WARNING] Tidak ditemukan user_unit untuk user_id: %s, unit_id: %d", userID, unitID)
	}
	return nil
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
