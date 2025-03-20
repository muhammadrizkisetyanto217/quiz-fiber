package category

import (
	"quiz-fiber/internals/features/quizzes/quiz_section/model"
	"time"

	"gorm.io/gorm"
)

type UnitModel struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Name                string         `gorm:"unique;not null" json:"name"`
	Status              string         `gorm:"type:varchar(10);default:'pending'" json:"status"`
	DescriptionShort    string         `gorm:"type:varchar(200);not null" json:"description_short"`
	DescriptionOverview string         `gorm:"type:text;not null" json:"description_overview"`
	TotalQuizzesSection int            `gorm:"default:0" json:"total_quizzes_section"` 
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ThemesOrLevelID     uint           `json:"themes_or_level_id"`
	CreatedBy           uint           `json:"created_by"`

	SectionQuizzes []model.SectionQuizzesModel `gorm:"foreignKey:UnitID" json:"section_quizzes"`
}

func (UnitModel) TableName() string {
	return "units"
}