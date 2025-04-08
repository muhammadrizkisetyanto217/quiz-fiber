package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Fungsi untuk memperbarui UserSectionQuizzesModel
// UserSectionQuizzesModel menyimpan daftar kuis yang telah diselesaikan dalam suatu section
type UserSubcategoryModel struct {
	ID                     uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                 uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	SubcategoryID          uint          `gorm:"column:subcategory_id;not null" json:"subcategory_id"`
	CompleteThemesOrLevels pq.Int64Array `gorm:"type:integer[]" json:"complete_themes_or_levels"`
	TotalThemesOrLevels    int           `gorm:"default:0" json:"total_themes_or_levels"`
	GradeResult            float64       `json:"grade_result"`
	CreatedAt              time.Time     `gorm:"default:current_timestamp" json:"created_at"`
}

func (UserSubcategoryModel) TableName() string {
	return "user_subcategory"
}
