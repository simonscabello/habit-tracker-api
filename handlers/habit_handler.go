package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/simonscabello/habit-tracker-api/database"
	"github.com/simonscabello/habit-tracker-api/middleware"
	"github.com/simonscabello/habit-tracker-api/models"
)

func CreateHabit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	var habit models.Habit
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if habit.Name == "" {
		http.Error(w, "O campo 'name' é obrigatório", http.StatusBadRequest)
		return
	}

	habit.UserID = userID

	if err := database.DB.Create(&habit).Error; err != nil {
		http.Error(w, "Erro ao salvar hábito", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)
}

func GetHabits(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	var habits []models.Habit
	if err := database.DB.Where("user_id = ?", userID).Find(&habits).Error; err != nil {
		http.Error(w, "Erro ao buscar hábitos", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	today := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	type HabitWithStatus struct {
		ID          uint      `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		DoneToday   bool      `json:"done_today"`
	}

	var response []HabitWithStatus

	for _, habit := range habits {
		var log models.HabitLog
		err := database.DB.
			Where("habit_id = ? AND user_id = ? AND DATE(completed_at) = ?", habit.ID, userID, today).
			First(&log).Error

		response = append(response, HabitWithStatus{
			ID:          habit.ID,
			Name:        habit.Name,
			Description: habit.Description,
			CreatedAt:   habit.CreatedAt,
			UpdatedAt:   habit.UpdatedAt,
			DoneToday:   err == nil,
		})
	}

	json.NewEncoder(w).Encode(response)
}

func GetHabitByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	id := chi.URLParam(r, "id")

	var habit models.Habit
	if err := database.DB.First(&habit, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		http.Error(w, "Hábito não encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(habit)
}

func UpdateHabit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	id := chi.URLParam(r, "id")

	var habit models.Habit
	if err := database.DB.First(&habit, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		http.Error(w, "Hábito não encontrado", http.StatusNotFound)
		return
	}

	var updated models.Habit
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	habit.Name = updated.Name
	habit.Description = updated.Description

	if err := database.DB.Save(&habit).Error; err != nil {
		http.Error(w, "Erro ao atualizar hábito", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(habit)
}

func DeleteHabit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	id := chi.URLParam(r, "id")

	var habit models.Habit
	if err := database.DB.First(&habit, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		http.Error(w, "Hábito não encontrado", http.StatusNotFound)
		return
	}

	if err := database.DB.Delete(&habit).Error; err != nil {
		http.Error(w, "Erro ao deletar hábito", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LogHabit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	habitID := chi.URLParam(r, "id")

	var habit models.Habit
	if err := database.DB.First(&habit, "id = ? AND user_id = ?", habitID, userID).Error; err != nil {
		http.Error(w, "Hábito não encontrado", http.StatusNotFound)
		return
	}

	// Data de hoje com hora zero e fuso horário local
	now := time.Now().In(time.Local)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	// Verifica se já existe log para hoje
	var existing models.HabitLog
	err := database.DB.
		Where("habit_id = ? AND user_id = ? AND DATE(completed_at) = ?", habitID, userID, today).
		First(&existing).Error

	if err == nil {
		// Já existe
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Hábito já registrado hoje",
		})
		return
	}

	log := models.HabitLog{
		HabitID:     habit.ID,
		UserID:      userID,
		CompletedAt: today,
	}

	if err := database.DB.Create(&log).Error; err != nil {
		http.Error(w, "Erro ao registrar hábito", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hábito registrado com sucesso para hoje",
	})
}
