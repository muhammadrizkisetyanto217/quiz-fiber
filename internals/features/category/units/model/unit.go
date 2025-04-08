package model

import (
	"time"

	"quiz-fiber/internals/features/quizzes/quizzes/model"

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
	TotalSectionQuizzes int            `gorm:"default:0" json:"total_section_quizzes"` // ✅ Kolom tambahan
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
// ✅ AfterCreate Hook untuk UnitModel
func (u *UnitModel) AfterCreate(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE themes_or_levels
		SET total_unit = total_unit + 1
		WHERE id = ?
	`, u.ThemesOrLevelID).Error
	return
}

// ✅ AfterDelete Hook untuk UnitModel
func (u *UnitModel) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE themes_or_levels
		SET total_unit = GREATEST(total_unit - 1, 0)
		WHERE id = ?
	`, u.ThemesOrLevelID).Error
	return
}
