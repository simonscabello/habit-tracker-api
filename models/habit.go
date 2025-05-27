package models

import "gorm.io/gorm"

type Habit struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      uint   `json:"-"`
}
