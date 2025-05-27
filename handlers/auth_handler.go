package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simonscabello/habit-tracker-api/database"
	"github.com/simonscabello/habit-tracker-api/middleware"
	"github.com/simonscabello/habit-tracker-api/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Email == "" || input.Password == "" {
		http.Error(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
		return
	}

	// Verifica se já existe um usuário com esse email
	var existing models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		http.Error(w, "Email já registrado", http.StatusBadRequest)
		return
	}

	// Criptografa a senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao gerar senha", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
		return
	}

	user.Password = "" // remove o hash da resposta
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Senha inválida", http.StatusUnauthorized)
		return
	}

	// Gerar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	user.Password = "" // nunca retornamos a senha
	json.NewEncoder(w).Encode(user)
}
