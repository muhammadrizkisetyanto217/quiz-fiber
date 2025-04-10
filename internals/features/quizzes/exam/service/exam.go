package service

import (
	"errors"
	"log"

	userUnitModel "quiz-fiber/internals/features/category/units/model"
	UpdateUserThemesOrLevelsIfUnitCompleted "quiz-fiber/internals/features/quizzes/quizzes/services"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ✅ Final: UpdateUserUnitFromExam
// func UpdateUserUnitFromExam(db *gorm.DB, userID uuid.UUID, examID uint, grade int) error {
// 	var unitID uint

// 	err := db.Table("exams").
// 		Select("unit_id").
// 		Where("id = ?", examID).
// 		Scan(&unitID).Error
// 	if err != nil || unitID == 0 {
// 		return err
// 	}

// 	var userUnit userUnitModel.UserUnitModel
// 	result := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit)

// 	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		userUnit = userUnitModel.UserUnitModel{
// 			UserID:      userID,
// 			UnitID:      unitID,
// 			GradeExam:   grade,
// 			GradeResult: grade,
// 		}
// 		err := db.Create(&userUnit).Error
// 		if err != nil {
// 			return err
// 		}
// 	} else if result.Error != nil {
// 		return result.Error
// 	} else {
// 		updates := map[string]interface{}{
// 			"grade_exam":   grade,
// 			"grade_result": grade,
// 		}
// 		if err := db.Model(&userUnit).Updates(updates).Error; err != nil {
// 			return err
// 		}
// 	}

// 	if grade > 65 {
// 		// Ambil themes_or_levels_id dari unit
// 		var themesOrLevelsID uint
// 		if err := db.Table("units").
// 			Select("themes_or_level_id").
// 			Where("id = ?", unitID).
// 			Scan(&themesOrLevelsID).Error; err != nil {
// 			log.Println("[ERROR] Failed to fetch themes_or_levels_id:", err)
// 		} else if themesOrLevelsID != 0 {
// 			// ⬇️ Tambahkan log ini di sini
// 			log.Println("[DEBUG] Trigger update themes_or_levels for user:", userID, "unitID:", unitID, "themesOrLevelsID:", themesOrLevelsID)

// 			if err := UpdateUserThemesOrLevelsIfUnitCompleted.UpdateUserThemesOrLevelsIfUnitCompleted(
// 				db, userID, int(unitID), int(themesOrLevelsID),
// 			); err != nil {
// 				log.Println("[ERROR] Failed to update themes_or_levels:", err)
// 			}
// 		} else {
// 			log.Println("[WARNING] themesOrLevelsID = 0, tidak akan update themes progress")
// 		}
// 	}

//		return nil
//	}
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
	var gradeResult int = 0

	// 1. Cek is_reading
	var isReading bool
	_ = db.Table("user_readings").
		Select("true").
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		Scan(&isReading).Error
	if isReading {
		gradeResult += 5
	}

	// 2. Cek is_evaluation
	var isEvaluation bool
	_ = db.Table("user_evaluations").
		Select("true").
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		Scan(&isEvaluation).Error
	if isEvaluation {
		gradeResult += 15
	}

	// 3. Cek apakah semua section_quizzes sudah dikerjakan
	var totalSections, completedSections int64
	_ = db.Table("section_quizzes").
		Where("unit_id = ?", unitID).
		Count(&totalSections).Error

	_ = db.Table("user_section_quizzes").
		Joins("JOIN section_quizzes ON user_section_quizzes.section_quiz_id = section_quizzes.id").
		Where("user_section_quizzes.user_id = ? AND section_quizzes.unit_id = ?", userID, unitID).
		Count(&completedSections).Error

	if totalSections > 0 && totalSections == completedSections {
		gradeResult += 30
	}

	// 4. Tambahkan 50 jika grade_exam == 100
	if grade == 100 {
		gradeResult += 50
	}

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		userUnit = userUnitModel.UserUnitModel{
			UserID:      userID,
			UnitID:      unitID,
			GradeExam:   grade,
			GradeResult: gradeResult,
		}
		err := db.Create(&userUnit).Error
		if err != nil {
			return err
		}
	} else if result.Error != nil {
		return result.Error
	} else {
		updates := map[string]interface{}{
			"grade_exam":   grade,
			"grade_result": gradeResult,
		}
		if err := db.Model(&userUnit).Updates(updates).Error; err != nil {
			return err
		}
	}

	if grade > 65 {
		var themesOrLevelsID uint
		if err := db.Table("units").
			Select("themes_or_level_id").
			Where("id = ?", unitID).
			Scan(&themesOrLevelsID).Error; err != nil {
			log.Println("[ERROR] Failed to fetch themes_or_levels_id:", err)
		} else if themesOrLevelsID != 0 {
			log.Println("[DEBUG] Trigger update themes_or_levels for user:", userID, "unitID:", unitID, "themesOrLevelsID:", themesOrLevelsID)

			if err := UpdateUserThemesOrLevelsIfUnitCompleted.UpdateUserThemesOrLevelsIfUnitCompleted(
				db, userID, int(unitID), int(themesOrLevelsID),
			); err != nil {
				log.Println("[ERROR] Failed to update themes_or_levels:", err)
			}
		} else {
			log.Println("[WARNING] themesOrLevelsID = 0, tidak akan update themes progress")
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
