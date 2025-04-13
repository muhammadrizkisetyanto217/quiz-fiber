package services

import (
	"errors"
	"log"

	// userCategoryModel "quiz-fiber/internals/features/category/category/model"
	// userSubcategoryModel "quiz-fiber/internals/features/category/subcategory/model"
	// userThemesOrLevelsModel "quiz-fiber/internals/features/category/themes_or_levels/model"
	userUnitModel "quiz-fiber/internals/features/category/units/model"
	quizzesModel "quiz-fiber/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func UpdateUserSectionIfQuizCompleted(db *gorm.DB, userID uuid.UUID, sectionID uint, quizID uint) error {
	log.Println("[SERVICE] UpdateUserSectionIfQuizCompleted - userID:", userID, "sectionID:", sectionID, "quizID:", quizID)

	// Hitung semua quiz yang aktif di section tersebut sebagai total
	var allQuizzes []quizzesModel.QuizModel
	if err := db.Where("section_quizzes_id = ? AND deleted_at IS NULL", sectionID).Find(&allQuizzes).Error; err != nil {
		log.Println("[ERROR] Failed to fetch quizzes for section:", err)
		return err
	}

	totalQuizIDs := pq.Int64Array{}
	for _, quiz := range allQuizzes {
		totalQuizIDs = append(totalQuizIDs, int64(quiz.ID))
	}

	var userSection quizzesModel.UserSectionQuizzesModel
	err := db.Where("user_id = ? AND section_quizzes_id = ?", userID, sectionID).First(&userSection).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userSection = quizzesModel.UserSectionQuizzesModel{
			UserID:           userID,
			SectionQuizzesID: sectionID,
			CompleteQuiz:     pq.Int64Array{int64(quizID)},
			TotalQuiz:        totalQuizIDs,
		}
		log.Println("[SERVICE] Creating new UserSectionQuizzesModel")
		return db.Create(&userSection).Error
	}

	// Cek apakah quizID sudah tercatat
	for _, id := range userSection.CompleteQuiz {
		if id == int64(quizID) {
			log.Println("[SERVICE] Quiz ID already recorded, skipping update")
			return nil
		}
	}

	userSection.CompleteQuiz = append(userSection.CompleteQuiz, int64(quizID))
	userSection.TotalQuiz = totalQuizIDs

	log.Println("[SERVICE] Updating existing UserSectionQuizzesModel")
	return db.Save(&userSection).Error
}

func UpdateUserUnitIfSectionCompleted(db *gorm.DB, userID uuid.UUID, unitID uint, sectionID uint) error {
	log.Println("[SERVICE] UpdateUserUnitIfSectionCompleted - userID:", userID, "unitID:", unitID, "sectionID:", sectionID)

	// Ambil semua section dari unit (untuk TotalSectionQuizzes)
	var allSections []quizzesModel.SectionQuizzesModel
	if err := db.Where("unit_id = ? AND deleted_at IS NULL", unitID).Find(&allSections).Error; err != nil {
		log.Println("[ERROR] Failed to fetch sections for unit:", err)
		return err
	}
	totalSectionIDs := pq.Int64Array{}
	for _, section := range allSections {
		totalSectionIDs = append(totalSectionIDs, int64(section.ID))
	}

	// Cek apakah semua quiz dalam sectionID sudah diselesaikan
	var userSection quizzesModel.UserSectionQuizzesModel
	err := db.Where("user_id = ? AND section_quizzes_id = ?", userID, sectionID).First(&userSection).Error
	if err != nil {
		log.Println("[ERROR] UserSectionQuizzesModel not found, skipping unit update")
		return nil
	}

	if len(userSection.CompleteQuiz) < len(userSection.TotalQuiz) {
		log.Println("[INFO] Section belum selesai, tidak update ke UserUnitModel")
		return nil
	}

	// Ambil user_unit
	var userUnit userUnitModel.UserUnitModel
	err = db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error
	if err != nil {
		log.Printf("[WARNING] user_unit belum tersedia untuk user_id=%s, unit_id=%d\n", userID.String(), unitID)
		return nil // tidak buat baru
	}

	// Jika sectionID belum masuk → tambahkan
	for _, id := range userUnit.CompleteSectionQuizzes {
		if id == int64(sectionID) {
			log.Println("[SERVICE] Section ID already recorded in UserUnitModel")
			return nil
		}
	}

	userUnit.CompleteSectionQuizzes = append(userUnit.CompleteSectionQuizzes, int64(sectionID))
	userUnit.TotalSectionQuizzes = totalSectionIDs

	log.Println("[SERVICE] Updating existing UserUnitModel")
	return db.Save(&userUnit).Error
}
