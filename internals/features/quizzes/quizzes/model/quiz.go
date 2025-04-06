package model

import (
	"time"

	"github.com/google/uuid"
)

type QuizModel struct {
	ID            int        `json:"id" gorm:"primaryKey"`
	Name          string     `json:"name_quizzes" gorm:"type:varchar(50);unique;not null;column:name_quizzes"`
	Status        string     `json:"status" gorm:"type:varchar(10);default:pending;check:status IN ('active', 'pending', 'archived')"`
	TotalQuestion int        `json:"total_question"`
	IconURL       string     `json:"icon_url" gorm:"type:varchar(100)"`
	CreatedAt     time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`

	SectionQuizID int       `json:"section_quizzes_id" gorm:"column:section_quizzes_id"`
	UnitID        int       `json:"unit_id"`
	CreatedBy     uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`
}

// TableName memastikan Gorm menggunakan tabel "quizzes"
func (QuizModel) TableName() string {
	return "quizzes"
}
