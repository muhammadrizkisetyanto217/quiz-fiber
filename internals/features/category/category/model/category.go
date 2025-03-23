package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CategoryModel struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	Name               string         `json:"name" gorm:"size:255;not null"`
	Status             string         `json:"status" gorm:"type:varchar(10);default:'pending';check:status IN ('active', 'pending', 'archived')"`
	DescriptionShort   string         `json:"description_short" gorm:"size:100"`
	DescriptionLong    string         `json:"description_long" gorm:"size:2000"`
	TotalSubcategories int            `json:"total_subcategories"`
	ImageURL           string         `json:"image_url" gorm:"size:100"`
	UpdateNews         datatypes.JSON `json:"update_news"`
	DifficultyID       uint           `json:"difficulty_id"`
	CreatedAt          time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`
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
