package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Fungsi untuk memperbarui UserSectionQuizzesModel
// UserSectionQuizzesModel menyimpan daftar kuis yang telah diselesaikan dalam suatu section
type UserSectionQuizzesModel struct {
	ID               uint          `gorm:"primaryKey" json:"id"`
	UserID           uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	SectionQuizzesID uint           `gorm:"not null" json:"section_quizzes_id"`
	CompleteQuiz     pq.Int64Array `gorm:"type:integer[];default:'{}'" json:"complete_quiz"`
	TotalQuiz        pq.Int64Array `gorm:"type:integer[];default:'{}'" json:"total_quiz"`
	CreatedAt        time.Time     `gorm:"autoCreateTime" json:"created_at"`
}

func (UserSectionQuizzesModel) TableName() string {
	return "user_section_quizzes"
}
