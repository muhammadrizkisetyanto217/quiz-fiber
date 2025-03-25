package model

import (
	"time"

	"quiz-fiber/internals/features/quizzes/quiz_section/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UnitModel struct {
	ID                  uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                string         `gorm:"type:varchar(50);unique;not null" json:"name"`
	Status              string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active','pending','archived')" json:"status"`
	DescriptionShort    string         `gorm:"type:varchar(200);not null" json:"description_short"`
	DescriptionOverview string         `gorm:"type:text;not null" json:"description_overview"`
	ImageURL            string         `gorm:"type:varchar(100)" json:"image_url"`
	UpdateNews          datatypes.JSON `gorm:"type:jsonb" json:"update_news"`
	CreatedAt           time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ThemesOrLevelID     uint           `gorm:"not null" json:"themes_or_level_id"`
	CreatedBy           uuid.UUID      `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`

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
