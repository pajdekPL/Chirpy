package main

import (
	"net/http"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/PajdekPL/Chirpy/internal/database"
	"github.com/google/uuid"
)

// @Summary      Get a specific chirp
// @Description  Retrieves a specific chirp by its ID
// @Tags         chirps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Chirp ID"
// @Success      200  {object}  Chirp
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /chirps/{id} [get]
func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:                 dbChirp.ID,
		CreatedAt:          dbChirp.CreatedAt,
		UpdatedAt:          dbChirp.UpdatedAt,
		UserID:             dbChirp.UserID,
		Body:               dbChirp.Body,
		ExpirationDatetime: dbChirp.ExpirationDatetime,
	})
}

// @Summary      Get all chirps
// @Description  Retrieves all chirps in the system
// @Tags         chirps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   Chirp
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /chirps [get]
func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	if filter == "expired" {
		cfg.handlerExpiredChirpsRetrieve(w, r)
		return
	}

	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:                 dbChirp.ID,
			CreatedAt:          dbChirp.CreatedAt,
			UpdatedAt:          dbChirp.UpdatedAt,
			UserID:             dbChirp.UserID,
			Body:               dbChirp.Body,
			AuthorName:         dbChirp.AuthorName,
			ExpirationDatetime: dbChirp.ExpirationDatetime,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerExpiredChirpsRetrieve(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.db.GetExpiredChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:                 dbChirp.ID,
			CreatedAt:          dbChirp.CreatedAt,
			UpdatedAt:          dbChirp.UpdatedAt,
			UserID:             dbChirp.UserID,
			Body:               dbChirp.Body,
			AuthorName:         dbChirp.AuthorName,
			ExpirationDatetime: dbChirp.ExpirationDatetime,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

// @Summary      Get chirps for the currently authenticated user
// @Description  Retrieves all chirps created by the authenticated user
// @Tags         chirps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        filter  query  string  expired  "filter=expired to filter by expired chirps"
// @Success      200  {array}   Chirp
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /user/chirps [get]
func (cfg *apiConfig) handlerChirpsUserRetrieve(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
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
	var userChirps []database.GetChirpsByUserRow
	if filter == "expired" {
		cfg.handlerExpiredChirpsUserRetrieve(w, r, userID)
		return
	}

	userChirps, err = cfg.db.GetChirpsByUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range userChirps {
		chirps = append(chirps, Chirp{
			ID:                 dbChirp.ID,
			CreatedAt:          dbChirp.CreatedAt,
			UpdatedAt:          dbChirp.UpdatedAt,
			UserID:             dbChirp.UserID,
			Body:               dbChirp.Body,
			AuthorName:         dbChirp.AuthorName,
			ExpirationDatetime: dbChirp.ExpirationDatetime,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerExpiredChirpsUserRetrieve(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {

	userChirps, err := cfg.db.GetExpiredChirpsByUser(r.Context(), userID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range userChirps {
		chirps = append(chirps, Chirp{
			ID:                 dbChirp.ID,
			CreatedAt:          dbChirp.CreatedAt,
			UpdatedAt:          dbChirp.UpdatedAt,
			UserID:             dbChirp.UserID,
			Body:               dbChirp.Body,
			AuthorName:         dbChirp.AuthorName,
			ExpirationDatetime: dbChirp.ExpirationDatetime,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
