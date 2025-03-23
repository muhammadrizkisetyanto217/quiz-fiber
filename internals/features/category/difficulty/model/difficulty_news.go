package model

import (
	"time"

	"gorm.io/gorm"
)

type DifficultyNews struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"`
	Title        string         `gorm:"type:varchar(255);not null"`
	Description  string         `gorm:"type:text;not null"`
	IsPublic     bool           `gorm:"default:true"`
	CreatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	DifficultyID uint           `gorm:"not null"`
}

func (DifficultyNews) TableName() string {
	return "difficulty_news"
}
