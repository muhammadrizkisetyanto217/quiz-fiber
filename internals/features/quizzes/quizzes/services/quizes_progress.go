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
	type QuizCompletion struct {
		ID      int `json:"id"`
		Score   int `json:"score"`
		Attempt int `json:"attempt"`
	}

	// Step 1: Ambil semua section_quizzes dalam unit
	var sections []quizzesModel.SectionQuizzesModel
	if err := db.Where("unit_id = ? AND deleted_at IS NULL", unitID).Find(&sections).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil section_quizzes: %v", err)
		return err
	}

	// Step 2: Cek apakah semua section sudah complete oleh user
	for _, section := range sections {
		var userSection quizzesModel.UserSectionQuizzesModel
		if err := db.Where("user_id = ? AND section_quizzes_id = ?", userID, section.ID).
			First(&userSection).Error; err != nil {
			log.Printf("[INFO] Section %d belum ada progress oleh user", section.ID)
			return nil // Belum lengkap
		}

		// Langsung pakai TotalQuiz dari model
		totalQuizIDs := userSection.TotalQuiz

		// Unmarshal complete_quiz
		var completedQuizData []QuizCompletion
		if err := json.Unmarshal(userSection.CompleteQuiz, &completedQuizData); err != nil {
			log.Printf("[ERROR] Gagal decode complete_quiz: %v", err)
			return err
		}

		// Cek apakah semua quiz sudah complete
		completedIDs := map[int]bool{}
		for _, q := range completedQuizData {
			completedIDs[q.ID] = true
		}

		for _, id := range totalQuizIDs {
			if !completedIDs[int(id)] {
				log.Printf("[INFO] Section %d belum lengkap, quiz ID %d belum dikerjakan", section.ID, id)
				return nil
			}
		}
	}

	// Step 3: Update user_unit jika semua section dalam unit sudah complete
	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_id = ? AND unit_id = ?", userID, unitID).
		First(&userUnit).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil user_unit: %v", err)
		return err
	}

	// Tambahkan sectionID ke complete_section_quizzes jika belum ada
	// Decode datatypes.JSON ke slice int64
	var completeSectionIDs []int64
	if len(userUnit.CompleteSectionQuizzes) > 0 {
		if err := json.Unmarshal(userUnit.CompleteSectionQuizzes, &completeSectionIDs); err != nil {
			log.Printf("[ERROR] Gagal decode complete_section_quizzes: %v", err)
			return err
		}
	}

	// Cek apakah sectionID sudah ada
	found := false
	for _, sid := range completeSectionIDs {
		if uint(sid) == sectionID {
			found = true
			break
		}
	}

	// Tambahkan jika belum ada
	if !found {
		completeSectionIDs = append(completeSectionIDs, int64(sectionID))

		// Encode ulang ke JSON
		updatedJSON, err := json.Marshal(completeSectionIDs)
		if err != nil {
			log.Printf("[ERROR] Gagal encode complete_section_quizzes: %v", err)
			return err
		}
		userUnit.CompleteSectionQuizzes = updatedJSON
	}

	userUnit.UpdatedAt = time.Now()
	if err := db.Save(&userUnit).Error; err != nil {
		log.Printf("[ERROR] Gagal update user_unit: %v", err)
		return err
	}

	log.Printf("[SERVICE] UserUnit berhasil diperbarui - userID: %s, unitID: %d, sectionID: %d", userID, unitID, sectionID)
	return nil
}
