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

	// ✳️ Hitung grade_result berdasarkan aktivitas
	var gradeResult int

	if userUnit.AttemptReading > 0 {
		gradeResult += 5
	}

	if userUnit.AttemptEvaluation > 0 {
		gradeResult += 15
	}

	var totalSections, completedSections int64
	_ = db.Table("section_quizzes").
		Where("unit_id = ?", unitID).
		Count(&totalSections).Error

	_ = db.Table("user_section_quizzes").
		Joins("JOIN section_quizzes ON user_section_quizzes.section_quizzes_id = section_quizzes.id").
		Where("user_section_quizzes.user_id = ? AND section_quizzes.unit_id = ?", userID, unitID).
		Count(&completedSections).Error

	if totalSections > 0 && totalSections == completedSections {
		gradeResult += 30
	}

	gradeResult += grade / 2

	// Insert atau Update user_unit
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		userUnit = userUnitModel.UserUnitModel{
			UserID:      userID,
			UnitID:      unitID,
			GradeExam:   grade,
			GradeResult: gradeResult,
			IsPassed:    gradeResult > 65,
		}
		if err := db.Create(&userUnit).Error; err != nil {
			return err
		}
	} else if result.Error != nil {
		return result.Error
	} else {
		// 🧠 Update hanya grade_exam jika nilai baru lebih tinggi
		updates := map[string]interface{}{
			"grade_result": gradeResult,
			"is_passed":    gradeResult > 65,
		}
		if grade > userUnit.GradeExam {
			updates["grade_exam"] = grade
		}
		if err := db.Model(&userUnit).Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}

// ✅ Final: CheckAndUnsetExamStatus
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

	if count == 0 {
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Updates(map[string]interface{}{
				"grade_exam":   0,
				"grade_result": 0,
			}).Error
	}

	return nil
}
