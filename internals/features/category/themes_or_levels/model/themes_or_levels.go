package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ThemesOrLevelsModel struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string         `gorm:"type:varchar(255)" json:"name"`
	Status           string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active','pending','archived')" json:"status"`
	DescriptionShort string         `gorm:"type:varchar(100)" json:"description_short"`
	DescriptionLong  string         `gorm:"type:varchar(2000)" json:"description_long"`
	TotalUnit        int            `json:"total_unit"`
	ImageURL         string         `gorm:"type:varchar(100)" json:"image_url"`
	UpdateNews       datatypes.JSON `gorm:"type:jsonb" json:"update_news"`
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	SubcategoriesID  int            `gorm:"column:subcategories_id" json:"subcategories_id"`
}

func (ThemesOrLevelsModel) TableName() string {
	return "themes_or_levels"
}

// Hook AfterSave untuk memperbarui total_themes_or_levels di SubcategoryModel setelah insert/update
func (t *ThemesOrLevelsModel) AfterSave(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE subcategories
		SET total_themes_or_levels = (
			SELECT COUNT(*) FROM themes_or_levels WHERE subcategories_id = ?
		)
		WHERE id = ?
	`, t.SubcategoriesID, t.SubcategoriesID).Error
	return
}

// Hook AfterDelete untuk memperbarui total_themes_or_levels di SubcategoryModel setelah delete
func (t *ThemesOrLevelsModel) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE subcategories
		SET total_themes_or_levels = (
			SELECT COUNT(*) FROM themes_or_levels WHERE subcategories_id = ?
		)
		WHERE id = ?
	`, t.SubcategoriesID, t.SubcategoriesID).Error
	return
}
