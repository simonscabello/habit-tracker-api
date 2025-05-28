package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/simonscabello/habit-tracker-api/config"
	"github.com/simonscabello/habit-tracker-api/database"
	"github.com/simonscabello/habit-tracker-api/models"
	"github.com/simonscabello/habit-tracker-api/router"
)

func main() {
	config.Load()
	database.Connect()

	loc, error := time.LoadLocation(os.Getenv("TZ"))
	if error != nil {
		panic("Fuso horário inválido")
	}
	time.Local = loc

	err := database.DB.AutoMigrate(&models.User{}, &models.Habit{}, &models.HabitLog{})

	if err != nil {
		log.Fatalf("Erro ao migrar tabela: %v", err)
	}

	r := router.SetupRoutes()

	port := fmt.Sprintf(":%s", config.Cfg.AppPort)
	log.Printf("Servidor iniciado na porta %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
