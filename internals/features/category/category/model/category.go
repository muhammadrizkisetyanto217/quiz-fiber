package model

import (
	"time"

	"gorm.io/gorm"
)

type CategoryModel struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Name             string    `json:"name"`
	DescriptionShort string    `json:"description_short"`
	DescriptionLong  string    `json:"description_long"`
	TotalSubcategories  int       `json:"total_subcategories"`
	Status           string    `json:"status" gorm:"type:varchar(10);check:status IN ('active', 'pending', 'archived')"`
	CreatedAt        time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `json:"updated_at"`
	DifficultyID     uint      `json:"difficulty_id"`
}



func (CategoryModel) TableName() string {
	return "categories"
}

// Hook AfterSave untuk memperbarui total_categories di DifficultyModel setelah insert/update kategori
func (c *CategoryModel) AfterSave(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE difficulties
		SET total_categories = (
			SELECT COUNT(*) FROM categories WHERE difficulty_id = ?
		)
		WHERE id = ?
	`, c.DifficultyID, c.DifficultyID).Error
	return
}

// Hook AfterDelete untuk memperbarui total_categories di DifficultyModel setelah delete kategori
func (c *CategoryModel) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Exec(`
		UPDATE difficulties
		SET total_categories = (
			SELECT COUNT(*) FROM categories WHERE difficulty_id = ?
		)
		WHERE id = ?
	`, c.DifficultyID, c.DifficultyID).Error
	return
}