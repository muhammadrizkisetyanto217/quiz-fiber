package services

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	userUnitModel "quiz-fiber/internals/features/category/units/model"
	quizzesModel "quiz-fiber/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QuizProgress struct {
	ID      uint `json:"id"`
	Attempt int  `json:"attempt"`
	Score   int  `json:"score"`
}

type SectionProgress struct {
	ID      uint `json:"id"`
	Score   int  `json:"score"`
	Attempt int  `json:"attempt"`
}

func UpdateUserSectionIfQuizCompleted(
	db *gorm.DB,
	userID uuid.UUID,
	sectionID uint,
	quizID uint,
	attempt int,
	percentageGrade int,
) error {
	log.Println("[SERVICE] UpdateUserSectionIfQuizCompleted - userID:", userID, "sectionID:", sectionID, "quizID:", quizID, "attempt:", attempt, "score:", percentageGrade)

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
	newProgress := []QuizProgress{{ID: quizID, Attempt: attempt, Score: percentageGrade}}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		progressJSON, _ := json.Marshal(newProgress)
		userSection = quizzesModel.UserSectionQuizzesModel{
			UserID:           userID,
			SectionQuizzesID: sectionID,
			CompleteQuiz:     datatypes.JSON(progressJSON),
			TotalQuiz:        totalQuizIDs,
			GradeResult:      percentageGrade,
		}
		log.Println("[SERVICE] Creating new UserSectionQuizzesModel")
		return db.Create(&userSection).Error
	}

	var progressList []QuizProgress
	if err := json.Unmarshal(userSection.CompleteQuiz, &progressList); err != nil {
		log.Println("[ERROR] Failed to parse existing complete_quiz:", err)
		return err
	}

	found := false
	for i, p := range progressList {
		if p.ID == quizID {
			if attempt > p.Attempt {
				progressList[i].Attempt = attempt
			}
			if percentageGrade > p.Score {
				progressList[i].Score = percentageGrade
			}
			found = true
			break
		}
	}
	if !found {
		progressList = append(progressList, QuizProgress{
			ID:      quizID,
			Attempt: attempt,
			Score:   percentageGrade,
		})
	}

	completedQuizIDs := map[uint]bool{}
	totalScore := 0
	for _, p := range progressList {
		completedQuizIDs[p.ID] = true
		totalScore += p.Score
	}
	if len(completedQuizIDs) == len(totalQuizIDs) && len(progressList) > 0 {
		userSection.GradeResult = totalScore / len(progressList)
		log.Println("[SERVICE] Semua quiz selesai - GradeResult:", userSection.GradeResult)
	}

	newJSON, _ := json.Marshal(progressList)
	userSection.CompleteQuiz = datatypes.JSON(newJSON)
	userSection.TotalQuiz = totalQuizIDs

	log.Println("[SERVICE] Updating UserSectionQuizzesModel")
	return db.Save(&userSection).Error
}

func UpdateUserUnitIfSectionCompleted(db *gorm.DB, userID uuid.UUID, unitID uint, sectionID uint) error {
	log.Println("[SERVICE] UpdateUserUnitIfSectionCompleted - userID:", userID, "unitID:", unitID, "sectionID:", sectionID)

	var allSections []quizzesModel.SectionQuizzesModel
	if err := db.Where("unit_id = ? AND deleted_at IS NULL", unitID).Find(&allSections).Error; err != nil {
		log.Println("[ERROR] Failed to fetch sections for unit:", err)
		return err
	}
	totalSectionIDs := pq.Int64Array{}
	for _, section := range allSections {
		totalSectionIDs = append(totalSectionIDs, int64(section.ID))
	}

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

	var userUnit userUnitModel.UserUnitModel
	err = db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error
	if err != nil {
		log.Printf("[WARNING] user_unit belum tersedia untuk user_id=%s, unit_id=%d\n", userID.String(), unitID)
		return nil
	}

	var progressList []SectionProgress
	_ = json.Unmarshal(userUnit.CompleteSectionQuizzes, &progressList)

	found := false
	for _, p := range progressList {
		if p.ID == sectionID {
			found = true
			break
		}
	}
	if found {
		log.Println("[SERVICE] Section ID already recorded in UserUnitModel")
		return nil
	}

	progressList = append(progressList, SectionProgress{
		ID:      sectionID,
		Score:   userSection.GradeResult,
		Attempt: 1,
	})

	// ✅ Hitung GradeQuiz jika semua section selesai
	var gradeQuiz int
	if len(progressList) == len(totalSectionIDs) {
		totalScore := 0
		for _, p := range progressList {
			totalScore += p.Score
		}
		gradeQuiz = totalScore / len(progressList)
		log.Printf("[SERVICE] Semua section complete. GradeQuiz di-set ke %d\n", gradeQuiz)
	}

	progressJSON, _ := json.Marshal(progressList)
	userUnit.CompleteSectionQuizzes = datatypes.JSON(progressJSON)
	userUnit.TotalSectionQuizzes = totalSectionIDs
	userUnit.GradeQuiz = gradeQuiz

	log.Println("[SERVICE] Updating existing UserUnitModel")
	return db.Model(&userUnit).Updates(map[string]interface{}{
		"complete_section_quizzes": userUnit.CompleteSectionQuizzes,
		"total_section_quizzes":    totalSectionIDs,
		"grade_quiz":               gradeQuiz,
		"updated_at":               time.Now(),
	}).Error
}
