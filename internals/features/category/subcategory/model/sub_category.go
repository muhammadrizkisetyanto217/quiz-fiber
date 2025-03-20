package category

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type SubcategoryModel struct {
	ID                           uint           `json:"id" gorm:"primaryKey"`
	Name                         string         `json:"name"`
	Status                       string         `json:"status" gorm:"type:VARCHAR(10);check:status IN ('active', 'pending', 'archived')"`
	DescriptionLong              string         `json:"description_long"`
	TotalThemesOrLevels          int            `json:"total_themes_or_levels"`
	UpdateNews                   json.RawMessage `json:"update_news" gorm:"type:jsonb"`
	ImageURL                     string         `json:"image_url"`
	CreatedAt                    time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt                    time.Time      `json:"updated_at"`
	DeletedAt                    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CategoriesID                 int            `json:"categories_id"`
}

func (SubcategoryModel) TableName() string {
	return "subcategories"
}