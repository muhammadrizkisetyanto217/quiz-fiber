package model

import (
	"time"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamModel struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	NameExams     string         `gorm:"size:50;not null" json:"name_exams" validate:"required,max=50"`
	Status        string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status" validate:"required,oneof=active pending archived"`
	Point         int            `gorm:"not null;default:30" json:"point" validate:"gte=0"`
	TotalQuestion *int           `json:"total_question" validate:"omitempty,gte=0"`
	IconURL       *string        `gorm:"size:100" json:"icon_url,omitempty" validate:"omitempty,url"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UnitID        uint           `json:"unit_id" validate:"required"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`
}

// TableName mengatur nama tabel agar sesuai dengan skema database
func (ExamModel) TableName() string {
	return "exams"
}

func (e *ExamModel) Validate() error {
	if e.NameExams == "" || len(e.NameExams) > 50 {
		return errors.New("name_exams is required and must be less than or equal to 50 characters")
	}
	if e.Status != "active" && e.Status != "pending" && e.Status != "archived" {
		return errors.New("status must be one of 'active', 'pending', or 'archived'")
	}
	if e.Point < 0 {
		return errors.New("point must be greater than or equal to 0")
	}
	if e.TotalQuestion != nil && *e.TotalQuestion < 0 {
		return errors.New("total_question must be greater than or equal to 0")
	}
	if e.UnitID == 0 {
		return errors.New("unit_id is required")
	}
	return nil
}
