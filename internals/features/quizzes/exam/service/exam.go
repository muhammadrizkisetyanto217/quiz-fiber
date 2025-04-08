package service

import (
	"errors"

	userUnitModel "quiz-fiber/internals/features/category/units/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateUserUnitFromExam(db *gorm.DB, userID uuid.UUID, examID uint, grade int) error {
	var unitID uint

	err := db.Table("exams").
		Select("unit_id").
		Where("id = ?", examID).
		Scan(&unitID).Error
	if err != nil || unitID == 0 {
		return err
	}

	var userUnit userUnitModel.UserUnitModel
	result := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		userUnit = userUnitModel.UserUnitModel{
			UserID:    userID,
			UnitID:    unitID,
			GradeExam: grade,
		}
		return db.Create(&userUnit).Error
	} else if result.Error != nil {
		return result.Error
	}

	// Update GradeExam-nya
	return db.Model(&userUnit).Update("grade_exam", grade).Error
}


func CheckAndUnsetExamStatus(db *gorm.DB, userID uuid.UUID, examID uint) error {
	var unitID uint
	err := db.Table("exams").
		Select("unit_id").
		Where("id = ?", examID).
		Scan(&unitID).Error
	if err != nil || unitID == 0 {
		return err
	}

	var count int64
	err = db.Table("user_exams").
		Joins("JOIN exams ON exams.id = user_exams.exam_id").
		Where("user_exams.user_id = ? AND exams.unit_id = ?", userID, unitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	// Jika tidak ada lagi exam → kosongkan GradeExam
	if count == 0 {
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Update("grade_exam", 0).Error
	}

	return nil
}
