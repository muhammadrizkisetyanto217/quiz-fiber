package model

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuizQuestionModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	QuestionText    string         `gorm:"type:varchar(200);not null" json:"question_text"`   // ✅ Ubah menjadi string biasa
	QuestionAnswer  pq.StringArray `gorm:"type:text[];not null" json:"question_answer"`       // ✅ Tetap dalam format array TEXT[]
	QuestionCorrect string         `gorm:"type:varchar(50);not null" json:"question_correct"` // ✅ Tambahkan jawaban yang benar
	TooltipsID      pq.Int64Array  `gorm:"type:int[]" json:"tooltips_id"`                     // 🔥 Array integer untuk menyimpan ID tooltips
	Status          string         `gorm:"type:varchar(10);not null;default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status"`
	ParagraphHelp   string         `gorm:"type:text;not null" json:"paragraph_help"`
	ExplainQuestion string         `gorm:"type:text;not null" json:"explain_question"`
	AnswerText      string         `gorm:"type:text;not null" json:"answer_text"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	QuizzesID       uint           `gorm:"not null;index" json:"quizzes_id"`
}

// TableName mengatur nama tabel agar sesuai dengan skema database
func (QuizQuestionModel) TableName() string {
	return "quizzes_questions"
}

func (q QuizQuestionModel) MarshalJSONQuizzes() ([]byte, error) {
	type Alias QuizQuestionModel
	return json.Marshal(&struct {
		TooltipsID []int64 `json:"tooltips_id"`
		*Alias
	}{
		TooltipsID: []int64(q.TooltipsID), // 🔥 Konversi `pq.Int64Array` ke `[]int64`
		Alias:      (*Alias)(&q),
	})
}
