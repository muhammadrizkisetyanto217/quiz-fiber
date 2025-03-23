package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SubcategoryModel struct {
	ID                  uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                string         `json:"name" gorm:"type:varchar(255)"`
	Status              string         `json:"status" gorm:"type:varchar(10);default:'pending';check:status IN ('active','pending','archived')"`
	DescriptionLong     string         `json:"description_long" gorm:"type:varchar(2000)"`
	TotalThemesOrLevels int            `json:"total_themes_or_levels"`
	ImageURL            string         `json:"image_url" gorm:"type:varchar(100)"`
	UpdateNews          datatypes.JSON `json:"update_news"` // pakai JSONB di PostgreSQL
	CreatedAt           time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt           *time.Time     `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CategoriesID        int            `json:"categories_id"`
}

func (SubcategoryModel) TableName() string {
	return "subcategories"
}

// Hook AfterSave untuk memperbarui total_subcategories di CategoryModel setelah insert/update
func (s *SubcategoryModel) AfterSave(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE categories
		SET total_subcategories = (
			SELECT COUNT(*) FROM subcategories WHERE categories_id = ?
		)
		WHERE id = ?
	`, s.CategoriesID, s.CategoriesID).Error
	return
}

// Hook AfterDelete untuk memperbarui total_subcategories di CategoryModel setelah delete
func (s *SubcategoryModel) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE categories
		SET total_subcategories = (
			SELECT COUNT(*) FROM subcategories WHERE categories_id = ?
		)
		WHERE id = ?
	`, s.CategoriesID, s.CategoriesID).Error
	return
}
