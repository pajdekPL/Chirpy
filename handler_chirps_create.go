package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/PajdekPL/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID                 uuid.UUID `json:"id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	UserID             uuid.UUID `json:"user_id"`
	Body               string    `json:"body"`
	AuthorName         string    `json:"author_name"`
	ExpirationDatetime time.Time `json:"expiration_datetime"`
}

type ChirpDataCreate struct {
	Body               string    `json:"body" example:"My super chirp!"`
	ExpirationDatetime time.Time `json:"expiration_datetime" example:"2023-12-31T23:59:59Z"`
}

// @Summary      Create a new chirp
// @Description  Creates a new chirp for the authenticated user. The chirp body must be less than 140 characters and cannot contain certain bad words.
// @Tags         chirps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chirp  body      ChirpDataCreate  true  "Chirp object"
// @Success      201    {object}  Chirp
// @Failure      400    {object}  ErrorResponse
// @Failure      401  	{object}  ErrorResponse
// @Failure      404  	{object}  ErrorResponse
// @Failure      500  	{object}  ErrorResponse
// @Router       /chirps [post]
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := ChirpDataCreate{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:               cleaned,
		ExpirationDatetime: params.ExpirationDatetime,
		UserID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:                 chirp.ID,
		CreatedAt:          chirp.CreatedAt,
		UpdatedAt:          chirp.UpdatedAt,
		Body:               chirp.Body,
		UserID:             chirp.UserID,
		ExpirationDatetime: chirp.ExpirationDatetime,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
