package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/simonscabello/habit-tracker-api/database"
	"github.com/simonscabello/habit-tracker-api/middleware"
	"github.com/simonscabello/habit-tracker-api/models"
)

func GetWeeklyProgress(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	location := time.Local
	now := time.Now().In(location)
	weekday := int(now.Weekday())

	sunday := time.Date(now.Year(), now.Month(), now.Day()-weekday, 0, 0, 0, 0, location)
	saturdayEnd := sunday.AddDate(0, 0, 7)

	progress := make([]int, 7)

	var logs []models.HabitLog
	err := database.DB.
		Where("user_id = ? AND completed_at >= ? AND completed_at < ?", userID, sunday, saturdayEnd).
		Find(&logs).Error

	if err != nil {
		http.Error(w, "Erro ao buscar progresso", http.StatusInternalServerError)
		return
	}

	for _, log := range logs {
		weekday := log.CompletedAt.In(location).Weekday()
		progress[int(weekday)]++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}
