package model

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
	// TotalQuizzesSection int            `gorm:"default:0" json:"total_quizzes_section"`
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

// Hook AfterSave untuk memperbarui total_unit di ThemesOrLevelModel setelah insert/update
func (u *UnitModel) AfterSave(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE themes_or_levels
		SET total_unit = (
			SELECT COUNT(*) FROM units WHERE themes_or_level_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, u.ThemesOrLevelID, u.ThemesOrLevelID).Error
	return
}

// Hook AfterDelete untuk memperbarui total_unit di ThemesOrLevelModel setelah delete
func (u *UnitModel) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE themes_or_levels
		SET total_unit = (
			SELECT COUNT(*) FROM units WHERE themes_or_level_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, u.ThemesOrLevelID, u.ThemesOrLevelID).Error
	return
}
