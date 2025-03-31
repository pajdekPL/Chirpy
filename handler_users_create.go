package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/PajdekPL/Chirpy/internal/database"
	"github.com/google/uuid"
)

type UserDataCreate struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type ReturnDataCreateUser struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

// @Summary      Create new user
// @Description  Creates a new user account with email, password and name
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userData  body      UserDataCreate  true  "User registration data"
// @Success      201       {object}  ReturnDataCreateUser
// @Failure      401  	   {object}  ErrorResponse
// @Failure      404  	   {object}  ErrorResponse
// @Failure      500  	   {object}  ErrorResponse
// @Router       /users [post]
func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	data := UserDataCreate{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem decoding data", err)
		return
	}
	if data.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Empty email", err)
		return
	}
	if data.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Empty password", err)
		return
	}
	if data.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Empty name", err)
		return
	}
	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem creating user", err)
		return
	}
	newUserData := database.CreateUserParams{
		Email:          data.Email,
		HashedPassword: hashedPassword,
		UserName:       data.Name,
	}
	user, err := cfg.db.CreateUser(context.Background(), newUserData)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			respondWithError(w, http.StatusBadRequest, "User already exists with the given email or name", err)
			return
		}
		respondWithError(w, http.StatusBadRequest, "Problem creating user", err)
		return
	}
	responseData := ReturnDataCreateUser{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		Name:        user.UserName,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, responseData)
}
