package model

import (
	"time"

	"quiz-fiber/internals/features/quizzes/evaluation/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserEvaluationModel struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uuid.UUID `gorm:"not null" json:"user_id"`
	EvaluationID    uint      `gorm:"not null;column:evaluation_id" json:"evaluation_id"`
	UnitID          uint      `gorm:"not null;column:unit_id" json:"unit_id"`
	Attempt         int       `gorm:"default:1;not null" json:"attempt"`
	PercentageGrade int       `gorm:"default:0;not null" json:"percentage_grade"`
	TimeDuration    int       `gorm:"default:0;not null" json:"time_duration"`
	Point           int       `gorm:"default:0;not null" json:"point"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserEvaluationModel) TableName() string {
	return "user_evaluations"
}

func (u *UserEvaluationModel) AfterCreate(tx *gorm.DB) error {
	return service.UpdateUserUnitFromEvaluation(tx, u.UserID, u.UnitID)
}

func (u *UserEvaluationModel) AfterDelete(tx *gorm.DB) error {
	return service.CheckAndUnsetEvaluationStatus(tx, u.UserID, u.UnitID)
}
