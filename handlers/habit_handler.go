package handlers

import (
	"encoding/json"
	"net/http"

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

	json.NewEncoder(w).Encode(habits)
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
