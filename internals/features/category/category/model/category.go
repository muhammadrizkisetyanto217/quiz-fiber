package model

import (
	"time"
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