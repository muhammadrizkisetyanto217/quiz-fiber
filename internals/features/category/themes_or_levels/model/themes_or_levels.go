package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UpdateNews struct {
	Notes        string `json:"notes"`
	LatestUpdate string `json:"latest_update"`
}

// ✅ Agar bisa disimpan ke jsonb
func (u UpdateNews) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// ✅ Agar bisa diambil dari jsonb
func (u *UpdateNews) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value for UpdateNews")
	}
	return json.Unmarshal(bytes, u)
}

type ThemesOrLevelModel struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Name             string     `json:"name"`
	Status           string     `gorm:"type:VARCHAR(10);check:status IN ('active', 'pending', 'archived')" json:"status"`
	DescriptionShort string     `json:"description_short"`
	DescriptionLong  string     `json:"description_long"`
	TotalUnit        int        `json:"total_unit"`
	UpdateNews       UpdateNews `gorm:"type:jsonb" json:"update_news"`
	ImageURL         string     `json:"image_url"`
	CreatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	DeletedAt        *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	SubcategoriesID  int        `gorm:"column:subcategories_id" json:"subcategories_id"`
}

func (ThemesOrLevelModel) TableName() string {
	return "themes_or_levels"
}

// Hook AfterSave untuk memperbarui total_themes_or_levels di SubcategoryModel setelah insert/update
func (t *ThemesOrLevelModel) AfterSave(tx *gorm.DB) (err error) {
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
func (t *ThemesOrLevelModel) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE subcategories
		SET total_themes_or_levels = (
			SELECT COUNT(*) FROM themes_or_levels WHERE subcategories_id = ?
		)
		WHERE id = ?
	`, t.SubcategoriesID, t.SubcategoriesID).Error
	return
}
