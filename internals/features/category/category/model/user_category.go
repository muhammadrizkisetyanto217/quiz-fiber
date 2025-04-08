package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Fungsi untuk memperbarui UserSectionQuizzesModel
// UserSectionQuizzesModel menyimpan daftar kuis yang telah diselesaikan dalam suatu section
type UserCategoryModel struct {
	ID                  uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID              uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	CategoryID          uint          `gorm:"column:category_id;not null" json:"category_id"`
	CompleteSubcategory pq.Int64Array `gorm:"type:integer[]" json:"complete_category"`
	TotalSubcategory    int           `gorm:"default:0" json:"total_category"`
	GradeResult         float64       `json:"grade_result"`
	CreatedAt           time.Time     `gorm:"default:current_timestamp" json:"created_at"`
}

func (UserCategoryModel) TableName() string {
	return "user_category"
}
