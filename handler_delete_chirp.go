package main

import (
	"context"
	"net/http"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "wrong chirp id - should be with uuid format", err)
		return
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
	chirp, err := cfg.db.GetChirp(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp doesn't exist", err)
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "action forbidden", err)
		return
	}
	err = cfg.db.DeleteChirp(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, "")
}
