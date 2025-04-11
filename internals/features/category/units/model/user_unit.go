package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserUnitModel struct {
	ID                     uint          `gorm:"primaryKey" json:"id"`
	UserID                 uuid.UUID     `gorm:"type:uuid;not null;index:idx_user_unit_user_id_unit_id" json:"user_id"`
	UnitID                 uint          `gorm:"not null;index:idx_user_unit_user_id_unit_id" json:"unit_id"`
	AttemptReading         int           `gorm:"default:0;not null" json:"attempt_reading"`
	AttemptEvaluation      int           `gorm:"default:0;not null" json:"attempt_evaluation"`
	CompleteSectionQuizzes pq.Int64Array `gorm:"type:integer[];default:'{}'" json:"complete_section_quizzes"`
	TotalSectionQuizzes    pq.Int64Array `gorm:"type:integer[];default:'{}'" json:"total_section_quizzes"`
	GradeExam              int           `gorm:"default:0" json:"grade_exam"`
	IsPassed               bool          `gorm:"default:false" json:"is_passed"`
	GradeResult            int           `gorm:"default:0" json:"grade_result"`
	CreatedAt              time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName untuk override nama tabel default
func (UserUnitModel) TableName() string {
	return "user_unit"
}
