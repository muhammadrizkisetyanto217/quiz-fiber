package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SectionQuizzesModel struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	NameQuizzes      string         `gorm:"size:50;not null" json:"name_quizzes"`
	Status           string         `gorm:"size:10;default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status"`
	MaterialsQuizzes string         `gorm:"type:text;not null" json:"materials_quizzes"`
	IconURL          string         `gorm:"size:100" json:"icon_url"`
	TotalQuizzes     int            `gorm:"default:0" json:"total_quizzes"` // <--- Tambahan disini
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy        uuid.UUID      `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`
	UnitID           uint           `gorm:"not null;constraint:OnDelete:CASCADE" json:"unit_id"`
	Quizzes          []QuizModel    `gorm:"foreignKey:SectionQuizID" json:"quizzes"`
}

func (SectionQuizzesModel) TableName() string {
	return "section_quizzes"
}

// ✅ AfterCreate Hook
func (s *SectionQuizzesModel) AfterCreate(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE units
		SET total_section_quizzes = total_section_quizzes + 1
		WHERE id = ?
	`, s.UnitID).Error
	return
}

// ✅ AfterDelete Hook
func (s *SectionQuizzesModel) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE units
		SET total_section_quizzes = GREATEST(total_section_quizzes - 1, 0)
		WHERE id = ?
	`, s.UnitID).Error
	return
}
