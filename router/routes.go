package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/simonscabello/habit-tracker-api/handlers"
	"github.com/simonscabello/habit-tracker-api/middleware"
)

func SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.CorsMiddleware)

	// Rotas públicas
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Habit Tracker API is running!"))
	})

	r.Post("/register", handlers.RegisterUser)
	r.Post("/login", handlers.LoginUser)

	// Rota protegida: /me
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Get("/me", handlers.GetMe)
		r.Get("/progress", handlers.GetWeeklyProgress)
	})

	// Rotas protegidas: /habits
	r.Route("/habits", func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		r.Get("/", handlers.GetHabits)
		r.Post("/", handlers.CreateHabit)
		r.Get("/{id}", handlers.GetHabitByID)
		r.Put("/{id}", handlers.UpdateHabit)
		r.Delete("/{id}", handlers.DeleteHabit)
		r.Post("/{id}/log", handlers.LogHabit)

	})

	return r
}
