package model

import (
    "time"

    "github.com/google/uuid"
)

type KajianAttendance struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
    Latitude   float64   `json:"latitude"`
    Longitude  float64   `json:"longitude"`
    Address    string    `json:"address"`
    AccessTime time.Time `json:"access_time"`
    Topic      string    `json:"topic"`   // originally "materi"
    Notes      string    `json:"notes"`   // originally "lainnya"
    CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (KajianAttendance) TableName() string {
	return "kajian_attendances"
}