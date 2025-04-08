package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserUnitModel struct {
	ID                     uint          `gorm:"primaryKey" json:"id"`
	UserID                 uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	UnitID                 uint           `gorm:"not null" json:"unit_id"`
	IsReading              bool          `gorm:"default:false" json:"is_reading"`
	IsEvaluation           bool          `gorm:"default:false" json:"is_evaluation"`
	CompleteSectionQuizzes pq.Int64Array `gorm:"type:integer[];default:'{}'" json:"complete_section_quizzes"`
	TotalSectionQuizzes    pq.Int64Array `gorm:"type:integer[];default:'{}'" json:"total_section_quizzes"`
	GradeExam              int           `gorm:"default:0" json:"grade_exam"`
	GradeResult            int           `gorm:"default:0" json:"grade_result"`
	CreatedAt              time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName untuk override nama tabel default
func (UserUnitModel) TableName() string {
	return "user_unit"
}
