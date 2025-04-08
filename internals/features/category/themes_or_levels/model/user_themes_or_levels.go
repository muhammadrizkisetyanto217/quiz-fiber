package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Fungsi untuk memperbarui UserSectionQuizzesModel
// UserSectionQuizzesModel menyimpan daftar kuis yang telah diselesaikan dalam suatu section
type UserThemesOrLevelsModel struct {
	ID               uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	ThemesOrLevelsID uint          `gorm:"column:themes_or_levels_id;not null" json:"themes_or_levels_id"`
	CompleteUnit     pq.Int64Array `gorm:"type:integer[]" json:"complete_unit"`
	TotalUnit        int           `gorm:"default:0" json:"total_unit"`
	GradeResult      float64       `json:"grade_result"`
	CreatedAt        time.Time     `gorm:"default:current_timestamp" json:"created_at"`
}

func (UserThemesOrLevelsModel) TableName() string {
	return "user_themes_or_levels"
}
