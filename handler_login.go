package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/PajdekPL/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type UserData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type ReturnUserData struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}
	data := UserData{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem decoding data", err)
		return
	}
	if data.Email == "" || data.Password == "" {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	userDb, err := cfg.db.GetUserByEmail(context.Background(), data.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	isValidPassword := auth.CheckPasswordHash(data.Password, userDb.HashedPassword)

	if isValidPassword != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	jwt, err := auth.MakeJWT(
		userDb.ID,
		cfg.jwtSecret,
		time.Hour,
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}

	refreshTokenDb, err := cfg.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userDb.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ReturnUserData{
		ID:           userDb.ID,
		CreatedAt:    userDb.CreatedAt,
		UpdatedAt:    userDb.UpdatedAt,
		Email:        userDb.Email,
		Token:        jwt,
		RefreshToken: refreshTokenDb.Token,
		IsChirpyRed:  userDb.IsChirpyRed,
	})
}

func getDurationFromExpiresInSeconds(expiresInSeconds string) time.Duration {
	seconds, err := strconv.Atoi(expiresInSeconds)
	if err != nil || seconds > 3600 || seconds <= 0 {
		duration, _ := time.ParseDuration("1h")
		return duration
	}
	duration, _ := time.ParseDuration(string(seconds) + "s")
	return duration
}
