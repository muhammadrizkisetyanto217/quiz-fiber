package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type EvaluationsQuestionModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	QuestionText    string         `gorm:"size:200;not null" json:"question_text"`
	QuestionAnswer  pq.StringArray `gorm:"type:text[];not null" json:"question_answer"`
	QuestionCorrect string         `gorm:"size:50;not null" json:"question_correct"`
	Status          string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status"`
	ParagraphHelp   string         `gorm:"type:text;not null" json:"paragraph_help"`
	ExplainQuestion string         `gorm:"type:text;not null" json:"explain_question"`
	AnswerText      string         `gorm:"type:text;not null" json:"answer_text"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	EvaluationID    uint           `gorm:"not null;index" json:"evaluation_id"`
}

// TableName memastikan nama tabel sesuai dengan skema database
func (EvaluationsQuestionModel) TableName() string {
	return "evaluations_questions"
}
