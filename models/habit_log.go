package models

import (
	"time"

	"gorm.io/gorm"
)

type HabitLog struct {
	gorm.Model
	HabitID     uint      `json:"habit_id"`
	UserID      uint      `json:"-"`
	CompletedAt time.Time `json:"completed_at"`
}
