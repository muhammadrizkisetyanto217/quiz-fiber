package category

import (
	"time"

)

type ThemesOrLevelModel struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `json:"name"`
	Status           string         `gorm:"type:VARCHAR(10);check:status IN ('active', 'pending', 'archived')" json:"status"`
	DescriptionShort string         `json:"description_short"`
	DescriptionLong  string         `json:"description_long"`
	TotalUnit        int            `json:"total_unit"`
	UpdateNews       string         `gorm:"type:jsonb" json:"update_news"`
	ImageURL         string         `json:"image_url"`
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at,omitempty"`
	DeletedAt        *time.Time     `gorm:"index" json:"deleted_at,omitempty"`
	SubcategoriesID int `gorm:"column:subcategories_id" json:"subcategories_id"`
}

func (ThemesOrLevelModel) TableName() string {
	return "themes_or_levels"
}