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

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type UserData struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
	}

	type ReturnUserData struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	data := UserData{}

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
	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem creating user", err)
		return
	}
	newUserData := database.CreateUserParams{
		Email:          data.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.db.CreateUser(context.Background(), newUserData)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			respondWithError(w, http.StatusBadRequest, "User already exists with the given email", err)
			return
		}
		respondWithError(w, http.StatusBadRequest, "Problem creating user", err)
		return
	}
	responseData := ReturnUserData{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, responseData)
}
