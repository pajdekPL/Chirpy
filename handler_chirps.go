package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/PajdekPL/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}

	chirpCensoredText, err := validateAndCensorChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}

	chirpDb := database.CreateChirpParams{
		UserID: userID,
		Body:   chirpCensoredText,
	}

	chirp, err := cfg.db.CreateChirp(context.Background(), chirpDb)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with saving chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpsResponse := []Chirp{}
	chirps, err := cfg.db.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with getting chirps", err)
		return
	}
	authorID := r.URL.Query().Get("author_id")
	sortDirection := "asc"
	sortDirectionParam := r.URL.Query().Get("sort")
	if sortDirectionParam == "desc" {
		sortDirection = "desc"
	}

	authorUUID := uuid.Nil
	if authorID != "" {
		authorUUID, err = uuid.Parse(authorID); if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author_id query parameter", err)
		return
	} 
	}

	for _, chirp := range chirps {
		if authorUUID != uuid.Nil && chirp.UserID != authorUUID {
			continue
		}
		chirpsResponse = append(chirpsResponse, Chirp{
			ID:        chirp.ID,
			Body:      chirp.Body,
			UpdatedAt: chirp.UpdatedAt,
			CreatedAt: chirp.CreatedAt,
			UserID:    chirp.UserID,
		})
	}
	if sortDirection == "desc"{
		slices.SortFunc(chirpsResponse, func(a Chirp, b Chirp) int {
			return b.UpdatedAt.Compare(a.UpdatedAt)
	
		})
	}
	// sort.Slice(chirps, func(i, j int) bool {
	// 	if sortDirection == "desc" {
	// 		return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
	// 	}
	// 	return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	// })
	respondWithJSON(w, http.StatusOK, chirpsResponse)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "wrong id - should be with uuid format", err)
		return
	}
	chirp, err := cfg.db.GetChirp(context.Background(), id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		Body:      chirp.Body,
		UpdatedAt: chirp.UpdatedAt,
		CreatedAt: chirp.CreatedAt,
		UserID:    chirp.UserID,
	})
}

func validateAndCensorChirp(chirpText string) (string, error) {
	maxChirpLength := 140
	if len(chirpText) > maxChirpLength {
		return "", fmt.Errorf("Chirp too long")

	}
	return getCensoredBody(chirpText), nil
}

func getCensoredBody(chirp string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(chirp, " ")

	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
