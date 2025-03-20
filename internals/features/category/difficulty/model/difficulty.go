package model

import (
	"time"

)

type DifficultyModel struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"size:255;not null" json:"name"`
	DescriptionShort  string    `gorm:"size:100" json:"description_short"`
	DescriptionLong   string    `gorm:"size:2000" json:"description_long"`
	TotalCategories   int       `json:"total_categories"`
	Status            string    `gorm:"size:10;default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}


func (DifficultyModel) TableName() string {
	return "difficulties"
}