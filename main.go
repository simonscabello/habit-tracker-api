package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/simonscabello/habit-tracker-api/config"
	"github.com/simonscabello/habit-tracker-api/database"
	"github.com/simonscabello/habit-tracker-api/models"
	"github.com/simonscabello/habit-tracker-api/router"
)

func main() {
	config.Load()
	database.Connect()

	err := database.DB.AutoMigrate(&models.User{}, &models.Habit{})

	if err != nil {
		log.Fatalf("Erro ao migrar tabela: %v", err)
	}

	r := router.SetupRoutes()

	port := fmt.Sprintf(":%s", config.Cfg.AppPort)
	log.Printf("Servidor iniciado na porta %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
