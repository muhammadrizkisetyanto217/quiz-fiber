package sectionquizzes

import (
	"encoding/json"
	"log"
	"os"

	"quiz-fiber/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SectionQuizSeed struct {
	NameQuizzes      string  `json:"name_quizzes"`
	Status           string  `json:"status"`
	MaterialsQuizzes string  `json:"materials_quizzes"`
	IconURL          string  `json:"icon_url"`
	UnitID           uint    `json:"unit_id"`
	CreatedBy        string  `json:"created_by"`
	TotalQuizzes     []int64 `json:"total_quizzes"`
}

func SeedSectionQuizzesFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []SectionQuizSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, seed := range seeds {
		var existing model.SectionQuizzesModel
		if err := db.Where("name_quizzes = ?", seed.NameQuizzes).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Section Quiz '%s' sudah ada, lewati...", seed.NameQuizzes)
			continue
		}

		newSection := model.SectionQuizzesModel{
			NameQuizzes:      seed.NameQuizzes,
			Status:           seed.Status,
			MaterialsQuizzes: seed.MaterialsQuizzes,
			IconURL:          seed.IconURL,
			TotalQuizzes:     pq.Int64Array(seed.TotalQuizzes),
			UnitID:           seed.UnitID,
			CreatedBy:        parseUUID(seed.CreatedBy),
		}

		if err := db.Create(&newSection).Error; err != nil {
			log.Printf("❌ Gagal insert '%s': %v", seed.NameQuizzes, err)
		} else {
			log.Printf("✅ Berhasil insert '%s'", seed.NameQuizzes)
		}
	}
}

// helper mandiri
func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		log.Fatalf("❌ Gagal parse UUID: %v", err)
	}
	return id
}
