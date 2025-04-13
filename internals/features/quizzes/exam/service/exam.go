package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	userUnitModel "quiz-fiber/internals/features/category/units/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateUserUnitFromExam(db *gorm.DB, userID uuid.UUID, examID uint, grade int) error {
	if grade < 0 || grade > 100 {
		return fmt.Errorf("nilai grade tidak valid: %d", grade)
	}

	var unitID uint
	err := db.Table("exams").
		Select("unit_id").
		Where("id = ?", examID).
		Scan(&unitID).Error
	if err != nil || unitID == 0 {
		log.Println("[ERROR] Gagal ambil unit_id dari exam_id:", examID)
		return err
	}

	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("[WARNING] user_unit tidak ditemukan saat UpdateUserUnitFromExam, user_id:", userID, "unit_id:", unitID)
		}
		return err
	}

	// Hitung tambahan poin dari aktivitas selain exam
	activityBonus := 0
	if userUnit.AttemptReading > 0 {
		activityBonus += 5
	}
	if userUnit.AttemptEvaluation > 0 {
		activityBonus += 15
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
		activityBonus += 30
	}

	// Final grade_result
	var gradeResult int
	if activityBonus == 0 {
		gradeResult = grade / 2
	} else {
		gradeResult = activityBonus + (grade / 2)
	}

	updates := map[string]interface{}{
		"grade_result": gradeResult,
		"is_passed":    gradeResult > 65,
		"updated_at":   time.Now(),
	}

	if grade > userUnit.GradeExam {
		updates["grade_exam"] = grade
	}

	return db.Model(&userUnit).Updates(updates).Error
}

// ✅ Final: CheckAndUnsetExamStatus
func CheckAndUnsetExamStatus(db *gorm.DB, userID uuid.UUID, examID uint) error {
	var unitID uint
	err := db.Table("exams").
		Select("unit_id").
		Where("id = ?", examID).
		Scan(&unitID).Error
	if err != nil || unitID == 0 {
		log.Println("[ERROR] Gagal ambil unit_id dari exam_id untuk reset status:", examID)
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
		log.Println("[INFO] Reset nilai exam dan result karena tidak ada user_exams tersisa, user_id:", userID, "unit_id:", unitID)
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Updates(map[string]interface{}{
				"grade_exam":   0,
				"grade_result": 0,
				"updated_at":   time.Now(),
			}).Error
	}

	return nil
}
