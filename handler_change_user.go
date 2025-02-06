package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/PajdekPL/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerChangeEmailPassword(w http.ResponseWriter, r *http.Request) {
	type UserData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type RespondData struct {
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"UpdatedAt"`
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bad token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bad token", err)
		return
	}
	data := UserData{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem decoding data", err)
		return
	}
	if data.Email == "" || data.Password == "" {
		respondWithError(w, http.StatusUnauthorized, "Incorrect new email or password", err)
		return
	}
	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "New password hashing problem", err)
		return
	}
	userDb, err := cfg.db.ChangeEmailAndPassword(context.Background(), database.ChangeEmailAndPasswordParams{
		ID:             userID,
		Email:          data.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}
	respondWithJSON(w, http.StatusOK, RespondData{
		Email:     userDb.Email,
		UpdatedAt: userDb.UpdatedAt,
		CreatedAt: userDb.CreatedAt,
	})
}
